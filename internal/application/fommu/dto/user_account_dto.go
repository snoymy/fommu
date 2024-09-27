package dto

import (
	"app/internal/core/types"
)

type UserAccountDTO struct {
    Email               *types.Nullable[string]
    CurrentPassword     *types.Nullable[string]
    NewPassword         *types.Nullable[string]
    Discoverable        *types.Nullable[bool]
}
