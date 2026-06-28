package main

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinearDependency(t *testing.T) {
	o := NewOrchestrator()

	wf, err := o.CreateWorkflow("linear-test", WorkflowConfig{
		MaxConcurrentTasks:     5,
		DefaultFailureStrategy: FailureStrategyFailFast,
	})
	require.NoError(t, err)

	taskA, err := o.AddTask(wf.ID, "task-A", map[string]any{"input": "a"})
	require.NoError(t, err)

	taskB, err := o.AddTask(wf.ID, "task-B", map[string]any{"input": "b"})
	require.NoError(t, err)

	taskC, err := o.AddTask(wf.ID, "task-C", map[string]any{"input": "c"})
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskA.ID, taskB.ID)
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskB.ID, taskC.ID)
	require.NoError(t, err)

	err = o.StartWorkflow(wf.ID)
	require.NoError(t, err)

	assert.Equal(t, WorkflowStatusRunning, wf.Status)

	readyTasks, err := o.GetReadyTasks(wf.ID)
	require.NoError(t, err)
	assert.Len(t, readyTasks, 1)
	assert.Equal(t, taskA.ID, readyTasks[0].ID)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.CompleteTask(taskA.ID, map[string]any{"output": "a"})
	require.NoError(t, err)

	readyTasks, err = o.GetReadyTasks(wf.ID)
	require.NoError(t, err)
	assert.Len(t, readyTasks, 1)
	assert.Equal(t, taskB.ID, readyTasks[0].ID)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.CompleteTask(taskB.ID, map[string]any{"output": "b"})
	require.NoError(t, err)

	readyTasks, err = o.GetReadyTasks(wf.ID)
	require.NoError(t, err)
	assert.Len(t, readyTasks, 1)
	assert.Equal(t, taskC.ID, readyTasks[0].ID)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.CompleteTask(taskC.ID, map[string]any{"output": "c"})
	require.NoError(t, err)

	wf, err = o.GetWorkflow(wf.ID)
	require.NoError(t, err)
	assert.Equal(t, WorkflowStatusCompleted, wf.Status)
}

func TestParallelExecution(t *testing.T) {
	o := NewOrchestrator()

	wf, err := o.CreateWorkflow("parallel-test", WorkflowConfig{
		MaxConcurrentTasks:     5,
		DefaultFailureStrategy: FailureStrategyFailFast,
	})
	require.NoError(t, err)

	taskA, err := o.AddTask(wf.ID, "task-A", map[string]any{"input": "a"})
	require.NoError(t, err)

	taskB, err := o.AddTask(wf.ID, "task-B", map[string]any{"input": "b"})
	require.NoError(t, err)

	taskC, err := o.AddTask(wf.ID, "task-C", map[string]any{"input": "c"})
	require.NoError(t, err)

	taskD, err := o.AddTask(wf.ID, "task-D", map[string]any{"input": "d"})
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskA.ID, taskB.ID)
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskA.ID, taskC.ID)
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskB.ID, taskD.ID)
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskC.ID, taskD.ID)
	require.NoError(t, err)

	err = o.StartWorkflow(wf.ID)
	require.NoError(t, err)

	assert.Equal(t, WorkflowStatusRunning, wf.Status)

	readyTasks, err := o.GetReadyTasks(wf.ID)
	require.NoError(t, err)
	assert.Len(t, readyTasks, 1)
	assert.Equal(t, taskA.ID, readyTasks[0].ID)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.CompleteTask(taskA.ID, map[string]any{"output": "a"})
	require.NoError(t, err)

	readyTasks, err = o.GetReadyTasks(wf.ID)
	require.NoError(t, err)
	assert.Len(t, readyTasks, 2)

	taskMap := make(map[uuid.UUID]bool)
	for _, task := range readyTasks {
		taskMap[task.ID] = true
	}
	assert.True(t, taskMap[taskB.ID])
	assert.True(t, taskMap[taskC.ID])

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.CompleteTask(taskB.ID, map[string]any{"output": "b"})
	require.NoError(t, err)

	readyTasks, err = o.GetReadyTasks(wf.ID)
	require.NoError(t, err)
	assert.Len(t, readyTasks, 0)

	err = o.CompleteTask(taskC.ID, map[string]any{"output": "c"})
	require.NoError(t, err)

	readyTasks, err = o.GetReadyTasks(wf.ID)
	require.NoError(t, err)
	assert.Len(t, readyTasks, 1)
	assert.Equal(t, taskD.ID, readyTasks[0].ID)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.CompleteTask(taskD.ID, map[string]any{"output": "d"})
	require.NoError(t, err)

	wf, err = o.GetWorkflow(wf.ID)
	require.NoError(t, err)
	assert.Equal(t, WorkflowStatusCompleted, wf.Status)
}

func TestDiamondDependency(t *testing.T) {
	o := NewOrchestrator()

	wf, err := o.CreateWorkflow("diamond-test", WorkflowConfig{
		MaxConcurrentTasks:     5,
		DefaultFailureStrategy: FailureStrategyFailFast,
	})
	require.NoError(t, err)

	taskA, err := o.AddTask(wf.ID, "task-A", map[string]any{"input": "a"})
	require.NoError(t, err)

	taskB, err := o.AddTask(wf.ID, "task-B", map[string]any{"input": "b"})
	require.NoError(t, err)

	taskC, err := o.AddTask(wf.ID, "task-C", map[string]any{"input": "c"})
	require.NoError(t, err)

	taskD, err := o.AddTask(wf.ID, "task-D", map[string]any{"input": "d"})
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskA.ID, taskB.ID)
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskA.ID, taskC.ID)
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskB.ID, taskD.ID)
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskC.ID, taskD.ID)
	require.NoError(t, err)

	err = o.StartWorkflow(wf.ID)
	require.NoError(t, err)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.CompleteTask(taskA.ID, map[string]any{"output": "a"})
	require.NoError(t, err)

	readyTasks, err := o.GetReadyTasks(wf.ID)
	require.NoError(t, err)
	assert.Len(t, readyTasks, 2)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.CompleteTask(taskB.ID, map[string]any{"output": "b"})
	require.NoError(t, err)

	readyTasks, err = o.GetReadyTasks(wf.ID)
	require.NoError(t, err)
	assert.Len(t, readyTasks, 0)

	err = o.CompleteTask(taskC.ID, map[string]any{"output": "c"})
	require.NoError(t, err)

	readyTasks, err = o.GetReadyTasks(wf.ID)
	require.NoError(t, err)
	assert.Len(t, readyTasks, 1)
	assert.Equal(t, taskD.ID, readyTasks[0].ID)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.CompleteTask(taskD.ID, map[string]any{"output": "d"})
	require.NoError(t, err)

	wf, err = o.GetWorkflow(wf.ID)
	require.NoError(t, err)
	assert.Equal(t, WorkflowStatusCompleted, wf.Status)
}

func TestCycleDetection(t *testing.T) {
	o := NewOrchestrator()

	wf, err := o.CreateWorkflow("cycle-test", WorkflowConfig{
		MaxConcurrentTasks:     5,
		DefaultFailureStrategy: FailureStrategyFailFast,
	})
	require.NoError(t, err)

	taskA, err := o.AddTask(wf.ID, "task-A", map[string]any{"input": "a"})
	require.NoError(t, err)

	taskB, err := o.AddTask(wf.ID, "task-B", map[string]any{"input": "b"})
	require.NoError(t, err)

	taskC, err := o.AddTask(wf.ID, "task-C", map[string]any{"input": "c"})
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskA.ID, taskB.ID)
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskB.ID, taskC.ID)
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskC.ID, taskA.ID)
	assert.Equal(t, ErrCycleDetected, err)
}

func TestFailureStrategyFailFast(t *testing.T) {
	o := NewOrchestrator()

	wf, err := o.CreateWorkflow("fail-fast-test", WorkflowConfig{
		MaxConcurrentTasks:     5,
		DefaultFailureStrategy: FailureStrategyFailFast,
	})
	require.NoError(t, err)

	taskA, err := o.AddTask(wf.ID, "task-A", map[string]any{"input": "a"})
	require.NoError(t, err)

	taskB, err := o.AddTask(wf.ID, "task-B", map[string]any{"input": "b"})
	require.NoError(t, err)

	taskC, err := o.AddTask(wf.ID, "task-C", map[string]any{"input": "c"})
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskA.ID, taskB.ID)
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskB.ID, taskC.ID)
	require.NoError(t, err)

	err = o.StartWorkflow(wf.ID)
	require.NoError(t, err)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.CompleteTask(taskA.ID, map[string]any{"output": "a"})
	require.NoError(t, err)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.FailTask(taskB.ID, "failed")
	require.NoError(t, err)

	wf, err = o.GetWorkflow(wf.ID)
	require.NoError(t, err)
	assert.Equal(t, WorkflowStatusFailed, wf.Status)

	taskC, err = o.GetTask(taskC.ID)
	require.NoError(t, err)
	assert.Equal(t, TaskStatusCancelled, taskC.Status)
}

func TestFailureStrategyRetry(t *testing.T) {
	o := NewOrchestrator()

	wf, err := o.CreateWorkflow("retry-test", WorkflowConfig{
		MaxConcurrentTasks:     5,
		DefaultFailureStrategy: FailureStrategyRetry,
	})
	require.NoError(t, err)

	taskA, err := o.AddTask(wf.ID, "task-A", map[string]any{"input": "a"},
		WithMaxRetries(2))
	require.NoError(t, err)

	err = o.StartWorkflow(wf.ID)
	require.NoError(t, err)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.FailTask(taskA.ID, "first failure")
	require.NoError(t, err)

	taskA, err = o.GetTask(taskA.ID)
	require.NoError(t, err)
	assert.Equal(t, TaskStatusReady, taskA.Status)
	assert.Equal(t, 1, taskA.RetryCount)

	readyTasks, err := o.GetReadyTasks(wf.ID)
	require.NoError(t, err)
	assert.Len(t, readyTasks, 1)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.FailTask(taskA.ID, "second failure")
	require.NoError(t, err)

	taskA, err = o.GetTask(taskA.ID)
	require.NoError(t, err)
	assert.Equal(t, TaskStatusReady, taskA.Status)
	assert.Equal(t, 2, taskA.RetryCount)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.FailTask(taskA.ID, "third failure")
	require.NoError(t, err)

	wf, err = o.GetWorkflow(wf.ID)
	require.NoError(t, err)
	assert.Equal(t, WorkflowStatusFailed, wf.Status)
}

func TestFailureStrategySkip(t *testing.T) {
	o := NewOrchestrator()

	wf, err := o.CreateWorkflow("skip-test", WorkflowConfig{
		MaxConcurrentTasks:     5,
		DefaultFailureStrategy: FailureStrategySkip,
	})
	require.NoError(t, err)

	taskA, err := o.AddTask(wf.ID, "task-A", map[string]any{"input": "a"})
	require.NoError(t, err)

	taskB, err := o.AddTask(wf.ID, "task-B", map[string]any{"input": "b"})
	require.NoError(t, err)

	taskC, err := o.AddTask(wf.ID, "task-C", map[string]any{"input": "c"})
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskA.ID, taskB.ID)
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskB.ID, taskC.ID)
	require.NoError(t, err)

	err = o.StartWorkflow(wf.ID)
	require.NoError(t, err)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.CompleteTask(taskA.ID, map[string]any{"output": "a"})
	require.NoError(t, err)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.FailTask(taskB.ID, "failed")
	require.NoError(t, err)

	taskB, err = o.GetTask(taskB.ID)
	require.NoError(t, err)
	assert.Equal(t, TaskStatusSkipped, taskB.Status)

	readyTasks, err := o.GetReadyTasks(wf.ID)
	require.NoError(t, err)
	assert.Len(t, readyTasks, 1)
	assert.Equal(t, taskC.ID, readyTasks[0].ID)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.CompleteTask(taskC.ID, map[string]any{"output": "c"})
	require.NoError(t, err)

	wf, err = o.GetWorkflow(wf.ID)
	require.NoError(t, err)
	assert.Equal(t, WorkflowStatusCompleted, wf.Status)
}

func TestConcurrencyLimit(t *testing.T) {
	o := NewOrchestrator()

	wf, err := o.CreateWorkflow("concurrency-test", WorkflowConfig{
		MaxConcurrentTasks:     3,
		DefaultFailureStrategy: FailureStrategyFailFast,
	})
	require.NoError(t, err)

	var tasks []*WorkflowTask
	for i := 0; i < 10; i++ {
		task, err := o.AddTask(wf.ID, "task-"+string(rune('A'+i)), map[string]any{"input": i})
		require.NoError(t, err)
		tasks = append(tasks, task)
	}

	err = o.StartWorkflow(wf.ID)
	require.NoError(t, err)

	readyTasks, err := o.GetReadyTasks(wf.ID)
	require.NoError(t, err)
	assert.Len(t, readyTasks, 10)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	assert.Equal(t, ErrMaxConcurrentTasks, err)

	runningCount, err := o.GetRunningCount(wf.ID)
	require.NoError(t, err)
	assert.Equal(t, 3, runningCount)

	err = o.CompleteTask(tasks[0].ID, map[string]any{"output": 0})
	require.NoError(t, err)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	runningCount, err = o.GetRunningCount(wf.ID)
	require.NoError(t, err)
	assert.Equal(t, 3, runningCount)

	err = o.CompleteTask(tasks[1].ID, map[string]any{"output": 1})
	require.NoError(t, err)

	err = o.CompleteTask(tasks[2].ID, map[string]any{"output": 2})
	require.NoError(t, err)

	err = o.CompleteTask(tasks[3].ID, map[string]any{"output": 3})
	require.NoError(t, err)

	for i := 4; i < 10; i++ {
		_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
		require.NoError(t, err)

		err = o.CompleteTask(tasks[i].ID, map[string]any{"output": i})
		require.NoError(t, err)
	}

	wf, err = o.GetWorkflow(wf.ID)
	require.NoError(t, err)
	assert.Equal(t, WorkflowStatusCompleted, wf.Status)
}

func TestCancelWorkflow(t *testing.T) {
	o := NewOrchestrator()

	wf, err := o.CreateWorkflow("cancel-test", WorkflowConfig{
		MaxConcurrentTasks:     5,
		DefaultFailureStrategy: FailureStrategyFailFast,
	})
	require.NoError(t, err)

	taskA, err := o.AddTask(wf.ID, "task-A", map[string]any{"input": "a"})
	require.NoError(t, err)

	taskB, err := o.AddTask(wf.ID, "task-B", map[string]any{"input": "b"})
	require.NoError(t, err)

	err = o.AddDependency(wf.ID, taskA.ID, taskB.ID)
	require.NoError(t, err)

	err = o.StartWorkflow(wf.ID)
	require.NoError(t, err)

	_, err = o.ClaimTask(wf.ID, uuid.New(), uuid.New())
	require.NoError(t, err)

	err = o.CancelWorkflow(wf.ID)
	require.NoError(t, err)

	wf, err = o.GetWorkflow(wf.ID)
	require.NoError(t, err)
	assert.Equal(t, WorkflowStatusCancelled, wf.Status)

	taskA, err = o.GetTask(taskA.ID)
	require.NoError(t, err)
	assert.Equal(t, TaskStatusCancelled, taskA.Status)

	taskB, err = o.GetTask(taskB.ID)
	require.NoError(t, err)
	assert.Equal(t, TaskStatusCancelled, taskB.Status)
}