package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/shan251197/booking-and-reservation/internal/config"
	"github.com/shan251197/booking-and-reservation/internal/forms"
	"github.com/shan251197/booking-and-reservation/internal/models"
	"github.com/shan251197/booking-and-reservation/internal/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})
}
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "hello,shan"
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})

}
func (m *Repository) Reservation(rw http.ResponseWriter, r *http.Request) {

	var emptyReservation models.Reservation
	data := make(map[string]interface{})

	data["reservation"] = emptyReservation

	render.RenderTemplate(rw, r, "make-reservation.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}
func (m *Repository) PostReservation(rw http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
	}

	form := forms.New(r.PostForm)

	// form.Has("first_name", r)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(rw, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(rw, r, "/reservation-summary", http.StatusSeeOther)

}
func (m *Repository) Generals(rw http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(rw, r, "generals.page.html", &models.TemplateData{})
}

func (m *Repository) Majors(rw http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(rw, r, "majors.page.html", &models.TemplateData{})
}
func (m *Repository) Availability(rw http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(rw, r, "search-availability.page.html", &models.TemplateData{})
}
func (m *Repository) PostAvailability(rw http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	csrfToken := r.Form.Get("csrf_token")
	rw.Write([]byte(fmt.Sprintf("start date %s and end date %s", start, end)))
	rw.Write([]byte(fmt.Sprintf("csrf token %s", csrfToken)))
}

type JsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) AvailabilityJSON(rw http.ResponseWriter, r *http.Request) {
	resp := JsonResponse{
		OK:      true,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		log.Println(err)
	}
	rw.Header().Set("content-type", "Application/json")
	rw.Write(out)
}
func (m *Repository) Contact(rw http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(rw, r, "contact.page.html", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(rw http.ResponseWriter, r *http.Request) {

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		log.Println("cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplate(rw, r, "reservation-summary.page.html", &models.TemplateData{
		Data: data,
	})
}
