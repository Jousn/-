package main

import (
	"time"

	"github.com/google/uuid"
)

type WorkflowStatus string

const (
	WorkflowStatusPending   WorkflowStatus = "pending"
	WorkflowStatusRunning   WorkflowStatus = "running"
	WorkflowStatusCompleted WorkflowStatus = "completed"
	WorkflowStatusFailed    WorkflowStatus = "failed"
	WorkflowStatusCancelled WorkflowStatus = "cancelled"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusReady      TaskStatus = "ready"
	TaskStatusRunning    TaskStatus = "running"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusCancelled  TaskStatus = "cancelled"
	TaskStatusSkipped    TaskStatus = "skipped"
)

type FailureStrategy string

const (
	FailureStrategyFailFast FailureStrategy = "fail_fast"
	FailureStrategyRetry    FailureStrategy = "retry"
	FailureStrategySkip     FailureStrategy = "skip"
)

type WorkflowConfig struct {
	MaxConcurrentTasks     int              `json:"max_concurrent_tasks"`
	DefaultFailureStrategy FailureStrategy  `json:"default_failure_strategy"`
	DefaultTimeoutSecs     int              `json:"default_timeout_secs"`
}

type Workflow struct {
	ID          uuid.UUID       `json:"id"`
	IssueID     uuid.UUID       `json:"issue_id"`
	Name        string          `json:"name"`
	Status      WorkflowStatus  `json:"status"`
	Config      WorkflowConfig  `json:"config"`
	CreatedAt   time.Time       `json:"created_at"`
	StartedAt   *time.Time      `json:"started_at"`
	CompletedAt *time.Time      `json:"completed_at"`
}

type WorkflowTask struct {
	ID             uuid.UUID       `json:"id"`
	WorkflowID     uuid.UUID       `json:"workflow_id"`
	TaskRef        string          `json:"task_ref"`
	Status         TaskStatus      `json:"status"`
	FailureStrategy FailureStrategy `json:"failure_strategy"`
	MaxRetries     int             `json:"max_retries"`
	RetryCount     int             `json:"retry_count"`
	TimeoutSecs    int             `json:"timeout_secs"`
	Input          map[string]any  `json:"input"`
	Output         map[string]any  `json:"output"`
	CreatedAt      time.Time       `json:"created_at"`
	StartedAt      *time.Time      `json:"started_at"`
	CompletedAt    *time.Time      `json:"completed_at"`
	ErrorMessage   string          `json:"error_message"`
	AgentID        *uuid.UUID      `json:"agent_id"`
	RuntimeID      *uuid.UUID      `json:"runtime_id"`
}

type TaskDependency struct {
	ID          uuid.UUID `json:"id"`
	WorkflowID  uuid.UUID `json:"workflow_id"`
	FromTaskID  uuid.UUID `json:"from_task_id"`
	ToTaskID    uuid.UUID `json:"to_task_id"`
}