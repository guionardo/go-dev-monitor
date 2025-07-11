package server

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/guionardo/go-dev-monitor/internal/api"
	"github.com/guionardo/go-dev-monitor/internal/config"
	"github.com/guionardo/go-dev-monitor/internal/logging"
	"github.com/guionardo/go-dev-monitor/internal/store"
)

type Server struct {
	config *config.Server
	store  *store.DataStore
}

var validate = validator.New(validator.WithRequiredStructEnabled())

const (
	AppDataHeader  = "X-App-Data"
	AppDataValue   = "go-dev-monitor"
	HostNameHeader = "X-App-Hostname"
)

func New() (*Server, error) {
	centralCfg := config.NewConfig()
	cfg := centralCfg.Server
	if cfg == nil {
		cfg = &config.Server{}
	}
	cfg.Reset()

	store, err := store.New(cfg.QueueSize, centralCfg.GetConfigDir())
	if err != nil {
		return nil, err
	}
	return &Server{
		config: cfg,
		store:  store,
	}, nil
}

func (s *Server) Run() {
	if !logging.IsDebug() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(cors.Default())
	logging.SetupGin(router)

	setupVueStatic(router)

	router.GET("/hc", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	router.POST("/data", s.PostData)
	router.GET("/data", s.GetData)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: router.Handler(),
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		// service connections
		logging.Info("Listening", slog.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("Failed to start listening", err, slog.String("address", srv.Addr))
			cancel()
		}
	}()

	defer cancel()
	<-ctx.Done()
	logging.Info("Finishing service...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelShutdown()

	err := srv.Shutdown(shutdownCtx)
	logging.Info("Finished.", slog.Any("error", err))
}

func (s *Server) GetData(c *gin.Context) {
	summary, err := s.store.Get()

	if err != nil {
		c.JSON(http.StatusBadGateway, api.ServerResponse{Message: "bad gateway", Error: err.Error()})
		return
	}
	var response = api.SummaryResponse{
		Origins: make(map[string][]api.LocalRepositoryResponse, len(summary)),
	}
	queryOrigin := c.Query("origin")

	for origin, localRepos := range summary {
		if len(queryOrigin) > 0 && queryOrigin != origin {
			continue
		}
		for _, localRepo := range localRepos {
			response.Origins[origin] = append(response.Origins[origin], api.ToLocalRepositoryResponse(localRepo))
		}
	}
	c.JSON(http.StatusOK, response)
}

func (s *Server) PostData(c *gin.Context) {
	if c.Request.Header.Get(AppDataHeader) != AppDataValue {
		c.JSON(http.StatusBadRequest, api.ServerResponse{Message: "bad request"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err == nil {
		requestData := &api.AgentRequest{}
		if err = requestData.UnmarshalJSON(body); err == nil {
			if err = validate.Struct(requestData); err == nil {
				s.store.BeginPosts(requestData.Hostname)
				for _, repo := range requestData.Repositories {
					if err = s.store.Post(requestData.Hostname, repo); err != nil {
						logging.Warn("Saving", slog.String("origin", repo.Origin), slog.Any("error", err))
						break
					}
					logging.Info("Saving", slog.String("origin", repo.Origin), slog.String("hostname", requestData.Hostname))
				}
				logging.Info("Saved", slog.Int("repos", len(requestData.Repositories)), slog.String("hostname", requestData.Hostname))
			}
		}
	}
	if err == nil {
		c.JSON(http.StatusAccepted, api.ServerResponse{Message: "accepted"})
	} else {
		c.JSON(http.StatusBadRequest,
			api.ServerResponse{Message: "bad request", Error: err.Error()})
	}
}
