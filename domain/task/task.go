package task

import (
	"time"
	
	"github.com/jrcartee/scheduling/types"
)

type DB interface {
	SelectAll() (Tasks, error)
	SelectOne(string, ...interface{}) (Task, error)
	Insert(*Task) error
	Update(*Task) error
}

type Task struct {
	ID           int            `db:"task_id" json:"id"`
	Schedule     int            `db:"schedule" json:"schedule"`
	Code         string         `db:"task_code" json:"code"`
	Name         string         `db:"task_name" json:"name"`
	Duration     int            `db:"duration" json:"duration"`
	Remaining    int            `db:"remaining" json:"remaining"`
	StartEarly   types.NullTime `db:"start_early" json:"start_early"`
	StartLate    types.NullTime `db:"start_late" json:"start_late"`
	StartActual  types.NullTime `db:"start_actual" json:"start_actual"`
	FinishEarly  types.NullTime `db:"finish_early" json:"finish_early"`
	FinishLate   types.NullTime `db:"finish_late" json:"finish_late"`
	FinishActual types.NullTime `db:"finish_actual" json:"finish_actual"`
}
type Tasks []Task


func (t Task) CurrentDrtn(date time.Time) int {
	if !t.StartActual.Valid {
		return 0
	}

	if t.FinishActual.Valid {
		date = t.FinishActual.Time
	}

	drtn := date.Sub(t.StartActual.Time)
	days := int(drtn.Hours() / 24)
	return days
}

func (t Task) ActualDrtn(date time.Time) int {
	if !t.StartActual.Valid {
		return t.Duration
	}
	return t.CurrentDrtn(date) + t.Remaining
}
