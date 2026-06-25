package sequence

import (
	"billionmail-core/api/sequence/v1"
	"billionmail-core/internal/service/public"
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

func (c *ControllerV1) SequenceList(ctx context.Context, req *v1.SequenceListReq) (res *v1.SequenceListRes, err error) {
	res = &v1.SequenceListRes{}

	model := g.DB().Model("bm_sequences")
	if req.Keyword != "" {
		model = model.WhereLike("name", "%"+req.Keyword+"%")
	}

	total, err := model.Count()
	if err != nil {
		return nil, err
	}

	var rows []map[string]interface{}
	err = model.Page(req.Page, req.PageSize).Order("create_time DESC").Scan(&rows)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.SequenceInfo, 0, len(rows))
	for _, row := range rows {
		info := rowToSequenceInfo(row)
		// get leads count
		count, _ := g.DB().Model("bm_sequence_leads").Where("sequence_id", info.Id).Count()
		info.LeadsCount = count
		list = append(list, info)
	}

	res.Data.Total = total
	res.Data.List = list
	res.SetSuccess(public.LangCtx(ctx, "Get sequence list successfully"))
	return
}

func (c *ControllerV1) SequenceFind(ctx context.Context, req *v1.SequenceFindReq) (res *v1.SequenceFindRes, err error) {
	res = &v1.SequenceFindRes{}
	var row map[string]interface{}
	err = g.DB().Model("bm_sequences").Where("id", req.Id).Scan(&row)
	if err != nil || row == nil {
		res.SetSuccess("not found")
		return
	}
	res.Data = rowToSequenceInfo(row)
	res.SetSuccess("ok")
	return
}

func (c *ControllerV1) SequenceCreate(ctx context.Context, req *v1.SequenceCreateReq) (res *v1.SequenceCreateRes, err error) {
	res = &v1.SequenceCreateRes{}
	stepsJSON, _ := json.Marshal(req.Steps)
	schedJSON, _ := json.Marshal(req.Schedule)
	now := int(time.Now().Unix())

	result, err := g.DB().Model("bm_sequences").Insert(g.Map{
		"name":         req.Name,
		"description":  req.Description,
		"status":       req.Status,
		"steps":        string(stepsJSON),
		"schedule":     string(schedJSON),
		"calendly_url": req.CalendlyUrl,
		"create_time":  now,
		"update_time":  now,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	res.Data.Id = int(id)
	res.SetSuccess(public.LangCtx(ctx, "Sequence created successfully"))
	return
}

func (c *ControllerV1) SequenceUpdate(ctx context.Context, req *v1.SequenceUpdateReq) (res *v1.SequenceUpdateRes, err error) {
	res = &v1.SequenceUpdateRes{}
	stepsJSON, _ := json.Marshal(req.Steps)
	schedJSON, _ := json.Marshal(req.Schedule)

	_, err = g.DB().Model("bm_sequences").Where("id", req.Id).Update(g.Map{
		"name":         req.Name,
		"description":  req.Description,
		"status":       req.Status,
		"steps":        string(stepsJSON),
		"schedule":     string(schedJSON),
		"calendly_url": req.CalendlyUrl,
		"update_time":  int(time.Now().Unix()),
	})
	if err != nil {
		return nil, err
	}
	res.SetSuccess(public.LangCtx(ctx, "Sequence updated successfully"))
	return
}

func (c *ControllerV1) SequenceDelete(ctx context.Context, req *v1.SequenceDeleteReq) (res *v1.SequenceDeleteRes, err error) {
	res = &v1.SequenceDeleteRes{}
	_, err = g.DB().Model("bm_sequences").Where("id", req.Id).Delete()
	if err != nil {
		return nil, err
	}
	res.SetSuccess(public.LangCtx(ctx, "Sequence deleted"))
	return
}

func (c *ControllerV1) SequencePause(ctx context.Context, req *v1.SequencePauseReq) (res *v1.SequencePauseRes, err error) {
	res = &v1.SequencePauseRes{}
	_, err = g.DB().Model("bm_sequences").Where("id", req.Id).Update(g.Map{"status": "paused", "update_time": int(time.Now().Unix())})
	if err != nil {
		return nil, err
	}
	res.SetSuccess("Sequence paused")
	return
}

func (c *ControllerV1) SequenceResume(ctx context.Context, req *v1.SequenceResumeReq) (res *v1.SequenceResumeRes, err error) {
	res = &v1.SequenceResumeRes{}
	_, err = g.DB().Model("bm_sequences").Where("id", req.Id).Update(g.Map{"status": "active", "update_time": int(time.Now().Unix())})
	if err != nil {
		return nil, err
	}
	res.SetSuccess("Sequence resumed")
	return
}

func (c *ControllerV1) SequenceLaunch(ctx context.Context, req *v1.SequenceLaunchReq) (res *v1.SequenceLaunchRes, err error) {
	res = &v1.SequenceLaunchRes{}

	// Enroll all contacts from the provided groups into this sequence
	enrolled := 0
	now := int(time.Now().Unix())
	for _, groupId := range req.LeadGroupIds {
		var contacts []struct {
			Id    int    `json:"id"`
			Email string `json:"email"`
		}
		_ = g.DB().Model("bm_contacts").Where("group_id", groupId).Fields("id, email").Scan(&contacts)
		for _, c := range contacts {
			// Avoid duplicates
			count, _ := g.DB().Model("bm_sequence_leads").
				Where("sequence_id", req.SequenceId).
				Where("contact_id", c.Id).Count()
			if count > 0 {
				continue
			}
			_, _ = g.DB().Model("bm_sequence_leads").Insert(g.Map{
				"sequence_id":    req.SequenceId,
				"contact_id":     c.Id,
				"email":          c.Email,
				"status":         "not_started",
				"current_step":   0,
				"next_send_time": now,
				"create_time":    now,
				"update_time":    now,
			})
			enrolled++
		}
	}

	// Activate the sequence
	_, _ = g.DB().Model("bm_sequences").Where("id", req.SequenceId).Update(g.Map{
		"status":      "active",
		"update_time": now,
	})

	res.Data.LeadsEnrolled = enrolled
	res.SetSuccess(public.LangCtx(ctx, "Sequence launched successfully"))
	return
}

func rowToSequenceInfo(row map[string]interface{}) *v1.SequenceInfo {
	info := &v1.SequenceInfo{}
	if row == nil {
		return info
	}
	if v, ok := row["id"]; ok {
		info.Id = int(g.NewVar(v).Int())
	}
	if v, ok := row["name"]; ok {
		info.Name = g.NewVar(v).String()
	}
	if v, ok := row["description"]; ok {
		info.Description = g.NewVar(v).String()
	}
	if v, ok := row["status"]; ok {
		info.Status = g.NewVar(v).String()
	}
	if v, ok := row["calendly_url"]; ok {
		info.CalendlyUrl = g.NewVar(v).String()
	}
	if v, ok := row["create_time"]; ok {
		info.CreateTime = g.NewVar(v).Int()
	}
	if v, ok := row["update_time"]; ok {
		info.UpdateTime = g.NewVar(v).Int()
	}
	if v, ok := row["steps"]; ok {
		_ = json.Unmarshal([]byte(g.NewVar(v).String()), &info.Steps)
	}
	if v, ok := row["schedule"]; ok {
		_ = json.Unmarshal([]byte(g.NewVar(v).String()), &info.Schedule)
	}
	return info
}
