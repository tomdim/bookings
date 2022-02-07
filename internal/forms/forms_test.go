package forms

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/test", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestNew(t *testing.T) {
	r := httptest.NewRequest("POST", "/test", nil)
	form := New(r.PostForm)
	switch v := interface{}(form).(type) {
	case *Form:
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not *Form, but is %T", v))
	}
}

func TestForm_Required(t *testing.T) {
	formData := url.Values{}
	form := New(formData)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("got valid when should be invalid")
	}

	formData.Add("a", "test")
	formData.Add("b", "test")
	formData.Add("c", "test")

	form = New(formData)
	form.Required("a", "b", "c")
	fmt.Println(form.Values)
	if !form.Valid() {
		t.Error("got invalid when should be valid")
	}
}

func TestForm_Has(t *testing.T) {
	formData := url.Values{}
	form := New(formData)

	has := form.Has("a")
	if has {
		t.Error("Has returned true when should return false")
	}

	formData.Add("a", "test")
	form.Values = formData
	has = form.Has("a")
	if !has {
		t.Error("Has returned false when should return true")
	}
}

func TestForm_MinLength(t *testing.T) {
	formData := url.Values{}
	form := New(formData)

	hasMinLength := form.MinLength("a", 3)
	if hasMinLength {
		t.Error("MinLength returned true when should return false")
	}

	formData.Add("a", "te")
	form.Values = formData

	hasMinLength = form.MinLength("a", 3)
	if hasMinLength {
		t.Error("MinLength returned true when should return false")
	}

	formData = url.Values{}
	formData.Add("a", "test")
	form.Values = formData

	hasMinLength = form.MinLength("a", 3)
	if !hasMinLength {
		t.Error("MinLength returned false when should return true")
	}
}

func TestForm_IsEmail(t *testing.T) {
	formData := url.Values{}
	form := New(formData)

	isEmail := form.IsEmail("a")
	if isEmail {
		t.Error("IsEmail returned true when should return false")
	}

	formData.Add("a", "te")
	form.Values = formData

	isEmail = form.IsEmail("a")
	if isEmail {
		t.Error("IsEmail returned true when should return false")
	}

	formData = url.Values{}
	formData.Add("a", "test@test.com")
	form.Values = formData

	isEmail = form.IsEmail("a")
	if !isEmail {
		t.Error("IsEmail returned false when should return true")
	}
}

//// IsEmail checks if a form field's value is a valid email
//func (f *Form) IsEmail(field string) bool {
//	if !govalidator.IsEmail(f.Get(field)) {
//		f.Errors.Add(field, "Invalid email address.")
//		return false
//	}
//	return true
//}
