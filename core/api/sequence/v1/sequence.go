package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"billionmail-core/utility/types/api_v1"
)

// SequenceStep represents a single step in an outreach sequence
type SequenceStep struct {
	Id         int    `json:"id"          dc:"Step ID"`
	Type       string `json:"type"        dc:"Step type: initial|followup1|followup2|reply"`
	Subject    string `json:"subject"     dc:"Email subject"`
	TemplateId int    `json:"template_id" dc:"Email template ID"`
	Addresser  string `json:"addresser"   dc:"From email address"`
	FullName   string `json:"full_name"   dc:"Sender display name"`
	DelayDays  int    `json:"delay_days"  dc:"Days to wait after previous step"`
	AiEnabled  int    `json:"ai_enabled"  dc:"Enable AI auto-reply (0/1)"`
	AiTone     string `json:"ai_tone"     dc:"AI tone: professional|friendly|concise"`
}

// SequenceSchedule holds schedule configuration
type SequenceSchedule struct {
	WorkingDays       []string `json:"working_days"        dc:"Active days e.g. [mon,tue,wed,thu,fri]"`
	SendWindowStart   string   `json:"send_window_start"   dc:"Send window start HH:MM"`
	SendWindowEnd     string   `json:"send_window_end"     dc:"Send window end HH:MM"`
	Timezone          string   `json:"timezone"            dc:"IANA timezone"`
	DailyLimit        int      `json:"daily_limit"         dc:"Max emails per mailbox per day"`
	RandomizeMinutes  int      `json:"randomize_minutes"   dc:"Random send delay variance in minutes"`
	StopOnReply       int      `json:"stop_on_reply"       dc:"Stop sequence when lead replies (0/1)"`
	TrackOpen         int      `json:"track_open"          dc:"Track email opens (0/1)"`
	TrackClick        int      `json:"track_click"         dc:"Track link clicks (0/1)"`
}

// SequenceInfo is the full sequence model
type SequenceInfo struct {
	Id           int              `json:"id"            dc:"Sequence ID"`
	Name         string           `json:"name"          dc:"Sequence name"`
	Description  string           `json:"description"   dc:"Description"`
	Status       string           `json:"status"        dc:"draft|active|paused"`
	Steps        []SequenceStep   `json:"steps"         dc:"Sequence steps"`
	Schedule     SequenceSchedule `json:"schedule"      dc:"Schedule settings"`
	CalendlyUrl  string           `json:"calendly_url"  dc:"Calendly booking URL"`
	LeadsCount   int              `json:"leads_count"   dc:"Total leads in sequence"`
	OpenRate     float64          `json:"open_rate"     dc:"Open rate %"`
	ReplyRate    float64          `json:"reply_rate"    dc:"Reply rate %"`
	CreateTime   int              `json:"create_time"   dc:"Create time"`
	UpdateTime   int              `json:"update_time"   dc:"Update time"`
}

// --- List ---
type SequenceListReq struct {
	g.Meta        `path:"/sequence/list" method:"get" tags:"Sequence" summary:"List sequences"`
	Authorization string `json:"authorization" in:"header"`
	Page          int    `json:"page"          d:"1"`
	PageSize      int    `json:"page_size"     d:"20"`
	Keyword       string `json:"keyword"`
}
type SequenceListRes struct {
	api_v1.StandardRes
	Data struct {
		Total int             `json:"total"`
		List  []*SequenceInfo `json:"list"`
	} `json:"data"`
}

// --- Find ---
type SequenceFindReq struct {
	g.Meta        `path:"/sequence/find" method:"get" tags:"Sequence" summary:"Get sequence by ID"`
	Authorization string `json:"authorization" in:"header"`
	Id            int    `json:"id" v:"required|min:1"`
}
type SequenceFindRes struct {
	api_v1.StandardRes
	Data *SequenceInfo `json:"data"`
}

// --- Create ---
type SequenceCreateReq struct {
	g.Meta        `path:"/sequence/create" method:"post" tags:"Sequence" summary:"Create a sequence"`
	Authorization string           `json:"authorization" in:"header"`
	Name          string           `json:"name"         v:"required|max-length:100"`
	Description   string           `json:"description"`
	Status        string           `json:"status"       d:"draft"`
	Steps         []SequenceStep   `json:"steps"`
	Schedule      SequenceSchedule `json:"schedule"`
	CalendlyUrl   string           `json:"calendly_url"`
}
type SequenceCreateRes struct {
	api_v1.StandardRes
	Data struct {
		Id int `json:"id"`
	} `json:"data"`
}

// --- Update ---
type SequenceUpdateReq struct {
	g.Meta        `path:"/sequence/update" method:"post" tags:"Sequence" summary:"Update a sequence"`
	Authorization string           `json:"authorization" in:"header"`
	Id            int              `json:"id"   v:"required|min:1"`
	Name          string           `json:"name" v:"required|max-length:100"`
	Description   string           `json:"description"`
	Status        string           `json:"status"`
	Steps         []SequenceStep   `json:"steps"`
	Schedule      SequenceSchedule `json:"schedule"`
	CalendlyUrl   string           `json:"calendly_url"`
}
type SequenceUpdateRes struct {
	api_v1.StandardRes
}

// --- Delete ---
type SequenceDeleteReq struct {
	g.Meta        `path:"/sequence/delete" method:"post" tags:"Sequence" summary:"Delete a sequence"`
	Authorization string `json:"authorization" in:"header"`
	Id            int    `json:"id" v:"required|min:1"`
}
type SequenceDeleteRes struct {
	api_v1.StandardRes
}

// --- Pause ---
type SequencePauseReq struct {
	g.Meta        `path:"/sequence/pause" method:"post" tags:"Sequence" summary:"Pause a sequence"`
	Authorization string `json:"authorization" in:"header"`
	Id            int    `json:"id" v:"required|min:1"`
}
type SequencePauseRes struct {
	api_v1.StandardRes
}

// --- Resume ---
type SequenceResumeReq struct {
	g.Meta        `path:"/sequence/resume" method:"post" tags:"Sequence" summary:"Resume a sequence"`
	Authorization string `json:"authorization" in:"header"`
	Id            int    `json:"id" v:"required|min:1"`
}
type SequenceResumeRes struct {
	api_v1.StandardRes
}

// --- Launch ---
type SequenceLaunchReq struct {
	g.Meta        `path:"/sequence/launch" method:"post" tags:"Sequence" summary:"Launch sequence with lead groups"`
	Authorization string `json:"authorization" in:"header"`
	SequenceId    int    `json:"sequence_id"    v:"required|min:1"`
	LeadGroupIds  []int  `json:"lead_group_ids" v:"required"`
}
type SequenceLaunchRes struct {
	api_v1.StandardRes
	Data struct {
		LeadsEnrolled int `json:"leads_enrolled"`
	} `json:"data"`
}
