package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	rhttp "github.com/mgjules/tic-tac-toe/http"
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

// TestHandleMove is a unit test
// to test the move handler
func (suite *ServerTestSuite) TestHandleMove() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	requestBody, err := json.Marshal(struct {
		GameID   string `json:"game_id"`
		Mark     uint8  `json:"mark"`
		Position uint8  `json:"position"`
	}{
		GameID:   "ababa",
		Mark:     1,
		Position: 5,
	})
	if err != nil {
		suite.Error(err)
	}

	c.Request, _ = http.NewRequest("POST", "/hola", bytes.NewBuffer(requestBody))
	suite.server.HandleMove()(c)

	suite.Equal(http.StatusOK, w.Code)
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}