package crm

import (
	"context"
	"billionmail-core/api/crm/v1"
)

type ICrmV1 interface {
	InboxList(ctx context.Context, req *v1.InboxListReq) (res *v1.InboxListRes, err error)
	AiReplySuggest(ctx context.Context, req *v1.AiReplySuggestReq) (res *v1.AiReplySuggestRes, err error)
	SendReply(ctx context.Context, req *v1.SendReplyReq) (res *v1.SendReplyRes, err error)
	MarkStatus(ctx context.Context, req *v1.MarkStatusReq) (res *v1.MarkStatusRes, err error)
	SequenceAnalytics(ctx context.Context, req *v1.SequenceAnalyticsReq) (res *v1.SequenceAnalyticsRes, err error)
}
