package server

import (
	"abuse/db/postgres"
	"abuse/db/redis"
	"abuse/pkg/mail"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run(envType *string) {
	// Validate envType
	if envType == nil || *envType == "" {
		logrus.Fatal("Environment type must be provided")
	}

	// Initialize logger first
	InitLogger()

	logrus.Infof("Running in %s environment", *envType)
	viper.SetConfigType("json")

	// Choose config name and log level based on environment
	var configName string
	var level logrus.Level

	switch *envType {
	case "production":
		configName = "config_production"
		level = logrus.InfoLevel
	case "staging":
		configName = "config_staging"
		level = logrus.InfoLevel
	case "dev":
		configName = "config_dev"
		level = logrus.DebugLevel
	case "local":
		configName = "config_local"
		level = logrus.DebugLevel
	default:
		logrus.Fatalf("Unknown environment type: %s", *envType)
	}

	viper.SetConfigName(configName)
	logrus.SetLevel(level)

	// Resolve $HOME properly
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logrus.Fatalf("Unable to resolve home directory: %v", err)
	}
	viper.AddConfigPath(filepath.Join(homeDir, ".sck"))

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Error reading config file: %v", err)
	}

	logrus.Infof("Using config file: %s", viper.ConfigFileUsed())
	// Initialize Postgres connection
	postgres := postgres.NewPostgres()
	redis := redis.NewRedis(envType)

	server := &Server{
		router:   mux.NewRouter(),
		postgres: postgres,
		user:     postgres,
		mail: mail.NewMail(
			viper.GetString("mail_id"),
			viper.GetString("mail_pass"),
			viper.GetString("app_pass"),
		),
		sessmanager: redis,
	}
	server.RegisterRoutes()

	logrus.Info("Server is starting...")
	if *envType == "staging" || *envType == "prod" {
		slog.Info("API server running on port 443")
		log.Fatal(
			http.ListenAndServeTLS(
				":443",
				viper.GetString("fullChainPath"),
				viper.GetString("privKeyPath"),
				server,
			),
		)

	} else {
		slog.Info("API server running on port 8194")
		log.Fatal(
			http.ListenAndServe(
				":8194",
				server,
			),
		)

	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// InitLogger initializes the logger with optional LOG_LEVEL override
func InitLogger() {
	logrus.SetOutput(os.Stdout)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Set initial log level based on LOG_LEVEL environment variable
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "info", "":
		logrus.SetLevel(logrus.InfoLevel)
	default:
		logrus.Warnf("Unknown LOG_LEVEL: %s, defaulting to info", logLevel)
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func (s *Server) respond(
	w http.ResponseWriter,
	data interface{},
	status int,
	err error,
) {
	// Set content type header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err == nil {
		resp := &ResponseMsg{
			Message: "success",
			Data:    data,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logrus.Error("Error in encoding the response", "error", err)
			return
		}
		return
	}
	resp := &ResponseMsg{
		Message: err.Error(),
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logrus.Error("Error in encoding the error response", "error", err)
		return
	}
}
