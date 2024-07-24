// auth/main.go
package main

import (
	"auth/internal/application"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

func main() {
	addrCfg := os.Getenv("PORT")

	mysqlCfg := mysql.Config{
		User:      os.Getenv("MYSQL_USER"),
		Passwd:    os.Getenv("MYSQL_PASSWORD"),
		Net:       "tcp",
		Addr:      os.Getenv("MYSQL_ADDR"),
		DBName:    os.Getenv("MYSQL_DBNAME"),
		ParseTime: true,
	}

	cfg := application.ConfigServer{Addr: addrCfg, MySQLDSN: mysqlCfg.FormatDSN()}

	server := application.New(cfg)

	if err := server.Run(); err != nil {
		fmt.Println(err)
	}
}
