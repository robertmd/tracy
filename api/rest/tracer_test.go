package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nccgroup/tracy/api/types"
)

// Testing addTracer with httptest. POST /tracers
func testAddTracer(t *testing.T) []RequestTestPair {
	var (
		tracerString     = "blahblah"
		URL              = "http://example.com"
		method           = "GET"
		addURL           = "http://127.0.0.1:8081/tracers"
		getURL           = "http://127.0.0.1:8081/tracers/1"
		rawRequest       = "GET / HTTP/1.1\\nHost: gorm.io\\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:58.0) Gecko/20100101 Firefox/58.0\\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,;q=0.8\\nAccept-Language: en-US,en;q=0.5\\nAccept-Encoding: gzip, deflate\\nConnection: keep-alive\\nPragma: no-cacheCache-Control: no-cache"
		addTracerPayload = fmt.Sprintf(`{"RawRequest": "%s", "RequestURL": "%s", "RequestMethod": "%s", "Tracers": [{"TracerPayload": "%s"}]}`, rawRequest, URL, method, tracerString)
	)

	/* ADDING A TRACER */
	/////////////////////
	/* Make the POST request. */
	addReq, err := http.NewRequest("POST", addURL, bytes.NewBuffer([]byte(addTracerPayload)))
	if err != nil {
		t.Fatalf("tried to build an HTTP request, but got the following error: %+v", err)
	}
	/* ADDING A TRACER */
	/////////////////////

	/* GETING A TRACER */
	/////////////////////
	getReq, err := http.NewRequest("GET", getURL, nil)
	if err != nil {
		t.Fatalf("tried to build an HTTP request but got the following error: %+v", err)
	}
	/* GETTING A TRACER */
	/////////////////////

	/* Create a mapping of the request/test and use the server helper to execute it. */
	tests := make([]RequestTestPair, 2)
	addReqTest := RequestTestPair{addReq, addTest}
	getReqTest := RequestTestPair{getReq, getTest}
	tests[0] = addReqTest
	tests[1] = getReqTest
	return tests
}

// getTest is the commonly used GET request test.
func getTest(rr *httptest.ResponseRecorder, t *testing.T) error {
	/* Return variable. */
	var err error

	if status := rr.Code; status != http.StatusOK {
		err = fmt.Errorf("GetTracer returned the wrong status code. Got %v, but wanted %v", status, http.StatusOK)
	} else {
		/* Validate the tracer was the first tracer inserted. */
		got := types.Tracer{}
		json.Unmarshal([]byte(rr.Body.String()), &got)

		/* This test only looks for the tracer just added. The ID should be 1. */
		if got.ID != 1 {
			err = fmt.Errorf("getTracer returned the wrong body in the response. Got ID of %+v, but expected %+v", got.ID, 1)
		}
	}

	/* Return nil to indicate no problems. */
	return err
}

// addTest is the commonly used POST request test.
func addTest(rr *httptest.ResponseRecorder, t *testing.T) error {
	/* Return variable. */
	var err error

	/* Make sure the status code is 200. */
	if status := rr.Code; status != http.StatusOK {
		err = fmt.Errorf("AddTracer returned the wrong status code: got %v, but wanted %v", status, http.StatusOK)
	} else {
		/* Make sure the body is a valid JSON object. */
		got := types.Tracer{}
		json.Unmarshal([]byte(rr.Body.String()), &got)

		/* Sanity checks to make sure the added tracer wasn't empty. */
		if got.ID != 1 {
			err = fmt.Errorf("The inserted tracer has the wrong ID. Expected 1, got: %d", got.ID)
		}
	}

	return err
}
