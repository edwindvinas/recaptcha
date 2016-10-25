// Package recaptcha is google golang module for google re-captcha.
//
// Installation
//
//   go get github.com/haisum/recaptcha
//
// Usage
//
// Usage example can be found in example/main.go file.
//
//
// Source code
//
// Available on github: http://github.com/haisum/recaptcha
//
// Author: Haisum (haisumbhatti@gmail.com)
//----------------------------------------------
// REV001 - Updated to use urlfetch
// Author: Edwin (edwin.d.vinas@gmail.com)
//----------------------------------------------

package recaptcha

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	//REV001
	"appengine/urlfetch"
	
)

// R type represents an object of Recaptcha and has public property Secret,
// which is secret obtained from google recaptcha tool admin interface
type R struct {
	Secret    string
	lastError []string
}

// Struct for parsing json in google's response
type googleResponse struct {
	Success    bool
	ErrorCodes []string `json:"error-codes"`
}

// url to post submitted re-captcha response to
var postURL = "https://www.google.com/recaptcha/api/siteverify"

// Verify method, verifies if current request have valid re-captcha response and returns true or false
// This method also records any errors in validation.
// These errors can be received by calling LastError() method.
func (r *R) Verify(req http.Request) bool {
	//REV001
	c := appengine.NewContext(req)
	r.lastError = make([]string, 1)
	response := req.FormValue("g-recaptcha-response")
	//REV001
	//client := &http.Client{Timeout: 20 * time.Second}
	client := urlfetch.Client(c)
	resp, err := client.PostForm(postURL,
		url.Values{"secret": {r.Secret}, "response": {response}})
	if err != nil {
		r.lastError = append(r.lastError, err.Error())
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.lastError = append(r.lastError, err.Error())
		return false
	}
	gr := new(googleResponse)
	err = json.Unmarshal(body, gr)
	if err != nil {
		r.lastError = append(r.lastError, err.Error())
		return false
	}
	if !gr.Success {
		r.lastError = append(r.lastError, gr.ErrorCodes...)
	}
	return gr.Success
}

// LastError returns errors occurred in last re-captcha validation attempt
func (r R) LastError() []string {
	return r.lastError
}
