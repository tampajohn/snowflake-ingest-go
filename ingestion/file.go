package ingestion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// IngestFileRequest wraps an ingestion.Manager and exposes injestion functionality
type IngestFileRequest struct {
	*Request `json:"-"`
	Files    []File `json:"files"`
}

// File contains the path information for the FileRequest
type File struct {
	Path string `json:"path"`
}

// IngestFileResponse is awesome
type IngestFileResponse struct {
	*http.Response
}

// AddFiles allows you to specify file paths and returns an IngestFileRequest
func (r *Request) AddFiles(paths ...string) *IngestFileRequest {
	files := make([]File, len(paths))
	for i, path := range paths {
		files[i] = File{
			Path: path,
		}
	}
	return &IngestFileRequest{
		Request: r,
		Files:   files,
	}
}

// DoIngest sends the prepared request to Snowpipe and returns an IngestFileResponse
func (r *IngestFileRequest) DoIngest(c context.Context) (*IngestFileResponse, error) {
	u, err := url.Parse(fmt.Sprintf(baseURL, r.Scheme, r.Account, r.Port))
	if err != nil {
		return nil, err
	}

	u.Path = fmt.Sprintf(insertFilePath, r.PipeName)
	q := u.Query()
	if r.RequestID != nil && len(*r.RequestID) > 0 {
		q.Add("requestId", *r.RequestID)
	}
	u.RawQuery = q.Encode()
	b, _ := json.Marshal(r)

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(b))
	req.Header.Add(authHeader, authHeaderPrefix+r.token)
	req.Header.Add(contentTypeHeader, applicationJSON)
	var resp *http.Response
	if r.client != nil {
		resp, err = r.client.Do(req)
	} else {
		return nil, fmt.Errorf("No client supplied.  Use ingestion.Manager.SetHTTPClient to configure a custom http.Client")
	}

	if err != nil {
		return &IngestFileResponse{resp}, err
	}
	return &IngestFileResponse{resp}, nil
}
