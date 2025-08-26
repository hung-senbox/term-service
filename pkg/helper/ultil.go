package helper

import (
	"context"
	"term-service/internal/gateway/dto"
	"term-service/pkg/constants"
)

func CurrentUserFromCtx(ctx context.Context) (*dto.CurrentUser, bool) {
	if cu, ok := ctx.Value(constants.CurrentUserKey).(*dto.CurrentUser); ok {
		return cu, true
	}
	return nil, false
}
