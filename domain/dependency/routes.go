package dependency

import (
	"database/sql"
	"github.com/jrcartee/router"
	"github.com/jrcartee/scheduling/httputil"
	"net/http"
)

func RegisterRoutes(r *router.Router, db *sql.DB) {
	tm := &TaskDependencyModel{TaskDependencyDB{db}}

	r.RegisterRoute("deps", router.Endpoints{
		http.MethodGet:  httputil.GenericIndex(tm),
		http.MethodPost: httputil.GenericInsert(tm, db),
	})

	r.RegisterRoute(`deps/{before:\d+}/{after:\d+}`, router.Endpoints{
		http.MethodGet:  httputil.GenericRetrieve(tm),
		http.MethodPost: httputil.GenericUpdate(tm, db),
	})
}
