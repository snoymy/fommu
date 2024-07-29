package usecase

import (
	"app/internal/application/activitypub/repo"
	"app/internal/core/appstatus"
	"app/internal/config"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

type VerifySignatureUsecase struct {
	userRepo repo.UsersRepo `injectable:""`
}

func NewVerifySignatureUsecase() *VerifySignatureUsecase {
	return &VerifySignatureUsecase{}
}

// VerifySignature verifies the signature of the request.
func (uc *VerifySignatureUsecase) Exec(ctx context.Context, r *http.Request) error {
    signatureHeader := r.Header.Get("Signature")
    if signatureHeader == "" {
        return fmt.Errorf("missing Signature header")
    }

    // Parse the Signature header into a map.
    sigParams := map[string]string{}
    for _, pair := range strings.Split(signatureHeader, ",") {
        kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
        if len(kv) == 2 {
            sigParams[kv[0]] = strings.Trim(kv[1], `"`)
        }
    }

    keyID := sigParams["keyId"]
    signature, err := base64.StdEncoding.DecodeString(sigParams["signature"])
    if err != nil {
        fmt.Printf("invalid base64 signature: %v\n", err)
        return fmt.Errorf("invalid base64 signature: %v", err)
    }

    headers := strings.Split(sigParams["headers"], " ")
    signingString := ""
    for _, header := range headers {
        var headerValue string
        if header == "(request-target)" {
            headerValue = fmt.Sprintf("(request-target): %s %s", strings.ToLower(r.Method), r.URL.RequestURI())
        } else if header == "host" {
            headerValue = fmt.Sprintf("host: %s", config.Fommu.Domain)
        } else {
            headerValue = fmt.Sprintf("%s: %s", header, r.Header.Get(header))
        }
        signingString += headerValue + "\n"
    }
    signingString = strings.TrimRight(signingString, "\n")
    byts, _ := httputil.DumpRequest(r, true)
    fmt.Println(string(byts))
    fmt.Println(signingString)

    // Fetch the public key.
    actorId := strings.Split(keyID, "#")[0]
    user, err := uc.userRepo.FindUserByActorId(ctx, actorId)
    if err != nil {
        return err
    }

    if user == nil {
        return appstatus.NotFound("User not found.")
    }

    pubKeyPem := user.PublicKey
    pubKey, err := parseKey(pubKeyPem)
    if err != nil {
        return appstatus.InvalidCredential("Failed to parse public key.")
    }

    if pubKeyPem == "" || pubKey == nil {
        return appstatus.InvalidCredential("Cannot get public key")
    }

    // Hash the signing string.
    hasher := sha256.New()
    hasher.Write([]byte(signingString))
    hashed := hasher.Sum(nil)

    // Verify the signature.
    err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed, signature)
    if err != nil {
        return appstatus.InvalidCredential("Signature verification failed.")
    }

    fmt.Println("succeed")

    return nil
}

func parseKey(key string) (*rsa.PublicKey, error) {
    block, _ := pem.Decode([]byte(key))
    if block == nil {
        return nil, fmt.Errorf("failed to decode PEM block")
    }

    pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, err
    }

    rsaPubKey, ok := pubKey.(*rsa.PublicKey)
    if !ok {
        return nil, fmt.Errorf("public key is not of type RSA")
    }

    return rsaPubKey, nil
}
