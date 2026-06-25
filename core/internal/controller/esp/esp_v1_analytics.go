package esp

import (
	"context"
	"time"

	v1 "billionmail-core/api/esp/v1"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) GetDeliveryStats(ctx context.Context, req *v1.GetDeliveryStatsReq) (res *v1.GetDeliveryStatsRes, err error) {
	res = &v1.GetDeliveryStatsRes{}

	dateFrom := req.DateFrom
	dateTo := req.DateTo
	if dateFrom == "" {
		dateFrom = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if dateTo == "" {
		dateTo = time.Now().Format("2006-01-02")
	}

	m := g.DB().Model("bm_delivery_stats").
		Where("date_str >= ?", dateFrom).
		Where("date_str <= ?", dateTo)
	if req.Domain != "" {
		m = m.Where("domain", req.Domain)
	}

	var byDay []*v1.DeliveryStats
	err = m.Fields("date_str, SUM(sent) as sent, SUM(delivered) as delivered, SUM(bounced_hard) as bounced_hard, SUM(bounced_soft) as bounced_soft, SUM(complained) as complained, SUM(opened) as opened, SUM(clicked) as clicked, SUM(unsubscribed) as unsubscribed").
		Group("date_str").
		Order("date_str ASC").
		Scan(&byDay)
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		return
	}
	if byDay == nil {
		byDay = []*v1.DeliveryStats{}
	}
	for _, d := range byDay {
		if d.Sent > 0 {
			d.DeliveryRate = float64(d.Delivered) / float64(d.Sent) * 100
			d.OpenRate = float64(d.Opened) / float64(d.Sent) * 100
			d.BounceRate = float64(d.BouncedHard+d.BouncedSoft) / float64(d.Sent) * 100
		}
	}

	var byDomain []*v1.DeliveryStats
	err = m.Fields("domain, SUM(sent) as sent, SUM(delivered) as delivered, SUM(bounced_hard) as bounced_hard, SUM(bounced_soft) as bounced_soft, SUM(complained) as complained, SUM(opened) as opened, SUM(clicked) as clicked").
		Group("domain").
		Order("sent DESC").
		Scan(&byDomain)
	if err != nil {
		byDomain = []*v1.DeliveryStats{}
		err = nil
	}
	if byDomain == nil {
		byDomain = []*v1.DeliveryStats{}
	}
	for _, d := range byDomain {
		if d.Sent > 0 {
			d.DeliveryRate = float64(d.Delivered) / float64(d.Sent) * 100
			d.OpenRate = float64(d.Opened) / float64(d.Sent) * 100
			d.BounceRate = float64(d.BouncedHard+d.BouncedSoft) / float64(d.Sent) * 100
		}
	}

	summary := &v1.DeliveryStats{}
	for _, d := range byDay {
		summary.Sent += d.Sent
		summary.Delivered += d.Delivered
		summary.BouncedHard += d.BouncedHard
		summary.BouncedSoft += d.BouncedSoft
		summary.Complained += d.Complained
		summary.Opened += d.Opened
		summary.Clicked += d.Clicked
	}
	if summary.Sent > 0 {
		summary.DeliveryRate = float64(summary.Delivered) / float64(summary.Sent) * 100
		summary.OpenRate = float64(summary.Opened) / float64(summary.Sent) * 100
		summary.BounceRate = float64(summary.BouncedHard+summary.BouncedSoft) / float64(summary.Sent) * 100
	}

	res.Data.Summary = summary
	res.Data.ByDay = byDay
	res.Data.ByDomain = byDomain
	return
}
