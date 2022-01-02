package main

import (
	"encoding/json"
	"httpfileserver/server"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListDirsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/listdirs?dir=.", nil)
	w := httptest.NewRecorder()
	server.ListDirsHandler(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Reading body expected error to be nil got %v", err)
	}
	// Check json
	var fileresp server.FilesResponse
	err = json.Unmarshal(data, &fileresp)

	if err != nil {
		t.Errorf("Pasing json response expected error to be nil got %v", err)
	}
	if fileresp.Folder != "." {
		t.Errorf("expected Folder to be '.' got %v", fileresp.Folder)
	}
}
