package cpm
import (
	// "fmt"
	"time"

	"github.com/jrcartee/scheduling/domain/task"
	"github.com/jrcartee/scheduling/domain/dependency"
)


type TaskDepDB struct {
	dependency.TaskDependencyDB
}

func (db TaskDepDB) PredecessorIdsForTask(tid int) ([]int, error) {
	q := "SELECT task_before FROM task_deps WHERE task_after=$1"
	rows, err := db.Query(q, tid)
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

func (db TaskDepDB) SuccessorIdsForTask(tid int) ([]int, error) {
	q := "SELECT task_before FROM task_deps WHERE task_before=$1"
	rows, err := db.Query(q, tid)
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

func (db TaskDepDB) SuccessorDataForTask(tid int) ([]dependencyTask, error) {	
	q := `SELECT (
			D.lag, D.type, T.start_actual, T.finish_actual, T.duration, T.remaining
			) FROM task_deps D JOIN tasks T ON D.task_after=T.task_id WHERE T.task_before=$1`
	rows, err := db.Query(q, tid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var set []dependencyTask
	var tmp dependencyTask
	for rows.Next() {
		err := rows.Scan(
			&tmp.Lag, &tmp.DType, 
			&tmp.StartActual, &tmp.FinishActual, 
			&tmp.Duration, &tmp.Remaining)
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


func (db TaskDepDB) findNextStart(t task.Task, d time.Time) (time.Time, error) {
	var date time.Time
	dtlist, err := db.SuccessorDataForTask(t.ID)
	if err != nil {
		return date, err
	}

	for _, dt := range dtlist {
		if date.IsZero() {
			date = dt.calculateNextStart(t, d)
		} else {
			date = maxTime(date, dt.calculateNextStart(t, d))
		}
	}

	return date, nil
}



type dependencyTask struct {
	dependency.TaskDependency
	task.Task
}

func (dt dependencyTask) calculateNextStart(b task.Task, dd time.Time) time.Time {
	var base time.Time
	var buffer int
	switch dt.DType {
	case 0:
		base = b.FinishEarly.Time
		buffer = dt.Lag + 1
	case 1:
		base = b.FinishEarly.Time
		buffer = dt.Lag - (dt.ActualDrtn(dd) - 1)
	case 2:
		base = b.StartEarly.Time
		buffer = dt.Lag
	case 3:
		base = b.StartEarly.Time
		buffer = dt.Lag - (dt.ActualDrtn(dd) - 1)
	}
	return base.Add(time.Hour * time.Duration(24 * buffer))
}
func (dt dependencyTask) calculatePreviousFinish(a task.Task, dd time.Time) time.Time {
	var base time.Time
	var buffer int
	switch dt.DType {
	case 0:
		base = a.StartLate.Time
		buffer = dt.Lag + 1
	case 1:
		base = a.FinishLate.Time
		buffer = dt.Lag
	case 2:
		base = a.StartLate.Time
		buffer = dt.Lag - (dt.ActualDrtn(dd) - 1)
	case 3:
		base = a.FinishLate.Time
		buffer = dt.Lag - (dt.ActualDrtn(dd) - 1)
	}
	return base.Add(time.Hour * time.Duration(24 * buffer))
}