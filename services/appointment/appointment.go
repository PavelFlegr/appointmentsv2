package appointment

import (
	"appointmentsv2/global"
	"appointmentsv2/services/auth"
	"errors"
	"log"
	"time"
)

type Appointment struct {
	Id           int    `db:"id"`
	Name         string `db:"name"`
	Instructions string `db:"instructions"`
	Token        string `db:"token"`
	Timezone     string `db:"timezone"`
}

type Schedule struct {
	Id            int    `db:"id"`
	StartDate     string `db:"start_date"`
	EndDate       string `db:"end_date"`
	StartTime     string `db:"start_time"`
	EndTime       string `db:"end_time"`
	Timezone      string `db:"timezone"`
	Length        int    `db:"length"`
	Spots         int    `db:"spots"`
	Status        string `db:"status"`
	AppointmentId int    `db:"appointment_id"`
}

type Slot struct {
	Id            int       `db:"id"`
	Start         time.Time `db:"start"`
	End           time.Time `db:"end"`
	Spots         int       `db:"spots"`
	Free          int       `db:"free"`
	AppointmentId int       `db:"appointment_id"`
	Token         string    `db:"token"`
}

type Reservation struct {
	Id            int    `db:"id"`
	FirstName     string `db:"first_name"`
	LastName      string `db:"last_name"`
	Email         string `db:"email"`
	SlotId        int    `db:"slot_id"`
	AppointmentId int    `db:"appointment_id"`
	Timezone      string `db:"timezone"`
	Token         string `db:"token"`
	Start         string
	End           string
}

type SlotGroup struct {
	Name  string
	Slots *[]Slot
}

func CreateAppointment(appointment *Appointment, userId int) error {
	err := global.Db.QueryRowx("insert into appointments (name, instructions, user_id, timezone) values ($1, $2, $3, $4) returning id, token", &appointment.Name, &appointment.Instructions, &userId, &appointment.Timezone).Scan(&appointment.Id, &appointment.Token)
	if err != nil {
		log.Println(err)
		return err
	}

	auth.AddPermission(userId, "appointment", appointment.Id, "manage")
	auth.AddPermission(userId, "appointment", appointment.Id, "edit")
	auth.AddPermission(userId, "appointment", appointment.Id, "read")

	return nil
}

func ListAppointments(userId int) ([]Appointment, error) {
	var appointments []Appointment
	err := global.Db.Select(&appointments, "select id, name, instructions, token, timezone from appointments inner join appointment_permissions sp on appointments.id = sp.entity_id where sp.user_id = $1 and sp.action = 'read'", userId)
	if err != nil {
		log.Println(err)
		return appointments, err
	}

	return appointments, nil
}

func DeleteAppointment(appointmentId int) error {
	_, err := global.Db.Exec("delete from appointments where id = $1", appointmentId)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetAppointment(appointmentId int) (Appointment, error) {
	var appointment Appointment
	err := global.Db.QueryRowx("select id, name, instructions, token, timezone from appointments where id = $1", appointmentId).StructScan(&appointment)
	if err != nil {
		log.Println(err)
		return appointment, err
	}

	return appointment, nil
}

func GetAppointmentByToken(token string) (Appointment, error) {
	var appointment Appointment
	err := global.Db.QueryRowx("select id, name, instructions, token, timezone from appointments where token = $1", token).StructScan(&appointment)
	if err != nil {
		log.Println(err)
		return appointment, err
	}

	return appointment, nil
}

func UpdateAppointment(appointment *Appointment) error {
	_, err := global.Db.Exec("update appointments set name = $1, instructions = $2 where id = $3", appointment.Name, appointment.Instructions, appointment.Id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func CreateSchedule(schedule *Schedule, userId int) error {
	err := global.Db.QueryRowx("insert into schedules (start_date, end_date, start_time, end_time, timezone, length, spots, status, appointment_id, user_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id",
		&schedule.StartDate, &schedule.EndDate, &schedule.StartTime, &schedule.EndTime, &schedule.Timezone, &schedule.Length, &schedule.Spots, &schedule.Status, &schedule.AppointmentId, &userId,
	).Scan(&schedule.Id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetSchedule(scheduleId int, appointmentId int) (Schedule, error) {
	var schedule Schedule
	err := global.Db.QueryRowx("select id, start_date, end_date, start_time, end_time, timezone, length, spots, status, appointment_id from schedules where id = $1 and appointment_id = $2", scheduleId, appointmentId).StructScan(&schedule)
	if err != nil {
		log.Println(err)
		return schedule, err
	}

	return schedule, nil
}

func ListSchedules(appointmentId int) ([]Schedule, error) {
	var schedules []Schedule
	err := global.Db.Select(&schedules, "select id, start_date, end_date, start_time, end_time, timezone, length, spots, status, appointment_id from schedules where appointment_id = $1", &appointmentId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return schedules, err
}

func DeleteSchedule(scheduleId int, appointmentId int) error {
	_, err := global.Db.Exec("delete from schedules where id = $1 and appointment_id = $2", &scheduleId, &appointmentId)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func SetScheduleStatus(status string, scheduleId int, appointmentId int) error {
	res, err := global.Db.Exec("update schedules set status = $1 where id = $2 and appointment_id = $3", status, scheduleId, appointmentId)
	if err != nil {
		log.Println(err)
		return err
	}
	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}
	if rows == 0 {
		return errors.New("no schedule found")
	}
	return nil
}

func UpdateSchedule(schedule *Schedule) error {
	_, err := global.Db.Exec("update schedules set start_date = $1, end_date = $2, start_time = $3, end_time = $4, length = $5, spots = $6 where id = $7 and appointment_id = $8",
		schedule.StartDate, schedule.EndDate, schedule.StartTime, schedule.EndTime, schedule.Length, schedule.Spots, schedule.Id, schedule.AppointmentId,
	)

	if err != nil {
		log.Println(err)
		return err
	}

	return err
}

func CreateSlot(slot *Slot) error {
	_, err := global.Db.Exec("insert into slots (start, \"end\", appointment_id, spots, free) values ($1, $2, $3, $4, $5)", slot.Start, slot.End, slot.AppointmentId, slot.Spots, slot.Free)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetSlot(slotId int, appointmentId int) (Slot, error) {
	var slot Slot
	err := global.Db.QueryRowx("select id, start, \"end\", appointment_id, spots, free, token from slots where id = $1 and appointment_id = $2", slotId, appointmentId).StructScan(&slot)
	if err != nil {
		log.Println(err)
		return slot, err
	}

	return slot, nil
}

func GetSlotByToken(token string) (Slot, error) {
	var slot Slot
	err := global.Db.QueryRowx("select id, start, \"end\", appointment_id, spots, free, token from slots where token = $1", token).StructScan(&slot)
	if err != nil {
		log.Println(err)
		return slot, err
	}

	return slot, nil
}

func ChangeSlotFree(slotId int, appointmentId, delta int) error {
	_, err := global.Db.Exec("update slots set free = free + $1 where id = $2 and appointment_id = $3", delta, slotId, appointmentId)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func ListAvailableSlots(token string) ([]Slot, error) {
	var slots []Slot
	err := global.Db.Select(&slots, "select s.id, s.start, s.\"end\", s.appointment_id, s.free, s.token from slots s join appointments a on a.id = s.appointment_id where a.token = $1 and s.free > 0 order by s.start", token)
	if err != nil {
		log.Println(err)
		return slots, err
	}

	return slots, nil
}

func CreateReservation(reservation *Reservation) error {
	err := global.Db.QueryRowx("insert into reservations (first_name, last_name, email, appointment_id, slot_id, timezone) values ($1, $2, $3, $4, $5, $6) returning id, token",
		reservation.FirstName, reservation.LastName, reservation.Email, reservation.AppointmentId, reservation.SlotId, reservation.Timezone,
	).Scan(&reservation.Id, &reservation.Token)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetReservation(reservationId int, appointmentId int) (Reservation, error) {
	var reservation Reservation
	err := global.Db.QueryRowx("select id, first_name, last_name, email, appointment_id, slot_id from reservations where id = $1 and appointment_id = $2", reservationId, appointmentId).StructScan(&reservation)
	if err != nil {
		log.Println(err)
		return reservation, err
	}

	return reservation, nil
}

func GetReservationByToken(token string) (Reservation, error) {
	var reservation Reservation
	err := global.Db.QueryRowx("select id, first_name, last_name, email, appointment_id, slot_id, token from reservations where token = $1", token).StructScan(&reservation)
	if err != nil {
		log.Println(err)
		return reservation, err
	}

	return reservation, nil
}

func ListReservations(appointmentId int) ([]Reservation, error) {
	var reservations []Reservation
	err := global.Db.Select(&reservations, "select r.id, r.first_name, r.last_name, r.email, r.appointment_id, r.slot_id, s.start, s.\"end\" from reservations r inner join slots s on s.id = r.slot_id where r.appointment_id = $1", appointmentId)
	if err != nil {
		log.Println(err)
		return reservations, err
	}

	return reservations, nil
}

func DeleteReservation(reservationId int, appointmentId int) error {
	_, err := global.Db.Exec("delete from reservations where id = $1 and appointment_id = $2", reservationId, appointmentId)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
