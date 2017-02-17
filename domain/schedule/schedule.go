package schedule

import (
	"github.com/jrcartee/scheduling/types"
)

type DB interface {
	SelectAll() (Schedules, error)
	SelectOne(string, ...interface{}) (Schedule, error)
	Insert(*Schedule) error
	Update(*Schedule) error
}

type Schedule struct {
	ID       int            `db:"schedule_id" json:"id"`
	Name     string         `db:"schedule_name" json:"name"`
	DataDate types.NullTime `db:"data_date" json:"data_date"`
}

type Schedules []Schedule
