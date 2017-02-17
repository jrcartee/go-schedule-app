package task

import (
	"net/http"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/jrcartee/scheduling/httputil"
	"github.com/jrcartee/scheduling/testutil"
	"github.com/jrcartee/scheduling/types"
	// "github.com/jrcartee/scheduling/domain/schedule"
)

type MockTaskDB struct{}

func (db MockTaskDB) SelectAll() (Tasks, error) {
	return nil, nil
}
func (db MockTaskDB) SelectOne(where string, args ...interface{}) (Task, error) {
	tmp := Task{}
	return tmp, nil
}
func (db MockTaskDB) Insert(s *Task) error {
	return nil
}
func (db MockTaskDB) Update(s *Task) error {
	return nil
}

func TestTaskInsert(t *testing.T) {
	testCases := map[string]testutil.HTTPTestCase{
		"Simple Insert": {
			URL:            "/tasks",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusOK,
			Body:           []byte(`{"schedule": 1, "name":"Test", "code":"t1"}`),
		},
		"Empty Name": {
			URL:            "/tasks",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusBadRequest,
			Body:           []byte(`{"schedule": 1, "name":"", "code":""}`),
		},
	}

	db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()


	m := TaskModel{&MockTaskDB{}}
	h := httputil.GenericInsert(m, db)
	rows := sqlmock.NewRows([]string{"schedule_id", "schedule_name", "data_date"}).AddRow(1, "MockSchedule", types.NullTime{Valid: false})

	for n, tc := range testCases {
	    mock.ExpectQuery(`^SELECT (.+) FROM schedules WHERE schedule_id=\$1`).WithArgs(1).WillReturnRows(rows)
		t.Run(n, testutil.HandlerTestCase(t, h, tc))
	}
}

func TestTaskUpdate(t *testing.T) {
	testCases := map[string]testutil.HTTPTestCase{
		"Simple Update": {
			URL:            "/tasks",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusOK,
			Body:           []byte(`{"name":"Test","code":"t2"}`),
		},
		"Empty Name": {
			URL:            "/tasks",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusBadRequest,
			Body:           []byte(`{"name":""}`),
		},
	}

	db, _, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

	m := TaskModel{&MockTaskDB{}}
	h := httputil.GenericUpdate(m, db)
	for n, tc := range testCases {
		t.Run(n, testutil.HandlerTestCase(t, h, tc))
	}
}
