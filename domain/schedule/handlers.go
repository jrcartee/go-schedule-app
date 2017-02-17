package schedule

import (
	"database/sql"
	"net/http"

	"github.com/jrcartee/scheduling/httputil"
)

type ScheduleController struct {
	db DB
}

func (sc ScheduleController) GetInstance(w http.ResponseWriter, r *http.Request) *Schedule {
	ctx := r.Context()
	pk := ctx.Value("id")
	s, err := sc.db.SelectOne("schedule_id=$1", pk)
	if err == sql.ErrNoRows {
		http.Error(w, "Not Found", http.StatusNotFound)
		return nil
	} else if err != nil {
		httputil.InternalError(w, err)
		return nil
	}
	return &s
}
