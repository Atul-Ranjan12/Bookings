package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Vaid(t *testing.T){
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("Form is invalid when it should have been valids")
	}
}

func TestForm_Required(t *testing.T){
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form should have been invalid as a, b and c is required")
	}

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "a")
	postData.Add("c", "a")

	r = httptest.NewRequest("POST", "/whatever", nil)
	r.PostForm = postData
	form = New(r.PostForm)

	form.Required("a", "b", "c")
	if !form.Valid(){
		t.Error("Added the required fields but the form is still invalid")
	}
}

func TestForm_IsEmailValid(t *testing.T){
	postData := url.Values{}
	postData.Add("email", "atul@atul.com")
	r := httptest.NewRequest("POST", "/whatever", nil)
	r.PostForm = postData
	form := New(r.PostForm)

	form.IsEmailValid("email")
	if !form.Valid() {
		t.Error("Got invalid email with a valid email address")
	}

	postData = url.Values{}
	postData.Add("email", "atulranjan")
	r = httptest.NewRequest("POST", "/whatever", nil)
	r.PostForm = postData
	form = New(r.PostForm)

	form.IsEmailValid("email")
	if form.Valid() {
		t.Error("Got Valid Email address with an invalid email")
	}
}

func createPostData(key string, value string) *http.Request{
	postData:= url.Values{}
	postData.Add(key, value)
	r := httptest.NewRequest("POST", "/whatever", nil)
	r.PostForm = postData
	return r
}

func TestForm_MinLength(t *testing.T){
	r := createPostData("first_name", "Atul")
	form := New(r.PostForm)

	form.MinLength("first_name", 3, r)
	if !form.Valid() {
		t.Error("Form not valid for valid name")
	}

	r = createPostData("last_name", "s")
	form = New(r.PostForm)

	form.MinLength("last_name", 3, r)
	if form.Valid() {
		t.Error("Form was valid for invalid input")
	}
}

func TestForm_Has(t *testing.T){
	r := createPostData("first_name", "Dummy")
	form := New(r.PostForm)

	form.Has("first_name")
	if !form.Valid() {
		t.Error("Got Invalid form even if it has the required field")
	}

	r = httptest.NewRequest("POST", "/whatever", nil)
	form = New(r.PostForm)

	form.Has("a")
	if form.Valid() {
		t.Error("Got Valid form even if it does not have required field")
	}
}