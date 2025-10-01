package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/4udiwe/download-task-service/internal/entity"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	mu    sync.RWMutex
	path  string
	tasks map[string]*entity.Task
}

func New(path string) *Storage {
	storage := &Storage{
		path:  path,
		tasks: make(map[string]*entity.Task),
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		logrus.Fatalf("cannot create storage directory %s: %v", dir, err)
	}

	storage.load()

	return storage
}

func (s *Storage) persist() error {
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}

	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	return enc.Encode(s.tasks)
}

func (s *Storage) load() {
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		logrus.Infof("Storage file not found, creating new: %s", s.path)
		if err := s.persist(); err != nil {
			logrus.Fatalf("cannot create new storage file %s: %v", s.path, err)
		}
		return
	}

	f, err := os.Open(s.path)
	if err != nil {
		logrus.Fatalf("Error opening storage file: %v", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(&s.tasks); err != nil {
		logrus.Warnf("Error decoding storage file, starting with empty storage: %v", err)
		s.tasks = make(map[string]*entity.Task)
	}
}

func (s *Storage) Save(task *entity.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks[task.ID] = task

	return s.persist()
}

func (s *Storage) Get(id string) (*entity.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]

	if !ok {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

func (s *Storage) List() ([]*entity.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]*entity.Task, 0, len(s.tasks))

	for _, t := range s.tasks {
		list = append(list, t)
	}

	return list, nil
}

func (s *Storage) Update(task *entity.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[task.ID]; !ok {
		return ErrTaskNotFound
	}

	s.tasks[task.ID] = task

	return s.persist()
}
