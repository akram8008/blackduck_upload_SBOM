package blackduck

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// -------- project
func (b *BlackDuckComponents) getOrCreateProject(ctx context.Context, log *log.Entry, projectName, versionName string) (string, error) {
	projectLink, err := b.getProject(ctx, projectName)
	if err == nil {
		log.Infof("project: '%s' exist\n", projectName)
		return projectLink, nil
	}
	if errors.As(err, &ErrProjectNotFound) {
		log.Infof("project: '%s' doesn't exist, will create a new one\n", projectName)
		return b.createProject(ctx, projectName, versionName)
	}
	return "", err
}

func (b *BlackDuckComponents) createProject(ctx context.Context, projectName, versionName string) (string, error) {
	if !b.ValidToken() {
		if err := b.authenticate(ctx); err != nil {
			return "", errors.Wrap(err, ErrAuthenticate.Error())
		}
	}
	request := CreateProjectRequest{
		Name: projectName,
		Version: VersionRequest{
			Name:         versionName,
			Phase:        "DEVELOPMENT",
			Distribution: "EXTERNAL",
		}}
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", errors.Wrap(err, ErrCreateRequest.Error())
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.serverURL+ProjectsApi, bytes.NewReader(requestBody))
	if err != nil {
		return "", errors.Wrap(err, ErrCreateRequest.Error())
	}
	req.Header.Add("Accept", HEADER_PROJECT_DETAILS_V6)
	req.Header.Add("Content-Type", HEADER_PROJECT_DETAILS_V6)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", b.bearerToken))

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, ErrSendRequest.Error())
	}
	if resp.StatusCode != http.StatusCreated {
		return "", errors.Wrap(errors.New(resp.Status), ErrResponseStatusNotCreated.Error())
	}

	projectLink := resp.Header.Get(HEADER_PROJECT_DETAILS_Location)
	if projectLink == "" {
		err = errors.Errorf("couldn't find the project link in header: %s", HEADER_PROJECT_DETAILS_Location)
		return "", errors.Wrapf(err, ErrReadResponseData.Error())
	}

	return projectLink, nil
}

func (b *BlackDuckComponents) getProject(ctx context.Context, projectName string) (string, error) {
	if !b.ValidToken() {
		if err := b.authenticate(ctx); err != nil {
			return "", errors.Wrap(err, ErrAuthenticate.Error())
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, b.serverURL+ProjectsApi, nil)
	if err != nil {
		return "", errors.Wrap(err, ErrCreateRequest.Error())
	}

	req.Header.Add("Accept", HEADER_PROJECT_DETAILS_V6)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", b.bearerToken))
	query := req.URL.Query()
	query.Add("q", fmt.Sprintf("name:%v", projectName))
	req.URL.RawQuery = query.Encode()

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, ErrSendRequest.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.Wrap(errors.New(resp.Status), ErrResponseStatusNotOK.Error())
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, ErrReadResponseData.Error())
	}
	defer resp.Body.Close()

	projects := Projects{}
	err = json.Unmarshal(responseBody, &projects)
	if err != nil {
		return "", errors.Wrap(err, ErrParseResponseData.Error())
	} else if projects.TotalCount == 0 {
		return "", ErrProjectNotFound
	}

	// even if more than one projects found, let's return the first one with exact project name match
	for _, project := range projects.Items {
		if project.Name == projectName {
			return project.Href, nil
		}
	}

	return "", ErrProjectNotFound
}

// -------- version
func (b *BlackDuckComponents) getOrCreateProjectVersion(ctx context.Context, log *log.Entry, projectLink string, versionName string) error {
	_, err := b.getProjectVersion(ctx, projectLink, versionName)
	if err == nil {
		log.Infof("project version: '%s' exist\n", versionName)
		return nil
	}
	if errors.As(err, &ErrProjectVersionNotFound) {
		log.Infof("project version: '%s' doesn't exist, will create a new one\n", versionName)
		_, err = b.createProjectVersion(ctx, projectLink, versionName)
		return err
	}
	return err
}

func (b *BlackDuckComponents) createProjectVersion(ctx context.Context, projectLink string, versionName string) (string, error) {
	if !b.ValidToken() {
		if err := b.authenticate(ctx); err != nil {
			return "", errors.Wrap(err, ErrAuthenticate.Error())
		}
	}
	request := VersionRequest{
		Name:         versionName,
		Phase:        "DEVELOPMENT",
		Distribution: "EXTERNAL",
	}
	requestBody, err := json.Marshal(request)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, projectLink+GetVersionsApi, bytes.NewReader(requestBody))
	if err != nil {
		return "", errors.Wrap(err, ErrCreateRequest.Error())
	}

	req.Header.Add("Accept", HEADER_PROJECT_DETAILS_V5)
	req.Header.Add("Content-Type", HEADER_PROJECT_DETAILS_V5)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", b.bearerToken))

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, ErrSendRequest.Error())
	}
	if resp.StatusCode != http.StatusCreated {
		return "", errors.Wrap(errors.New(resp.Status), ErrResponseStatusNotOK.Error())
	}

	versionLink := resp.Header.Get(HEADER_PROJECT_DETAILS_Location)
	if versionLink == "" {
		err = errors.Errorf("couldn't find the version link in header: %s", HEADER_PROJECT_DETAILS_Location)
		return "", errors.Wrapf(err, ErrReadResponseData.Error())
	}

	return versionLink, nil
}

func (b *BlackDuckComponents) getProjectVersion(ctx context.Context, projectLink string, versionName string) (string, error) {
	if !b.ValidToken() {
		if err := b.authenticate(ctx); err != nil {
			return "", errors.Wrap(err, ErrAuthenticate.Error())
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, projectLink+GetVersionsApi, nil)
	if err != nil {
		return "", errors.Wrap(err, ErrCreateRequest.Error())
	}

	req.Header.Add("Accept", HEADER_PROJECT_DETAILS_V5)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", b.bearerToken))

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, ErrSendRequest.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.Wrap(errors.New(resp.Status), ErrResponseStatusNotOK.Error())
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, ErrReadResponseData.Error())
	}

	versions := ProjectVersions{}
	err = json.Unmarshal(responseBody, &versions)
	if err != nil {
		return "", errors.Wrap(err, ErrParseResponseData.Error())
	} else if versions.TotalCount == 0 {
		return "", ErrProjectVersionNotFound
	}

	for _, version := range versions.Items {
		if version.Name == versionName {
			return version.Href, nil
		}
	}

	return "", ErrProjectVersionNotFound
}
