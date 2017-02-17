package dependency

import (
	"database/sql"
	"fmt"
	"github.com/jrcartee/scheduling/database"
)


type TaskDependencyDB struct {
	database.Queryer
}

func (db TaskDependencyDB) SelectAll() (TaskDependencies, error) {
	rows, err := db.Query("SELECT * FROM task_deps")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	set, err := ScanTaskDependencies(rows)
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (db TaskDependencyDB) Select(where string, args ...interface{}) (TaskDependencies, error) {
	q := "SELECT * FROM task_deps WHERE " + where
	rows, err := db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	set, err := ScanTaskDependencies(rows)
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (db TaskDependencyDB) SelectOne(where string, args ...interface{}) (TaskDependency, error) {
	row := db.QueryRow("SELECT * FROM task_deps WHERE "+where, args...)
	return ScanTaskDependency(row)
}

func (db TaskDependencyDB) Insert(t *TaskDependency) error {
	f := database.FormatForInsert([]string{
		"schedule", "lag", "type", "task_before", "task_after",
	})
	q := fmt.Sprintf("INSERT INTO task_deps %s;", f)

	_, err := db.Exec(q, t.Schedule, t.Lag, t.DType, t.Before, t.After)
	return err
}

func (db TaskDependencyDB) Update(t *TaskDependency) error {
	fields := []string{
		"lag", "type",
	}
	f := database.FormatForUpdate(fields)
	p1 := len(fields) + 1
	p2 := p1 + 1
	q := fmt.Sprintf("UPDATE task_deps SET %s WHERE task_before=$%d AND task_after=$%d;", f, p1, p2)

	_, err := db.Exec(q, t.Lag, t.DType, t.Before, t.After)
	return err

}

func (db TaskDependencyDB) SelectIdsForSchedule(sid int) ([][]int, error) {
	q := "SELECT task_before, task_after FROM task_deps WHERE schedule_id=$1"
	rows, err := db.Query(q, sid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var set [][]int
	var btmp, atmp int
	for rows.Next() {
		err := rows.Scan(&btmp, &atmp)
		if err != nil {
			return nil, err
		}
		set = append(set, []int{btmp, atmp})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return set, nil
}


func ScanTaskDependency(r database.Scannable) (TaskDependency, error) {
	var td TaskDependency
	err := r.Scan(
		&td.Schedule, &td.Before, &td.After, &td.Lag, &td.DType,
	)
	if err != nil {
		return td, err
	}
	return td, nil
}

func ScanTaskDependencies(r *sql.Rows) (TaskDependencies, error) {
	var set TaskDependencies
	var err error
	for r.Next() {
		tmp, err := ScanTaskDependency(r)
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