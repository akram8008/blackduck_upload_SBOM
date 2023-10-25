package blackduck

import (
	"os"
)

const (
	TokenAuthenticateApi            = "/api/tokens/authenticate"
	UploadSBOMApi                   = "/api/scan/data"
	HEADER_USER_V4                  = "application/vnd.blackducksoftware.user-4+json"
	HEADER_PROJECT_DETAILS_V6       = "application/vnd.blackducksoftware.project-detail-6+json"
	HEADER_PROJECT_DETAILS_Location = "Location"
	HEADER_PROJECT_DETAILS_V5       = "application/vnd.blackducksoftware.project-detail-5+json"

	ProjectsApi    = "/api/projects"
	GetVersionsApi = "/versions"
)

// ResponseAuth defines the response to BlackDuck authenticate
type ResponseAuth struct {
	BearerToken                 string `json:"bearerToken,omitempty"`
	BearerExpiresInMilliseconds int64  `json:"expiresInMilliseconds,omitempty"`
}

// RequestUploadSBOM defines the request to BlackDuck uploadSBOM Api
type RequestUploadSBOM struct {
	File           *os.File
	FileName       string
	ProjectName    string
	ProjectVersion string
}

// CreateProjectRequest defines the request to BlackDuck uploadSBOM Api
type CreateProjectRequest struct {
	Name    string         `json:"name"`
	Version VersionRequest `json:"versionRequest"`
}
type VersionRequest struct {
	Name         string `json:"versionName"`
	Phase        string `json:"phase"`
	Distribution string `json:"distribution"`
}

// Projects defines the response to a BlackDuck project API request
type Projects struct {
	TotalCount int       `json:"totalCount,omitempty"`
	Items      []Project `json:"items,omitempty"`
}

// Project defines a BlackDuck project
type Project struct {
	Name     string `json:"name,omitempty"`
	Metadata `json:"_meta,omitempty"`
}

// ProjectVersions defines the response to a BlackDuck project version API request
type ProjectVersions struct {
	TotalCount int              `json:"totalCount,omitempty"`
	Items      []ProjectVersion `json:"items,omitempty"`
}

// ProjectVersion defines a version of a BlackDuck project
type ProjectVersion struct {
	Name     string `json:"versionName,omitempty"`
	Metadata `json:"_meta,omitempty"`
}

// Metadata defines BlackDuck metadata for e.g. projects
type Metadata struct {
	Href  string `json:"href,omitempty"`
	Links []Link `json:"links,omitempty"`
}

// Link defines BlackDuck links to e.g. versions of projects
type Link struct {
	Rel  string `json:"rel,omitempty"`
	Href string `json:"href,omitempty"`
}
