package main

import (
	"appointmentsv2/global"
	"appointmentsv2/handlers/auth"
	"appointmentsv2/handlers/manage"
	"appointmentsv2/templates"
	"context"

	"fmt"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/securecookie"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	log.SetFlags(log.Llongfile | log.Ltime | log.Ldate | log.Lmsgprefix)

	global.Conf = &global.Config{
		Port:    80,
		ConnStr: "postgresql://postgres:postgres@localhost/appointments?sslmode=disable",
		HashKey: "iisdzb75MHcyAo4xPDSLsPVakm9saAs7",
		Host:    "localhost",
	}
	toml.DecodeFile("./conf.toml", global.Conf)
	var db *sqlx.DB
	db, err := sqlx.Open("postgres", global.Conf.ConnStr)
	if err != nil {
		log.Fatalln(err)
	}
	global.Db = db
	global.Sc = securecookie.New([]byte(global.Conf.HashKey), nil)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	RegisterRoutes(r)

	fmt.Printf("starting on port %v\n", global.Conf.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", global.Conf.Port), r)
}

func RegisterRoutes(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		templates.Index().Render(context.Background(), w)
	})

	r.Get("/manage/appointment", manage.ListAppointments)
	r.Post("/manage/appointment", manage.CreateAppointment)
	r.Delete("/manage/appointment/{appointmentId}", manage.RemoveAppointment)
	r.Get("/manage/appointment/{appointmentId}", manage.GetAppointment)
	r.Put("/manage/appointment/{appointmentId}", manage.UpdateAppointment)

	r.Get("/manage/appointment/{appointmentId}/schedule/{scheduleId}", manage.GetSchedule)
	r.Post("/manage/appointment/{appointmentId}/schedule", manage.AddSchedule)
	r.Put("/manage/appointment/{appointmentId}/schedule/{scheduleId}", manage.UpdateSchedule)
	r.Delete("/manage/appointment/{appointmentId}/schedule/{scheduleId}", manage.DeleteSchedule)
	r.Post("/manage/appointment/{appointmentId}/schedule/{scheduleId}/activate", manage.ActivateSchedule)

	r.Get("/manage/appointment/{appointmentId}/reservation", manage.ListReservations)
	r.Delete("/manage/appointment/{appointmentId}/reservation/{reservationId}", manage.DeleteReservation)

	r.Get("/reserve/{token}", manage.ListAvailableSlots)
	r.Get("/reserve/slot/{token}", manage.GetSlot)
	r.Post("/reserve/slot/{token}", manage.ReserveSlot)
	r.Get("/cancel/{token}", manage.CancelReservation)

	r.Post("/login", auth.Login)
	r.Get("/logout", auth.Logout)
	r.Post("/register", auth.Register)
}
