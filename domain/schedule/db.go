package schedule

import (
	"database/sql"
	"fmt"
	"github.com/jrcartee/scheduling/database"
)

var Fields = []string{
	"schedule_name",
	"data_date",
}

func ScanSchedule(r database.Scannable) (Schedule, error) {
	var s Schedule
	err := r.Scan(
		&s.ID,
		&s.Name,
		&s.DataDate,
	)
	if err != nil {
		return s, err
	}
	return s, nil
}

func ScanSchedules(r *sql.Rows) (Schedules, error) {
	var set Schedules
	var err error
	for r.Next() {
		tmp, err := ScanSchedule(r)
		if err != nil {
			return nil, err
		}
		set = append(set, tmp)
	}
	err = r.Err()
	if err != nil {
		return nil, err
	}
	return set, nil
}

type ScheduleDB struct {
	database.Queryer
}

func (db ScheduleDB) SelectAll() (Schedules, error) {
	rows, err := db.Query("SELECT * FROM schedules")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	set, err := ScanSchedules(rows)
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (db ScheduleDB) SelectOne(where string, args ...interface{}) (Schedule, error) {
	row := db.QueryRow("SELECT * FROM schedules WHERE "+where, args...)
	return ScanSchedule(row)
}

func (db ScheduleDB) Insert(s *Schedule) error {
	f := database.FormatForInsert(Fields)
	q := fmt.Sprintf("INSERT INTO schedules %s RETURNING schedule_id;", f)

	ddate, _ := s.DataDate.Value()
	var id int
	err := db.QueryRow(q, s.Name, ddate).Scan(&id)
	if err != nil {
		return err
	}
	s.ID = id
	return nil
}

func (db ScheduleDB) Update(s *Schedule) error {
	f := database.FormatForUpdate(Fields)
	argNum := len(Fields) + 1
	q := fmt.Sprintf("UPDATE schedules SET %s WHERE schedule_id=$%d;", f, argNum)
	ddate, _ := s.DataDate.Value()
	_, err := db.Exec(q, s.Name, ddate, s.ID)
	return err
}
