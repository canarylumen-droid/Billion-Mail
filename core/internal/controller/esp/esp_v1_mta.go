package esp

import (
	"context"
	"time"

	v1 "billionmail-core/api/esp/v1"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) AddMtaServer(ctx context.Context, req *v1.AddMtaServerReq) (res *v1.AddMtaServerRes, err error) {
	res = &v1.AddMtaServerRes{}

	sshPort := req.SshPort
	if sshPort <= 0 {
		sshPort = 22
	}
	sshUser := req.SshUser
	if sshUser == "" {
		sshUser = "root"
	}
	now := time.Now().Unix()

	result, err := g.DB().Model("bm_mta_servers").Data(g.Map{
		"name":       req.Name,
		"host":       req.Host,
		"ssh_port":   sshPort,
		"ssh_user":   sshUser,
		"status":     "active",
		"region":     req.Region,
		"ip_count":   1,
		"created_at": now,
		"updated_at": now,
	}).Insert()
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		return
	}

	id, _ := result.LastInsertId()
	res.Data = &v1.MtaServer{
		Id:        id,
		Name:      req.Name,
		Host:      req.Host,
		SshPort:   sshPort,
		SshUser:   sshUser,
		Status:    "active",
		Region:    req.Region,
		IpCount:   1,
		CreatedAt: now,
	}
	return
}

func (c *ControllerV1) ListMtaServers(ctx context.Context, req *v1.ListMtaServersReq) (res *v1.ListMtaServersRes, err error) {
	res = &v1.ListMtaServersRes{}
	err = g.DB().Model("bm_mta_servers").
		Fields("id,name,host,ssh_port,ssh_user,status,region,ip_count,created_at").
		Scan(&res.Data)
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
	}
	if res.Data == nil {
		res.Data = []*v1.MtaServer{}
	}
	return
}

func (c *ControllerV1) DeleteMtaServer(ctx context.Context, req *v1.DeleteMtaServerReq) (res *v1.DeleteMtaServerRes, err error) {
	res = &v1.DeleteMtaServerRes{}
	_, err = g.DB().Model("bm_mta_servers").Where("id", req.Id).Delete()
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
	}
	return
}
