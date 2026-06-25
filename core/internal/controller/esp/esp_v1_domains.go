package esp

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net"
	"strings"
	"time"

	v1 "billionmail-core/api/esp/v1"
	"github.com/gogf/gf/v2/frame/g"
)

func generateDkimKeyPair() (privateKeyPEM string, publicKeyDNS string, err error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	privBytes := x509.MarshalPKCS1PrivateKey(key)
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})
	privateKeyPEM = string(privPEM)

	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return
	}
	publicKeyDNS = base64.StdEncoding.EncodeToString(pubBytes)
	return
}

func buildDnsRecords(domain, dkimSelector, dkimPublicKey string) []*v1.DnsRecordEntry {
	return []*v1.DnsRecordEntry{
		{
			Type:  "TXT",
			Host:  domain,
			Value: "v=spf1 include:relay.yourdomain.com ~all",
			TTL:   "3600",
			Note:  "SPF record — authorizes your relay server to send on behalf of this domain",
		},
		{
			Type:  "TXT",
			Host:  fmt.Sprintf("%s._domainkey.%s", dkimSelector, domain),
			Value: fmt.Sprintf("v=DKIM1; k=rsa; p=%s", dkimPublicKey),
			TTL:   "3600",
			Note:  "DKIM public key — proves emails are signed by you",
		},
		{
			Type:  "TXT",
			Host:  fmt.Sprintf("_dmarc.%s", domain),
			Value: fmt.Sprintf("v=DMARC1; p=none; rua=mailto:dmarc@%s; ruf=mailto:dmarc@%s; fo=1", domain, domain),
			TTL:   "3600",
			Note:  "DMARC policy — start with p=none, upgrade to p=quarantine after confirming alignment",
		},
	}
}

func (c *ControllerV1) AddSendingDomain(ctx context.Context, req *v1.AddSendingDomainReq) (res *v1.AddSendingDomainRes, err error) {
	res = &v1.AddSendingDomainRes{}

	domain := strings.ToLower(strings.TrimSpace(req.Domain))
	selector := "bm1"
	now := time.Now().Unix()

	privKey, pubKey, keyErr := generateDkimKeyPair()
	if keyErr != nil {
		res.Code = 1
		res.Message = "Failed to generate DKIM keys: " + keyErr.Error()
		return
	}

	spfRec := fmt.Sprintf("v=spf1 include:relay.yourdomain.com ~all")
	dmarcRec := fmt.Sprintf("v=DMARC1; p=none; rua=mailto:dmarc@%s; fo=1", domain)

	result, err := g.DB().Model("bm_sending_domains").Data(g.Map{
		"domain":           domain,
		"status":           "pending",
		"dkim_selector":    selector,
		"dkim_private_key": privKey,
		"dkim_public_key":  pubKey,
		"spf_record":       spfRec,
		"dmarc_record":     dmarcRec,
		"spf_verified":     false,
		"dkim_verified":    false,
		"dmarc_verified":   false,
		"created_at":       now,
		"updated_at":       now,
	}).Insert()
	if err != nil {
		res.Code = 1
		res.Message = "Domain already exists or DB error: " + err.Error()
		return
	}

	id, _ := result.LastInsertId()
	d := &v1.SendingDomain{
		Id:            id,
		Domain:        domain,
		Status:        "pending",
		DkimSelector:  selector,
		DkimPublicKey: pubKey,
		SpfRecord:     spfRec,
		DmarcRecord:   dmarcRec,
		CreatedAt:     now,
	}
	res.Data.Domain = d
	res.Data.DnsRecords = buildDnsRecords(domain, selector, pubKey)
	return
}

func (c *ControllerV1) ListSendingDomains(ctx context.Context, req *v1.ListSendingDomainsReq) (res *v1.ListSendingDomainsRes, err error) {
	res = &v1.ListSendingDomainsRes{}
	err = g.DB().Model("bm_sending_domains").
		Fields("id,domain,status,dkim_selector,dkim_public_key,spf_record,dmarc_record,spf_verified,dkim_verified,dmarc_verified,last_checked_at,created_at").
		Scan(&res.Data)
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
	}
	if res.Data == nil {
		res.Data = []*v1.SendingDomain{}
	}
	return
}

func checkTxtRecord(domain, expectedContains string) bool {
	records, err := net.LookupTXT(domain)
	if err != nil {
		return false
	}
	for _, r := range records {
		if strings.Contains(r, expectedContains) {
			return true
		}
	}
	return false
}

func (c *ControllerV1) VerifyDomain(ctx context.Context, req *v1.VerifyDomainReq) (res *v1.VerifyDomainRes, err error) {
	res = &v1.VerifyDomainRes{}

	var d v1.SendingDomain
	err = g.DB().Model("bm_sending_domains").Where("id", req.Id).Scan(&d)
	if err != nil || d.Id == 0 {
		res.Code = 1
		res.Message = "Domain not found"
		return
	}

	spfOk := checkTxtRecord(d.Domain, "v=spf1")
	dkimHost := fmt.Sprintf("%s._domainkey.%s", d.DkimSelector, d.Domain)
	dkimOk := checkTxtRecord(dkimHost, "v=DKIM1")
	dmarcHost := "_dmarc." + d.Domain
	dmarcOk := checkTxtRecord(dmarcHost, "v=DMARC1")

	status := "pending"
	if spfOk && dkimOk && dmarcOk {
		status = "verified"
	} else if spfOk || dkimOk {
		status = "partial"
	}

	now := time.Now().Unix()
	_, err = g.DB().Model("bm_sending_domains").Where("id", req.Id).Data(g.Map{
		"spf_verified":    spfOk,
		"dkim_verified":   dkimOk,
		"dmarc_verified":  dmarcOk,
		"status":          status,
		"last_checked_at": now,
		"updated_at":      now,
	}).Update()
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		return
	}

	d.SpfVerified = spfOk
	d.DkimVerified = dkimOk
	d.DmarcVerified = dmarcOk
	d.Status = status
	d.LastCheckedAt = now

	res.Data.Domain = &d
	res.Data.DnsRecords = buildDnsRecords(d.Domain, d.DkimSelector, d.DkimPublicKey)
	return
}

func (c *ControllerV1) DeleteSendingDomain(ctx context.Context, req *v1.DeleteSendingDomainReq) (res *v1.DeleteSendingDomainRes, err error) {
	res = &v1.DeleteSendingDomainRes{}
	_, err = g.DB().Model("bm_sending_domains").Where("id", req.Id).Delete()
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
	}
	return
}
