package httpclient

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/snoymy/activitypub"
)

type ActivitypubClientImpl struct { }

func NewActivitypubClientImpl() ActivitypubClient {
    return &ActivitypubClientImpl{}
}

func (c *ActivitypubClientImpl) FetchWebfinger(ctx context.Context, domain string, username string) ([]interface{}, error) {
    urls := []string{
        fmt.Sprintf("https://%s/.well-known/webfinger?resource=acct:%s@%s", domain, username, domain),
        fmt.Sprintf("https://www.%s/.well-known/webfinger?resource=acct:%s@%s", domain, username, domain),
        fmt.Sprintf("http://%s/.well-known/webfinger?resource=acct:%s@%s", domain, username, domain),
        fmt.Sprintf("http://www.%s/.well-known/webfinger?resource=acct:%s@%s", domain, username, domain),
    }

    var body []byte = nil
    for _, url := range urls {
        res, err := http.Get(url)
        if err != nil {
            continue
        }
        if res.StatusCode != http.StatusOK {
            continue
        }

        body, err = io.ReadAll(res.Body)
        if err != nil {
            return nil, err
        }
        break
    }

    if body == nil {
        return nil, nil
    }

    var info map[string]interface{}
    err := json.Unmarshal(body, &info)
    if err != nil {
        return nil, err
    }

    links, ok := info["links"].([]interface{})
    if !ok {
        return nil, err
    }

    return links, nil
}

func (c *ActivitypubClientImpl) FetchActor(ctx context.Context, url string) (*activitypub.Actor, error) {
    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Accept", "application/activity+json")

    client := &http.Client{
        Timeout: 30 * time.Second,
    }

    res, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    if res.StatusCode != http.StatusOK {
        return nil, nil 
    }

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    actor := &activitypub.Actor{}

    err = json.Unmarshal(body, &actor)
    if err != nil {
        return nil, err
    }

    var nestedScheme struct {
        Tag []*activitypub.Object `json:"tag"`
        Attachment []*activitypub.Object `json:"attachment"`
    }

    err = json.Unmarshal(body, &nestedScheme)
    if err != nil {
        return nil, err
    }

    for _, item := range nestedScheme.Tag {
        actor.Tag.Append(item)
    }

    for _, item := range nestedScheme.Attachment {
        actor.Attachment.Append(item)
    }

    return actor, nil
}

func (c *ActivitypubClientImpl) FetchOrderedCollectionPage(ctx context.Context, url string, page int) (*activitypub.OrderedCollectionPage, error) {
    queryString := ""
    if page > 0 {
        queryString = fmt.Sprintf("?page=%d", page)
    }
    url = url + queryString

    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Accept", "application/activity+json")
    client := &http.Client{
        Timeout: 30 * time.Second,
    }

    res, err := client.Do(req)
    if err != nil {
        return nil, err
    }

    if res.StatusCode != http.StatusOK {
        return nil, nil 
    }

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }
    
    //var temp types.JsonObject
    collectionPage := &activitypub.OrderedCollectionPage{}
    err = json.Unmarshal(body, &collectionPage)
    if err != nil {
        return nil, err
    }

    return collectionPage, nil
}

func (c *ActivitypubClientImpl) PublishActivity(ctx context.Context, targetUrl string, privateKey string, keyId string, activity *activitypub.Activity) error {
    message, err := json.Marshal(activity)
    if err != nil {
        return err
    }
    fmt.Println(string(message))
    fmt.Println(targetUrl)

    req, err := http.NewRequest(http.MethodPost, targetUrl, bytes.NewReader(message))
    if err != nil {
        return err
    }

    headers := []string{"(request-target)", "host", "date", "digest", "content-type"}
    parsedURL, err := url.Parse(targetUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return err
	}

	// Optionally, if you need just the hostname without the port:
	hostName := parsedURL.Hostname()
	fmt.Println("Host without port:", hostName)

    req.Header.Set("Content-Type", "application/activity+json")
    req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
    req.Header.Set("Host", hostName)
    signature, err := c.createSignature(req, headers, privateKey, keyId)
	if err != nil {
		fmt.Println("Error creating signature:", err)
		return err
	}

	// Set the Signature header in the request
	req.Header.Set("Signature", signature)
    client := &http.Client{
        Timeout: 30 * time.Second,
    }

    res, err := client.Do(req)
    if err != nil {
        return err
    }

    if res.StatusCode >= http.StatusBadRequest {
        errorMessage, err := io.ReadAll(res.Body)
        if err != nil {
            return err
        }

        return errors.New(fmt.Sprintf("error: %d, message: %s", res.StatusCode, string(errorMessage)))
    }
    
    fmt.Println("Send Succeed")
    return nil
}

func (c *ActivitypubClientImpl) createSignature(r *http.Request, headers []string, privateKeyPem string, keyID string) (string, error) {
	// Calculate Digest if included in headers
	if contains(headers, "digest") && r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read body: %v", err)
		}
		// Restore the body since ioutil.ReadAll drains it
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Calculate the SHA-256 digest
		digest := sha256.Sum256(bodyBytes)
		digestBase64 := base64.StdEncoding.EncodeToString(digest[:])
		r.Header.Set("Digest", fmt.Sprintf("SHA-256=%s", digestBase64))
	}

	// Construct the signing string
	signingString := ""
	for _, header := range headers {
		var headerValue string
		if header == "(request-target)" {
			headerValue = fmt.Sprintf("(request-target): %s %s", strings.ToLower(r.Method), r.URL.RequestURI())
		} else {
			fmt.Printf("dsfdsfdsfds %s: %s", header, r.Header.Get(header))
			headerValue = fmt.Sprintf("%s: %s", header, r.Header.Get(header))
		}
		signingString += headerValue + "\n"
	}
	signingString = strings.TrimRight(signingString, "\n")

	// Hash the signing string
	hasher := sha256.New()
	hasher.Write([]byte(signingString))
	hashed := hasher.Sum(nil)

	// Parse the private key
	block, _ := pem.Decode([]byte(privateKeyPem))
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block")
	}

    privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	// Sign the hashed value
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", fmt.Errorf("failed to sign: %v", err)
	}

	// Encode the signature in base64
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	// Construct the Signature header value
	signatureHeader := fmt.Sprintf(`keyId="%s",algorithm="rsa-sha256",headers="%s",signature="%s"`, keyID, strings.Join(headers, " "), signatureBase64)

    fmt.Println(signatureHeader)

	return signatureHeader, nil
}

// Helper function to check if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
