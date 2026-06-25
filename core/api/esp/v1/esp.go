package v1

import (
	"billionmail-core/utility/types/api_v1"
	"github.com/gogf/gf/v2/frame/g"
)

// ─── SMTP Credentials ────────────────────────────────────────────────────────

type SmtpCredential struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	DailyLimit  int    `json:"daily_limit"`
	SentToday   int    `json:"sent_today"`
	Plan        string `json:"plan"`
	IpPoolId    int64  `json:"ip_pool_id"`
	IpPoolName  string `json:"ip_pool_name"`
	SmtpHost    string `json:"smtp_host"`
	SmtpPort    int    `json:"smtp_port"`
	CreatedAt   int64  `json:"created_at"`
}

type CreateCredentialReq struct {
	g.Meta        `path:"/esp/credentials/create" method:"post" tags:"ESP" summary:"Create SMTP credential"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Name          string `json:"name" v:"required|max-length:100"`
	Description   string `json:"description" v:"max-length:255"`
	DailyLimit    int    `json:"daily_limit" d:"1000"`
	Plan          string `json:"plan" d:"free" v:"in:free,starter,pro"`
	IpPoolId      int64  `json:"ip_pool_id"`
}
type CreateCredentialRes struct {
	api_v1.StandardRes
	Data *SmtpCredential `json:"data"`
}

type ListCredentialsReq struct {
	g.Meta        `path:"/esp/credentials/list" method:"get" tags:"ESP" summary:"List SMTP credentials"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
}
type ListCredentialsRes struct {
	api_v1.StandardRes
	Data []*SmtpCredential `json:"data"`
}

type DeleteCredentialReq struct {
	g.Meta        `path:"/esp/credentials/delete" method:"post" tags:"ESP" summary:"Delete SMTP credential"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Id            int64  `json:"id" v:"required|min:1"`
}
type DeleteCredentialRes struct {
	api_v1.StandardRes
}

type RegeneratePasswordReq struct {
	g.Meta        `path:"/esp/credentials/regenerate" method:"post" tags:"ESP" summary:"Regenerate credential password"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Id            int64  `json:"id" v:"required|min:1"`
}
type RegeneratePasswordRes struct {
	api_v1.StandardRes
	Data struct {
		Password string `json:"password"`
	} `json:"data"`
}

// ─── Sending Domains ──────────────────────────────────────────────────────────

type SendingDomain struct {
	Id             int64  `json:"id"`
	Domain         string `json:"domain"`
	Status         string `json:"status"`
	DkimSelector   string `json:"dkim_selector"`
	DkimPublicKey  string `json:"dkim_public_key"`
	SpfRecord      string `json:"spf_record"`
	DmarcRecord    string `json:"dmarc_record"`
	SpfVerified    bool   `json:"spf_verified"`
	DkimVerified   bool   `json:"dkim_verified"`
	DmarcVerified  bool   `json:"dmarc_verified"`
	LastCheckedAt  int64  `json:"last_checked_at"`
	CreatedAt      int64  `json:"created_at"`
}

type DnsRecordEntry struct {
	Type  string `json:"type"`
	Host  string `json:"host"`
	Value string `json:"value"`
	TTL   string `json:"ttl"`
	Note  string `json:"note"`
}

type AddSendingDomainReq struct {
	g.Meta        `path:"/esp/domains/add" method:"post" tags:"ESP" summary:"Add sending domain"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Domain        string `json:"domain" v:"required|max-length:255"`
}
type AddSendingDomainRes struct {
	api_v1.StandardRes
	Data struct {
		Domain     *SendingDomain    `json:"domain"`
		DnsRecords []*DnsRecordEntry `json:"dns_records"`
	} `json:"data"`
}

type ListSendingDomainsReq struct {
	g.Meta        `path:"/esp/domains/list" method:"get" tags:"ESP" summary:"List sending domains"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
}
type ListSendingDomainsRes struct {
	api_v1.StandardRes
	Data []*SendingDomain `json:"data"`
}

type VerifyDomainReq struct {
	g.Meta        `path:"/esp/domains/verify" method:"post" tags:"ESP" summary:"Verify domain DNS records"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Id            int64  `json:"id" v:"required|min:1"`
}
type VerifyDomainRes struct {
	api_v1.StandardRes
	Data struct {
		Domain     *SendingDomain    `json:"domain"`
		DnsRecords []*DnsRecordEntry `json:"dns_records"`
	} `json:"data"`
}

type DeleteSendingDomainReq struct {
	g.Meta        `path:"/esp/domains/delete" method:"post" tags:"ESP" summary:"Delete sending domain"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Id            int64  `json:"id" v:"required|min:1"`
}
type DeleteSendingDomainRes struct {
	api_v1.StandardRes
}

// ─── IP Pools ─────────────────────────────────────────────────────────────────

type IpPool struct {
	Id          int64    `json:"id"`
	Name        string   `json:"name"`
	PoolType    string   `json:"pool_type"`
	Ips         []string `json:"ips"`
	Status      string   `json:"status"`
	WarmDay     int      `json:"warm_day"`
	DailyLimit  int      `json:"daily_limit"`
	SentToday   int      `json:"sent_today"`
	Description string   `json:"description"`
	CreatedAt   int64    `json:"created_at"`
}

type CreateIpPoolReq struct {
	g.Meta        `path:"/esp/pools/create" method:"post" tags:"ESP" summary:"Create IP pool"`
	Authorization string   `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Name          string   `json:"name" v:"required|max-length:100"`
	PoolType      string   `json:"pool_type" d:"shared" v:"in:shared,dedicated"`
	Ips           []string `json:"ips"`
	Description   string   `json:"description" v:"max-length:255"`
	DailyLimit    int      `json:"daily_limit" d:"500"`
}
type CreateIpPoolRes struct {
	api_v1.StandardRes
	Data *IpPool `json:"data"`
}

type ListIpPoolsReq struct {
	g.Meta        `path:"/esp/pools/list" method:"get" tags:"ESP" summary:"List IP pools"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
}
type ListIpPoolsRes struct {
	api_v1.StandardRes
	Data []*IpPool `json:"data"`
}

type DeleteIpPoolReq struct {
	g.Meta        `path:"/esp/pools/delete" method:"post" tags:"ESP" summary:"Delete IP pool"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Id            int64  `json:"id" v:"required|min:1"`
}
type DeleteIpPoolRes struct {
	api_v1.StandardRes
}

// ─── Delivery Analytics ───────────────────────────────────────────────────────

type DeliveryStats struct {
	DateStr      string  `json:"date_str"`
	Domain       string  `json:"domain"`
	Sent         int     `json:"sent"`
	Delivered    int     `json:"delivered"`
	BouncedSoft  int     `json:"bounced_soft"`
	BouncedHard  int     `json:"bounced_hard"`
	Complained   int     `json:"complained"`
	Opened       int     `json:"opened"`
	Clicked      int     `json:"clicked"`
	Unsubscribed int     `json:"unsubscribed"`
	DeliveryRate float64 `json:"delivery_rate"`
	OpenRate     float64 `json:"open_rate"`
	BounceRate   float64 `json:"bounce_rate"`
}

type GetDeliveryStatsReq struct {
	g.Meta        `path:"/esp/analytics/stats" method:"get" tags:"ESP" summary:"Get delivery analytics"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	DateFrom      string `json:"date_from"`
	DateTo        string `json:"date_to"`
	Domain        string `json:"domain"`
}
type GetDeliveryStatsRes struct {
	api_v1.StandardRes
	Data struct {
		Summary  *DeliveryStats    `json:"summary"`
		ByDay    []*DeliveryStats  `json:"by_day"`
		ByDomain []*DeliveryStats  `json:"by_domain"`
	} `json:"data"`
}

// ─── Suppression List ─────────────────────────────────────────────────────────

type SuppressionEntry struct {
	Id        int64  `json:"id"`
	Email     string `json:"email"`
	Reason    string `json:"reason"`
	Domain    string `json:"domain"`
	CreatedAt int64  `json:"created_at"`
}

type ListSuppressionReq struct {
	g.Meta        `path:"/esp/suppression/list" method:"get" tags:"ESP" summary:"List suppression entries"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Page          int    `json:"page" d:"1"`
	PageSize      int    `json:"page_size" d:"20"`
	Reason        string `json:"reason"`
	Keyword       string `json:"keyword"`
}
type ListSuppressionRes struct {
	api_v1.StandardRes
	Data struct {
		Total int                 `json:"total"`
		List  []*SuppressionEntry `json:"list"`
	} `json:"data"`
}

type AddSuppressionReq struct {
	g.Meta        `path:"/esp/suppression/add" method:"post" tags:"ESP" summary:"Add to suppression list"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Email         string `json:"email" v:"required|email"`
	Reason        string `json:"reason" d:"manual" v:"in:hard_bounce,soft_bounce,complaint,unsubscribe,manual"`
}
type AddSuppressionRes struct {
	api_v1.StandardRes
}

type DeleteSuppressionReq struct {
	g.Meta        `path:"/esp/suppression/delete" method:"post" tags:"ESP" summary:"Remove from suppression list"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Id            int64  `json:"id" v:"required|min:1"`
}
type DeleteSuppressionRes struct {
	api_v1.StandardRes
}

// ─── MTA Servers ──────────────────────────────────────────────────────────────

type MtaServer struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Host      string `json:"host"`
	SshPort   int    `json:"ssh_port"`
	SshUser   string `json:"ssh_user"`
	Status    string `json:"status"`
	Region    string `json:"region"`
	IpCount   int    `json:"ip_count"`
	CreatedAt int64  `json:"created_at"`
}

type AddMtaServerReq struct {
	g.Meta        `path:"/esp/mta/add" method:"post" tags:"ESP" summary:"Register MTA server"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Name          string `json:"name" v:"required|max-length:100"`
	Host          string `json:"host" v:"required|max-length:255"`
	SshPort       int    `json:"ssh_port" d:"22"`
	SshUser       string `json:"ssh_user" d:"root" v:"max-length:100"`
	Region        string `json:"region" v:"max-length:50"`
}
type AddMtaServerRes struct {
	api_v1.StandardRes
	Data *MtaServer `json:"data"`
}

type ListMtaServersReq struct {
	g.Meta        `path:"/esp/mta/list" method:"get" tags:"ESP" summary:"List MTA servers"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
}
type ListMtaServersRes struct {
	api_v1.StandardRes
	Data []*MtaServer `json:"data"`
}

type DeleteMtaServerReq struct {
	g.Meta        `path:"/esp/mta/delete" method:"post" tags:"ESP" summary:"Delete MTA server"`
	Authorization string `json:"authorization" in:"header" dc:"Authorization" v:"required"`
	Id            int64  `json:"id" v:"required|min:1"`
}
type DeleteMtaServerRes struct {
	api_v1.StandardRes
}
