package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/tomdim/bookings/internal/config"
	"github.com/tomdim/bookings/internal/models"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

type myWriter struct{}

func TestMain(m *testing.M) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// Change this to true when in production
	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

func (mw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (mw *myWriter) WriteHeader(i int) {
}

func (mw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
