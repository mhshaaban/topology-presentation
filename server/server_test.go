package server_test

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/owulveryck/topology-presentation/server"

	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

var (
	testServer *httptest.Server
	reader     io.Reader //Ignore this for now
	baseWsUrl  string
)

func init() {
	router := server.NewRouter()
	go server.AllHubs.Run()
	testServer = httptest.NewServer(router) //Creating new server with the user handlers

	baseWsUrl = fmt.Sprintf("%s/serveWs/", testServer.URL) //Grab the address for the endpoint

}

func TestServeWs(t *testing.T) {
	tsURL, err := url.Parse(testServer.URL)
	if err != nil {
		t.Error(err)
	}
	httpURL := url.URL{Scheme: tsURL.Scheme, Host: tsURL.Host, Path: "/serveWs/"}
	// Try to connect to a socket without an ID
	request, err := http.NewRequest("GET", httpURL.String(), nil)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	// We don't serve the baseurl, a tag is mandatory
	if res.StatusCode != 404 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}

	//Try with a valid tag
	httpURL.Path = "/serveWs/1234'"
	request, err = http.NewRequest("GET", httpURL.String(), nil)

	res, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	// We shall get a bad request as we are expected a websocket
	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	// Now test the websocket
	wsURL := url.URL{Scheme: "ws", Host: tsURL.Host, Path: "/serveWs/1234"}
	c, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		t.Errorf("Cannot connect to the websocket %v", err)

	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				t.Errorf("read: %v", err)
			}
			t.Logf("recv: %s", message)
		}

	}()
	message := &server.Node{}
	b, err := json.Marshal(message)
	if err != nil {
		t.Error(err)
	}
	//reader = strings.NewReader(string(b))
	err = c.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		t.Errorf("write:", err)
	}
	<-time.After(5 * time.Second)
}
