package schedule

import (
	"github.com/jrcartee/scheduling/database"
	"github.com/jrcartee/scheduling/types"
	"testing"
	"time"
)

func TestScheduleDB(t *testing.T) {
	db, teardown := database.Setup(database.DefaultConfig)
	defer teardown()

	t.Run("SelectAll", func(t *testing.T) {
		sdb := ScheduleDB{db}
		s, err := sdb.SelectAll()
		CheckError(t, err)

		if s == nil || len(s) == 0 {
			t.Error("no schedules")
		}
	})

	t.Run("SelectOne", func(t *testing.T) {
		sdb := ScheduleDB{db}
		empty := Schedule{}

		s, err := sdb.SelectOne("schedule_id=$1", 1)
		CheckError(t, err)

		if s == empty {
			t.Error("Empty schedule")
		}
	})

	t.Run("Insert", func(t *testing.T) {
		tx, err := db.Begin()
		CheckError(t, err)
		defer tx.Rollback()

		sdb := ScheduleDB{tx}
		s := &Schedule{
			Name: "Testing Schedules",
			DataDate: types.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		}
		err = sdb.Insert(s)
		CheckError(t, err)

		if s.ID == 0 {
			t.Error("ID not set")
		}
	})

	t.Run("Update", func(t *testing.T) {
		tx, err := db.Begin()
		CheckError(t, err)
		defer tx.Rollback()

		sdb := ScheduleDB{tx}
		s1, err := sdb.SelectOne("schedule_id=$1", 1)
		CheckError(t, err)

		s1.Name = "Different Name"
		err = sdb.Update(&s1)
		CheckError(t, err)

		s2, err := sdb.SelectOne("schedule_id=$1", 1)
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
