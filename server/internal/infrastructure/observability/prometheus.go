package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Job-Metriken
	JobsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "borg_jobs_total",
			Help: "Total number of tasks in database",
		},
	)

	JobsByStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "borg_jobs_by_status",
			Help: "Number of tasks by status",
		},
		[]string{"status"},
	)

	// Storage-Metriken
	StorageFilesCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "borg_storage_files_count",
			Help: "Number of files in storage directory",
		},
	)

	StorageSizeBytes = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "borg_storage_size_bytes",
			Help: "Total size of files in storage directory in bytes",
		},
	)

	// HTTP-Metriken (falls noch nicht vorhanden)
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "borg_http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "borg_http_request_duration_seconds",
			Help:    "HTTP request duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)
