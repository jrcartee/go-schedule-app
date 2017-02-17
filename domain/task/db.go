package task

import (
	"database/sql"
	"fmt"
	// "strings"
	"github.com/jrcartee/scheduling/database"
)

type TaskDB struct {
	database.Queryer
}

func (db TaskDB) SelectAll() (Tasks, error) {
	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	set, err := ScanTasks(rows)
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (db TaskDB) Select(where string, args ...interface{}) (Tasks, error) {
	q := fmt.Sprintf("SELECT * FROM tasks WHERE %s", where)
	rows, err := db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	set, err := ScanTasks(rows)
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (db TaskDB) SelectOne(where string, args ...interface{}) (Task, error) {
	row := db.QueryRow("SELECT * FROM tasks WHERE "+where, args...)
	return ScanTask(row)
}

func (db TaskDB) Insert(t *Task) error {
	f := database.FormatForInsert([]string{
		"schedule", "task_code", "task_name",
		"duration", "remaining",
		"start_early", "start_late", "start_actual",
		"finish_early", "finish_late", "finish_actual",
	})
	q := fmt.Sprintf("INSERT INTO tasks %s RETURNING task_id;", f)
	se, _ := t.StartEarly.Value()
	sl, _ := t.StartLate.Value()
	sa, _ := t.StartActual.Value()
	fe, _ := t.FinishEarly.Value()
	fl, _ := t.FinishLate.Value()
	fa, _ := t.FinishActual.Value()
	var id int

	err := db.QueryRow(q,
		t.Schedule, t.Code, t.Name,
		t.Duration, t.Remaining,
		se, sl, sa, fe, fl, fa).Scan(&id)
	if err != nil {
		return err
	}
	t.ID = id
	return nil
}

func (db TaskDB) Update(t *Task) error {
	fields := []string{
		"task_code", "task_name",
		"duration", "remaining",
		"start_early", "start_late", "start_actual",
		"finish_early", "finish_late", "finish_actual",
	}
	f := database.FormatForUpdate(fields)
	argNum := len(fields) + 1
	q := fmt.Sprintf("UPDATE tasks SET %s WHERE task_id=$%d;", f, argNum)
	se, _ := t.StartEarly.Value()
	sl, _ := t.StartLate.Value()
	sa, _ := t.StartActual.Value()
	fe, _ := t.FinishEarly.Value()
	fl, _ := t.FinishLate.Value()
	fa, _ := t.FinishActual.Value()

	_, err := db.Exec(q,
		t.Code, t.Name,
		t.Duration, t.Remaining,
		se, sl, sa, fe, fl, fa, t.ID)
	return err
}

func (db TaskDB) SelectIdsForSchedule(sid int) ([]int, error) {
	q := "SELECT task_id FROM tasks WHERE schedule_id=$1"
	rows, err := db.Query(q, sid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var set []int
	var tmp int
	for rows.Next() {
		err := rows.Scan(&tmp)
		if err != nil {
			return nil, err
		}
		set = append(set, tmp)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return set, nil
}


func (db TaskDB) SelectRootIdsForSchedule(sid int) ([]int, error) {
	q := `SELECT T.task_id FROM tasks T 
			LEFT JOIN task_deps D 
			ON T.schedule=$1 AND D.schedule=$1 
			AND T.task_id=D.task_after 
			GROUP BY T.task_id 
			HAVING COUNT(D.task_after)=0;`
	rows, err := db.Query(q, sid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var set []int
	var tmp int
	for rows.Next() {
		err := rows.Scan(&tmp)
		if err != nil {
			return nil, err
		}
		set = append(set, tmp)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return set, nil
}



func ScanTask(r database.Scannable) (Task, error) {
	var s Task
	err := r.Scan(
		&s.ID, &s.Schedule, &s.Code, &s.Name,
		&s.Duration, &s.Remaining,
		&s.StartEarly, &s.StartLate, &s.StartActual,
		&s.FinishEarly, &s.FinishLate, &s.FinishActual,
	)
	if err != nil {
		return s, err
	}
	return s, nil
}

func ScanTasks(r *sql.Rows) (Tasks, error) {
	var set Tasks
	var err error
	for r.Next() {
		tmp, err := ScanTask(r)
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
