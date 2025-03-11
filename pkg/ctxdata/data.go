package ctxdata

import "context"

/**
@Author: loser
@Description: get the data from context
*/

func GetUid(ctx context.Context) string {
	if u, ok := ctx.Value(IdentifyKey).(string); ok {
		return u
	}
	return ""
}
