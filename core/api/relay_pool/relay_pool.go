package relay_pool

import (
        "context"

        v1 "billionmail-core/api/relay_pool/v1"
)

type IRelayPoolV1 interface {
        ListProviders(ctx context.Context, req *v1.ListProvidersReq) (res *v1.ListProvidersRes, err error)
        UpdateProvider(ctx context.Context, req *v1.UpdateProviderReq) (res *v1.UpdateProviderRes, err error)
        CreateProvider(ctx context.Context, req *v1.CreateProviderReq) (res *v1.CreateProviderRes, err error)
        DeleteProvider(ctx context.Context, req *v1.DeleteProviderReq) (res *v1.DeleteProviderRes, err error)
        TestProvider(ctx context.Context, req *v1.TestProviderReq) (res *v1.TestProviderRes, err error)
        TestPool(ctx context.Context, req *v1.TestPoolReq) (res *v1.TestPoolRes, err error)
        ResetCounters(ctx context.Context, req *v1.ResetCountersReq) (res *v1.ResetCountersRes, err error)
}
