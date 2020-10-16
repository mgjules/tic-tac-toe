package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	rhttp "github.com/JulesMike/ringier-test/http"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	server *rhttp.Server
}

func (suite *ServerTestSuite) SetupTest() {
	suite.server = rhttp.NewServer(false)
	suite.server.Routes()
}

// TestRoutes is an integration test
// to test availability of the routes
func (suite *ServerTestSuite) TestRoutes() {
	cases := []struct {
		name   string
		method string
		path   string
		code   int
	}{
		{"route exist", http.MethodGet, "/", http.StatusOK},
		{"invalid method", http.MethodPost, "/", http.StatusNotFound},
		{"route does not exist", http.MethodGet, "/does-not-exist", http.StatusNotFound},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			suite.server.Router.ServeHTTP(w, req)

			suite.Equal(tc.code, w.Code)
		})
	}
}

// TestHandleHealthCheck is a unit test
// to test the health-check route
func (suite *ServerTestSuite) TestHandleHealthCheck() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	suite.server.HandleHealthCheck()(c)

	suite.Equal(http.StatusOK, w.Code)
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
