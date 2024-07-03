package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"snippetbox.abdulalsh.com/internal/assert"
	"testing"
)

func TestPing(t *testing.T) {
	//init a new response recorder
	rr := httptest.NewRecorder()
	//init a new dummy request
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	ping(rr, r)
	//get the result from http.response recorder to get the http.response generated by the ping handler
	rs := rr.Result()
	//check the status code written by the ping handler is 200
	assert.Equal(t, rs.StatusCode, http.StatusOK)
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")

}
