package listener

import (
	"app/internal/application/activitypub/usecase"
	"app/internal/log"
	"context"
)

type ProcessActivityListener struct {
    followUser *usecase.ProcessFollowActivityUsecase `injectable:""`
}

func NewProcessActivityListener() *ProcessActivityListener {
    return &ProcessActivityListener{}
}

func (l *ProcessActivityListener) Handler(activityId string, activityType string) {
    ctx := context.Background()
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    log.Debug(ctx, activityType)
    switch activityType {
        case "Follow": l.followUser.Exec(ctx, activityId)
    }
}
