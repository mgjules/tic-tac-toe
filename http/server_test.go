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

func (s *ServerTestSuite) SetupTest() {
	s.mockRepository = repository.NewMock()

	s.server = rhttp.NewServer(false)
	s.server.Routes(s.mockRepository)
}

// TestRoutes is an integration test
// to test availability of the routes
func (s *ServerTestSuite) TestRoutes() {
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
		clonedTc := tc
		s.Run(tc.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(clonedTc.method, clonedTc.path, nil)
			s.server.Router.ServeHTTP(w, req)

			s.Equal(clonedTc.code, w.Code)
		})
	}
}

// TestHandleHealthCheck is a unit test
// to test the health-check route
func (s *ServerTestSuite) TestHandleHealthCheck() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s.server.HandleHealthCheck()(c)

	s.Equal(http.StatusOK, w.Code)
}

// TestHandleMove is a unit test
// to test the move handler
func (s *ServerTestSuite) TestHandleMove() {
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
		clonedTc := tc
		s.Run(tc.name, func() {
			var err error
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			requestBody, err := json.Marshal(clonedTc.request)
			if err != nil {
				s.Error(err)
			}

			c.Request, err = http.NewRequest("POST", "/move", bytes.NewBuffer(requestBody))
			if err != nil {
				s.Error(err)
			}

			s.server.HandleMove(s.mockRepository)(c)

			s.Equal(clonedTc.code, w.Code)
		})
	}
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
