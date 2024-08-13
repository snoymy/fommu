package usecases_test

import (
	"app/domain/usecase"
	repository_mock "app/mock/repository"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestValidEmail(t *testing.T) {
    userRepo := &repository_mock.UserRepoMock{}
    userRepo.On("FindUserByUsername", mock.Anything).Return(nil, nil)
    userRepo.On("FindUserByEmail", mock.Anything).Return(nil, nil)
    userRepo.On("CreateUser", mock.Anything).Return(nil, nil)

    uc := usecases.NewSignupUsecase(userRepo)
    err := uc.Exec(context.Background(), dto.SignupDTO{
        Username: "username",
        Password: "password1",
        Email: "test@email.com",
    })
    
    assert.NoError(t, err)
}

// func TestInvalidEmail(t *testing.T) {
//     uc := usecases.NewCreateUserUsecase()
//     err := uc.Execute(dto.CreateUserDTO{
//         Username: "username",
//         Password: "password1",
//         Email: "testemail.com",
//     })
//     
//     assert.Error(t, err)
// }
// 
// func TestValidUsername(t *testing.T) {
//     uc := usecases.NewCreateUserUsecase()
//     err := uc.Execute(dto.CreateUserDTO{
//         Username: "username",
//         Password: "password1",
//         Email: "test@email.com",
//     })
//     
//     assert.NoError(t, err)
// }
// 
// func TestInvalidUsername(t *testing.T) {
//     uc := usecases.NewCreateUserUsecase()
//     err := uc.Execute(dto.CreateUserDTO{
//         Username: "user@name",
//         Password: "password1",
//         Email: "test@email.com",
//     })
//     
//     assert.Error(t, err)
// }
