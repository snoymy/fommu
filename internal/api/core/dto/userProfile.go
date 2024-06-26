package dto

import (
	"app/internal/types"
	"time"
)

type UserProfileDTO struct {
    Displayname         *types.Nullable[string]
    NamePrefix          *types.Nullable[string]
    NameSuffix          *types.Nullable[string]
    Gender              *types.Nullable[string]
    Pronoun             *types.Nullable[string]
    DateOfBirth         *types.Nullable[time.Time]
    ShowDateOfBirth     *types.Nullable[bool]
    Bio                 *types.Nullable[string]
    Discoverable        *types.Nullable[bool]
    AutoApproveFollower *types.Nullable[bool]
    Avatar              *types.Nullable[string]
    Banner              *types.Nullable[string]
    Preference          *types.Nullable[types.JsonObject]
}
