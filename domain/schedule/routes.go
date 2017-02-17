package schedule

import (
	"database/sql"
	"github.com/jrcartee/router"
	"github.com/jrcartee/scheduling/httputil"
	"net/http"
)

func RegisterRoutes(r *router.Router, db *sql.DB) {
	// sc := &ScheduleController{ &ScheduleDB{db} }
	sm := &ScheduleModel{ScheduleDB{db}}

	r.RegisterRoute("schedules", router.Endpoints{
		http.MethodGet:  httputil.GenericIndex(sm),
		http.MethodPost: httputil.GenericInsert(sm, db),
	})

	r.RegisterRoute(`schedules/{id:\d+}`, router.Endpoints{
		http.MethodGet:  httputil.GenericRetrieve(sm),
		http.MethodPost: httputil.GenericUpdate(sm, db),
	})
}
