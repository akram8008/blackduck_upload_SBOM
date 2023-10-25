package blackduck

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/pkg/errors"
)

func (b *BlackDuckComponents) UploadSBOM(ctx context.Context, reqParams RequestUploadSBOM) error {
	l := b.log.WithField("request_id", ctx.Value("request_id"))

	projectLink, err := b.getOrCreateProject(ctx, l, reqParams.ProjectName, reqParams.ProjectVersion)
	if err != nil {
		customErr := errors.Wrap(err, "failed to get or create a project")
		l.Error(customErr)
		return customErr
	}

	if err = b.getOrCreateProjectVersion(ctx, l, projectLink, reqParams.ProjectVersion); err != nil {
		customErr := errors.Wrap(err, "failed to get or create a project version")
		l.Error(customErr)
		return customErr
	}

	l.Infoln("uploading SBOM file has started...")
	if err = b.uploadSBOM(ctx, reqParams); err != nil {
		customErr := errors.Wrap(err, "failed to upload SBOM file")
		l.Error(customErr)
		return customErr
	}
	l.Infoln("SBOM file has successfully uploaded")

	return nil
}

func (b *BlackDuckComponents) uploadSBOM(ctx context.Context, reqParams RequestUploadSBOM) error {
	if !b.ValidToken() {
		if err := b.authenticate(ctx); err != nil {
			return errors.Wrap(err, ErrAuthenticate.Error())
		}
	}

	body, bound, err := createBody(reqParams)
	if err != nil {
		return errors.Wrap(err, ErrCreateRequest.Error())
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.serverURL+UploadSBOMApi, body)
	if err != nil {
		return errors.Wrap(err, ErrCreateRequest.Error())
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", b.bearerToken))
	req.Header.Add("Content-Type", bound)
	query := req.URL.Query()
	query.Set("projectName", reqParams.ProjectName)
	query.Set("versionName", reqParams.ProjectVersion)
	req.URL.RawQuery = query.Encode()

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, ErrSendRequest.Error())
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.Wrap(errors.New(resp.Status), ErrResponseStatusNotCreated.Error())
	}

	return nil
}

func createBody(reqParams RequestUploadSBOM) (*bytes.Buffer, string, error) {
	body := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(body)
	fileHeader := make(textproto.MIMEHeader)
	fileHeader.Set("Content-Type", "application/vnd.cyclonedx")
	fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s" `, reqParams.FileName))
	writer, err := multipartWriter.CreatePart(fileHeader)
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(writer, reqParams.File)
	if err != nil {
		return nil, "", err
	}

	multipartWriter.Close()

	return body, multipartWriter.FormDataContentType(), nil
}
