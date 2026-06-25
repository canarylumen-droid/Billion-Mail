// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sequence

import (
	"context"

	"billionmail-core/api/sequence/v1"
)

type ISequenceV1 interface {
	SequenceList(ctx context.Context, req *v1.SequenceListReq) (res *v1.SequenceListRes, err error)
	SequenceFind(ctx context.Context, req *v1.SequenceFindReq) (res *v1.SequenceFindRes, err error)
	SequenceCreate(ctx context.Context, req *v1.SequenceCreateReq) (res *v1.SequenceCreateRes, err error)
	SequenceUpdate(ctx context.Context, req *v1.SequenceUpdateReq) (res *v1.SequenceUpdateRes, err error)
	SequenceDelete(ctx context.Context, req *v1.SequenceDeleteReq) (res *v1.SequenceDeleteRes, err error)
	SequencePause(ctx context.Context, req *v1.SequencePauseReq) (res *v1.SequencePauseRes, err error)
	SequenceResume(ctx context.Context, req *v1.SequenceResumeReq) (res *v1.SequenceResumeRes, err error)
	SequenceLaunch(ctx context.Context, req *v1.SequenceLaunchReq) (res *v1.SequenceLaunchRes, err error)
}
