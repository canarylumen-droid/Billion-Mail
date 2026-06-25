package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"billionmail-core/utility/types/api_v1"
)

type LeadInfo struct {
	Id             int               `json:"id"`
	Email          string            `json:"email"`
	FirstName      string            `json:"first_name"`
	LastName       string            `json:"last_name"`
	Company        string            `json:"company"`
	Title          string            `json:"title"`
	Phone          string            `json:"phone"`
	Linkedin       string            `json:"linkedin"`
	Website        string            `json:"website"`
	Attributes     map[string]string `json:"attributes"`
	Status         string            `json:"status"`
	SequenceId     int               `json:"sequence_id"`
	SequenceName   string            `json:"sequence_name"`
	SequenceStep   string            `json:"sequence_step"`
	LastContacted  int               `json:"last_contacted"`
	NextSend       int               `json:"next_send"`
	GroupId        int               `json:"group_id"`
	CreateTime     int               `json:"create_time"`
}

type LeadListReq struct {
	g.Meta        `path:"/leads/list" method:"get" tags:"Leads" summary:"List leads"`
	Authorization string `json:"authorization" in:"header"`
	Page          int    `json:"page"       d:"1"`
	PageSize      int    `json:"page_size"  d:"20"`
	Keyword       string `json:"keyword"`
	Status        string `json:"status"`
	GroupId       int    `json:"group_id"`
	SequenceId    int    `json:"sequence_id"`
}
type LeadListRes struct {
	api_v1.StandardRes
	Data struct {
		Total int         `json:"total"`
		List  []*LeadInfo `json:"list"`
	} `json:"data"`
}

type LeadImportReq struct {
	g.Meta        `path:"/leads/import" method:"post" tags:"Leads" summary:"Import leads from CSV"`
	Authorization string `json:"authorization" in:"header"`
	FileData      string `json:"file_data"  v:"required"`
	FileType      string `json:"file_type"  d:"csv"`
	GroupId       int    `json:"group_id"`
	GroupName     string `json:"group_name"`
	Overwrite     int    `json:"overwrite"  d:"0"`
	SequenceId    int    `json:"sequence_id"`
}
type LeadImportRes struct {
	api_v1.StandardRes
	Data struct {
		Imported int `json:"imported"`
		Skipped  int `json:"skipped"`
		GroupId  int `json:"group_id"`
	} `json:"data"`
}

type LeadDeleteReq struct {
	g.Meta        `path:"/leads/delete" method:"post" tags:"Leads" summary:"Delete a lead"`
	Authorization string `json:"authorization" in:"header"`
	Id            int    `json:"id" v:"required|min:1"`
}
type LeadDeleteRes struct {
	api_v1.StandardRes
}

type LeadBatchDeleteReq struct {
	g.Meta        `path:"/leads/batch_delete" method:"post" tags:"Leads" summary:"Batch delete leads"`
	Authorization string `json:"authorization" in:"header"`
	Ids           []int  `json:"ids" v:"required"`
}
type LeadBatchDeleteRes struct {
	api_v1.StandardRes
}

type LeadAssignSequenceReq struct {
	g.Meta        `path:"/leads/assign_sequence" method:"post" tags:"Leads" summary:"Assign leads to sequence"`
	Authorization string `json:"authorization" in:"header"`
	Ids        []int `json:"ids"         v:"required"`
	SequenceId int   `json:"sequence_id" v:"required|min:1"`
}
type LeadAssignSequenceRes struct {
	api_v1.StandardRes
}

type LeadStatsReq struct {
	g.Meta        `path:"/leads/stats" method:"get" tags:"Leads" summary:"Get lead stats"`
	Authorization string `json:"authorization" in:"header"`
	GroupId       int    `json:"group_id"`
}
type LeadStatsRes struct {
	api_v1.StandardRes
	Data struct {
		Total        int `json:"total"`
		NotStarted   int `json:"not_started"`
		InSequence   int `json:"in_sequence"`
		Replied      int `json:"replied"`
		Interested   int `json:"interested"`
		Bounced      int `json:"bounced"`
		Unsubscribed int `json:"unsubscribed"`
	} `json:"data"`
}
