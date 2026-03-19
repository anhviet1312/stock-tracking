package utils

import (
	"context"

	storyW "codebase/internal/models"
)

func AuthenticatedUserClaimFromContext(ctx context.Context) (*storyW.UserClaim, bool) {
	v, ok := ctx.Value(storyW.ContextUserClaimKey).(*storyW.UserClaim)
	return v, ok
}
