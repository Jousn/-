package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCycleDetected           = errors.New("cycle detected in task dependencies")
	ErrWorkflowNotFound        = errors.New("workflow not found")
	ErrTaskNotFound            = errors.New("task not found")
	ErrTaskNotReady            = errors.New("task is not ready")
	ErrTaskAlreadyExists       = errors.New("task with same ref already exists")
	ErrDependencyAlreadyExists = errors.New("dependency already exists")
	ErrNoReadyTasks            = errors.New("no ready tasks available")
	ErrMaxConcurrentTasks      = errors.New("max concurrent tasks reached")
)

type Orchestrator struct {
	mu               sync.RWMutex
	workflows        map[uuid.UUID]*Workflow
	tasks            map[uuid.UUID]*WorkflowTask
	dependencies     map[uuid.UUID]*TaskDependency
	workflowTasks    map[uuid.UUID][]uuid.UUID
	taskDependencies map[uuid.UUID][]uuid.UUID
	taskDependents   map[uuid.UUID][]uuid.UUID
	runningCount     map[uuid.UUID]int
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		workflows:        make(map[uuid.UUID]*Workflow),
		tasks:            make(map[uuid.UUID]*WorkflowTask),
		dependencies:     make(map[uuid.UUID]*TaskDependency),
		workflowTasks:    make(map[uuid.UUID][]uuid.UUID),
		taskDependencies: make(map[uuid.UUID][]uuid.UUID),
		taskDependents:   make(map[uuid.UUID][]uuid.UUID),
		runningCount:     make(map[uuid.UUID]int),
	}
}

func (o *Orchestrator) CreateWorkflow(name string, config WorkflowConfig) (*Workflow, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	wf := &Workflow{
		ID:        uuid.New(),
		Name:      name,
		Status:    WorkflowStatusPending,
		Config:    config,
		CreatedAt: time.Now(),
	}

	o.workflows[wf.ID] = wf
	o.workflowTasks[wf.ID] = []uuid.UUID{}
	o.runningCount[wf.ID] = 0

	return wf, nil
}

func (o *Orchestrator) GetWorkflow(id uuid.UUID) (*Workflow, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	wf, ok := o.workflows[id]
	if !ok {
		return nil, ErrWorkflowNotFound
	}
	return wf, nil
}

func (o *Orchestrator) AddTask(workflowID uuid.UUID, taskRef string, input map[string]any, options ...TaskOption) (*WorkflowTask, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	wf, ok := o.workflows[workflowID]
	if !ok {
		return nil, ErrWorkflowNotFound
	}

	for _, taskID := range o.workflowTasks[workflowID] {
		task := o.tasks[taskID]
		if task.TaskRef == taskRef {
			return nil, ErrTaskAlreadyExists
		}
	}

	task := &WorkflowTask{
		ID:             uuid.New(),
		WorkflowID:     workflowID,
		TaskRef:        taskRef,
		Status:         TaskStatusPending,
		FailureStrategy: wf.Config.DefaultFailureStrategy,
		MaxRetries:     0,
		RetryCount:     0,
		TimeoutSecs:    wf.Config.DefaultTimeoutSecs,
		Input:          input,
		Output:         nil,
		CreatedAt:      time.Now(),
	}

	for _, opt := range options {
		opt(task)
	}

	o.tasks[task.ID] = task
	o.workflowTasks[workflowID] = append(o.workflowTasks[workflowID], task.ID)
	o.taskDependencies[task.ID] = []uuid.UUID{}
	o.taskDependents[task.ID] = []uuid.UUID{}

	return task, nil
}

type TaskOption func(*WorkflowTask)

func WithFailureStrategy(strategy FailureStrategy) TaskOption {
	return func(t *WorkflowTask) {
		t.FailureStrategy = strategy
	}
}

func WithMaxRetries(maxRetries int) TaskOption {
	return func(t *WorkflowTask) {
		t.MaxRetries = maxRetries
	}
}

func WithTimeout(secs int) TaskOption {
	return func(t *WorkflowTask) {
		t.TimeoutSecs = secs
	}
}

func (o *Orchestrator) AddDependency(workflowID, fromTaskID, toTaskID uuid.UUID) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.workflowContainsTask(workflowID, fromTaskID) || !o.workflowContainsTask(workflowID, toTaskID) {
		return ErrTaskNotFound
	}

	for _, dep := range o.taskDependencies[toTaskID] {
		if dep == fromTaskID {
			return ErrDependencyAlreadyExists
		}
	}

	if o.hasCycle(workflowID, fromTaskID, toTaskID) {
		return ErrCycleDetected
	}

	dep := &TaskDependency{
		ID:         uuid.New(),
		WorkflowID: workflowID,
		FromTaskID: fromTaskID,
		ToTaskID:   toTaskID,
	}

	o.dependencies[dep.ID] = dep
	o.taskDependencies[toTaskID] = append(o.taskDependencies[toTaskID], fromTaskID)
	o.taskDependents[fromTaskID] = append(o.taskDependents[fromTaskID], toTaskID)

	return nil
}

func (o *Orchestrator) workflowContainsTask(workflowID, taskID uuid.UUID) bool {
	tasks, ok := o.workflowTasks[workflowID]
	if !ok {
		return false
	}
	for _, t := range tasks {
		if t == taskID {
			return true
		}
	}
	return false
}

func (o *Orchestrator) hasCycle(workflowID, fromTaskID, toTaskID uuid.UUID) bool {
	visited := make(map[uuid.UUID]bool)

	var dfs func(node uuid.UUID) bool
	dfs = func(node uuid.UUID) bool {
		if node == fromTaskID {
			return true
		}
		if visited[node] {
			return false
		}
		visited[node] = true

		for _, dependent := range o.taskDependents[node] {
			if dfs(dependent) {
				return true
			}
		}

		return false
	}

	return dfs(toTaskID)
}

func (o *Orchestrator) StartWorkflow(workflowID uuid.UUID) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	wf, ok := o.workflows[workflowID]
	if !ok {
		return ErrWorkflowNotFound
	}

	if wf.Status != WorkflowStatusPending {
		return fmt.Errorf("workflow is not in pending state: %s", wf.Status)
	}

	if o.hasCycleInWorkflow(workflowID) {
		return ErrCycleDetected
	}

	now := time.Now()
	wf.Status = WorkflowStatusRunning
	wf.StartedAt = &now

	o.updateReadyTasks(workflowID)

	return nil
}

func (o *Orchestrator) hasCycleInWorkflow(workflowID uuid.UUID) bool {
	visited := make(map[uuid.UUID]bool)
	inStack := make(map[uuid.UUID]bool)

	tasks := o.workflowTasks[workflowID]
	for _, taskID := range tasks {
		if !visited[taskID] {
			if o.dfsCycle(taskID, visited, inStack) {
				return true
			}
		}
	}
	return false
}

func (o *Orchestrator) dfsCycle(node uuid.UUID, visited, inStack map[uuid.UUID]bool) bool {
	visited[node] = true
	inStack[node] = true

	for _, dependent := range o.taskDependents[node] {
		if !visited[dependent] {
			if o.dfsCycle(dependent, visited, inStack) {
				return true
			}
		} else if inStack[dependent] {
			return true
		}
	}

	inStack[node] = false
	return false
}

func (o *Orchestrator) updateReadyTasks(workflowID uuid.UUID) {
	tasks := o.workflowTasks[workflowID]
	for _, taskID := range tasks {
		task := o.tasks[taskID]
		if task.Status == TaskStatusPending {
			if o.isTaskReady(taskID) {
				task.Status = TaskStatusReady
			}
		}
	}
}

func (o *Orchestrator) isTaskReady(taskID uuid.UUID) bool {
	deps := o.taskDependencies[taskID]
	for _, dep := range deps {
		depTask := o.tasks[dep]
		if depTask.Status != TaskStatusCompleted && depTask.Status != TaskStatusSkipped {
			return false
		}
	}
	return true
}

func (o *Orchestrator) CancelWorkflow(workflowID uuid.UUID) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	wf, ok := o.workflows[workflowID]
	if !ok {
		return ErrWorkflowNotFound
	}

	if wf.Status == WorkflowStatusCompleted || wf.Status == WorkflowStatusFailed || wf.Status == WorkflowStatusCancelled {
		return nil
	}

	now := time.Now()
	wf.Status = WorkflowStatusCancelled
	wf.CompletedAt = &now

	tasks := o.workflowTasks[workflowID]
	for _, taskID := range tasks {
		task := o.tasks[taskID]
		if task.Status != TaskStatusCompleted && task.Status != TaskStatusFailed {
			task.Status = TaskStatusCancelled
			task.CompletedAt = &now
		}
	}

	return nil
}

func (o *Orchestrator) ClaimTask(workflowID, agentID, runtimeID uuid.UUID) (*WorkflowTask, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	wf, ok := o.workflows[workflowID]
	if !ok {
		return nil, ErrWorkflowNotFound
	}

	if wf.Status != WorkflowStatusRunning {
		return nil, fmt.Errorf("workflow is not running: %s", wf.Status)
	}

	if o.runningCount[workflowID] >= wf.Config.MaxConcurrentTasks {
		return nil, ErrMaxConcurrentTasks
	}

	tasks := o.workflowTasks[workflowID]
	var readyTasks []*WorkflowTask
	for _, taskID := range tasks {
		task := o.tasks[taskID]
		if task.Status == TaskStatusReady {
			readyTasks = append(readyTasks, task)
		}
	}

	if len(readyTasks) == 0 {
		return nil, ErrNoReadyTasks
	}

	task := readyTasks[0]
	now := time.Now()
	task.Status = TaskStatusRunning
	task.StartedAt = &now
	task.AgentID = &agentID
	task.RuntimeID = &runtimeID

	o.runningCount[workflowID]++

	return task, nil
}

func (o *Orchestrator) CompleteTask(taskID uuid.UUID, output map[string]any) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	task, ok := o.tasks[taskID]
	if !ok {
		return ErrTaskNotFound
	}

	if task.Status != TaskStatusRunning {
		return fmt.Errorf("task is not running: %s", task.Status)
	}

	now := time.Now()
	task.Status = TaskStatusCompleted
	task.Output = output
	task.CompletedAt = &now

	o.runningCount[task.WorkflowID]--
	o.updateReadyTasks(task.WorkflowID)
	o.checkWorkflowCompletion(task.WorkflowID)

	return nil
}

func (o *Orchestrator) FailTask(taskID uuid.UUID, errorMessage string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	task, ok := o.tasks[taskID]
	if !ok {
		return ErrTaskNotFound
	}

	if task.Status != TaskStatusRunning {
		return fmt.Errorf("task is not running: %s", task.Status)
	}

	task.Status = TaskStatusFailed
	task.ErrorMessage = errorMessage

	o.runningCount[task.WorkflowID]--

	return o.handleTaskFailure(task)
}

func (o *Orchestrator) handleTaskFailure(task *WorkflowTask) error {
	switch task.FailureStrategy {
	case FailureStrategyFailFast:
		return o.failFast(task.WorkflowID)
	case FailureStrategyRetry:
		if task.RetryCount < task.MaxRetries {
			task.RetryCount++
			task.Status = TaskStatusPending
			task.StartedAt = nil
			task.CompletedAt = nil
			task.ErrorMessage = ""
			o.updateReadyTasks(task.WorkflowID)
			return nil
		}
		return o.failFast(task.WorkflowID)
	case FailureStrategySkip:
		now := time.Now()
		task.Status = TaskStatusSkipped
		task.CompletedAt = &now
		o.updateReadyTasks(task.WorkflowID)
		o.checkWorkflowCompletion(task.WorkflowID)
		return nil
	default:
		return o.failFast(task.WorkflowID)
	}
}

func (o *Orchestrator) failFast(workflowID uuid.UUID) error {
	wf := o.workflows[workflowID]
	now := time.Now()
	wf.Status = WorkflowStatusFailed
	wf.CompletedAt = &now

	tasks := o.workflowTasks[workflowID]
	for _, taskID := range tasks {
		task := o.tasks[taskID]
		if task.Status != TaskStatusCompleted && task.Status != TaskStatusFailed {
			task.Status = TaskStatusCancelled
			task.CompletedAt = &now
		}
	}

	return nil
}

func (o *Orchestrator) checkWorkflowCompletion(workflowID uuid.UUID) {
	wf := o.workflows[workflowID]
	if wf.Status != WorkflowStatusRunning {
		return
	}

	tasks := o.workflowTasks[workflowID]
	allDone := true
	for _, taskID := range tasks {
		task := o.tasks[taskID]
		if task.Status != TaskStatusCompleted && task.Status != TaskStatusSkipped {
			allDone = false
			break
		}
	}

	if allDone {
		now := time.Now()
		wf.Status = WorkflowStatusCompleted
		wf.CompletedAt = &now
	}
}

func (o *Orchestrator) GetReadyTasks(workflowID uuid.UUID) ([]*WorkflowTask, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if _, ok := o.workflows[workflowID]; !ok {
		return nil, ErrWorkflowNotFound
	}

	var readyTasks []*WorkflowTask
	tasks := o.workflowTasks[workflowID]
	for _, taskID := range tasks {
		task := o.tasks[taskID]
		if task.Status == TaskStatusReady {
			readyTasks = append(readyTasks, task)
		}
	}

	return readyTasks, nil
}

func (o *Orchestrator) GetRunningCount(workflowID uuid.UUID) (int, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if _, ok := o.workflows[workflowID]; !ok {
		return 0, ErrWorkflowNotFound
	}

	return o.runningCount[workflowID], nil
}

func (o *Orchestrator) GetTask(taskID uuid.UUID) (*WorkflowTask, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	task, ok := o.tasks[taskID]
	if !ok {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (o *Orchestrator) GetTasks(workflowID uuid.UUID) ([]*WorkflowTask, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if _, ok := o.workflows[workflowID]; !ok {
		return nil, ErrWorkflowNotFound
	}

	var tasks []*WorkflowTask
	for _, taskID := range o.workflowTasks[workflowID] {
		tasks = append(tasks, o.tasks[taskID])
	}

	return tasks, nil
}