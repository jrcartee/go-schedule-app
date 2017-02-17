package cpm

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jrcartee/scheduling/types"
	"github.com/jrcartee/scheduling/database"
	"github.com/jrcartee/scheduling/httputil"


	"github.com/jrcartee/scheduling/domain/schedule"
	"github.com/jrcartee/scheduling/domain/task"
	"github.com/jrcartee/scheduling/domain/dependency"
)

type CPMRequest struct {
	Schedule int `json:"schedule"`
	DataDate types.NullTime `json:"data_date"`
}

type CPMState struct {
	dataDate time.Time
	visited map[int]struct{}
	depDB TaskDepDB
	taskDB task.TaskDB
}

func CreateCPMHandler(db database.Queryer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rdata := new(CPMRequest)
		err := httputil.DecodeJSONRequest(r, rdata)
		if err != nil {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
		}

		if !rdata.DataDate.Valid {
			http.Error(w, "Must provide data_date", http.StatusBadRequest)			
		}

		sdb := schedule.ScheduleDB{db}
		s, err := sdb.SelectOne("schedule_id=$1", rdata.Schedule)
		if err != nil {
			fmt.Println(err)			
			http.Error(w, "Schedule not found", http.StatusBadRequest)			
		}

		tdb := task.TaskDB{db}
		tddb := TaskDepDB{dependency.TaskDependencyDB{db}}

		// tlist, err := tdb.Select("schedule=$1", s.ID)
		// if err != nil {
		// 	fmt.Println(err)
		// 	http.Error(w, "Cannot retrieve tasks", http.StatusBadRequest)
		// }

		// tdlist, err := tddb.Select("schedule=$1", s.ID)
		// if err != nil {
		// 	fmt.Println(err)
		// 	http.Error(w, "Cannot retrieve dependencies", http.StatusBadRequest)
		// }

		rootIDs, err := tdb.SelectRootIdsForSchedule(s.ID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Cannot retrieve root tasks", http.StatusBadRequest)
		}

		// tlist, err := tdb.Select("schedule=$1", s.ID)

		visited := make(map[int]struct{})

		state := CPMState{
			dataDate: rdata.DataDate.Time,
			visited: visited,
			taskDB: tdb,
			depDB: tddb,
		}

		for _, rid := range rootIDs {
			forwardPass(rid, state)			
		}


	}
}

func forwardPass(tid int, state CPMState) error {
	if _, found := state.visited[tid]; found {
		return nil
	}
	predIDs, err := state.depDB.PredecessorIdsForTask(tid)
	if err != nil {
		return err
	}
	if len(predIDs) != 0 && !allKeysInMap(predIDs, state.visited) {
		return nil
	}

	tdata, err := state.taskDB.SelectOne("task_id=$1", tid)
	if err != nil {
		return err
	}

	var start time.Time
	if tdata.StartActual.Valid {
		// already started
		start = tdata.StartActual.Time
	} else if len(predIDs) == 0 {
		// hasn't started & no predecessors, start ASAP
		start = state.dataDate
	} else {
		start, err = state.depDB.findNextStart(tdata, state.dataDate)
		if err != nil {
			return err
		}
		start = maxTime(start, state.dataDate)
	}
	tdata.StartEarly = types.NullTime{Valid: true, Time: start}

	var finish time.Time
	if tdata.FinishActual.Valid {
		finish = tdata.FinishActual.Time		
	} else {
		numDays := tdata.ActualDrtn(state.dataDate)
		finish = start.Add(time.Hour * time.Duration(24 * numDays))
	}
	tdata.FinishEarly = types.NullTime{Valid: true, Time: finish}

	err = state.taskDB.Update(&tdata)
	if err != nil {
		return err
	}


	return nil
}
