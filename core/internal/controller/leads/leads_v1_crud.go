package leads

import (
	"billionmail-core/api/leads/v1"
	"billionmail-core/internal/service/public"
	"context"
	"encoding/csv"
	"encoding/json"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) LeadList(ctx context.Context, req *v1.LeadListReq) (res *v1.LeadListRes, err error) {
	res = &v1.LeadListRes{}
	model := g.DB().Model("bm_leads bl").
		LeftJoin("bm_sequences bs", "bl.sequence_id = bs.id")

	if req.Keyword != "" {
		model = model.WhereOrLike("bl.email", "%"+req.Keyword+"%").
			WhereOrLike("bl.first_name", "%"+req.Keyword+"%").
			WhereOrLike("bl.company", "%"+req.Keyword+"%")
	}
	if req.Status != "" {
		model = model.Where("bl.status", req.Status)
	}
	if req.GroupId > 0 {
		model = model.Where("bl.group_id", req.GroupId)
	}
	if req.SequenceId > 0 {
		model = model.Where("bl.sequence_id", req.SequenceId)
	}

	total, err := model.Count()
	if err != nil {
		return nil, err
	}

	var rows []map[string]interface{}
	err = model.Fields("bl.*, bs.name as sequence_name").
		Page(req.Page, req.PageSize).
		Order("bl.create_time DESC").
		Scan(&rows)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.LeadInfo, 0, len(rows))
	for _, row := range rows {
		lead := &v1.LeadInfo{
			Id:            g.NewVar(row["id"]).Int(),
			Email:         g.NewVar(row["email"]).String(),
			FirstName:     g.NewVar(row["first_name"]).String(),
			LastName:      g.NewVar(row["last_name"]).String(),
			Company:       g.NewVar(row["company"]).String(),
			Title:         g.NewVar(row["title"]).String(),
			Phone:         g.NewVar(row["phone"]).String(),
			Status:        g.NewVar(row["status"]).String(),
			SequenceId:    g.NewVar(row["sequence_id"]).Int(),
			SequenceName:  g.NewVar(row["sequence_name"]).String(),
			SequenceStep:  g.NewVar(row["sequence_step"]).String(),
			LastContacted: g.NewVar(row["last_contacted"]).Int(),
			GroupId:       g.NewVar(row["group_id"]).Int(),
			CreateTime:    g.NewVar(row["create_time"]).Int(),
		}
		if attrs, ok := row["attributes"]; ok {
			_ = json.Unmarshal([]byte(g.NewVar(attrs).String()), &lead.Attributes)
		}
		list = append(list, lead)
	}

	res.Data.Total = total
	res.Data.List = list
	res.SetSuccess("ok")
	return
}

func (c *ControllerV1) LeadImport(ctx context.Context, req *v1.LeadImportReq) (res *v1.LeadImportRes, err error) {
	res = &v1.LeadImportRes{}
	now := int(time.Now().Unix())

	// Create group if name provided
	groupId := req.GroupId
	if req.GroupName != "" {
		result, gErr := g.DB().Model("bm_contact_groups").Insert(g.Map{
			"name":        req.GroupName,
			"description": "Imported leads",
			"create_time": now,
			"update_time": now,
		})
		if gErr == nil {
			id, _ := result.LastInsertId()
			groupId = int(id)
		}
	}
	res.Data.GroupId = groupId

	// Parse CSV
	reader := csv.NewReader(strings.NewReader(req.FileData))
	records, err := reader.ReadAll()
	if err != nil || len(records) < 2 {
		res.SetSuccess("No valid records found")
		return
	}

	headers := records[0]
	// Map headers to field indices
	fieldMap := map[string]int{}
	for i, h := range headers {
		fieldMap[strings.ToLower(strings.TrimSpace(h))] = i
	}

	imported := 0
	skipped := 0
	for _, record := range records[1:] {
		getField := func(names ...string) string {
			for _, n := range names {
				if idx, ok := fieldMap[n]; ok && idx < len(record) {
					return strings.TrimSpace(record[idx])
				}
			}
			return ""
		}

		email := getField("email", "e-mail")
		if email == "" || !strings.Contains(email, "@") {
			skipped++
			continue
		}

		// Build attributes from remaining columns
		attrs := map[string]string{}
		knownFields := map[string]bool{"email": true, "first_name": true, "firstname": true, "last_name": true, "lastname": true, "company": true, "title": true, "phone": true, "linkedin": true, "website": true}
		for _, h := range headers {
			key := strings.ToLower(strings.TrimSpace(h))
			if !knownFields[key] {
				if idx, ok := fieldMap[key]; ok && idx < len(record) {
					attrs[key] = record[idx]
				}
			}
		}
		attrsJSON, _ := json.Marshal(attrs)

		// Check duplicate
		count, _ := g.DB().Model("bm_leads").Where("email", email).Count()
		if count > 0 && req.Overwrite == 0 {
			skipped++
			continue
		}

		rowData := g.Map{
			"email":        email,
			"first_name":   getField("first_name", "firstname", "first"),
			"last_name":    getField("last_name", "lastname", "last"),
			"company":      getField("company", "organization", "org"),
			"title":        getField("title", "role", "position", "jobtitle"),
			"phone":        getField("phone", "telephone", "mobile"),
			"linkedin":     getField("linkedin", "linkedin_url"),
			"website":      getField("website", "url"),
			"attributes":   string(attrsJSON),
			"status":       "not_started",
			"group_id":     groupId,
			"sequence_id":  req.SequenceId,
			"create_time":  now,
			"update_time":  now,
		}

		if count > 0 && req.Overwrite == 1 {
			_, _ = g.DB().Model("bm_leads").Where("email", email).Update(rowData)
		} else {
			_, _ = g.DB().Model("bm_leads").Insert(rowData)
		}
		imported++
	}

	res.Data.Imported = imported
	res.Data.Skipped = skipped
	res.SetSuccess(public.LangCtx(ctx, "Leads imported successfully"))
	return
}

func (c *ControllerV1) LeadDelete(ctx context.Context, req *v1.LeadDeleteReq) (res *v1.LeadDeleteRes, err error) {
	res = &v1.LeadDeleteRes{}
	_, err = g.DB().Model("bm_leads").Where("id", req.Id).Delete()
	if err != nil {
		return nil, err
	}
	res.SetSuccess("Deleted")
	return
}

func (c *ControllerV1) LeadBatchDelete(ctx context.Context, req *v1.LeadBatchDeleteReq) (res *v1.LeadBatchDeleteRes, err error) {
	res = &v1.LeadBatchDeleteRes{}
	_, err = g.DB().Model("bm_leads").WhereIn("id", req.Ids).Delete()
	if err != nil {
		return nil, err
	}
	res.SetSuccess("Deleted")
	return
}

func (c *ControllerV1) LeadAssignSequence(ctx context.Context, req *v1.LeadAssignSequenceReq) (res *v1.LeadAssignSequenceRes, err error) {
	res = &v1.LeadAssignSequenceRes{}
	_, err = g.DB().Model("bm_leads").WhereIn("id", req.Ids).Update(g.Map{
		"sequence_id": req.SequenceId,
		"status":      "in_sequence",
		"update_time": int(time.Now().Unix()),
	})
	if err != nil {
		return nil, err
	}
	res.SetSuccess("Assigned")
	return
}

func (c *ControllerV1) LeadStats(ctx context.Context, req *v1.LeadStatsReq) (res *v1.LeadStatsRes, err error) {
	res = &v1.LeadStatsRes{}
	model := g.DB().Model("bm_leads")
	if req.GroupId > 0 {
		model = model.Where("group_id", req.GroupId)
	}

	total, _ := model.Count()
	res.Data.Total = total

	statuses := []string{"not_started", "in_sequence", "replied", "interested", "bounced", "unsubscribed"}
	for _, s := range statuses {
		count, _ := g.DB().Model("bm_leads").Where("status", s).Count()
		switch s {
		case "not_started":
			res.Data.NotStarted = count
		case "in_sequence":
			res.Data.InSequence = count
		case "replied":
			res.Data.Replied = count
		case "interested":
			res.Data.Interested = count
		case "bounced":
			res.Data.Bounced = count
		case "unsubscribed":
			res.Data.Unsubscribed = count
		}
	}
	res.SetSuccess("ok")
	return
}
