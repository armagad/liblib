package books

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
)

// func TestReadBodyAsBook(t *testing.T) {
// 	testcases := []struct {
// 		in   *strings.Reader
// 		want Book
// 	}{
// 		{strings.NewReader(`{"Title":"The Two Towers","Author":"","Abridged":false,"Edition":"","ISBN":""}`),
// 			Book{Title: "The Two Towers"}},
// 	}
// 	for _, tc := range testcases {
// 		book, err := ReadBodyAsBook(tc.in)
// 		if err != nil {
// 			t.Errorf(err.Error())
// 		}
// 		if book.Title != tc.want.Title {
// 			t.Errorf("JSON: %q, want %q", book, tc.want)
// 		}
// 	}
// }

func TestHandler(t *testing.T) {
	var UNIT_TEST_ADDR = "127.0.0.1:2020"
	api := NewApi("books", "", []DecoratorService{})
	mux := http.NewServeMux()
	mux.HandleFunc(api.CollectionRoot, api.Handler)
	var srv http.Server
	srv.Handler = mux

	l, _ := net.Listen("tcp", UNIT_TEST_ADDR)
	defer l.Close()
	go srv.Serve(l)
	defer srv.Shutdown(context.Background())

	var client http.Client
	var body *strings.Reader
	var err error

	// Create
	body = strings.NewReader(`{"Title":"The Long and Winding Road","Author":"McCartney-Lennon"}`)
	_, err = client.Post("http://"+UNIT_TEST_ADDR+api.CollectionRoot, "application/json", body)
	if err != nil {
		panic(err)
	}
	// Create
	body = strings.NewReader(`{"Title":"The Two Towers","Author":"Tolkien","Abridged":false}`)
	_, err = client.Post("http://"+UNIT_TEST_ADDR+api.CollectionRoot, "application/json", body)
	if err != nil {
		panic(err)
	}
	res, err := client.Get("http://" + UNIT_TEST_ADDR + api.CollectionRoot)
	if err != nil {
		panic(err)
	}
	b, _ := io.ReadAll(res.Body)
	if len(b) != 187 {
		t.Errorf("List: %d, want %d", len(b), 187)
	}

	// Update
	body = strings.NewReader(`{"Title":"The Road","Author":"McCarthy"}`)
	req, err := http.NewRequest(http.MethodPatch, "http://"+UNIT_TEST_ADDR+api.CollectionRoot+"1", body)
	if err != nil {
		panic(err)
	}
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	res, err = client.Get("http://" + UNIT_TEST_ADDR + api.CollectionRoot)
	if err != nil {
		panic(err)
	}
	b, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if len(b) != 162 {
		t.Errorf("List after update: %d, want %d", len(b), 162)
	}

	// Get
	res, err = client.Get("http://" + UNIT_TEST_ADDR + api.CollectionRoot + "1")
	if err != nil {
		panic(err)
	}
	b, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if len(b) != 77 {
		t.Errorf("Get: %d, want %d", len(b), 77)
	}

	// Delete
	req, err = http.NewRequest(http.MethodDelete, "http://"+UNIT_TEST_ADDR+api.CollectionRoot+"1", nil)
	if err != nil {
		panic(err)
	}
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	res, err = client.Get("http://" + UNIT_TEST_ADDR + api.CollectionRoot)
	if err != nil {
		panic(err)
	}
	b, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if len(b) != 84 {
		t.Errorf("List after delete: %d, want %d", len(b), 84)
	}

	// Get
	res, err = client.Get("http://" + UNIT_TEST_ADDR + api.CollectionRoot + "1")
	if err != nil {
		panic(err)
	}
	b, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if len(b) != 0 {
		t.Errorf("Get: %d, want %d", len(b), 0)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("got %q, want 404 Not Found", res.Status)
	}

	// Errors
	_, err = client.Post("http://"+UNIT_TEST_ADDR+api.CollectionRoot+"bad_path", "application/json", body)
	if err != nil {
		panic(err)
	}

	body = strings.NewReader(`tle":"The Long and Winding Road","Author":"McCartney-Lennon"}`)
	res, err = client.Post("http://"+UNIT_TEST_ADDR+api.CollectionRoot, "application/json", body)
	if err != nil {
		panic(err)
	}
	b, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if len(b) == 0 {
		t.Errorf("Create: 0, want >0")
	}
	if res.StatusCode != http.StatusUnsupportedMediaType {
		t.Errorf("got %q, want 415 Unsupported Media Type", res.Status)
	}

	// Update
	body = strings.NewReader(`SCALAR`)
	req, err = http.NewRequest(http.MethodPatch, "http://"+UNIT_TEST_ADDR+api.CollectionRoot+"2", body)
	if err != nil {
		panic(err)
	}
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	// Update
	body = strings.NewReader(`{"Title":"The Road","Author":"McCarthy"}`)
	req, err = http.NewRequest(http.MethodPatch, "http://"+UNIT_TEST_ADDR+api.CollectionRoot+"19236", body)
	if err != nil {
		panic(err)
	}
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	body = strings.NewReader(`{"Title":"The Road","Author":"McCarthy"}`)
	req, err = http.NewRequest(http.MethodPut, "http://"+UNIT_TEST_ADDR+api.CollectionRoot, body)
	if err != nil {
		panic(err)
	}
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	// Delete
	req, err = http.NewRequest(http.MethodDelete, "http://"+UNIT_TEST_ADDR+api.CollectionRoot+"14986", nil)
	if err != nil {
		panic(err)
	}
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
}

func TestErrors(t *testing.T) {
}
