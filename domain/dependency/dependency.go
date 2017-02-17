package dependency

type DB interface {
	SelectAll() (TaskDependencies, error)
	SelectOne(string, ...interface{}) (TaskDependency, error)
	Insert(*TaskDependency) error
	Update(*TaskDependency) error
}

var DependencyTypes = map[int]string{
	0: "fs",
	1: "ff",
	2: "ss",
	3: "sf",
}

type TaskDependency struct {
	Schedule int `db:"schedule" json:"schedule"`
	Lag      int `db:"lag" json:"lag"`
	DType    int `db:"type" json:"type"`
	Before   int `db:"task_before" json:"task_before"`
	After    int `db:"task_after" json:"task_after"`
}
type TaskDependencies []TaskDependency

