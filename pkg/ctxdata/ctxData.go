package ctxdata

import (
	"context"
	"encoding/json"
)

// CtxKeyJwtUserId get uid from ctx
var CtxKeyJwtUserId = "jwtUserId"

func GetUidFromCtx(ctx context.Context) string {
	var userId string
	if jsonUid, ok := ctx.Value(CtxKeyJwtUserId).(json.Number); ok {
		userId = jsonUid.String()
	}
	return userId
}
