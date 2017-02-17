package schedule

import (
	"net/http"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/jrcartee/scheduling/httputil"
	"github.com/jrcartee/scheduling/testutil"
)

type MockScheduleDB struct{}

func (db MockScheduleDB) SelectAll() (Schedules, error) {
	return nil, nil
}
func (db MockScheduleDB) SelectOne(where string, args ...interface{}) (Schedule, error) {
	tmp := Schedule{}
	return tmp, nil
}
func (db MockScheduleDB) Insert(s *Schedule) error {
	return nil
}
func (db MockScheduleDB) Update(s *Schedule) error {
	return nil
}

func TestScheduleInsert(t *testing.T) {
	testCases := map[string]testutil.HTTPTestCase{
		"Simple Insert": {
			URL:            "/schedules",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusOK,
			Body:           []byte(`{"name":"Test","data_date":"01Jan2017"}`),
		},
		"Empty Name": {
			URL:            "/schedules",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusBadRequest,
			Body:           []byte(`{"name":"","data_date":""}`),
		},
	}


	db, _, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()


	m := ScheduleModel{&MockScheduleDB{}}
	h := httputil.GenericInsert(m, db)
	for n, tc := range testCases {
		t.Run(n, testutil.HandlerTestCase(t, h, tc))
	}
}

func TestScheduleUpdate(t *testing.T) {
	testCases := map[string]testutil.HTTPTestCase{
		"Simple Update": {
			URL:            "/schedules",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusOK,
			Body:           []byte(`{"name":"Test","data_date":"01Jan2017"}`),
		},
		"Empty Update": {
			URL:            "/schedules",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusBadRequest,
			Body:           []byte(`{"name":"","data_date":""}`),
		},
	}


	db, _, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()


	m := ScheduleModel{&MockScheduleDB{}}
	h := httputil.GenericUpdate(m, db)
	for n, tc := range testCases {
		t.Run(n, testutil.HandlerTestCase(t, h, tc))
	}
}
