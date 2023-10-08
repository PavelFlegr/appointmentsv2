package auth

import (
	"appointmentsv2/global"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Id       int    `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func CreateUser(user *User) error {
	log.Println("test")
	err := global.Db.QueryRowx("insert into users (email, password) values ($1, $2) returning id", &user.Email, &user.Password).Scan(&user.Id)
	if err != nil {
		log.Println(err)
	}

	return err
}

func GetUserByEmail(email string) (User, error) {
	var user User
	err := global.Db.QueryRowx("select id, email, password from users where email = $1", email).StructScan(&user)
	if err != nil {
		log.Println(err)
	}

	return user, err
}

func AddPermission(userId int, entityType string, entityId int, action string) {
	query := fmt.Sprintf("insert into %v_permissions (user_id, action, entity_id) VALUES ($1, $2, $3)", entityType)
	_, err := global.Db.Exec(query, userId, action, entityId)
	if err != nil {
		log.Println(err)
	}
}

func HasPermission(userId int, entityType string, entityId int, action string) bool {
	query := fmt.Sprintf("select exists(select 1 from %v_permissions where user_id = $1 and entity_id = $2 and action = $3)", entityType)
	var found bool
	err := global.Db.QueryRow(query, userId, entityId, action).Scan(&found)
	if err != nil {
		log.Println(err)
	}
	return found
}

func RemovePermission(userId int, entityType string, entityId int, action string) {
	query := fmt.Sprintf("delete from %v_permissions where user_id = $1 and entity_id = $2 and action = $3", entityType)
	_, err := global.Db.Exec(query, userId, entityId, action)
	if err != nil {
		log.Println(err)
	}
}

func CheckAuth(r *http.Request) (int, error) {
	cookie, err := r.Cookie("userId")
	if err != nil {
		log.Println("Authentication failed", err)
		return 0, err
	}
	var value int
	err = global.Sc.Decode("userId", cookie.Value, &value)
	if err != nil {
		log.Println("Authentication failed", err)
	}

	return value, err
}

func RequireAuth(r *http.Request, w http.ResponseWriter) int {
	userId, err := CheckAuth(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		panic(err)
	}
	return userId
}

func RequirePermission(r *http.Request, w http.ResponseWriter, entityType string, entityId int, action string) int {
	userId := RequireAuth(r, w)
	if !HasPermission(userId, entityType, entityId, action) {
		w.WriteHeader(http.StatusForbidden)
		panic(fmt.Sprintln("user %v can't %v on %v %v", userId, action, entityType, entityId))
	}

	return userId
}
