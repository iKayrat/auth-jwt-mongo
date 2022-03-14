package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestNewToken(t *testing.T) {

	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, time.Second*15)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	err = client.Connect(ctx)
	require.NoError(t, err)

	defer func() {
		err := client.Disconnect(ctx)
		require.NoError(t, err)
	}()

	testserver := NewServer(client)

	ts := httptest.NewServer(testserver.Router)
	defer ts.Close()

	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", fmt.Sprintf(ts.URL+"/rest/token/"), nil)
	require.NoError(t, err)
	q := req.URL.Query()
	q.Add("id", "4b1fb443-c758-4552-900d-8f85def681fe")
	req.URL.RawQuery = q.Encode()

	testserver.Router.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)

	log.Println("err", err)
	log.Println("ts.Url", ts.URL)
	log.Println("req.body", req.Body)

	w2, r2 := testInvalidId(t, ts)
	testserver.Router.ServeHTTP(w2, r2)
	require.Equal(t, 400, w2.Code)
}

func testInvalidId(t *testing.T, ts *httptest.Server) (w *httptest.ResponseRecorder, r *http.Request) {
	w = httptest.NewRecorder()
	r, err := http.NewRequest("GET", fmt.Sprintf(ts.URL+"/rest/token/"), nil)
	require.NoError(t, err)
	q := r.URL.Query()
	q.Add("id", "4b1fb443-c758-4552-900d-8f85def681f")
	r.URL.RawQuery = q.Encode()

	return
}

func TestNewRefreshToken(t *testing.T) {

	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, time.Second*15)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	err = client.Connect(ctx)
	require.NoError(t, err)

	defer func() {
		err := client.Disconnect(ctx)
		require.NoError(t, err)
	}()

	testserver := NewServer(client)

	ts := httptest.NewServer(testserver.Router)
	defer ts.Close()

	w := httptest.NewRecorder()

	type Body struct {
		Refreshtoken string `json:"refresh_token"`
	}

	testbody := Body{
		Refreshtoken: "$2a$10$YnKP8F7HaI4Q0LQOntxrBu4GlVUep5OOIJagMRvJ1H7i3oIhCQ6h6",
	}

	jsonbody, err := json.Marshal(testbody)
	log.Println(string(jsonbody))
	require.NoError(t, err)

	req, err := http.NewRequest("POST", ts.URL+"/rest/token/refresh", bytes.NewBuffer(jsonbody))
	require.NoError(t, err)

	testserver.Router.ServeHTTP(w, req)

	log.Println("err", err)
	log.Println("ts.Url", ts.URL)
	log.Println("req.body", req.Body)

	require.Equal(t, 200, w.Code)
}
