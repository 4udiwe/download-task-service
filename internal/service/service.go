package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/4udiwe/download-task-service/internal/entity"
)

// TaskService - сервис управления задачами скачивания
type TaskService struct {
	storage TaskStorage
	queue   chan *entity.Task
	workers int

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// New создаёт сервис: workers воркеров, контекст parentCtx (можно nil)
func New(st TaskStorage, workers int, parentCtx context.Context) *TaskService {
	if parentCtx == nil {
		parentCtx = context.Background()
	}
	ctx, cancel := context.WithCancel(parentCtx)

	s := &TaskService{
		storage: st,
		queue:   make(chan *entity.Task, 100),
		workers: workers,
		ctx:     ctx,
		cancel:  cancel,
	}

	// стартуем воркеров и добавляем их в waitgroup
	for i := range workers {
		s.wg.Add(1)
		go func(id int) {
			defer s.wg.Done()
			s.worker(id)
		}(i)
	}

	// восстановление незавершённых задач (попытается поставить в очередь)
	s.resumePendingTasks()

	return s
}

func (s *TaskService) worker(id int) {
	for {
		select {
		case <-s.ctx.Done():
			// при отмене контекста выходим
			return
		case task := <-s.queue:
			if task == nil {
				// на случай, если кто-то закрыл канал
				return
			}
			s.processTask(task)
		}
	}
}

func (s *TaskService) processTask(task *entity.Task) {
	// пометим в "running" и сразу попытаемся сохранить
	task.Status = entity.TaskRunning
	if err := s.storage.Update(task); err != nil {
		fmt.Printf("warning: storage.Update (set running) task=%s: %v\n", task.ID, err)
	}

	for _, f := range task.Files {
		if f.Status == entity.FileDone {
			continue
		}

		f.Status = entity.FileRunning
		if err := s.storage.Update(task); err != nil {
			fmt.Printf("warning: storage.Update (file running) task=%s file=%s: %v\n", task.ID, f.URL, err)
		}

		// используем контекст сервиса, чтобы можно было отменять загрузки при Shutdown
		pathOnDisk, err := downloadFile(s.ctx, f.URL, task.ID)
		if err != nil {
			f.Status = entity.FileFailed
			f.Error = err.Error()
			fmt.Printf("download error task=%s url=%s: %v\n", task.ID, f.URL, err)
		} else {
			f.Status = entity.FileDone
			f.Path = pathOnDisk
		}

		if err := s.storage.Update(task); err != nil {
			fmt.Printf("warning: storage.Update (after file) task=%s: %v\n", task.ID, err)
		}
	}

	// проверка статуса всей задачи
	allDone := true
	for _, f := range task.Files {
		if f.Status != entity.FileDone {
			allDone = false
			break
		}
	}
	if allDone {
		task.Status = entity.TaskDone
	} else {
		// остаёмся в TaskRunning — при отмене останется in_progress/running
		task.Status = entity.TaskRunning
	}

	if err := s.storage.Update(task); err != nil {
		fmt.Printf("warning: storage.Update (final) task=%s: %v\n", task.ID, err)
	}
}

func (s *TaskService) resumePendingTasks() {
	tasks, err := s.storage.List()
	if err != nil {
		fmt.Printf("warning: storage.List on resume: %v\n", err)
		return
	}
	for _, t := range tasks {
		if t.Status == entity.TaskPending || t.Status == entity.TaskRunning {
			// пытаемся положить в очередь немедленно, иначе запускаем отложенную попытку
			select {
			case s.queue <- t:
			default:
				// очередь заполнена — ставим в goroutine (будет ждать места)
				go func(tt *entity.Task) { s.queue <- tt }(t)
			}
		}
	}
}

// CreateTask сохраняет задачу и пытается поставить её в очередь (но не блокирует)
func (s *TaskService) CreateTask(urls []string) (*entity.Task, error) {
	task := &entity.Task{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
		Status:    entity.TaskPending,
		Files:     make([]*entity.File, 0, len(urls)),
	}
	for _, u := range urls {
		task.Files = append(task.Files, &entity.File{
			URL:    u,
			Status: entity.FilePending,
		})
	}

	if err := s.storage.Save(task); err != nil {
		return nil, ErrCannotSaveTask
	}

	// non-blocking enqueue: если очередь занята — сохраняем задачу как pending,
	// она будет подобрана при рестарте или когда появится место.
	select {
	case s.queue <- task:
	default:
		// ничего делать не нужно — задача уже сохранена и будет обработана позже
	}

	return task, nil
}

func (s *TaskService) GetTask(id string) (*entity.Task, error) {
	return s.storage.Get(id)
}

func (s *TaskService) ListTasks() ([]*entity.Task, error) {
	return s.storage.List()
}

// Shutdown — отменяет контекст и ждёт, пока воркеры завершат текущие работы
func (s *TaskService) Shutdown() {
	// отменяем контекст — все downloadFile, сделанные с этим контекстом, должны прерваться
	s.cancel()
	// ждём, пока все воркеры закончат обработку
	s.wg.Wait()
}
