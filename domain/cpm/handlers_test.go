package cpm

import (
	"net/http"
	"testing"

	// "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/jrcartee/scheduling/database"
	"github.com/jrcartee/scheduling/testutil"
	// "github.com/jrcartee/scheduling/types"
	// "github.com/jrcartee/scheduling/domain/schedule"
)

func TestTaskInsert(t *testing.T) {
	testCases := map[string]testutil.HTTPTestCase{
		"Happy": {
			URL:            "/critical-path",
			Method:         http.MethodPost,
			ExpectedStatus: http.StatusOK,
			Body:           []byte(`{"schedule": 1, "data_date":"01JAN2017"}`),
		},
	}

	db, closeDB := database.Setup(database.DefaultConfig)
    defer closeDB()

    // TODO: Wrap these tests in a transaction that is rolled back
    //		Or mock the db calls
    
    h := CreateCPMHandler(db)
	for n, tc := range testCases {
		t.Run(n, testutil.HandlerTestCase(t, h, tc))
	}
}