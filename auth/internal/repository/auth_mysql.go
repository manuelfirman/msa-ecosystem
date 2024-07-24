package repository

import (
	"auth/internal/domain"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

// AuthMySQL represents the MySQL implementation of the AuthRepository
type AuthMySQL struct {
	DB *sql.DB
}

// NewAuthMySQL creates a new instance of AuthMySQL
func NewAuthMySQL(db *sql.DB) *AuthMySQL {
	return &AuthMySQL{DB: db}
}

func (am *AuthMySQL) Login(auth *domain.Auth) (err error) {
	err = am.DB.QueryRow("SELECT id, password FROM users WHERE username=?", auth.Username).Scan(&auth.ID, &auth.Password)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1146:
				err = ErrRepositoryTableNotFound
				return
			default:
				return
			}
		}
	}

	if auth.ID == 0 {
		err = ErrRepositoryEntryNotFound
	}

	return
}

func (am *AuthMySQL) Register(auth *domain.Auth) (err error) {
	result, err := am.DB.Exec("INSERT INTO users (username, password, unique_id) VALUES (?, ?, ?)", auth.Username, auth.Password, auth.UniqueID)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				return ErrRepositoryDuplicateEntry
			default:
				return err
			}
		}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	auth.ID = int(id)
	return
}

// serializeUser converts a User struct to a JSON string
// func serializeUser(userReq UserJSON) User {
// 	return User{
// 		Username:  userReq.Username,
// 		Email:     userReq.Email,
// 		FirstName: userReq.FirstName,
// 		LastName:  userReq.LastName,
// 	}
// }

// func publishUserCreated(user User) {
// 	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer ch.Close()

// 	q, err := ch.QueueDeclare(
// 		"user_created", // name
// 		false,          // durable
// 		false,          // delete when unused
// 		false,          // exclusive
// 		false,          // no-wait
// 		nil,            // arguments
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	body, err := json.Marshal(user)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	err = ch.Publish(
// 		"",     // exchange
// 		q.Name, // routing key
// 		false,  // mandatory
// 		false,  // immediate
// 		amqp.Publishing{
// 			ContentType: "application/json",
// 			Body:        body,
// 		})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
