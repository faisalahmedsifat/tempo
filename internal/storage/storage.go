package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"tempo/internal/types"
)

type Storage struct {
	filepath string
	mutex    sync.RWMutex
	jobs     map[string]types.Job
}

func NewStorage(dataDir string) (*Storage, error) {
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %v", err)
		}
		dataDir = filepath.Join(homeDir, ".tempo")
	}

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}

	filepath := filepath.Join(dataDir, "jobs.json")
	storage := &Storage{
		filepath: filepath,
		jobs:     make(map[string]types.Job),
	}

	// Load existing jobs
	if err := storage.load(); err != nil {
		return nil, fmt.Errorf("failed to load jobs: %v", err)
	}

	return storage, nil
}

func (s *Storage) AddJob(job types.Job) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.jobs[job.ID] = job
	return s.save()
}

func (s *Storage) GetJob(id string) (types.Job, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	job, exists := s.jobs[id]
	return job, exists
}

func (s *Storage) GetAllJobs() []types.Job {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	jobs := make([]types.Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

func (s *Storage) RemoveJob(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.jobs[id]; !exists {
		return fmt.Errorf("job '%s' not found", id)
	}

	delete(s.jobs, id)
	return s.save()
}

func (s *Storage) RemoveAllJobs() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.jobs = make(map[string]types.Job)
	return s.save()
}

func (s *Storage) load() error {
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, that's okay
			return nil
		}
		return err
	}

	var jobs []types.Job
	if err := json.Unmarshal(data, &jobs); err != nil {
		return fmt.Errorf("failed to unmarshal jobs: %v", err)
	}

	// Convert slice to map
	s.jobs = make(map[string]types.Job)
	for _, job := range jobs {
		s.jobs[job.ID] = job
	}

	return nil
}

func (s *Storage) save() error {
	// Convert map to slice
	jobs := make([]types.Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}

	data, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal jobs: %v", err)
	}

	if err := os.WriteFile(s.filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write jobs file: %v", err)
	}

	return nil
}
