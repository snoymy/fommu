package dto

import (
	"app/internal/types"
)

type UserAccountDTO struct {
    Email               *types.Nullable[string]
    CurrentPassword     *types.Nullable[string]
    NewPassword         *types.Nullable[string]
}
