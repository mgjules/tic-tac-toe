package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mgjules/tic-tac-toe/game"
	"github.com/mgjules/tic-tac-toe/game/repository"
	rhttp "github.com/mgjules/tic-tac-toe/http"
	"github.com/stretchr/testify/suite"
)

const fixturesDir = "./fixtures/"

type ServerTestSuite struct {
	suite.Suite
	server         *rhttp.Server
	mockRepository game.Repository
}

func (s *ServerTestSuite) SetupTest() {
	s.mockRepository = repository.NewMock()

	s.server = rhttp.NewServer(false)
	s.server.Routes(s.mockRepository, "", "")
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
		{"index", http.MethodGet, "/", http.StatusOK},
		{"handle move", http.MethodPost, "/move", http.StatusBadRequest},
		{"route does not exist", http.MethodGet, "/does-not-exist", http.StatusNotFound},
	}

	for _, tc := range cases {
		clonedTc := tc
		s.Run(tc.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequestWithContext(context.Background(), clonedTc.method, clonedTc.path, nil)
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

	cases := []struct {
		name        string
		gameID      string
		requestFile string
		textToCheck string
	}{
		{"X does forbidden move", "1", "x_forbidden", "already occupied"},
		{"X does not wait for turn", "2", "x_no_wait", "please wait"},
		{"X won", "3", "x_won", "X won"},
		{"O won", "4", "o_won", "O won"},
		{"tie", "5", "tie", "tie"},
	}

	for _, tc := range cases {
		clonedTc := tc
		s.Run(tc.name, func() {
			jsonFile, err := os.Open(fixturesDir + clonedTc.requestFile + ".json")
			if err != nil {
				s.FailNow(err.Error())
			}

			byteValue, err := ioutil.ReadAll(jsonFile)
			if err != nil {
				s.FailNow(err.Error())
			}

			var content struct {
				Requests []Request
			}
			if err := json.Unmarshal(byteValue, &content); err != nil {
				s.FailNow(err.Error())
			}

			var responseBody string
			for _, request := range content.Requests {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				request.GameID = clonedTc.gameID

				requestJSON, err := json.Marshal(request)
				if err != nil {
					s.FailNow(err.Error())
				}

				c.Request, err = http.NewRequest("POST", "/move", bytes.NewBuffer(requestJSON))
				if err != nil {
					s.FailNow(err.Error())
				}

				s.server.HandleMove(s.mockRepository)(c)

				responseBody = w.Body.String()
			}

			s.Contains(responseBody, clonedTc.textToCheck)
		})
	}
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
