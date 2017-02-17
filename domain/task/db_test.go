package task

import (
	// "time"
	"github.com/jrcartee/scheduling/database"
	"testing"
	// "github.com/jrcartee/scheduling/types"
)

func TestTaskDB(t *testing.T) {
	db, teardown := database.Setup(database.DefaultConfig)
	defer teardown()

	t.Run("SelectAll", func(t *testing.T) {
		tdb := TaskDB{db}
		s, err := tdb.SelectAll()
		CheckError(t, err)

		if s == nil || len(s) == 0 {
			t.Error("no tasks")
		}
	})

	t.Run("SelectOne", func(t *testing.T) {
		tdb := TaskDB{db}
		empty := Task{}

		s, err := tdb.SelectOne("task_id=$1", 1)
		CheckError(t, err)

		if s == empty {
			t.Error("Empty task")
		}
	})

	t.Run("Insert", func(t *testing.T) {
		tx, err := db.Begin()
		CheckError(t, err)
		defer tx.Rollback()

		tdb := TaskDB{tx}
		s := &Task{
			Schedule: 1,
			Name:     "Testing Tasks",
			Code:     "tt-1",
		}
		err = tdb.Insert(s)
		CheckError(t, err)

		if s.ID == 0 {
			t.Error("ID not set")
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := db.Begin()
		CheckError(t, err)
		defer tx.Rollback()

		tdb := TaskDB{tx}
		s1, err := tdb.SelectOne("task_id=$1", 1)
		CheckError(t, err)

		s1.Name = "Different Name"
		err = tdb.Update(&s1)
		CheckError(t, err)

		s2, err := tdb.SelectOne("task_id=$1", 1)
		CheckError(t, err)
		if s2.Name != s1.Name {
			t.Errorf("Names don't match: %s %s", s1.Name, s2.Name)
		}
	})

}

func CheckError(t *testing.T, e error) {
	if e != nil {
		t.Fatalf("error not nil:\n %s", e)
	}
}
