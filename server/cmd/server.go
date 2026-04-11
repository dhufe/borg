package main

import (
	"fmt"
	"lath/borg/internal"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	DefaultResponse = "Borg server version %s is running"
	FileStorePath   = "/borg/file-store"
)

var version = os.Getenv("BORG_VERSION")

type fileAnalysis struct {
	// Summary describes the overall verification result.
	Summary internal.Summary `json:"summary"`
	// Merged feature sets ...
	FeatureSets []internal.FeatureSet `json:"featureSets"`
	// ToolResults is a list of complete responses from all tools, mapped by
	// tool name.
	ToolResults []internal.ToolResult `json:"toolResults"`
	// DurationInMs represents the duration of the analysis in milliseconds.
	DurationInMs int64 `json:"durationInMs"`
}

func main() {
	log.Printf(DefaultResponse, version)
	initServer()
	router := gin.Default()
	router.MaxMultipartMemory = 3000 << 20 // 3 GiB
	// Allow cors to integrate Borg in other applications.
	router.ForwardedByClientIP = true
	if err := router.SetTrustedProxies([]string{"*"}); err != nil {
		log.Printf("Failed to set trusted proxies: %v", err)
		return
	}
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type"}
	corsConfig.AllowMethods = []string{"GET", "POST"}
	// It's important that the cors configuration is used before declaring the routes.
	router.Use(cors.New(corsConfig))
	router.GET("api", getDefaultResponse)
	router.GET("api/version", getVersion)
	router.POST("api/analyze", analyzeFile)
	err := router.Run()
	if err != nil {
		return
	}
}

func initServer() {
	internal.ParseConfig()
}

func getDefaultResponse(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf(DefaultResponse, version))
}

func getVersion(c *gin.Context) {
	c.String(http.StatusOK, version)
}

func analyzeFile(c *gin.Context) {
	start := time.Now()
	file, err := c.FormFile("file")
	// no file received
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "no file received",
		})
		return
	}
	// generate unique file name for storing
	filename := uuid.New().String() + "_" + file.Filename
	fileStorePath := filepath.Join(FileStorePath, filename)
	err = c.SaveUploadedFile(file, fileStorePath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "unable to save file",
		})
		return
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println(err)
		}
	}(fileStorePath)

	identResults := internal.RunIdentificationTools(filename)
	triggeredResults := internal.RunTriggeredTools(filename, identResults)
	toolResults := internal.CombineToolResults(identResults, triggeredResults)
	mergedSets := internal.MergeFeatureSets(toolResults)
	if len(mergedSets) == 0 {
		mergedSets = make([]internal.FeatureSet, 0)
	}
	tr := internal.GetSortedToolResults(identResults, triggeredResults)
	fileAnalysis := fileAnalysis{
		Summary:      internal.GetSummary(mergedSets, tr),
		FeatureSets:  mergedSets,
		ToolResults:  tr,
		DurationInMs: time.Since(start).Milliseconds(),
	}
	c.JSON(http.StatusOK, fileAnalysis)
}
