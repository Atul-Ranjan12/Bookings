package handlers

import (
	"encoding/json"
	"errors"
	"github/Atul-Ranjan12/booking/internal/config"
	"github/Atul-Ranjan12/booking/internal/driver"
	"github/Atul-Ranjan12/booking/internal/forms"
	"github/Atul-Ranjan12/booking/internal/helpers"
	"github/Atul-Ranjan12/booking/internal/models"
	"github/Atul-Ranjan12/booking/internal/render"
	"github/Atul-Ranjan12/booking/internal/repository"
	"github/Atul-Ranjan12/booking/internal/repository/dbrepo"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := map[string]string{}
	stringMap["test"] = "Hello Again"

	stringMap["remote_ip"] = m.App.Session.GetString(r.Context(), "remote_ip")

	render.Template(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("type assertion to string failed"))
		return
	}
	// Casting start date and end date to the right format
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	// stroring room information
	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
	}
	res.Room.RoomName = room.RoomName

	// Store reservation in the session
	m.App.Session.Put(r.Context(), "reservation", res)

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// Handles Posting of a Reservation Form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("can't get reservation from the session in PostReservation"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.MinLength("last_name", 3, r)
	form.IsEmailValid("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// Write information to the database
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Fill in the restrictions table in the database
	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}
	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Put the reservation back in the context
	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Availibility(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availibility.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// Post Availability
func (m *Repository) PostAvailibility(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
	}

	rooms, err := m.DB.SearchAvailibilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		// No rooms are available
		log.Println("Reached inside the if statement to handle no available rooms")
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	// Storing the Reservation in the session
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	RoomID    string `json:"room_id"`
}

// JSON for availability
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	// Get string values for start date and end date
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	rid := r.Form.Get("room_id")

	//Check if dates are parsing
	log.Println("CHECK: Dates and room ID", sd, ed, rid)

	startDate, _ := time.Parse("2006-01-02", sd)
	endDate, _ := time.Parse("2006-01-02", ed)
	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	// Check if the room is available
	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		log.Println("Error searching for the availibility function")
		return
	}

	// Populate the JSON Response
	response := jsonResponse{
		OK:        available,
		Message:   "Available",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    rid,
	}

	out, err := json.MarshalIndent(response, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
	}

	log.Println(string(out))

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Reservation Summary:
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.ErrorLog.Println("Cannot get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("type assertion to string failed"))
		return
	}
	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// Handler to redirect to the make-reservatin page and pass query string parameters
func (m *Repository) BookNow(w http.ResponseWriter, r *http.Request) {
	// Get id startdate and enddate from query string
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	// Parse the dates
	var layout string = "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	// get the room as well
	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Create a reservation instance
	var res models.Reservation
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate
	res.Room = room

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
