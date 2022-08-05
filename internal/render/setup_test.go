package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/shan251197/booking-and-reservation/internal/config"
	"github.com/shan251197/booking-and-reservation/internal/models"
)

var session *scs.SessionManager

var testApp config.AppConfig

func TestMain(m *testing.M) {

	gob.Register(models.Reservation{})

	testApp.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = false
	session.Cookie.SameSite = http.SameSiteLaxMode

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct{}

func (w *myWriter) Header() http.Header {
	var h http.Header

	return h
}

func (w *myWriter) Write(b []byte) (int, error) {
	length := len(b)

	return length, nil
}

func (w *myWriter) WriteHeader(statusCode int) {

}
