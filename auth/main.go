// auth/main.go
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/streadway/amqp"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

type UserJSON struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Auth struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var db *sql.DB

func initDB() {
	var err error
	connStr := mysql.Config{
		User:   "root",
		Passwd: "root",
		Net:    "tcp",
		Addr:   "auth_db:3306",
		DBName: "ms_auth",
	}

	db, err = sql.Open("mysql", connStr.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
}

func publishUserCreated(user User) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"user_created", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	body, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Fatal(err)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	var userBody UserJSON
	err := json.NewDecoder(r.Body).Decode(&userBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userBody.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	var auth Auth
	auth.Password = string(hashedPassword)
	auth.Username = userBody.Username
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", auth.Username, auth.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := serializeUser(userBody)
	publishUserCreated(user)
	w.WriteHeader(http.StatusCreated)
}

func login(w http.ResponseWriter, r *http.Request) {
	var creds Auth
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user Auth
	err = db.QueryRow("SELECT id, password FROM users WHERE username=?", creds.Username).Scan(&user.ID, &user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    tokenString,
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
}

func main() {
	initDB()
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	log.Println("Auth service listening on :5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

// serializeUser converts a User struct to a JSON string
func serializeUser(userReq UserJSON) User {
	return User{
		Username:  userReq.Username,
		Email:     userReq.Email,
		FirstName: userReq.FirstName,
		LastName:  userReq.LastName,
	}
}
