package dependency

import (
	"database/sql"
	"github.com/jrcartee/scheduling/database"
	"github.com/jrcartee/scheduling/domain/task"
	"github.com/jrcartee/scheduling/domain/schedule"
)

func validateCommon(td TaskDependency) map[string]string {
	verrs := make(map[string]string)
	if td.DType < 0 || td.DType > 3 {
		verrs["type"] = "Invalid type provided"
	}

	return verrs
}


func validateSchedule(td TaskDependency, db database.Queryer) map[string]string {
	var err error
	verrs := make(map[string]string)

	sdb := schedule.ScheduleDB{db}
	_, err = sdb.SelectOne("schedule_id=$1", td.Schedule)
	if err == sql.ErrNoRows {
		verrs["schedule"] = "Schedule doesn't exist"
	} else if err != nil {
		verrs["schedule"] = "Error occured while verifying schedule"
	}
	return verrs
}

func validateTasks(td TaskDependency, db database.Queryer) map[string]string {
	var err error
	verrs := make(map[string]string)

	tdb := task.TaskDB{db}
	t, err := tdb.SelectOne("task_id=$1", td.Before)
	if err == sql.ErrNoRows {
		verrs["task_before"] = "Task doesn't exist"
	} else if err != nil {		
		verrs["task_before"] = "Error occured while verifying Task"
	} else if t.Schedule != td.Schedule {
		verrs["task_before"] = "Not found in this schedule"		
	}

	t, err = tdb.SelectOne("task_id=$1", td.After)
	if err == sql.ErrNoRows {
		verrs["task_after"] = "Task doesn't exist"
	} else if err != nil {
		verrs["task_after"] = "Error occured while verifying Task"
	} else if t.Schedule != td.Schedule {
		verrs["task_after"] = "Not found in this schedule"		
	}

	if len(verrs) == 0 {
		tddb := TaskDependencyDB{db}
		_, err := tddb.SelectOne("task_before=$1 AND task_after=$2", td.Before, td.After)
		if err != nil && err != sql.ErrNoRows {
			verrs["all"] = "Error occured while checking for prexisting dependency"
		} else if err != sql.ErrNoRows {
			verrs["all"] = "A dependency exists between these tasks"
		}

	}

	return verrs
}


func causesCycle(data TaskDependency, db database.Queryer) (bool, error) {
	if data.Before == data.After {
		return true, nil
	}

	tddb := TaskDependencyDB{db}
	ids, err := tddb.SelectIdsForSchedule(data.Schedule)
	if err != nil {
		return false, err
	}

	edges := append(ids, []int{data.Before, data.After})
	return detectCycle(edges), nil
}


func detectCycle(edges [][]int) bool {
	// build index of node -> children
    nodes := make(map[int][]int, 0)
    for _, e := range edges {
    	if _, exists := nodes[e[0]]; exists {
    		nodes[e[0]] = append(nodes[e[0]], e[1])
    	} else {
    		nodes[e[0]] = []int{e[1]}
    	}
    }

    visited := make(map[int]bool, 0)
    // not visited == not in map
    // visit in progress == false in map
    // visit complete == true in map


    detectedCycle := false

    var visitNode func(int)
    visitNode = func(n int) {
    	if done, found := visited[n]; found && done {
    		return
    	}

    	// mark node as started
        visited[n] = false

        // visit nodes children
        for _, child := range nodes[n] {
        	done, found := visited[child]
	        if found && !done {
	        	// revisited unfinished path
	            detectedCycle = true
	            break
	        } else if !found {
	        	// unvisited node, recurse
	        	visitNode(child)
	        }
        }

        // visited all child edges
        visited[n] = true
    }

    for k := range nodes {
    	if !detectedCycle {    		
	    	visitNode(k)
    	}
    }

    return detectedCycle
}