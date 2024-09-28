package server

import (
    "context"
    "fmt"
    "net/http"
    _ "os"

    "github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
    "github.com/sirupsen/logrus"

    "github.com/adityaladwa/todo-app/internal/handlers"
    "github.com/adityaladwa/todo-app/internal/db"
    "github.com/adityaladwa/todo-app/pkg/logger"
    "github.com/spf13/viper"
)

type Server struct {
    Config   *Config
    Router   *chi.Mux
    Logger   *logrus.Logger
    Queries  *db.Queries
    DBPool   *pgx.Conn
}

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
    Port int `mapstructure:"port"`
}

// DatabaseConfig holds database-related configuration.
type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    User     string `mapstructure:"user"`
    Password string `mapstructure:"password"`
    DBName   string `mapstructure:"dbname"`
    SSLMode  string `mapstructure:"sslmode"`
}

// LogConfig holds logging-related configuration.
type LogConfig struct {
    Level string `mapstructure:"level"`
}

// LoadConfig reads configuration from the specified path using Viper.
func LoadConfig(path string) (*Config, error) {
    viper.SetConfigFile(path)
    viper.SetConfigType("yaml")

    // Set default values
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("database.sslmode", "disable")
    viper.SetDefault("log.level", "info")

    // Read the configuration file
    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("error reading config file: %w", err)
    }

    // Unmarshal the configuration into the Config struct
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("unable to decode config into struct: %w", err)
    }

    return &config, nil
}

func NewServer(configPath string) (*Server, error) {
    // Load configuration
    config, err := LoadConfig(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }

    // Initialize logger
    log := logger.NewLogger(config.Log.Level)

    // Connect to the database using pgxpool
    dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
        config.Database.User,
        config.Database.Password,
        config.Database.Host,
        config.Database.Port,
        config.Database.DBName,
        config.Database.SSLMode,
    )

    pool, err := pgx.Connect(context.Background(), dsn)
    if err != nil {
        log.Fatal("Cannot connect to database", err)
    }

    queries := db.New(pool)

    // Initialize router
    router := chi.NewRouter()

    server := &Server{
        Config:  config,
        Router:  router,
        Logger:  log,
        Queries: queries,
        DBPool:  pool,
    }

    server.setupRoutes()

    return server, nil
}



func (s *Server) setupRoutes() {
    s.Router.Route("/todos", func(r chi.Router) {
        handler := handlers.TodoHandler{
            Queries: s.Queries,
            Logger:  s.Logger,
        }

        r.Post("/", handler.CreateTodo)
        r.Get("/", handler.ListTodos)
        r.Get("/{id}", handler.GetTodo)
        r.Put("/{id}", handler.UpdateTodo)
        r.Delete("/{id}", handler.DeleteTodo)
    })
}

func (s *Server) Start() {
    addr := fmt.Sprintf(":%d", s.Config.Server.Port)
    s.Logger.Infof("Starting server on %s", addr)
    if err := http.ListenAndServe(addr, s.Router); err != nil {
        s.Logger.Fatal("Server failed", err)
    }
}

func (s *Server) Shutdown() {
    err := s.DBPool.Close(context.Background())
    if err != nil {
        s.Logger.Fatal(err) 
    }
    // Implement other graceful shutdown steps if needed
}
