package ingestion

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type FileTestSuite struct {
	suite.Suite
	Request       *Request
	TestFile      string
	ExpectedFiles []File
}

type mockRequester struct {
	http.Client
	mock.Mock
}

func (m *mockRequester) Do(r *http.Request) (*http.Response, error) {
	switch r.URL.Path {
	case "/v1/data/pipes/test_pipe/insertFiles":
		return &http.Response{StatusCode: http.StatusOK}, nil
	default:
		return &http.Response{StatusCode: http.StatusNotFound}, errorNotFound
	}
}

func (suite *FileTestSuite) SetupTest() {
	m, _ := NewManager(
		"TEST_ACCOUNT",
		"TEST_USER",
		"TEST_PIPE",
		nil,
	)
	m.client = &mockRequester{}
	suite.Request = m.NewRequest(nil)
	suite.TestFile = "TEST_FILE"
	suite.ExpectedFiles = []File{
		File{Path: suite.TestFile},
	}
}

func (suite *FileTestSuite) TestAddFiles() {
	suite.SetupTest()
	r := suite.Request.AddFiles(suite.TestFile)
	assert.Equal(suite.T(), suite.ExpectedFiles, r.Files)
}

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(FileTestSuite))
}

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"net/url"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestAddFiles(t *testing.T) {
// 	expected := &IngestFileRequest{
// 		Manager: &Manager{},
// 		Files: []File{File{
// 			Path: "testfile.csv",
// 		}},
// 	}
// 	m := &Manager{}
// 	actual := m.AddFiles("testfile.csv")
// 	assert.Equal(t, expected, actual)
// }

// func TestDoIngest(t *testing.T) {
// 	m := &Manager{
// 		client:   &mockRequester{},
// 		Account:  "test_account",
// 		PipeName: "test_pipe",
// 		UserName: "test_user",
// 		Scheme:   "https",
// 		Port:     443,
// 	}

// 	r := &IngestFileRequest{
// 		Manager: m,
// 		Files: []File{
// 			File{
// 				Path: "test.csv",
// 			},
// 		},
// 	}

// 	expected := &IngestFileResponse{
// 		&http.Response{
// 			StatusCode: 200,
// 		},
// 	}

// 	actual, err := r.DoIngest(context.Background())

// 	assert.Nil(t, err)
// 	assert.Equal(t, expected, actual)
// }

// func TestDoIngestBadScheme(t *testing.T) {
// 	m := &Manager{
// 		client:   &mockRequester{},
// 		Account:  "test_account",
// 		PipeName: "test_pipe",
// 		UserName: "test_user",
// 		Scheme:   "not_scheme",
// 		Port:     443,
// 	}

// 	r := &IngestFileRequest{
// 		Manager: m,
// 		Files: []File{
// 			File{
// 				Path: "test.csv",
// 			},
// 		},
// 	}

// 	actual, err := r.DoIngest(context.Background())

// 	assert.Nil(t, actual)
// 	assert.Equal(t, &url.Error{Op: "parse", URL: "not_scheme://test_account.us-east-1.snowflakecomputing.com:443", Err: fmt.Errorf("first path segment in URL cannot contain colon")}, err)
// }

// func TestDoIngestNoClient(t *testing.T) {
// 	m := &Manager{
// 		client:   nil,
// 		Account:  "test_account",
// 		PipeName: "test_pipe",
// 		UserName: "test_user",
// 		Scheme:   "https",
// 		Port:     443,
// 	}

// 	r := &IngestFileRequest{
// 		Manager: m,
// 		Files: []File{
// 			File{
// 				Path: "test.csv",
// 			},
// 		},
// 	}

// 	actual, err := r.DoIngest(context.Background())

// 	assert.Nil(t, actual)
// 	assert.Equal(t, noClientError, err)
// }

// func TestDoIngestBadRequest(t *testing.T) {
// 	m := &Manager{
// 		client:   &mockRequester{},
// 		Account:  "test_account",
// 		PipeName: "bad_pipe",
// 		UserName: "test_user",
// 		Scheme:   "https",
// 		Port:     443,
// 	}

// 	r := &IngestFileRequest{
// 		Manager: m,
// 		Files: []File{
// 			File{
// 				Path: "test.csv",
// 			},
// 		},
// 	}

// 	actual, err := r.DoIngest(context.Background())
// 	assert.Equal(t, notFoundError, err)
// 	assert.Equal(t, &IngestFileResponse{&http.Response{StatusCode: http.StatusNotFound}}, actual)
// }
