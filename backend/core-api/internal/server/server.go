package server

import (
	"log"
	"net/http"

	"github.com/Poted/raitometer/backend/core-api/internal/database"
	"github.com/Poted/raitometer/backend/core-api/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	db     *sqlx.DB
	router *chi.Mux
}

func New(db *sqlx.DB) *Server {
	s := &Server{
		db:     db,
		router: chi.NewRouter(),
	}

	projectStore := database.NewPostgresProjectStore(db)
	userStore := database.NewPostgresUserStore(db)
	calculatorStore := database.NewPostgresCalculatorStore(db)

	h := handlers.New(db, projectStore, userStore, calculatorStore)

	s.registerRoutes(h)

	return s
}

func (s *Server) registerRoutes(h *handlers.Handlers) {
	s.router.Get("/healthcheck", h.HealthCheckHandler)
	s.router.Mount("/users", s.userRoutes(h))

	s.router.Group(func(r chi.Router) {
		r.Use(h.AuthMiddleware)

		r.Mount("/projects", s.projectRoutes(h))
		r.Mount("/calculators", s.calculatorRoutes(h))
		r.Mount("/modules", s.moduleRoutes(h))

	})
}

func (s *Server) userRoutes(h *handlers.Handlers) http.Handler {
	r := chi.NewRouter()
	r.Post("/register", h.RegisterUserHandler)
	r.Post("/login", h.LoginUserHandler)
	return r
}

func (s *Server) projectRoutes(h *handlers.Handlers) http.Handler {
	r := chi.NewRouter()

	r.Get("/", h.GetAllProjectsHandler)
	r.Post("/", h.CreateProjectHandler)
	r.Get("/{projectID}", h.GetProjectHandler)
	r.Put("/{projectID}", h.UpdateProjectHandler)
	r.Delete("/{projectID}", h.DeleteProjectHandler)

	r.Post("/{projectID}/calculator", h.CreateCalculatorHandler)

	return r
}

func (s *Server) calculatorRoutes(h *handlers.Handlers) http.Handler {
	r := chi.NewRouter()

	r.Post("/{calculatorID}/modules", h.CreateModuleHandler)
	r.Get("/{calculatorID}", h.GetFullCalculatorHandler)

	return r
}

func (s *Server) moduleRoutes(h *handlers.Handlers) http.Handler {
	r := chi.NewRouter()

	r.Post("/{moduleID}/items", h.CreateItemHandler)

	return r
}

func (s *Server) Start(port string) error {
	log.Printf("server starting on port %s", port)

	return http.ListenAndServe(port, s.router)
}
