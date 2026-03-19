package models

type Ctxkey string

const (
	ContextJWTKey              = "jwt_auth_user"
	ContextUserClaimKey Ctxkey = "auth_user"
)
