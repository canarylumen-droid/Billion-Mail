package delivery

import (
	"context"

	"billionmail-core/api/delivery/v1"
)

type IDeliveryV1 interface {
	ListBackends(ctx context.Context, req *v1.ListBackendsReq) (res *v1.ListBackendsRes, err error)
	CreateBackend(ctx context.Context, req *v1.CreateBackendReq) (res *v1.CreateBackendRes, err error)
	UpdateBackend(ctx context.Context, req *v1.UpdateBackendReq) (res *v1.UpdateBackendRes, err error)
	DeleteBackend(ctx context.Context, req *v1.DeleteBackendReq) (res *v1.DeleteBackendRes, err error)
	TestBackend(ctx context.Context, req *v1.TestBackendReq) (res *v1.TestBackendRes, err error)
	SetDefaultBackend(ctx context.Context, req *v1.SetDefaultBackendReq) (res *v1.SetDefaultBackendRes, err error)
}
