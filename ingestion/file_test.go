package ingestion

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type FileTestSuite struct {
	suite.Suite
	Request                        *Request
	RequestID                      string
	TestFile                       string
	TestManager                    *Manager
	ExpectedFiles                  []File
	ExpectedOkResponse             *IngestFileResponse
	ExpectedContinueResponseWithID *IngestFileResponse
	BadSchemeManager               *Manager
	ErrorExpectedMissingSchema     error
}

type mockRequester struct {
	http.Client
	mock.Mock
}

func (m *mockRequester) Do(r *http.Request) (*http.Response, error) {
	switch r.URL.Query().Get("requestId") {
	case "TEST_ID":
		return &http.Response{StatusCode: http.StatusContinue}, nil
	}
	switch r.URL.Path {
	case "/v1/data/pipes/TEST_PIPE/insertFiles":
		return &http.Response{StatusCode: http.StatusOK}, nil
	}

	return &http.Response{StatusCode: http.StatusNotFound}, errorNotFound
}

func (suite *FileTestSuite) SetupTest() {
	m, _ := NewManager(
		"TEST_ACCOUNT",
		"TEST_USER",
		"TEST_PIPE",
		nil,
	)
	m.client = &mockRequester{}
	suite.TestManager = m
	suite.BadSchemeManager = &Manager{
		Account:  "TEST_ACCOUNT",
		PipeName: "TEST_PIPE",
		Port:     443,
		Scheme:   "",
		UserName: "TEST_USER",
		client:   &mockRequester{},
	}
	suite.Request = m.NewRequest(nil)
	suite.TestFile = "TEST_FILE"
	suite.ExpectedFiles = []File{
		File{Path: suite.TestFile},
	}
	suite.ExpectedOkResponse = &IngestFileResponse{
		&http.Response{
			StatusCode: http.StatusOK,
		},
	}
	suite.RequestID = "TEST_ID"
	suite.ErrorExpectedMissingSchema = &url.Error{Op: "parse", URL: "://TEST_ACCOUNT.us-east-1.snowflakecomputing.com:443", Err: fmt.Errorf("missing protocol scheme")}
	suite.ExpectedContinueResponseWithID = &IngestFileResponse{
		&http.Response{
			StatusCode: http.StatusContinue,
		},
	}
}

func (suite *FileTestSuite) TestAddFiles() {
	suite.SetupTest()
	r := suite.Request.AddFiles(suite.TestFile)
	assert.Equal(suite.T(), suite.ExpectedFiles, r.Files)
}

func (suite *FileTestSuite) TestDoIngest() {
	suite.SetupTest()
	actual, err := suite.Request.AddFiles(suite.TestFile).DoIngest(context.Background())

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.ExpectedOkResponse, actual)
}

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileTestSuite))
}

func (suite *FileTestSuite) TestDoIngestBadScheme() {

	actual, err := suite.BadSchemeManager.NewRequest(nil).AddFiles(suite.TestFile).DoIngest(context.Background())

	assert.Nil(suite.T(), actual)
	assert.Equal(suite.T(), suite.ErrorExpectedMissingSchema, err)
}

func (suite *FileTestSuite) TestDoIngestWithRequestID() {

	actual, err := suite.TestManager.NewRequest(&suite.RequestID).AddFiles(suite.TestFile).DoIngest(context.Background())

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.ExpectedContinueResponseWithID, actual)
}
