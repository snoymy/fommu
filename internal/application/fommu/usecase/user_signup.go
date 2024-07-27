package usecase

import (
	"app/internal/core/entity"
	"app/internal/application/fommu/repo"
	"app/internal/application/fommu/validator"
	"app/internal/core/appstatus"
	"app/internal/config"
	"app/internal/log"
	"app/internal/core/types"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"time"

	"github.com/google/uuid"
)

type SignupUsecase struct {
    userRepository repo.UsersRepo `injectable:""`
}

func NewSignupUsecase() *SignupUsecase {
    return &SignupUsecase{}
}

func (uc *SignupUsecase) Exec(ctx context.Context, email string, username string, password string) error {
    log.EnterMethod(ctx)
    defer log.ExitMethod(ctx)

    var (
        existUser *entity.UserEntity
        err error
    )

    // validate username
    log.Info(ctx, "Validate username.")
    if err := validator.ValidateUsername(username); err != nil {
        return appstatus.BadUsername(err.Error())
    }

    // check if username is used
    log.Info(ctx, "Check username if used.")
    existUser, err = uc.userRepository.FindUserByUsername(ctx, username, config.Fommu.Domain)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Somethig went wrong")
    }
    if existUser != nil {
        log.Info(ctx, "Username already used.")
        return appstatus.BadUsername("Username already used.")
    }

    // validate email
    log.Info(ctx, "Validate email.")
    if err := validator.ValidateEmail(email); err != nil {
        log.Info(ctx, "Email validation failed: " + err.Error())
        return appstatus.BadEmail(err.Error())
    }

    // check if email is used
    log.Info(ctx, "Check email if used.")
    existUser, err = uc.userRepository.FindUserByEmail(ctx, email, config.Fommu.Domain)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Somethig went wrong")
    }
    if existUser != nil {
        log.Info(ctx, "Email already used.")
        return appstatus.BadEmail("Email already used.")
    }

    // validate password
    log.Info(ctx, "Validate password.")
    if err := validator.ValidatePassword(password); err != nil {
        log.Info(ctx, "Password validation failed: " + err.Error())
        return appstatus.BadPassword(err.Error())
    }
    
    log.Info(ctx, "Create user entity.")
    user := entity.NewUserEntity()
    // set id
    user.ID = uuid.New().String()
    // set email
    user.Email.Set(email)
    // set username
    user.Username = username

    // hash password
    // set password_hash
    log.Info(ctx, "Hashing Password.")
    user.PasswordHash.Set(uc.hashPassword(password))

    // set display name
    user.Displayname = username
    // set url
    // user.URL, err = url.JoinPath(config.Fommu.URL, "users", user.Username)
    // if err != nil {
    //     return appstatus.InternalServerError(err.Error())
    // }
    // set remote
    user.Remote = false
    // set discoverable
    user.Discoverable = true 
    // set auto_approve_follower
    user.AutoApproveFollower = false

    user.Domain = config.Fommu.Domain
    user.Preference.SetNull()
    user.Attachment.Set(types.JsonArray{})
    user.Tag.Set(types.JsonArray{})

    // generate key
    log.Info(ctx, "Generate public/private key.")
    const bitSize = 4096
    privateKeyByte, publicKeyByte, err := uc.generateKeyPair(bitSize)
    if err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Somethig went wrong")
    }

    // set private key
    user.PrivateKey.Set(string(privateKeyByte))
    // set public key
    user.PublicKey = string(publicKeyByte)

    user.JoinAt.Set(time.Now().UTC())
    // set create date
    user.CreateAt = time.Now().UTC()
    // set active
    user.Status = "active"

    // write user to db
    log.Info(ctx, "Write user to database.")
    if err := uc.userRepository.CreateUser(ctx, user); err != nil {
        log.Error(ctx, "Error: " + err.Error())
        return appstatus.InternalServerError("Somethig went wrong")
    }

    return nil
}

func (uc *SignupUsecase) hashPassword(password string) string {
    h := sha256.New()
    h.Write([]byte(password))
    passwordHash := string(base64.StdEncoding.EncodeToString(h.Sum([]byte(password))))

    return passwordHash
}

func (uc *SignupUsecase) generateKeyPair(bitSize int) ([]byte, []byte, error) {
    privateKey, err := uc.generatePrivateKey(bitSize)
    if err != nil {
        return nil, nil, err
    }

    publicKeyByte, err := uc.generatePublicKey(privateKey)
    if err != nil {
        return nil, nil, err
    }

    privateKeyByte := uc.encodePrivateKeyToPEM(privateKey)

    return privateKeyByte, publicKeyByte, nil
}


func (uc *SignupUsecase) generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (uc *SignupUsecase) generatePublicKey(privatekey *rsa.PrivateKey) ([]byte, error) {
	// Get ASN.1 DER format
	pubDER, err := x509.MarshalPKIXPublicKey(&privatekey.PublicKey)
    if err != nil {
        return nil, err
    }

	// pem.Block
	pubBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubDER,
	}

	// Private key in PEM format
	publicPEM := pem.EncodeToMemory(&pubBlock)

	return publicPEM, nil
}

func (uc *SignupUsecase) encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}
