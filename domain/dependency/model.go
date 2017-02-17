package dependency

import (
	"net/http"
	
	"github.com/jrcartee/scheduling/database"
	"github.com/jrcartee/scheduling/httputil"
)

type TaskDependencyModel struct {
	db DB
}

func (m TaskDependencyModel) New() interface{} {
	return new(TaskDependency)
}

func (m TaskDependencyModel) Instance(w http.ResponseWriter, r *http.Request) interface{} {
	tc := TaskDependencyController{m.db}
	return tc.GetInstance(w, r)
}

func (m TaskDependencyModel) SelectAll() (interface{}, error) {
	return m.db.SelectAll()
}

func (m TaskDependencyModel) Insert(data interface{}) error {
	return m.db.Insert(data.(*TaskDependency))
}

func (m TaskDependencyModel) Update(data interface{}) error {
	return m.db.Update(data.(*TaskDependency))
}

func (m TaskDependencyModel) ValidateInsert(data interface{}, db database.Queryer) httputil.ValidationErrors {
	
	t := data.(*TaskDependency)
	errs := httputil.MergeErrors(
		validateCommon(*t),
		validateSchedule(*t, db),
		validateTasks(*t, db),
	)

	_, sError := errs["schedule"]
	_, tbError := errs["task_before"]
	_, taError := errs["task_after"]
	if !sError && !tbError && !taError {
		cycle, err := causesCycle(*t, db)
		if err != nil {
			errs["all"] = "Error while checking for cycles"			
		} else if cycle {
			errs["all"] = "This combination would cause a cycle"
		}
	}
	return errs
}

func (m TaskDependencyModel) ValidateUpdate(data interface{}, db database.Queryer) httputil.ValidationErrors {
	t := data.(*TaskDependency)
	return validateCommon(*t)
}


