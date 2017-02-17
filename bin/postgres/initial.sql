CREATE TABLE schedules (
	schedule_id		SERIAL PRIMARY KEY,
	schedule_name	VARCHAR(50),
	data_date		DATE
);
CREATE TABLE tasks (
	task_id			SERIAL PRIMARY KEY,
	schedule		INT NOT NULL REFERENCES schedules(schedule_id) ON DELETE CASCADE,
	task_code		VARCHAR(15) NOT NULL,
	task_name		VARCHAR(50),
	duration		INT NOT NULL DEFAULT 1,
	remaining		INT NOT NULL DEFAULT 1,
	start_early		DATE,
	start_late		DATE,
	start_actual	DATE,
	finish_early	DATE,
	finish_late		DATE,
	finish_actual	DATE
);

CREATE TABLE task_deps (
	schedule 		INT NOT NULL REFERENCES schedules(schedule_id) On DELETE CASCADE,
	task_before		INT NOT NULL REFERENCES tasks(task_id) ON DELETE CASCADE,
	task_after		INT NOT NULL REFERENCES tasks(task_id) ON DELETE CASCADE,
	lag				INT,
	type			INT,
	PRIMARY KEY (task_before, task_after)
);

INSERT INTO schedules (schedule_id, schedule_name, data_date)
VALUES 
	(1, 'Sample Schedule', CURRENT_DATE)
;

SELECT setval('schedules_schedule_id_seq', 1);

INSERT INTO tasks (task_id, schedule, task_code, task_name, duration)
VALUES
	(1, 1, 'begin', 'prep for tasks', 2),
	(2, 1, 'a1', 'first task of a', 3),
	(3, 1, 'a2', 'second task of a', 2),
	(4, 1, 'b1', 'first task of b', 1),
	(5, 1, 'b2', 'second task of b', 5),
	(6, 1, 'end', 'finish work', 2)
;

SELECT setval('tasks_task_id_seq', 6);

INSERT INTO task_deps (schedule, lag, type, task_before, task_after)
VALUES
	(1, 0, 0, 1, 2),
	(1, 1, 0, 1, 4),
	(1, 4, 0, 2, 3),
	(1, 0, 0, 5, 6),
	(1, 2, 0, 3, 6)
;