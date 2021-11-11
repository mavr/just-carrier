package http_server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var startAt time.Time

const versionAPI = "0.1"

func init() {
	startAt = time.Now()
}

// APIServiceConfig configuration struct for api service
type APIServiceConfig struct {
	AppRevision  string
	AppVersion   string
	AppDebugMode bool

	ServPort int
}
type apiService struct {
	log  *zap.Logger

	s *http.Server
}

// NewAPIService create new api service
func NewAPIService(c APIServiceConfig, log *zap.Logger) *apiService {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/health", healthCheckHandle())

	apiGR := r.Group("/api")
	apiGR.GET("/status", statusHandle(c.AppVersion, c.AppRevision))

	return &apiService{
		log: log,
		s: &http.Server{
			Handler: r,
			Addr:    fmt.Sprintf(":%d", c.ServPort),
		},
	}
}

func (s *apiService) Run(ctx context.Context) error {
	go func() {
		_ = s.run(ctx)
	}()

	s.log.Info("api server started")

	return nil
}

func (s *apiService) run(ctx context.Context) error {
	go func(ctx context.Context) {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		if err := s.s.Shutdown(ctxShutdown); err != nil {
			s.log.Error("api server shutdown", zap.Error(err))
		}
	}(ctx)

	if err := s.s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	s.log.Info("api server stopped")
	return nil
}

func healthCheckHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthCheck{Status: healthStatusOK})
	}
}

func statusHandle(version, revision string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, StatusResponse{
			Build:      revision,
			Version:    version,
			VersionAPI: versionAPI,
			Uptime:     time.Now().Sub(startAt).Round(time.Second).String(),
		})
	}
}
