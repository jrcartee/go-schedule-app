package schedule

import (
	"net/http"
	"github.com/jrcartee/scheduling/database"
)

type ScheduleModel struct {
	db DB
}

func (m ScheduleModel) New() interface{} {
	return new(Schedule)
}

func (m ScheduleModel) Instance(w http.ResponseWriter, r *http.Request) interface{} {
	sc := ScheduleController{m.db}
	return sc.GetInstance(w, r)
}

func (m ScheduleModel) SelectAll() (interface{}, error) {
	return m.db.SelectAll()
}

func (m ScheduleModel) Insert(data interface{}) error {
	return m.db.Insert(data.(*Schedule))
}

func (m ScheduleModel) Update(data interface{}) error {
	return m.db.Update(data.(*Schedule))
}

func (m ScheduleModel) ValidateInsert(data interface{}, db database.Queryer) map[string]string {
	s := data.(*Schedule)
	return validate(s)
}

func (m ScheduleModel) ValidateUpdate(data interface{}, db database.Queryer) map[string]string {
	s := data.(*Schedule)
	return validate(s)
}



func validate(s *Schedule) map[string]string {
	verrs := make(map[string]string)
	if len(s.Name) == 0 {
		verrs["name"] = "Name is required"
	}
	return verrs
}