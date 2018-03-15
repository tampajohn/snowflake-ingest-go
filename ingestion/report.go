package ingestion

import (
	"fmt"
	"net/http"
	"net/url"
)

// IngestReportResponse is awesome
type IngestReportResponse struct {
	*http.Response
}

// DoReport sends the prepared request to Snowpipe and returns an IngestFileResponse
// You can provide an optional beginMark returned in a previous IngestionReportResponse
func (r *Request) DoReport(beginMark ...string) (*IngestReportResponse, error) {
	fmt.Println(r.token)
	u, err := url.Parse(fmt.Sprintf(baseURL, r.Scheme, r.Account, r.Port))

	if err != nil {
		return nil, err
	}

	u.Path = fmt.Sprintf(insertReportPath, r.PipeName)
	q := u.Query()
	if r.RequestID != nil && len(*r.RequestID) > 0 {
		q.Add("requestId", *r.RequestID)
	}
	if len(beginMark) > 0 {
		q.Add("beginMark", beginMark[0])
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Authorization", "Bearer "+r.token)
	req.Header.Add("Content-Type", "application/json")

	var resp *http.Response
	if r.client != nil {
		resp, err = r.client.Do(req)
	} else {
		return nil, errorNoClient
	}

	if err != nil {
		fmt.Println(err)
		return &IngestReportResponse{resp}, err
	}
	return &IngestReportResponse{resp}, nil
}
