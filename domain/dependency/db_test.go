package dependency

import (
	"github.com/jrcartee/scheduling/database"
	"testing"
)

func TestTaskDependencyDB(t *testing.T) {
	db, teardown := database.Setup(database.DefaultConfig)
	defer teardown()

	t.Run("SelectAll", func(t *testing.T) {
		tdb := TaskDependencyDB{db}
		s, err := tdb.SelectAll()
		CheckError(t, err)

		if s == nil || len(s) == 0 {
			t.Error("no task deps")
		}
	})

	t.Run("SelectOne", func(t *testing.T) {
		tdb := TaskDependencyDB{db}
		empty := TaskDependency{}

		s, err := tdb.SelectOne("task_before=$1 AND task_after=$2", 1, 2)
		CheckError(t, err)

		if s == empty {
			t.Error("Empty dependency")
		}
	})

	t.Run("Insert", func(t *testing.T) {
		tx, err := db.Begin()
		CheckError(t, err)
		defer tx.Rollback()

		tdb := TaskDependencyDB{tx}
		td := &TaskDependency{
			Schedule: 1,
			Before: 1,
			After: 3,
			Lag: 5,
			DType: 1,
		}
		err = tdb.Insert(td)
		CheckError(t, err)
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := db.Begin()
		CheckError(t, err)
		defer tx.Rollback()

		tdb := TaskDependencyDB{tx}
		td1, err := tdb.SelectOne("task_before=$1 AND task_after=$2", 1, 2)
		CheckError(t, err)

		td1.Lag = 99
		err = tdb.Update(&td1)
		CheckError(t, err)

		td2, err := tdb.SelectOne("task_before=$1 AND task_after=$2", 1, 2)
		CheckError(t, err)
		if td2.Lag != td1.Lag {
			t.Errorf("Lag doesn't match: %s %s", td1.Lag, td2.Lag)
		}
	})

}

func CheckError(t *testing.T, e error) {
	if e != nil {
		t.Fatalf("error not nil:\n %s", e)
	}
}
