package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tomdim/bookings/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close() // ts will live until the end of execution of testHandlers function

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_PostReservation(t *testing.T) {
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, "2050-01-01")
	endDate, _ := time.Parse(layout, "2050-01-02")

	// reservation to put in session
	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	// post data
	reqBody := "first_name=John"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// test case reservation is not in the session (reset everything)
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for missing post body request
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid form data
	reqBody = "first_name=J"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=j@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned unexpected response cod for invalid form data: got %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// test for failure in insert reservation into db
	reqBody = "first_name=John"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=j@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	reservation.RoomID = 2
	reservation.Room.ID = 2
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned unexpected response code for failed reservation insertion: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for failure in insert room restriction into db
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	reservation.RoomID = 1000
	reservation.Room.ID = 1000
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned unexpected response code for failed room restriction insertion: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// test case reservation is not in the session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation summary handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// test case reservation summary reservation is not in the session
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_AvailabilityJSON(t *testing.T) {
	var j jsonResponse
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	// invalid (missing) form test case
	// create request
	req, _ := http.NewRequest("POST", "/search-availability-json", nil)

	// get context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get response recorder
	rr := httptest.NewRecorder()

	// make request to handler func
	handler.ServeHTTP(rr, req)

	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json response")
	}
	if j.OK {
		t.Error("json response ok - invalid post form: expected false but got true")
	}
	if j.Message != "Internal server error" {
		t.Errorf("json response message - invalid post form: expected `%s`, got `%s`", "Internal server error", j.Message)
	}

	// invalid start date test case
	reqBody := "start=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	// create request
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get response recorder
	rr = httptest.NewRecorder()

	// make request to handler func
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json response")
	}
	if j.OK {
		t.Error("json response ok - invalid start date test case: expected false but got true")
	}
	if j.Message != "Invalid start date format" {
		t.Errorf("json response message - invalid start date test case: expected `%s`, got `%s`", "Invalid start date format", j.Message)
	}

	// invalid end date test case
	reqBody = "start=2050-05-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	// create request
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get response recorder
	rr = httptest.NewRecorder()

	// make request to handler func
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json response")
	}
	if j.OK {
		t.Error("json response ok - invalid end date test case: expected false but got true")
	}
	if j.Message != "Invalid end date format" {
		t.Errorf("json response message - invalid end date test case: expected `%s`, got `%s`", "Invalid end date format", j.Message)
	}

	// invalid room id test case
	reqBody = "start=2050-05-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-05-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")

	// create request
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get response recorder
	rr = httptest.NewRecorder()

	// make request to handler func
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json response")
	}
	if j.OK {
		t.Error("json response ok - invalid room id test case: expected false but got true")
	}
	if j.Message != "Invalid room id format" {
		t.Errorf("json response message - invalid room id test case: expected `%s`, got `%s`", "Invalid room id format", j.Message)
	}

	// error connecting to db test case
	reqBody = "start=2050-05-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-05-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000")

	// create request
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get response recorder
	rr = httptest.NewRecorder()

	// make request to handler func
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json response")
	}
	if j.OK {
		t.Error("json response ok - db connection error test case: expected false but got true")
	}
	if j.Message != "Error connecting to database" {
		t.Errorf("json response message - invalid post form: expected `%s`, got `%s`", "Error connecting to database", j.Message)
	}

	// happy path test case
	reqBody = "start=2050-05-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-05-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	// create request
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get response recorder
	rr = httptest.NewRecorder()

	// make request to handler func
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json response")
	}
	if !j.OK {
		t.Error("json response ok - happy path test case: expected true but got false")
	}
	if j.Message != "" {
		t.Errorf("json response message - happy path test case: expected `%s`, got `%s`", "", j.Message)
	}
	if j.StartDate != "2050-05-01" {
		t.Errorf("json response start date - happy path test case: expected `%s`, got `%s`", "2050-05-01", j.StartDate)
	}
	if j.EndDate != "2050-05-02" {
		t.Errorf("json response end date - happy path test case: expected `%s`, got `%s`", "2050-05-02", j.EndDate)
	}
	if j.RoomID != "2" {
		t.Errorf("json response room id - happy path test case: expected `%s`, got `%s`", "2", j.RoomID)
	}
}

func TestRepository_PostAvailability(t *testing.T) {
	// happy path test case
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")

	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PostAvailability handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	// invalid (missing) post form data
	req, _ = http.NewRequest("POST", "/search-availability", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// invalid start date test case
	reqBody = "start=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// invalid end date test case
	reqBody = "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=invalid")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// db connection error test case
	reqBody = "start=2000-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostAvailability handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// no rooms found (no availability) test case
	reqBody = "start=2050-02-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-02-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_BookRoom(t *testing.T) {

	req, _ := http.NewRequest("GET", "/book-room?id=2&s=2050-01-01&e=2050-01-02", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Book room handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// invalid room id test case
	req, _ = http.NewRequest("GET", "/book-room?id=invalid&s=2050-01-01&e=2050-01-02", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Book room handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// invalid start date test case
	req, _ = http.NewRequest("GET", "/book-room?id=1&s=invalid&e=2050-01-02", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Book room handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// invalid end date test case
	req, _ = http.NewRequest("GET", "/book-room?id=1&s=2050-01-01&e=invalid", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Book room handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// db connection error test case
	req, _ = http.NewRequest("GET", "/book-room?id=3&s=2050-01-01&e=2050-01-02", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Book room handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_ChooseRoom(t *testing.T) {
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, "2050-01-01")
	endDate, _ := time.Parse(layout, "2050-01-02")

	// reservation to put in session
	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	req, _ := http.NewRequest("GET", "/choose-room", nil)
	req.RequestURI = "/choose-room/2"
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Choose room handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusSeeOther)
	}

	// invalid room id on request uri test case
	req, _ = http.NewRequest("GET", "/choose-room", nil)
	req.RequestURI = "/choose-room/invalid"
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Choose room handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// reservation not found on session test case
	req, _ = http.NewRequest("GET", "/choose-room", nil)
	req.RequestURI = "/choose-room/2"
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Choose room handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// db connection error test case
	req, _ = http.NewRequest("GET", "/choose-room", nil)
	req.RequestURI = "/choose-room/3"
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Choose room handler returned unexpected response code: got %d, expected %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
