package scheduler

import (
	"context"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
)

// Task represents a scheduled task.
type Task interface {
	// Name returns the task name.
	Name() string
	// Run executes the task.
	Run(ctx context.Context) error
}

// Scheduler manages and runs scheduled tasks.
type Scheduler struct {
	cron   *cron.Cron
	tasks  map[string]cron.EntryID
	mu     sync.RWMutex
	logger *log.Logger
}

// NewScheduler creates a new scheduler instance.
func NewScheduler(logger *log.Logger) *Scheduler {
	return &Scheduler{
		cron:   cron.New(cron.WithSeconds()), // Support seconds in cron expression
		tasks:  make(map[string]cron.EntryID),
		logger: logger,
	}
}

// AddTask schedules a task with a cron expression.
// Cron expression format: "second minute hour day month weekday"
// Examples:
//   - "0 * * * * *" - Every minute
//   - "0 */5 * * * *" - Every 5 minutes
//   - "0 0 * * * *" - Every hour
//   - "0 0 0 * * *" - Every day at midnight
//   - "0 0 9 * * MON-FRI" - Every weekday at 9 AM
func (s *Scheduler) AddTask(cronExpr string, task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entryID, err := s.cron.AddFunc(cronExpr, func() {
		ctx := context.Background()
		s.logger.Printf("Running scheduled task: %s", task.Name())

		if err := task.Run(ctx); err != nil {
			s.logger.Printf("Task %s failed: %v", task.Name(), err)
		} else {
			s.logger.Printf("Task %s completed successfully", task.Name())
		}
	})

	if err != nil {
		return err
	}

	s.tasks[task.Name()] = entryID
	s.logger.Printf("Scheduled task '%s' with cron expression: %s", task.Name(), cronExpr)
	return nil
}

// Start starts the scheduler.
func (s *Scheduler) Start() {
	s.cron.Start()
	s.logger.Printf("Scheduler started with %d tasks", len(s.tasks))
}

// Stop stops the scheduler.
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	s.logger.Printf("Scheduler stopped")
}

// RemoveTask removes a scheduled task by name.
func (s *Scheduler) RemoveTask(taskName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entryID, exists := s.tasks[taskName]
	if !exists {
		return nil
	}

	s.cron.Remove(entryID)
	delete(s.tasks, taskName)
	s.logger.Printf("Removed task: %s", taskName)
	return nil
}

// ListTasks returns a list of all scheduled task names.
func (s *Scheduler) ListTasks() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]string, 0, len(s.tasks))
	for name := range s.tasks {
		tasks = append(tasks, name)
	}
	return tasks
}
