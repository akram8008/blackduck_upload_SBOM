package blackduck

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type BlackDuckComponents struct {
	serverURL  string
	token      string
	httpClient HttpClient
	mutex      *sync.Mutex
	log        *log.Logger

	bearerToken                 string
	lastAuthentication          time.Time
	bearerExpiresInMilliseconds int64
}

func NewBlackDuckComponents(serverURL, token string) *BlackDuckComponents {
	b := BlackDuckComponents{
		serverURL:  serverURL,
		token:      token,
		httpClient: &http.Client{Timeout: time.Minute},
		mutex:      &sync.Mutex{},
		log:        log.New(),
	}

	return &b
}

func (b *BlackDuckComponents) ValidToken() bool {
	b.mutex.Lock()
	expiryTime := b.lastAuthentication.Add(time.Millisecond * time.Duration(b.bearerExpiresInMilliseconds))
	b.mutex.Unlock()

	return time.Now().Before(expiryTime)
}

func (b *BlackDuckComponents) authenticate(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.serverURL+TokenAuthenticateApi, nil)
	if err != nil {
		customErr := errors.Wrap(err, ErrCreateRequest.Error())
		return customErr
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %v", b.token))
	req.Header.Add("Accept", HEADER_USER_V4)

	resp, err := b.httpClient.Do(req)
	if err != nil {
		customErr := errors.Wrap(err, ErrSendRequest.Error())
		return customErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(errors.New(resp.Status), ErrResponseStatusNotOK.Error())
	}

	responseBody, err := io.ReadAll(resp.Body)
	var responseAuth ResponseAuth
	err = json.Unmarshal(responseBody, &responseAuth)
	if err != nil {
		return errors.Wrap(err, ErrParseResponseData.Error())
	}

	b.mutex.Lock()
	b.lastAuthentication = time.Now()
	b.bearerToken = responseAuth.BearerToken
	b.bearerExpiresInMilliseconds = responseAuth.BearerExpiresInMilliseconds
	b.mutex.Unlock()

	return nil
}
