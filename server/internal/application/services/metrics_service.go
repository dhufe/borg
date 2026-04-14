package services

import (
	"context"
	"lath/borg/internal/domain/model"
	"log"
	"os"
	"path/filepath"

	"lath/borg/internal/domain/interfaces"
	"lath/borg/internal/infrastructure/observability"
)

type MetricsService struct {
	repo            interfaces.JobRepository
	fileStoragePath string
}

func NewMetricsService(repo interfaces.JobRepository, fileStoragePath string) *MetricsService {
	return &MetricsService{
		repo:            repo,
		fileStoragePath: fileStoragePath,
	}
}

// UpdateMetrics aktualisiert alle Business-Metriken aus Repository und Filesystem
func (s *MetricsService) UpdateMetrics(ctx context.Context) error {
	// Tasks aus Repository
	if err := s.updateTaskMetrics(ctx); err != nil {
		log.Printf("Failed to update task metrics: %v", err)
		// Nicht returnen, damit andere Metriken trotzdem aktualisiert werden
	}

	// Datei-Storage Metriken
	if err := s.updateStorageMetrics(); err != nil {
		log.Printf("Failed to update storage metrics: %v", err)
	}

	return nil
}

func (s *MetricsService) updateTaskMetrics(ctx context.Context) error {
	// Gesamtanzahl
	totalCount, err := s.repo.Count(ctx)
	if err != nil {
		return err
	}
	metrics.JobsTotal.Set(float64(totalCount))

	// Alle Status auf einmal holen
	statusCounts, err := s.repo.CountByAllStatuses(ctx)
	if err != nil {
		return err
	}

	// Metriken für alle möglichen Status setzen
	for _, status := range model.AllJobStatuses() {
		count := statusCounts[status] // 0 wenn nicht vorhanden
		metrics.JobsByStatus.WithLabelValues(string(status)).Set(float64(count))
	}

	return nil
}

func (s *MetricsService) updateStorageMetrics() error {
	// Anzahl Dateien im Storage
	fileCount := 0
	totalSize := int64(0)

	err := filepath.Walk(s.fileStoragePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileCount++
			totalSize += info.Size()
		}
		return nil
	})
	if err != nil {
		return err
	}

	metrics.StorageFilesCount.Set(float64(fileCount))
	metrics.StorageSizeBytes.Set(float64(totalSize))

	return nil
}
