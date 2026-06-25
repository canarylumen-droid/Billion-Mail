package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"billionmail-core/utility/types/api_v1"
)

type InboxMessage struct {
	Id          int    `json:"id"`
	FromEmail   string `json:"from_email"`
	FromName    string `json:"from_name"`
	Subject     string `json:"subject"`
	BodyText    string `json:"body_text"`
	BodyHtml    string `json:"body_html"`
	ReceivedAt  int    `json:"received_at"`
	Status      string `json:"status"`
	LeadId      int    `json:"lead_id"`
	SequenceId  int    `json:"sequence_id"`
	ThreadId    string `json:"thread_id"`
	HasCalendly int    `json:"has_calendly"`
}

// --- Inbox List ---
type InboxListReq struct {
	g.Meta        `path:"/crm/inbox" method:"get" tags:"CRM" summary:"List inbox messages"`
	Authorization string `json:"authorization" in:"header"`
	Page          int    `json:"page"      d:"1"`
	PageSize      int    `json:"page_size" d:"20"`
	Status        string `json:"status"`
	Keyword       string `json:"keyword"`
}
type InboxListRes struct {
	api_v1.StandardRes
	Data struct {
		Total int             `json:"total"`
		List  []*InboxMessage `json:"list"`
	} `json:"data"`
}

// --- AI Reply Suggestion ---
type AiReplySuggestReq struct {
	g.Meta        `path:"/crm/ai_reply" method:"post" tags:"CRM" summary:"AI reply suggestion"`
	Authorization string `json:"authorization" in:"header"`
	MessageId     int    `json:"message_id"   v:"required|min:1"`
	Tone          string `json:"tone"         d:"professional"`
	CalendlyUrl   string `json:"calendly_url"`
}
type AiReplySuggestRes struct {
	api_v1.StandardRes
	Data struct {
		Subject     string `json:"subject"`
		Body        string `json:"body"`
		HasCalendly bool   `json:"has_calendly"`
	} `json:"data"`
}

// --- Send Reply ---
type SendReplyReq struct {
	g.Meta        `path:"/crm/send_reply" method:"post" tags:"CRM" summary:"Send reply to a lead"`
	Authorization string `json:"authorization" in:"header"`
	MessageId     int    `json:"message_id" v:"required|min:1"`
	Subject       string `json:"subject"    v:"required"`
	Body          string `json:"body"       v:"required"`
	FromEmail     string `json:"from_email" v:"required|email"`
}
type SendReplyRes struct {
	api_v1.StandardRes
}

// --- Mark Status ---
type MarkStatusReq struct {
	g.Meta        `path:"/crm/mark_status" method:"post" tags:"CRM" summary:"Mark lead status"`
	Authorization string `json:"authorization" in:"header"`
	LeadId        int    `json:"lead_id" v:"required|min:1"`
	Status        string `json:"status"  v:"required"`
}
type MarkStatusRes struct {
	api_v1.StandardRes
}

// --- Sequence Analytics ---
type SequenceAnalyticsReq struct {
	g.Meta        `path:"/crm/analytics" method:"get" tags:"CRM" summary:"Sequence analytics"`
	Authorization string `json:"authorization" in:"header"`
	SequenceId    int    `json:"sequence_id"`
	StartTime     int    `json:"start_time"`
	EndTime       int    `json:"end_time"`
}
type SequenceAnalyticsRes struct {
	api_v1.StandardRes
	Data struct {
		TotalSent      int     `json:"total_sent"`
		Delivered      int     `json:"delivered"`
		Opens          int     `json:"opens"`
		Clicks         int     `json:"clicks"`
		Replies        int     `json:"replies"`
		Bounces        int     `json:"bounces"`
		Unsubscribes   int     `json:"unsubscribes"`
		OpenRate       float64 `json:"open_rate"`
		ReplyRate      float64 `json:"reply_rate"`
		ClickRate      float64 `json:"click_rate"`
		StepBreakdown  []StepStat `json:"step_breakdown"`
	} `json:"data"`
}

type StepStat struct {
	Step    string  `json:"step"`
	Sent    int     `json:"sent"`
	Opens   int     `json:"opens"`
	Replies int     `json:"replies"`
	Rate    float64 `json:"rate"`
}
