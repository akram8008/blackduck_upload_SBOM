package blackduck

import "github.com/pkg/errors"

var (
	ErrProjectNotFound        = errors.New("project not found")
	ErrProjectVersionNotFound = errors.New("project version not found")

	ErrAuthenticate = errors.New("failed to get a new token")

	ErrCreateRequest            = errors.New("failed to create a request")
	ErrSendRequest              = errors.New("failed to send the request")
	ErrResponseStatusNotOK      = errors.New("response status is not OK(200)")
	ErrResponseStatusNotCreated = errors.New("response status is not Created(201)")
	ErrReadResponseData         = errors.New("failed to read the response data")
	ErrParseResponseData        = errors.New("failed to parse the response data")
)
