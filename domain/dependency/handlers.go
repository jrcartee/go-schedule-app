package dependency

import (
	"database/sql"
	"net/http"

	"github.com/jrcartee/scheduling/httputil"
)

type TaskDependencyController struct {
	db DB
}

func (tc TaskDependencyController) GetInstance(w http.ResponseWriter, r *http.Request) *TaskDependency {
	ctx := r.Context()
	before, after := ctx.Value("before"), ctx.Value("after")
	td, err := tc.db.SelectOne("task_before=$1 AND task_after=$2", before, after)
	if err == sql.ErrNoRows {
		http.Error(w, "Not Found", http.StatusNotFound)
		return nil
	} else if err != nil {
		httputil.InternalError(w, err)
		return nil
	}
	return &td
}
