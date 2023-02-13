package render

import (
	"encoding/gob"
	"github/Atul-Ranjan12/booking/internal/config"
	"github/Atul-Ranjan12/booking/internal/models"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
)

var testApp config.AppConfig;
var session *scs.SessionManager

func TestMain(m *testing.M) {
	// Primitive types for the Session
	gob.Register(models.Reservation{})

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

type myTestWriter struct {}

func (ww *myTestWriter) Header() http.Header{
	var h http.Header
	return h
}

func (ww *myTestWriter) Write(b []byte) (int , error){
	length := len(b)
	return length, nil
}

func (ww *myTestWriter) WriteHeader(statusCode int){
	
}