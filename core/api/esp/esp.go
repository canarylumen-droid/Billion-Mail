package esp

import (
	"context"

	"billionmail-core/api/esp/v1"
)

type IEspV1 interface {
	// SMTP Credentials
	CreateCredential(ctx context.Context, req *v1.CreateCredentialReq) (res *v1.CreateCredentialRes, err error)
	ListCredentials(ctx context.Context, req *v1.ListCredentialsReq) (res *v1.ListCredentialsRes, err error)
	DeleteCredential(ctx context.Context, req *v1.DeleteCredentialReq) (res *v1.DeleteCredentialRes, err error)
	RegeneratePassword(ctx context.Context, req *v1.RegeneratePasswordReq) (res *v1.RegeneratePasswordRes, err error)

	// Sending Domains
	AddSendingDomain(ctx context.Context, req *v1.AddSendingDomainReq) (res *v1.AddSendingDomainRes, err error)
	ListSendingDomains(ctx context.Context, req *v1.ListSendingDomainsReq) (res *v1.ListSendingDomainsRes, err error)
	VerifyDomain(ctx context.Context, req *v1.VerifyDomainReq) (res *v1.VerifyDomainRes, err error)
	DeleteSendingDomain(ctx context.Context, req *v1.DeleteSendingDomainReq) (res *v1.DeleteSendingDomainRes, err error)

	// IP Pools
	CreateIpPool(ctx context.Context, req *v1.CreateIpPoolReq) (res *v1.CreateIpPoolRes, err error)
	ListIpPools(ctx context.Context, req *v1.ListIpPoolsReq) (res *v1.ListIpPoolsRes, err error)
	DeleteIpPool(ctx context.Context, req *v1.DeleteIpPoolReq) (res *v1.DeleteIpPoolRes, err error)

	// Analytics
	GetDeliveryStats(ctx context.Context, req *v1.GetDeliveryStatsReq) (res *v1.GetDeliveryStatsRes, err error)

	// Suppression
	ListSuppression(ctx context.Context, req *v1.ListSuppressionReq) (res *v1.ListSuppressionRes, err error)
	AddSuppression(ctx context.Context, req *v1.AddSuppressionReq) (res *v1.AddSuppressionRes, err error)
	DeleteSuppression(ctx context.Context, req *v1.DeleteSuppressionReq) (res *v1.DeleteSuppressionRes, err error)

	// MTA Servers
	AddMtaServer(ctx context.Context, req *v1.AddMtaServerReq) (res *v1.AddMtaServerRes, err error)
	ListMtaServers(ctx context.Context, req *v1.ListMtaServersReq) (res *v1.ListMtaServersRes, err error)
	DeleteMtaServer(ctx context.Context, req *v1.DeleteMtaServerReq) (res *v1.DeleteMtaServerRes, err error)
}
