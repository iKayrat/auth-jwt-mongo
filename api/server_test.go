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

	log.Println("err", err)
	log.Println("ts.Url", ts.URL)
	log.Println("req.body", req.Body)

	require.Equal(t, 200, w.Code)

	w2 := httptest.NewRecorder()
	req2, err := http.NewRequest("GET", fmt.Sprintf(ts.URL+"/rest/token/"), nil)
	require.NoError(t, err)
	q2 := req.URL.Query()
	q.Add("id", "4b1fb443-c758-4552-900d-8f85def681f")
	req.URL.RawQuery = q2.Encode()

	testserver.Router.ServeHTTP(w2, req2)
	require.Equal(t, 400, w2.Code)

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
	// parent := context.Background()
	// ctx, cancel := context.WithTimeout(parent, time.Second*15)

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

	// queryParam := "?id=4b1fb443-c758-4552-900d-8f85def681fe"
	// req, err := http.NewRequest("GET", fmt.Sprintf(ts.URL+"/rest/token/"+queryParam), nil)
	req, err := http.NewRequest("POST", ts.URL+"/rest/token/refresh", bytes.NewBuffer(jsonbody))
	require.NoError(t, err)

	testserver.Router.ServeHTTP(w, req)

	log.Println("err", err)
	log.Println("ts.Url", ts.URL)
	log.Println("req.body", req.Body)

	// body := req.Body
	// err = json.Unmarshal([]byte(req.Body), &body)
	// log.Println(err)
	// require.NoError(t, err)
	// require.NoError(t, err)

	// log.Println("req.body", resp.Body)
	// log.Println("req.body", w.Result().Request.Body)

	require.Equal(t, 200, w.Code)
}
