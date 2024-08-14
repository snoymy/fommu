package repoimpl_test

import (
	"app/config"
	"app/config/database"
	"app/domain/entity"
	"app/infa/repoimpl"
	"app/internal/core/entities"
	"app/internal/utils"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
    config.Init()
    conn := database.NewConnection()
    defer conn.Close(context.Background())

    userRepo := repoimpl.NewUserRepoImpl(conn) 

    user := entity.NewUserEntity()
    user.ID = uuid.New().String()
    user.Email = "test@email.com"
    user.Username = "username"
    h := sha256.New()
    h.Write([]byte("password"))
    user.PasswordHash = string(base64.StdEncoding.EncodeToString(h.Sum([]byte("password"))))
    user.Displayname = "username"
    user.URL = "www.baseurl.com/users/" + user.Username
    user.Remote = false
    user.Discoverable = true 
    user.AutoApproveFollower = false

    const bitSize = 4096
    privateKeyByte, publicKeyByte, err := utils.GenerateKeyPair(bitSize)
    assert.NoError(t, err)

    user.PrivateKey = string(privateKeyByte)
    user.PublicKey = string(publicKeyByte)
    user.CreateAt = time.Now().UTC()
    user.Status = entities.UserStatusActive
    
    err = userRepo.CreateUser(context.Background(), user)

    assert.NoError(t, err)
}

func TestFindUserByName(t *testing.T) {
    config.Init()
    conn := database.NewConnection()
    defer conn.Close(context.Background())

    userRepo := repoimpl.NewUserRepoImpl(conn) 

    user, err := userRepo.FindUserByName(context.Background(), "username")
    if !assert.NoError(t, err) {
        t.FailNow()
    }
    if !assert.NotNil(t, user) {
        t.FailNow()
    }
    assert.Equal(t, user.Username, "username")
}

func TestFindUserByEmail(t *testing.T) {
    config.Init()
    conn := database.NewConnection()
    defer conn.Close(context.Background())

    userRepo := repoimpl.NewUserRepoImpl(conn) 

    user, err := userRepo.FindUserByEmail(context.Background(), "test@email.com")
    if !assert.NoError(t, err) {
        t.FailNow()
    }
    if !assert.NotNil(t, user) {
        t.FailNow()
    }
    assert.Equal(t, user.Email, "test@email.com")
}
