package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	STATE_INIT      = "INIT"
	STATE_STARTED   = "STARTED"
	STATE_COMPLETED = "COMPLETED"
)

var workflowStorageInstance *workflowStorage

func InitDB(db *sql.DB) {
	workflowStorageInstance = &workflowStorage{db: db}
	workflowStorageInstance.setupTables()
}

type workflowStorage struct {
	db *sql.DB
}

type storeStartExecutionReq struct {
	ExecutionId  string
	RunId        string
	StartPayload string
	TaskId       string
}

func (ws *workflowStorage) setupTables() {
	db := ws.db

	dropTables := `
        drop table if exists workflow_execution;
        drop table if exists workflow_task_queue;
    `
	_, err := db.Exec(dropTables)
	if err != nil {
		log.Fatalf("%q: %s\n", err, dropTables)
	}

	workflowExecutionDDL := `create table if not exists workflow_execution (
        run_id text primary key,
        execution_id text,
        execution_state text not null,
        start_payload text); `
	_, err = db.Exec(workflowExecutionDDL)
	if err != nil {
		log.Fatalf("%q: %s\n", err, workflowExecutionDDL)
	}

	workflowTaskDDL := `create table if not exists workflow_task_queue (
        task_id text primary key,
        task_payload text,
        task_state text,
        execution_id text,
        run_id text,
        create_timestamp int,
        update_timestamp int
    ); `
	_, err = db.Exec(workflowTaskDDL)
	if err != nil {
		log.Fatalf("%q: %s\n", err, workflowTaskDDL)
	}
}

type WorkflowTaskEntity struct {
	TaskId      string `db:"task_id"`
	TaskPayload string `db:"task_payload"`
	TaskState   string `db:"task_state"`
	ExecutionId string `db:"execution_id"`
	RunId       string `db:"run_id"`
}

func (ws *workflowStorage) pollPendingWorkflowTask() (*WorkflowTaskEntity, error) {
	db := ws.db
	stmt := `SELECT task_id, task_payload, task_state, execution_id, run_id FROM workflow_task_queue WHERE task_state = ? ORDER BY create_timestamp ASC LIMIT 1;`
	rows, err := db.Query(stmt, STATE_INIT)
	if err != nil {
		return nil, err
	}
	log.Printf("rows %s \n", rows)

	taskEntities := make([]WorkflowTaskEntity, 0, 1)
	err = sqlx.StructScan(rows, &taskEntities)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		} else {
			return nil, nil
		}
	}
	if len(taskEntities) > 0 {
		res := taskEntities[0]
		return &res, nil
	}
	return nil, nil
}

func (ws *workflowStorage) storeStartExecution(req storeStartExecutionReq) error {

	db := ws.db
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// store workflow execution state
	nowSec := time.Now().Unix()
	stmt, err := tx.Prepare("INSERT INTO workflow_execution(execution_id, run_id, execution_state, start_payload) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(req.ExecutionId, req.RunId, STATE_INIT, req.StartPayload)
	if err != nil {
		log.Fatal(err)
	}

	// store a worklfow task
	stmt, err = tx.Prepare("INSERT INTO workflow_task_queue (task_id, task_payload, task_state, execution_id, run_id, create_timestamp, update_timestamp) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		req.TaskId, req.StartPayload, STATE_INIT, req.ExecutionId, req.RunId, nowSec, nowSec)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
