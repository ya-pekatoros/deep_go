package main

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	tasks     priorityQueue
	positions map[int]*scheduledTask
}

type scheduledTask struct {
	task  Task
	index int
}

type priorityQueue []*scheduledTask

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	if pq[i].task.Priority == pq[j].task.Priority {
		return pq[i].task.Identifier < pq[j].task.Identifier
	}

	return pq[i].task.Priority > pq[j].task.Priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x any) {
	item := x.(*scheduledTask)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	lastIndex := len(old) - 1
	item := old[lastIndex]
	item.index = -1
	*pq = old[:lastIndex]

	return item
}

func NewScheduler() Scheduler {
	return Scheduler{
		positions: make(map[int]*scheduledTask),
	}
}

func (s *Scheduler) AddTask(task Task) {
	if s.positions == nil {
		s.positions = make(map[int]*scheduledTask)
	}

	if s.ChangeTaskPriority(task.Identifier, task.Priority) {
		return
	}

	s.scheduleTask(task)
}

func (s *Scheduler) scheduleTask(task Task) {
	item := &scheduledTask{
		task: task,
	}

	heap.Push(&s.tasks, item)
	s.positions[task.Identifier] = item
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) bool {
	if item, ok := s.positions[taskID]; ok {
		item.task.Priority = newPriority
		heap.Fix(&s.tasks, item.index)
		return true
	}

	return false
}

func (s *Scheduler) GetTask() Task {
	if len(s.tasks) == 0 {
		return Task{}
	}

	item := heap.Pop(&s.tasks).(*scheduledTask)
	delete(s.positions, item.task.Identifier)

	return item.task
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, Task{Identifier: task1.Identifier, Priority: 100}, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}

func TestChangeUnknownTaskPriorityDoesNothing(t *testing.T) {
	task := Task{Identifier: 1, Priority: 10}

	scheduler := NewScheduler()
	scheduler.AddTask(task)
	scheduler.ChangeTaskPriority(2, 100)

	assert.Equal(t, task, scheduler.GetTask())
}

func TestGetTaskFromEmptyScheduler(t *testing.T) {
	scheduler := NewScheduler()

	assert.Equal(t, Task{}, scheduler.GetTask())
}

func TestChangeTaskPriorityDown(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 30}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 10}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.ChangeTaskPriority(1, 5)

	assert.Equal(t, task2, scheduler.GetTask())
	assert.Equal(t, task3, scheduler.GetTask())
	assert.Equal(t, Task{Identifier: task1.Identifier, Priority: 5}, scheduler.GetTask())
}

func TestEqualPriorityTasksAreReturnedByIdentifier(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 10}
	task3 := Task{Identifier: 3, Priority: 10}

	scheduler := NewScheduler()
	scheduler.AddTask(task3)
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)

	assert.Equal(t, task1, scheduler.GetTask())
	assert.Equal(t, task2, scheduler.GetTask())
	assert.Equal(t, task3, scheduler.GetTask())
}
