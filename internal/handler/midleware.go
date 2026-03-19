package handler

import (
	"context"
	"fmt"

	storyW "codebase/internal/models"
	"codebase/pkg/errorx"
	"codebase/pkg/httpx-echo"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func WithAutheticatedJWTTokenData() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := c.Get(storyW.ContextJWTKey).(*jwt.Token)
			if !ok {
				return httpx.RestAbort(c, nil, errorx.Wrap(fmt.Errorf("cannot get jwt token data"), errorx.Authn))
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return httpx.RestAbort(c, nil, errorx.Wrap(fmt.Errorf("cannot get jwt token data to map claim"), errorx.Authn))
			}

			idStr, _ := claims["id"].(string)
			id, err := uuid.Parse(idStr)
			if err != nil {
				return httpx.RestAbort(c, nil, errorx.Wrap(fmt.Errorf("invalid user id in token"), errorx.Authn))
			}

			email, _ := claims["email"].(string)
			username, _ := claims["username"].(string)

			var firstName *string
			if fn, ok := claims["first_name"].(string); ok {
				firstName = &fn
			}

			var lastName *string
			if ln, ok := claims["last_name"].(string); ok {
				lastName = &ln
			}

			isActive, _ := claims["is_active"].(bool)

			userClaim := &storyW.UserClaim{
				ID:        id,
				Email:     email,
				Username:  username,
				FirstName: firstName,
				LastName:  lastName,
				IsActive:  isActive,
			}

			ctx := context.WithValue(c.Request().Context(), storyW.ContextUserClaimKey, userClaim)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
