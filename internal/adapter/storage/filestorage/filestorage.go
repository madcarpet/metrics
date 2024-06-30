package filestorage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
)

type FileStorage struct {
	*memstorage.MemStorage
	file     *os.File
	writer   *bufio.Writer
	scanner  *bufio.Scanner
	syncmode bool
	mutex    sync.Mutex
}

func NewFileStorage(fileName string, syncmode bool) (*FileStorage, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	memStorage := memstorage.NewMemStorage()
	return &FileStorage{
		MemStorage: memStorage,
		file:       file,
		writer:     bufio.NewWriter(file),
		scanner:    bufio.NewScanner(file),
		syncmode:   syncmode,
		mutex:      sync.Mutex{},
	}, nil
}

func (fs *FileStorage) UpdateMetric(ctx context.Context, m entity.Metric) error {
	err := fs.MemStorage.UpdateMetric(ctx, m)
	if err != nil {
		return err
	}
	if fs.syncmode {
		err := fs.ExportToFile()
		if err != nil {
			return err
		}
	}
	return nil
}

func (fs *FileStorage) ExportToFile() error {
	metrics, err := fs.GetAllMetrics(context.Background())
	if err != nil {
		return err
	}
	for _, m := range metrics {
		data, err := json.Marshal(m)
		if err != nil {
			return err
		}
		if _, err := fs.writer.Write(data); err != nil {
			return err
		}
		if err := fs.writer.WriteByte('\n'); err != nil {
			return err
		}
	}
	fs.mutex.Lock()
	fs.file.Truncate(0)
	fs.file.Seek(0, 0)
	fs.writer.Flush()
	fs.mutex.Unlock()
	return nil
}

func (fs *FileStorage) ImportFromFile() error {
	for fs.scanner.Scan() {
		data := fs.scanner.Bytes()
		metric := entity.Metric{}
		err := json.Unmarshal(data, &metric)
		if err != nil {
			return err
		}
		err = fs.UpdateMetric(context.Background(), metric)
		if err != nil {
			return err
		}
	}
	if err := fs.scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage) Close() error {
	err := fs.file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage) IsConnected(ctx context.Context) error {
	return fmt.Errorf("DB type without connection support")
}
