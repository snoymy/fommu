package usecases_test

import (
	"app/domain/usecase"
	"app/internal/api/core/entities"
	"app/internal/utils"
	repository_mock "app/mock/repository"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignin(t *testing.T) {
    user := entity.NewUserEntity()
    user.ID = uuid.New().String()
    user.Email = "user001@testemail.com"
    user.PasswordHash = utils.HashPassword("password1")

    userRepo := &repository_mock.UserRepoMock{}
    sessionRepo := &repository_mock.SessionRepoMock{}
    sessionRepo.On("CreateSession", mock.Anything).Return(nil, nil)
    userRepo.On("FindUserByEmail", mock.Anything).Return(user, nil)

    uc := usecases.NewSigninUsecase(userRepo, sessionRepo)
    session, err := uc.Exec(context.Background(), "user001@testemail.com", "password1")

    j, _ := json.MarshalIndent(session, "", "\t")
    fmt.Println(string(j))
    
    assert.NoError(t, err)
}
