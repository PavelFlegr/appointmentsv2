package manage

import (
	"appointmentsv2/global"
	"appointmentsv2/services/appointment"
	"appointmentsv2/services/auth"
	appointmentTemplates "appointmentsv2/templates/appointment"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func getIntParam(name string, r *http.Request, w http.ResponseWriter) int {
	val, err := strconv.Atoi(chi.URLParam(r, name))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	return val
}

func ListAppointments(w http.ResponseWriter, r *http.Request) {
	userId, authErr := auth.CheckAuth(r)
	if authErr != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	appointments, err := appointment.ListAppointments(userId)
	if err != nil {
		log.Println(err)
		return
	}

	err = appointmentTemplates.Appointments(appointments).Render(context.Background(), w)
	if err != nil {
		log.Println(err)
	}
}

func CreateAppointment(w http.ResponseWriter, r *http.Request) {
	userId := auth.RequireAuth(r, w)
	name := r.PostFormValue("name")
	timezone := r.PostFormValue("timezone")
	appt := appointment.Appointment{
		Name:     name,
		Timezone: timezone,
	}
	appointment.CreateAppointment(&appt, userId)
	appointmentTemplates.AppointmentItem(appt).Render(context.Background(), w)
}

func UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	appointmentId := getIntParam("appointmentId", r, w)
	auth.RequirePermission(r, w, "appointment", appointmentId, "edit")

	name := r.PostFormValue("name")
	instructions := r.PostFormValue("instructions")

	appt := appointment.Appointment{
		Id:           appointmentId,
		Name:         name,
		Instructions: instructions,
	}
	err := appointment.UpdateAppointment(&appt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, name)
}

func RemoveAppointment(w http.ResponseWriter, r *http.Request) {
	appointmentId := getIntParam("appointmentId", r, w)
	auth.RequirePermission(r, w, "appointment", appointmentId, "manage")

	err := appointment.DeleteAppointment(appointmentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func GetAppointment(w http.ResponseWriter, r *http.Request) {
	appointmentId := getIntParam("appointmentId", r, w)
	auth.RequirePermission(r, w, "appointment", appointmentId, "read")

	appt, err := appointment.GetAppointment(appointmentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	apptLoc, _ := time.LoadLocation(appt.Timezone)

	schedules, err := appointment.ListSchedules(appointmentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i := range schedules {
		schedule := &schedules[i]
		loc, _ := time.LoadLocation(schedule.Timezone)
		start, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", schedule.StartDate, schedule.StartTime), loc)
		end, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", schedule.EndDate, schedule.EndTime), loc)

		schedule.StartDate = start.In(apptLoc).Format("2. 1. 2006")
		schedule.EndDate = end.In(apptLoc).Format("2. 1. 2006")
		schedule.StartTime = start.In(apptLoc).Format("15:04")
		schedule.EndTime = end.In(apptLoc).Format("15:04 MST")
	}

	appointmentTemplates.Appointment(appt, schedules).Render(context.Background(), w)
}

func AddSchedule(w http.ResponseWriter, r *http.Request) {
	appointmentId := getIntParam("appointmentId", r, w)
	userId := auth.RequirePermission(r, w, "appointment", appointmentId, "edit")
	startDate := r.PostFormValue("startDate")
	endDate := r.PostFormValue("endDate")
	startTime := r.PostFormValue("startTime")
	endTime := r.PostFormValue("endTime")
	timezone := r.PostFormValue("timezone")
	length, err := strconv.Atoi(r.PostFormValue("length"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var spots int
	spots, err = strconv.Atoi(r.PostFormValue("spots"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	schedule := appointment.Schedule{
		StartDate:     startDate,
		EndDate:       endDate,
		StartTime:     startTime,
		EndTime:       endTime,
		Timezone:      timezone,
		Length:        length,
		Spots:         spots,
		AppointmentId: appointmentId,
		Status:        "new",
	}

	err = appointment.CreateSchedule(&schedule, userId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	appointmentTemplates.ScheduleItem(schedule).Render(context.Background(), w)
}

func GetSchedule(w http.ResponseWriter, r *http.Request) {
	appointmentId := getIntParam("appointmentId", r, w)
	auth.RequirePermission(r, w, "appointment", appointmentId, "read")
	scheduleId := getIntParam("scheduleId", r, w)

	schedule, err := appointment.GetSchedule(scheduleId, appointmentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	appointmentTemplates.Schedule(schedule).Render(context.Background(), w)

}

func DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	appointmentId := getIntParam("appointmentId", r, w)
	auth.RequirePermission(r, w, "appointment", appointmentId, "edit")
	scheduleId := getIntParam("scheduleId", r, w)

	err := appointment.DeleteSchedule(scheduleId, appointmentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func ActivateSchedule(w http.ResponseWriter, r *http.Request) {
	appointmentId := getIntParam("appointmentId", r, w)
	auth.RequirePermission(r, w, "appointment", appointmentId, "edit")
	scheduleId := getIntParam("scheduleId", r, w)

	schedule, err := appointment.GetSchedule(scheduleId, appointmentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	loc, err := time.LoadLocation(schedule.Timezone)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = appointment.SetScheduleStatus("active", scheduleId, appointmentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	schedule.Status = "active"

	start, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", schedule.StartDate, schedule.StartTime), loc)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	end, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", schedule.EndDate, schedule.EndTime), loc)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	minutes := (end.Hour()-start.Hour())*60 + end.Minute() - start.Minute()
	diff := end.Sub(start)
	days := int(diff.Truncate(time.Duration(time.Hour*24)).Hours() / 24)

	for day := 0; day < days; day += 1 {
		dayOffset := time.Duration(24*day) * time.Hour
		for minute := 0; minute < minutes-1; minute += schedule.Length {
			minuteOffset := time.Duration(minute) * time.Minute
			slot := appointment.Slot{
				AppointmentId: appointmentId,
				Start:         start.Add(dayOffset).Add(minuteOffset),
				End:           start.Add(dayOffset).Add(minuteOffset).Add(time.Minute * time.Duration(schedule.Length)),
				Spots:         schedule.Spots,
				Free:          schedule.Spots,
			}

			err := appointment.CreateSlot(&slot)
			if err != nil {
				log.Println(err)
			}
		}
	}

	appointmentTemplates.ScheduleItem(schedule).Render(context.Background(), w)
}

func UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	appointmentId := getIntParam("appointmentId", r, w)
	auth.RequirePermission(r, w, "appointment", appointmentId, "edit")
	scheduleId := getIntParam("scheduleId", r, w)
	startDate := r.PostFormValue("startDate")
	endDate := r.PostFormValue("endDate")
	startTime := r.PostFormValue("startTime")
	endTime := r.PostFormValue("endTime")
	length, err := strconv.Atoi(r.PostFormValue("length"))
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Length must be a number")
		return
	}
	var spots int
	spots, err = strconv.Atoi(r.PostFormValue("spots"))
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Spots must be a number")
		return
	}

	schedule := appointment.Schedule{
		Id:            scheduleId,
		StartDate:     startDate,
		EndDate:       endDate,
		StartTime:     startTime,
		EndTime:       endTime,
		Length:        length,
		Spots:         spots,
		AppointmentId: appointmentId,
	}

	err = appointment.UpdateSchedule(&schedule)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Failed to update schedule")
		return
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/manage/appointment/%v", appointmentId))
}

func ListAvailableSlots(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	slots, err := appointment.ListAvailableSlots(token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var groups []appointment.SlotGroup
	group := &appointment.SlotGroup{}
	for _, slot := range slots {
		date := slot.Start.Local().Format("2. 1. 2006")
		if date != group.Name {
			group = &appointment.SlotGroup{
				Name:  date,
				Slots: &[]appointment.Slot{},
			}
			groups = append(groups, *group)
		}
		*group.Slots = append(*group.Slots, slot)
	}

	appointmentTemplates.Reserve(groups).Render(context.Background(), w)
}

func GetSlot(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	slot, err := appointment.GetSlotByToken(token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	appointment, err := appointment.GetAppointment(slot.AppointmentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	appointmentTemplates.ReserveSlot(slot, appointment.Token).Render(context.Background(), w)
}

func ReserveSlot(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	firstName := r.PostFormValue("firstName")
	lastName := r.PostFormValue("lastName")
	email := r.PostFormValue("email")
	timezone := r.PostFormValue("timezone")

	slot, err := appointment.GetSlotByToken(token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if slot.Free <= 0 {
		log.Println("no free spot in da slot")
		fmt.Fprint(w, "Sorry this spot is full. try a different one")
		return
	}

	slot.Free -= 1
	err = appointment.ChangeSlotFree(slot.Id, slot.AppointmentId, -1)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reservation := appointment.Reservation{
		AppointmentId: slot.AppointmentId,
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		SlotId:        slot.Id,
		Timezone:      timezone,
	}

	err = appointment.CreateReservation(&reservation)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	appt, err := appointment.GetAppointment(reservation.AppointmentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	style := `box-shadow:inset 0px 1px 0px 0px #f5978e;
	background:linear-gradient(to bottom, #f24537 5%, #c62d1f 100%);
	background-color:#f24537;
	border-radius:6px;
	border:1px solid #d02718;
	display:inline-block;
	cursor:pointer;
	color:#ffffff;
	font-family:Arial;
	font-size:15px;
	font-weight:bold;
	padding:6px 24px;
	text-decoration:none;
	text-shadow:0px 1px 0px #810e05;`

	if err != nil {
		log.Println(err)
	}
	loc, err := time.LoadLocation(reservation.Timezone)
	if err != nil {
		log.Println(err)
	}

	message := fmt.Sprintf("Your reservation for %v is registered. You can cancel it by clicking <a style=\"%v\" href=\"%v/cancel/%v\">here</a><div style=\"min-height: 300px\">%v</div>",
		slot.Start.In(loc).Format("2. 1. 2006 15:04 MST"), style, global.Conf.Host, reservation.Token, appt.Instructions,
	)

	go SendMail(reservation.Email, fmt.Sprintf("Reservation %v created", appt.Name), message)
}

func SendMail(to string, subject string, message string) {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %v\r\n", global.Conf.Email.From)
	msg += fmt.Sprintf("To: %v\r\n", to)
	msg += fmt.Sprintf("Subject: %v\r\n", subject)
	msg += fmt.Sprintf("\r\n%v\r\n", message)
	auth := smtp.PlainAuth("", global.Conf.Email.User, global.Conf.Email.Password, global.Conf.Email.Host)
	err := smtp.SendMail(fmt.Sprintf("%v:%v", global.Conf.Email.Host, global.Conf.Email.Port), auth, global.Conf.Email.From, []string{to}, []byte(msg))
	if err != nil {
		log.Println(err)
	}
}

func ListReservations(w http.ResponseWriter, r *http.Request) {
	appointmentId := getIntParam("appointmentId", r, w)
	auth.RequirePermission(r, w, "appointment", appointmentId, "read")

	appt, _ := appointment.GetAppointment(appointmentId)
	loc, _ := time.LoadLocation(appt.Timezone)

	reservations, err := appointment.ListReservations(appointmentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i := range reservations {
		reservation := &reservations[i]
		start, _ := time.Parse(time.RFC3339, reservation.Start)
		end, _ := time.Parse(time.RFC3339, reservation.End)
		reservation.Start = start.In(loc).Format("2. 1. 2006 15:04 MST")
		reservation.End = end.In(loc).Format("2. 1. 2006 15:04 MST")
	}

	appointmentTemplates.Reservations(reservations).Render(context.Background(), w)
}

func DeleteReservation(w http.ResponseWriter, r *http.Request) {
	appointmentId := getIntParam("appointmentId", r, w)
	reservationId := getIntParam("reservationId", r, w)

	auth.RequirePermission(r, w, "appointment", appointmentId, "edit")

	reservation, err := appointment.GetReservation(reservationId, appointmentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = appointment.DeleteReservation(reservationId, appointmentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = appointment.ChangeSlotFree(reservation.SlotId, appointmentId, 1)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func CancelReservation(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	reservation, err := appointment.GetReservationByToken(token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = appointment.DeleteReservation(reservation.Id, reservation.AppointmentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = appointment.ChangeSlotFree(reservation.SlotId, reservation.AppointmentId, 1)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
