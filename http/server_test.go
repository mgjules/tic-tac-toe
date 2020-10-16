package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mgjules/tic-tac-toe/game"
	"github.com/mgjules/tic-tac-toe/game/repository"
	rhttp "github.com/mgjules/tic-tac-toe/http"
	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	server         *rhttp.Server
	mockRepository game.Repository
}

func (suite *ServerTestSuite) SetupTest() {
	suite.mockRepository = repository.NewMock()

	suite.server = rhttp.NewServer(false)
	suite.server.Routes(suite.mockRepository)
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
	type Request struct {
		GameID   string `json:"game_id"`
		Mark     uint8  `json:"mark"`
		Position uint8  `json:"position"`
	}

	const gameID = "123123"
	const X = 1
	const O = 2

	cases := []struct {
		name    string
		request Request
		code    int
	}{
		{"correct move by X", Request{gameID, X, 1}, http.StatusOK},
		{"forbidden move by O", Request{gameID, O, 1}, http.StatusForbidden},
		{"out-of-bound move by O", Request{gameID, O, 9}, http.StatusBadRequest},
		{"wait-your-turn move by X", Request{gameID, X, 2}, http.StatusForbidden},
		{"correct move by O", Request{gameID, O, 0}, http.StatusOK},
		{"forbidden move by X", Request{gameID, X, 0}, http.StatusForbidden},
		{"correct move by X", Request{gameID, X, 2}, http.StatusOK},
		{"correct move by O", Request{gameID, O, 6}, http.StatusOK},
		{"correct move by X", Request{gameID, X, 5}, http.StatusOK},
		{"winning move by O", Request{gameID, O, 3}, http.StatusOK},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			var err error
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			requestBody, err := json.Marshal(tc.request)
			if err != nil {
				suite.Error(err)
			}

			c.Request, err = http.NewRequest("POST", "/move", bytes.NewBuffer(requestBody))
			if err != nil {
				suite.Error(err)
			}

			suite.server.HandleMove(suite.mockRepository)(c)

			suite.Equal(tc.code, w.Code)
		})
	}
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
