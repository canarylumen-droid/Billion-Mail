package leads

import (
	"context"
	"billionmail-core/api/leads/v1"
)

type ILeadsV1 interface {
	LeadList(ctx context.Context, req *v1.LeadListReq) (res *v1.LeadListRes, err error)
	LeadImport(ctx context.Context, req *v1.LeadImportReq) (res *v1.LeadImportRes, err error)
	LeadDelete(ctx context.Context, req *v1.LeadDeleteReq) (res *v1.LeadDeleteRes, err error)
	LeadBatchDelete(ctx context.Context, req *v1.LeadBatchDeleteReq) (res *v1.LeadBatchDeleteRes, err error)
	LeadAssignSequence(ctx context.Context, req *v1.LeadAssignSequenceReq) (res *v1.LeadAssignSequenceRes, err error)
	LeadStats(ctx context.Context, req *v1.LeadStatsReq) (res *v1.LeadStatsRes, err error)
}
