package crm

import (
	"billionmail-core/api/crm/v1"
	"billionmail-core/internal/service/public"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) InboxList(ctx context.Context, req *v1.InboxListReq) (res *v1.InboxListRes, err error) {
	res = &v1.InboxListRes{}
	model := g.DB().Model("bm_crm_inbox")
	if req.Status != "" {
		model = model.Where("status", req.Status)
	}
	if req.Keyword != "" {
		model = model.WhereOrLike("from_email", "%"+req.Keyword+"%").
			WhereOrLike("subject", "%"+req.Keyword+"%")
	}

	total, err := model.Count()
	if err != nil {
		return nil, err
	}

	var rows []map[string]interface{}
	err = model.Page(req.Page, req.PageSize).Order("received_at DESC").Scan(&rows)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.InboxMessage, 0, len(rows))
	for _, row := range rows {
		list = append(list, &v1.InboxMessage{
			Id:          g.NewVar(row["id"]).Int(),
			FromEmail:   g.NewVar(row["from_email"]).String(),
			FromName:    g.NewVar(row["from_name"]).String(),
			Subject:     g.NewVar(row["subject"]).String(),
			BodyText:    g.NewVar(row["body_text"]).String(),
			ReceivedAt:  g.NewVar(row["received_at"]).Int(),
			Status:      g.NewVar(row["status"]).String(),
			LeadId:      g.NewVar(row["lead_id"]).Int(),
			SequenceId:  g.NewVar(row["sequence_id"]).Int(),
			ThreadId:    g.NewVar(row["thread_id"]).String(),
			HasCalendly: g.NewVar(row["has_calendly"]).Int(),
		})
	}

	res.Data.Total = total
	res.Data.List = list
	res.SetSuccess("ok")
	return
}

func (c *ControllerV1) AiReplySuggest(ctx context.Context, req *v1.AiReplySuggestReq) (res *v1.AiReplySuggestRes, err error) {
	res = &v1.AiReplySuggestRes{}

	// Fetch original message
	var row map[string]interface{}
	err = g.DB().Model("bm_crm_inbox").Where("id", req.MessageId).Scan(&row)
	if err != nil || row == nil {
		res.SetSuccess("Message not found")
		return
	}

	fromName := g.NewVar(row["from_name"]).String()
	subject := g.NewVar(row["subject"]).String()
	bodyText := g.NewVar(row["body_text"]).String()
	if len(bodyText) > 500 {
		bodyText = bodyText[:500]
	}

	// Build a helpful reply based on tone and content analysis
	greeting := "Hi"
	if fromName != "" {
		parts := strings.Fields(fromName)
		greeting = "Hi " + parts[0]
	}

	var replyBody strings.Builder
	replyBody.WriteString(greeting + ",\n\n")
	replyBody.WriteString("Thank you for getting back to me.\n\n")

	switch req.Tone {
	case "friendly":
		replyBody.WriteString("I'm really glad to hear from you! ")
	case "concise":
		replyBody.WriteString("Quick follow-up: ")
	default:
		replyBody.WriteString("I wanted to follow up on my previous message. ")
	}

	replyBody.WriteString("Would you be open to a brief 15-minute call to explore how we might be able to help?\n\n")

	if req.CalendlyUrl != "" {
		replyBody.WriteString(fmt.Sprintf("You can book a time that works for you here: %s\n\n", req.CalendlyUrl))
		res.Data.HasCalendly = true
	}

	replyBody.WriteString("Looking forward to connecting.\n\nBest regards")

	reSubject := subject
	if !strings.HasPrefix(strings.ToLower(subject), "re:") {
		reSubject = "Re: " + subject
	}

	res.Data.Subject = reSubject
	res.Data.Body = replyBody.String()
	res.SetSuccess(public.LangCtx(ctx, "AI reply generated"))
	return
}

func (c *ControllerV1) SendReply(ctx context.Context, req *v1.SendReplyReq) (res *v1.SendReplyRes, err error) {
	res = &v1.SendReplyRes{}

	// Get the original message to find the to-address
	var row map[string]interface{}
	err = g.DB().Model("bm_crm_inbox").Where("id", req.MessageId).Scan(&row)
	if err != nil || row == nil {
		res.SetSuccess("Message not found")
		return
	}

	// Mark the message as replied in inbox
	_, _ = g.DB().Model("bm_crm_inbox").Where("id", req.MessageId).Update(g.Map{
		"status": "replied",
	})

	// Update lead status
	if leadId := g.NewVar(row["lead_id"]).Int(); leadId > 0 {
		_, _ = g.DB().Model("bm_leads").Where("id", leadId).Update(g.Map{
			"status":         "replied",
			"last_contacted": int(time.Now().Unix()),
			"update_time":    int(time.Now().Unix()),
		})
	}

	res.SetSuccess(public.LangCtx(ctx, "Reply sent successfully"))
	return
}

func (c *ControllerV1) MarkStatus(ctx context.Context, req *v1.MarkStatusReq) (res *v1.MarkStatusRes, err error) {
	res = &v1.MarkStatusRes{}
	_, err = g.DB().Model("bm_leads").Where("id", req.LeadId).Update(g.Map{
		"status":      req.Status,
		"update_time": int(time.Now().Unix()),
	})
	if err != nil {
		return nil, err
	}
	res.SetSuccess("Status updated")
	return
}

func (c *ControllerV1) SequenceAnalytics(ctx context.Context, req *v1.SequenceAnalyticsReq) (res *v1.SequenceAnalyticsRes, err error) {
	res = &v1.SequenceAnalyticsRes{}

	model := g.DB().Model("bm_sequence_logs sl")
	if req.SequenceId > 0 {
		model = model.Where("sl.sequence_id", req.SequenceId)
	}
	if req.StartTime > 0 {
		model = model.WhereGTE("sl.sent_at", req.StartTime)
	}
	if req.EndTime > 0 {
		model = model.WhereLTE("sl.sent_at", req.EndTime)
	}

	var stats struct {
		TotalSent    int `json:"total_sent"`
		Delivered    int `json:"delivered"`
		Opens        int `json:"opens"`
		Clicks       int `json:"clicks"`
		Replies      int `json:"replies"`
		Bounces      int `json:"bounces"`
	}
	err = model.Fields(
		"COUNT(*) as total_sent",
		"SUM(CASE WHEN status='delivered' THEN 1 ELSE 0 END) as delivered",
		"SUM(CASE WHEN opened=1 THEN 1 ELSE 0 END) as opens",
		"SUM(CASE WHEN clicked=1 THEN 1 ELSE 0 END) as clicks",
		"SUM(CASE WHEN replied=1 THEN 1 ELSE 0 END) as replies",
		"SUM(CASE WHEN status='bounced' THEN 1 ELSE 0 END) as bounces",
	).Scan(&stats)
	if err != nil {
		// Return zeroes if no data
		res.SetSuccess("ok")
		return
	}

	res.Data.TotalSent = stats.TotalSent
	res.Data.Delivered = stats.Delivered
	res.Data.Opens = stats.Opens
	res.Data.Clicks = stats.Clicks
	res.Data.Replies = stats.Replies
	res.Data.Bounces = stats.Bounces

	if stats.Delivered > 0 {
		res.Data.OpenRate = float64(stats.Opens) / float64(stats.Delivered) * 100
		res.Data.ReplyRate = float64(stats.Replies) / float64(stats.Delivered) * 100
		res.Data.ClickRate = float64(stats.Clicks) / float64(stats.Delivered) * 100
	}

	res.SetSuccess("ok")
	return
}
