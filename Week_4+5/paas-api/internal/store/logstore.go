package store

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dreule28/Week_4/paas-api/internal/model"
)

const DefaultLogStorePath = "data/logs.jsonl"

type LogStore interface {
	Append(entry model.LogEntry) (model.LogEntry, error)
	List(q model.LogQuery) ([]model.LogEntry, error)
}

type FileLogStore struct {
	path    string
	mu      sync.RWMutex
	entries []model.LogEntry
	seq     uint64
}

func NewFileLogStore(path string) (*FileLogStore, error) {
	if strings.TrimSpace(path) == "" {
		path = DefaultLogStorePath
	}

	s := &FileLogStore{path: path}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *FileLogStore) load() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	f, err := os.OpenFile(s.path, os.O_RDONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var entry model.LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		s.entries = append(s.entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	s.seq = uint64(len(s.entries))
	return nil
}

func (s *FileLogStore) Append(entry model.LogEntry) (model.LogEntry, error) {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	if entry.ID == "" {
		n := atomic.AddUint64(&s.seq, 1)
		entry.ID = entry.Timestamp.Format("20060102T150405.000000000Z07:00") + "-" + strings.ReplaceAll(strings.ToLower(entry.Type), " ", "_") + "-" + itoa(n)
	}

	payload, err := json.Marshal(entry)
	if err != nil {
		return model.LogEntry{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return model.LogEntry{}, err
	}
	f, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return model.LogEntry{}, err
	}
	defer f.Close()

	if _, err := f.Write(append(payload, '\n')); err != nil {
		return model.LogEntry{}, err
	}

	s.entries = append(s.entries, entry)
	return entry, nil
}

func (s *FileLogStore) List(q model.LogQuery) ([]model.LogEntry, error) {
	limit := q.Limit
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}

	instanceID := strings.TrimSpace(q.InstanceID)
	logType := strings.ToLower(strings.TrimSpace(q.Type))

	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]model.LogEntry, 0, min(limit, len(s.entries)))
	for i := len(s.entries) - 1; i >= 0 && len(out) < limit; i-- {
		entry := s.entries[i]
		if instanceID != "" && entry.InstanceID != instanceID {
			continue
		}
		if logType != "" && strings.ToLower(entry.Type) != logType {
			continue
		}
		out = append(out, entry)
	}

	return out, nil
}

func itoa(v uint64) string {
	return strconv.FormatUint(v, 10)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
