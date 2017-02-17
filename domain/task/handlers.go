package task

import (
	"database/sql"
	"net/http"

	"github.com/jrcartee/scheduling/httputil"
)

type TaskController struct {
	db DB
}

func (tc TaskController) GetInstance(w http.ResponseWriter, r *http.Request) *Task {
	ctx := r.Context()
	pk := ctx.Value("id")
	t, err := tc.db.SelectOne("task_id=$1", pk)
	if err == sql.ErrNoRows {
		http.Error(w, "Not Found", http.StatusNotFound)
		return nil
	} else if err != nil {
		httputil.InternalError(w, err)
		return nil
	}
	return &t
}
