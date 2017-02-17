package cpm

import (
	"database/sql"
	"github.com/jrcartee/router"
	"net/http"
)

func RegisterRoutes(r *router.Router, db *sql.DB) {

	r.RegisterRoute(`critical-path`, router.Endpoints{
		http.MethodPost: CreateCPMHandler(db),
	})
}
