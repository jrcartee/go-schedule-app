package task

import (
	"database/sql"
	"github.com/jrcartee/router"
	"github.com/jrcartee/scheduling/httputil"
	"net/http"
)

func RegisterRoutes(r *router.Router, db *sql.DB) {
	m := &TaskModel{TaskDB{db}}
	r.RegisterRoute("tasks", router.Endpoints{
		http.MethodGet:  httputil.GenericIndex(m),
		http.MethodPost: httputil.GenericInsert(m, db),
	})

	r.RegisterRoute(`tasks/{id:\d+}`, router.Endpoints{
		http.MethodGet:  httputil.GenericRetrieve(m),
		http.MethodPost: httputil.GenericUpdate(m, db),
	})
}
