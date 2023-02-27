package main

import (
	"github/Atul-Ranjan12/booking/internal/config"
	"github/Atul-Ranjan12/booking/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quaters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)

	mux.Get("/search-availability", handlers.Repo.Availibility)
	mux.Post("/search-availability", handlers.Repo.PostAvailibility)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)

	// Route when the link is clicked
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)

	// Route to book now page
	mux.Get("/book-room", handlers.Repo.BookNow)

	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	// Route to the user login page
	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)

	// Route to log the user out
	mux.Get("/user/logout", handlers.Repo.Logout)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		// mux.Use(Auth)

		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		mux.Get("/reservations-all", handlers.Repo.AdminAllReservations)
		mux.Get("/reservations-calendar", handlers.Repo.AdminReservationCalender)
	})

	return mux
}
