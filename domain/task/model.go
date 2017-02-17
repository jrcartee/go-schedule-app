package task

import (
	"net/http"
	"database/sql"
	"github.com/jrcartee/scheduling/database"

	"github.com/jrcartee/scheduling/domain/schedule"
)

type TaskModel struct {
	db DB
}

func (m TaskModel) New() interface{} {
	return new(Task)
}

func (m TaskModel) Instance(w http.ResponseWriter, r *http.Request) interface{} {
	tc := TaskController{m.db}
	return tc.GetInstance(w, r)
}

func (m TaskModel) SelectAll() (interface{}, error) {
	return m.db.SelectAll()
}

func (m TaskModel) Insert(data interface{}) error {
	return m.db.Insert(data.(*Task))
}

func (m TaskModel) Update(data interface{}) error {
	return m.db.Update(data.(*Task))
}



func (m TaskModel) ValidateInsert(data interface{}, db database.Queryer) map[string]string {
	t := data.(*Task)
	verrs := m.ValidateUpdate(t, db)

	sdb := schedule.ScheduleDB{db}
	_, err := sdb.SelectOne("schedule_id=$1", t.Schedule)			
	if err == sql.ErrNoRows {
		verrs["schedule"] = "Schedule doesn't exist"
	} else if err != nil {
		verrs["schedule"] = "An error occured while verifying schedule"			
	}
	return verrs
}

func (m TaskModel) ValidateUpdate(data interface{}, db database.Queryer) map[string]string {
	t := data.(*Task)
	verrs := make(map[string]string)
	if len(t.Name) == 0 {
		verrs["name"] = "Name is required"
	}
	return verrs
}