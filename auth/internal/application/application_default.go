package application

import (
	"auth/internal/handler"
	"auth/internal/repository"
	"auth/internal/service"
	"database/sql"
	"net/http"
)

// ConfigServer is the configuration for the server
type ConfigServer struct {
	// Addr is the address to listen on
	Addr string
	// MySQLDSN is the DSN for the MySQL database
	MySQLDSN string
}

// New creates a new instance of the server
func New(cfg ConfigServer) *ApplicationDefault {
	// default config
	defaultCfg := ConfigServer{
		Addr:     ":5000",
		MySQLDSN: "",
	}
	if cfg.Addr != "" {
		defaultCfg.Addr = cfg.Addr
	}
	if cfg.MySQLDSN != "" {
		defaultCfg.MySQLDSN = cfg.MySQLDSN
	}

	return &ApplicationDefault{
		addr:     defaultCfg.Addr,
		mysqlDSN: defaultCfg.MySQLDSN,
	}
}

type ApplicationDefault struct {
	// addr is the address to listen on
	addr string
	// mysqlDSN is the DSN for the MySQL database
	mysqlDSN string
}

// Run runs the server
func (s *ApplicationDefault) Run() (err error) {
	// dependencies
	// - database: connection
	db, err := sql.Open("mysql", s.mysqlDSN)
	if err != nil {
		return
	}
	defer db.Close()
	// - database: ping
	err = db.Ping()
	if err != nil {
		return
	}

	// - router
	router := http.NewServeMux()

	// endpoints
	buildAuthRouter(router, db)

	// run
	err = http.ListenAndServe(s.addr, router)
	return
}

// *buildAuthRouter builds the router for the products endpoints
func buildAuthRouter(router *http.ServeMux, db *sql.DB) {
	// instance dependences
	rp := repository.NewAuthMySQL(db)
	sv := service.NewAuthDefault(rp)
	hd := handler.NewAuthDefault(sv)

	router.HandleFunc("/auth/login", methodHandler(hd.Login(), "POST"))
	router.HandleFunc("/auth/register", methodHandler(hd.Register(), "POST"))
	router.HandleFunc("/auth/verify", methodHandler(hd.Verify(), "POST"))
}

// methodHandler es un wrapper que verifica el método HTTP
func methodHandler(h http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(w, r)
	}
}
