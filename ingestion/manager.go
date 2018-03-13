package ingestion

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	baseURL                    = "%s://%s.us-east-1.snowflakecomputing.com:%d"
	insertFilePath             = "/v1/data/pipes/%s/insertFiles"
	insertReportPath           = "/v1/data/pipes/%s/insertReport"
	loadHistoryScanURLTemplate = "/v1/data/pipes/%s/loadHistoryScan"

	requestIDParam         = "requestId"
	stageParam             = "stage"
	recentSecondsParam     = "recentSeconds"
	historyBeginMarkParam  = "beginMark"
	historyRangeStartParam = "startTimeInclusive"
	historyRangeEndParam   = "endTimeExclusive"

	issuer     = "iss"
	expireTime = "exp"
	issueTime  = "iat"
	subject    = "sub"

	authHeader        = "Authorization"
	authHeaderPrefix  = "Bearer "
	contentTypeHeader = "Content-Type"
	applicationJSON   = "application/json"

	lifetime = 59 * time.Minute
)

var (
	errorNoClient = errors.New("No client supplied.  Use WithClient to configure a custom http.Client")
	errorNotFound = errors.New("Not Found")
)

// Requester is an interface for making HTTP Calls and can be fulfilled by
// the standard http.Client type.
type Requester interface {
	Do(req *http.Request) (*http.Response, error)
}

// Manager is the Base Snowpipe IngestionManager
type Manager struct {
	Account  string
	UserName string
	PipeName string
	Port     int
	Scheme   string
	token    string
	client   Requester
}

// Request is the base request used by Ingestion and Reporting
type Request struct {
	*Manager
	RequestID *string
}

// NewManager returns a new ingest manager instance.  Given a private key,
// username, pipename, and account, it ensures all requests to Snowflake Ingest
// have a valid JWT token specified.
func NewManager(account, userName, pipeName string, privateKey *[]byte) (*Manager, error) {

	m := &Manager{
		Account:  strings.ToUpper(account),
		UserName: strings.ToUpper(userName),
		PipeName: pipeName,
		Scheme:   "https",
		Port:     443,
		client:   http.DefaultClient,
	}

	if privateKey != nil {
		b, _ := pem.Decode(*privateKey)
		if _, err := x509.ParsePKCS1PrivateKey(b.Bytes); err != nil {
			return m, err
		}

		t := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			issuer:     fmt.Sprintf("%s.%s", strings.ToUpper(account), strings.ToUpper(userName)),
			expireTime: time.Now().UTC().UnixNano(),
			issueTime:  time.Now().UTC().UnixNano(),
		})

		key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(*privateKey))
		if err != nil {

			fmt.Println(err)
			return m, err
		}

		tokenString, err := t.SignedString(key)
		if err != nil {
			return m, err
		}
		m.token = tokenString
	}

	return m, nil
}

// SetToken exposes the ability to specify a token on the manager without requiring
// sharing the contents of your private key with the package.
func (m *Manager) SetToken(token string) {
	m.token = token
}

// SetHTTPClient exposes the ability to specify a custom http.Client to use when
// making Snowflake Ingest Requests.
func (m *Manager) SetHTTPClient(client *http.Client) {
	m.client = client
}

//NewRequest returns an instance of ingestion.Request
func (m *Manager) NewRequest(requestID *string) *Request {
	r := &Request{
		Manager:   m,
		RequestID: requestID,
	}
	return r
}
