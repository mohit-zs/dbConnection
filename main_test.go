package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_Handler(t *testing.T) {

	config := mysqlConfig{
		dbDriver: "mysql",
		user:     "root",
		password: "insert_password",
		dbName:   "animal",
	}
	var err error
	db, err = connectToMySQL(config)
	if err != nil {
		log.Println(err)
	}

	testcases := []struct {
		method string
		body   []byte

		expectedStatusCode int
		expectedResponse   []byte
	}{
		{"GET", nil, http.StatusOK, []byte(`[{"ID":1,"name":"Hippo","age":10},{"ID":2,"name":"Ele","age":20},{"ID":3,"name":"Zebra","age":6}]`)},
		{"POST", []byte(`{"ID": 4, "Name":"Cat","Age":3}`), http.StatusOK, []byte(`success`)},
		//{"DELETE", nil, http.StatusOK, []byte(`Deletion Successful`)},

		//{"DELETE", nil, http.StatusMethodNotAllowed, nil},
	}

	for _, v := range testcases {
		req := httptest.NewRequest(v.method, "/animal/", bytes.NewReader(v.body))
		w := httptest.NewRecorder()

		h := http.HandlerFunc(handler)
		h.ServeHTTP(w, req)

		if w.Code != v.expectedStatusCode {
			t.Errorf("Expected %v Got %v", v.expectedStatusCode, w.Code)
		}

		expected := bytes.NewBuffer(v.expectedResponse)
		if !reflect.DeepEqual(expected, w.Body) {
			t.Errorf("Expected %v\tGot %v", expected.String(), w.Body.String())
		}
	}
}

//[{"ID":1,"name":"Hippo","age":10},{"ID":2,"name":"Ele","age":20},{"ID":3,"name":"Zebra","age":6}]
//[{"ID":1,"name":"Hippo","age":10},{"ID":2,"name":"Ele","age":20},{"ID":3,"name":"Zebra","age":6}]
//[{"ID":1,"name":"Hippo","age":10},{"ID":2,"name":"Ele","age":20},{"ID":3,"name":"Zebra","age":6}]
