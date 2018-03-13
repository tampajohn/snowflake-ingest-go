package ingestion

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ManagerTestSuite struct {
	suite.Suite
	DefaultManager     *Manager
	Account            string
	PipeName           string
	UserName           string
	ExpectedToken      string
	ExpectedHTTPClient *http.Client
	ExpectedRequest    *Request
	ExpectedRequestID  string
}

func (suite *ManagerTestSuite) SetupTest() {
	suite.Account = "TEST_ACCOUNT"
	suite.PipeName = "TEST_PIPE"
	suite.UserName = "TEST_USER"
	suite.DefaultManager = &Manager{
		Account:  suite.Account,
		client:   http.DefaultClient,
		PipeName: suite.PipeName,
		Port:     443,
		Scheme:   "https",
		UserName: suite.UserName,
	}
	suite.ExpectedToken = "token"
	suite.ExpectedHTTPClient = &http.Client{}
	suite.ExpectedRequestID = "some-id"
	suite.ExpectedRequest = &Request{
		Manager:   suite.DefaultManager,
		RequestID: &suite.ExpectedRequestID,
	}
	suite.ExpectedToken = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjA5MTMzMzE1OTQ0NzkxNDYsImlhdCI6MTUyMDkxMzMzMTU5NDQ3OTI2NSwiaXNzIjoiQlJJTkdIVUIuSk9ITiJ9.Qk-gPFF4gRAq9_HD1m0b8kdHGJbUTmpIqrANKiLUoSEw8jIR_ZfbwqUtyl5F1iqr-U03fWa5vxYJfF6vHiGfcgVrIjbMcjUOjb4Hf1Mf5fwG0dOshKTTEe5J322zINl14fGqLP3ze4t9vShv4fEsfuwIdAsQZwH25dsUgE94W3ALjGDUwWozOV-aeLqoahFO-1OXt7veb7zsA8lDzkihqQojZQ4lW8lOR6RIGhWECK0lnD_EYnEfn_aeSoId4gsWUlrkdGxrCBwz1JaOB6nAZYmGzpwuMHf4rNCxU3AC7Tpj0mF2yQzMPvJsGxCRnOde8lQQ-buXRqz2ClUv-lW7tw`
}

func (suite *ManagerTestSuite) TestNewManager() {
	m, err := NewManager(suite.Account, suite.UserName, suite.PipeName, nil)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.DefaultManager, m)
}

func (suite *ManagerTestSuite) TestSetToken() {
	m, err := NewManager(suite.Account, suite.UserName, suite.PipeName, nil)
	assert.Nil(suite.T(), err)
	m.SetToken(suite.ExpectedToken)
	assert.Equal(suite.T(), suite.ExpectedToken, m.token)
}

func (suite *ManagerTestSuite) TestSetClient() {
	m, err := NewManager(suite.Account, suite.UserName, suite.PipeName, nil)
	assert.Nil(suite.T(), err)
	m.SetHTTPClient(suite.ExpectedHTTPClient)
	assert.Equal(suite.T(), suite.ExpectedHTTPClient, m.client)
}

func (suite *ManagerTestSuite) TestNewRequest() {
	m, err := NewManager(suite.Account, suite.UserName, suite.PipeName, nil)
	assert.Nil(suite.T(), err)
	r := m.NewRequest(&suite.ExpectedRequestID)
	assert.Equal(suite.T(), suite.ExpectedRequest, r)
}
func TestManagerSuite(t *testing.T) {
	suite.Run(t, new(ManagerTestSuite))
}
