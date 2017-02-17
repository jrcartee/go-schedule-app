package dependency

import (
	"net/http"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/jrcartee/scheduling/httputil"
	"github.com/jrcartee/scheduling/testutil"
	"github.com/jrcartee/scheduling/types"
)

type MockTaskDependencyDB struct{}

func (db MockTaskDependencyDB) SelectAll() (TaskDependencies, error) {
	return nil, nil
}
func (db MockTaskDependencyDB) SelectOne(where string, args ...interface{}) (TaskDependency, error) {
	tmp := TaskDependency{}
	return tmp, nil
}
func (db MockTaskDependencyDB) Insert(s *TaskDependency) error {
	return nil
}
func (db MockTaskDependencyDB) Update(s *TaskDependency) error {
	return nil
}


func TestTaskDependencyInsert(t *testing.T) {

	/*
		TODO: FIX THESE TESTS
		 Seperate cases so that empty calls can be mocked successfully.

		 figure out WHY WOULD THIS BE NON-DETERMINISTIC?!
		 After splitting test cases into separate test funcs is it still non-deterministic?
		 Maybe i should bail on this sqlmock shin-dig and find a better way 
		 	to validate and/or stub validation queries
	*/

	testCases := map[string]testutil.HTTPTestCase{
		"Happy": {
			URL:            "/deps",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusOK,
			Body:           []byte(`{"schedule":1,"task_before":1,"task_after":2}`),
		},
		"Bad schedule": {
			URL:            "/deps",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusBadRequest,
			Body:           []byte(`{"schedule":0,"task_before":1,"task_after":2}`),
		},
		"Bad task_before": {
			URL:            "/deps",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusBadRequest,
			Body:           []byte(`{"schedule":1,"task_before":0,"task_after":2}`),
		},
		"Bad task_after": {
			URL:            "/deps",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusBadRequest,
			Body:           []byte(`{"schedule":1,"task_before":1,"task_after":0}`),
		},
		"Bad type": {
			URL:            "/deps",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusBadRequest,
			Body:           []byte(`{"schedule":1,"task_before":1,"task_after":2,"type":5}`),
		},
	}

	db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()


	m := TaskDependencyModel{&MockTaskDependencyDB{}}
	h := httputil.GenericInsert(m, db)
	sRows := sqlmock.NewRows([]string{"schedule_id", "schedule_name", "data_date"}).
		AddRow(1, "MockSchedule", types.NullTime{Valid: false})
	tFields := []string{
		"task_id",
		"schedule",
		"task_code",
		"task_name",
		"duration",
		"remaining",
		"start_early",
		"start_late",
		"start_actual",
		"finish_early",
		"finish_late",
		"finish_actual",
	}
	null := types.NullTime{Valid:false}
	tRow1 := sqlmock.NewRows(tFields).
		AddRow(1, 1, "t1", "MockTask1", 1, 1, null,null,null,null,null,null)
	tRow2 := sqlmock.NewRows(tFields).
		AddRow(2, 1, "t2", "MockTask2", 1, 1, null,null,null,null,null,null)


	for n, tc := range testCases {
	    mock.ExpectQuery(`^SELECT (.+) FROM schedules WHERE schedule_id=\$1`).
	    	WithArgs(1).
	    	WillReturnRows(sRows)
	    mock.ExpectQuery(`^SELECT (.+) FROM tasks WHERE task_id=\$1`).
	    	WithArgs(1).
	    	WillReturnRows(tRow1)
	    mock.ExpectQuery(`^SELECT (.+) FROM tasks WHERE task_id=\$1`).
	    	WithArgs(2).
	    	WillReturnRows(tRow2)
		t.Run(n, testutil.HandlerTestCase(t, h, tc))
	}
}





func TestTaskDependencyUpdate(t *testing.T) {
	testCases := map[string]testutil.HTTPTestCase{
		"Happy case": {
			URL:            "/deps",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusOK,
			Body:           []byte(`{"lag": 1, "type": 2}`),
		},
		"Bad type": {
			URL:            "/deps",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusBadRequest,
			Body:           []byte(`{"lag": 1, "type": 6}`),
		},
	}

	db, _, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

	m := TaskDependencyModel{&MockTaskDependencyDB{}}
	h := httputil.GenericUpdate(m, db)
	for n, tc := range testCases {
		t.Run(n, testutil.HandlerTestCase(t, h, tc))
	}
}
