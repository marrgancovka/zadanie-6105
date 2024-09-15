package app

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	handlerBid "zadanie-6105/internal/pkg/bids/delivery/http"
	repoBid "zadanie-6105/internal/pkg/bids/repo"
	usecaseBid "zadanie-6105/internal/pkg/bids/usecase"
	"zadanie-6105/internal/pkg/middleware"
	handlerTender "zadanie-6105/internal/pkg/tenders/delivery/http"
	repoTender "zadanie-6105/internal/pkg/tenders/repo"
	usecaseTender "zadanie-6105/internal/pkg/tenders/usecase"
	repoUser "zadanie-6105/internal/pkg/users/repo"
)

type App struct {
	log *logrus.Logger
}

func NewApp(log *logrus.Logger) *App {
	return &App{log: log}
}

func (a *App) Start() error {
	_ = godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_CONN"))
	if err != nil {
		a.log.Error("failed to connect database ", err.Error())
	}
	if err = db.Ping(); err != nil {
		a.log.Error("failed to ping database ", err.Error())
	}
	defer db.Close()

	m, err := migrate.New("file://schema", os.Getenv("POSTGRES_CONN"))
	if err != nil {
		a.log.Error("failed to create migrate: ", err.Error())
	}
	if err == nil {
		if err = m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	}

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	srv := &http.Server{
		Addr:         os.Getenv("SERVER_ADDRESS"),
		Handler:      r,
		ReadTimeout:  50 * time.Second,
		WriteTimeout: 50 * time.Second,
	}

	r.HandleFunc("/ping", ping).Methods(http.MethodGet)

	uRepo := repoUser.NewRepository(db)
	md := middleware.NewMiddleware(uRepo)

	tRepo := repoTender.NewRepository(db)
	tUsecase := usecaseTender.NewUsecase(tRepo)
	tHandler := handlerTender.NewHandler(tUsecase)

	r.HandleFunc("/tenders", tHandler.GetTendersList).Methods(http.MethodGet)
	r.HandleFunc("/tenders/new", tHandler.CreateNewTender).Methods(http.MethodPost)
	r.Handle("/tenders/my", md.UserExistsMiddleware(http.HandlerFunc(tHandler.GetUserTenders))).Methods(http.MethodGet)
	r.Handle("/tenders/{tenderId}/status", md.UserExistsMiddleware(http.HandlerFunc(tHandler.GetTenderStatus))).Methods(http.MethodGet)
	r.Handle("/tenders/{tenderId}/status", md.UserExistsMiddleware(http.HandlerFunc(tHandler.EditTenderStatus))).Methods(http.MethodPut)
	r.Handle("/tenders/{tenderId}/edit", md.UserExistsMiddleware(http.HandlerFunc(tHandler.EditTender))).Methods(http.MethodPatch)

	bRepo := repoBid.NewRepository(db)
	bUsecase := usecaseBid.NewBidUsecase(bRepo)
	bHandler := handlerBid.NewHandler(bUsecase)

	r.HandleFunc("/bids/new", bHandler.CreateNewBid).Methods(http.MethodPost)
	r.Handle("/bids/my", md.UserExistsMiddleware(http.HandlerFunc(bHandler.GetUserBids))).Methods(http.MethodGet)
	r.Handle("/bids/{tenderId}/list", md.UserExistsMiddleware(http.HandlerFunc(bHandler.GetTenderBids))).Methods(http.MethodGet)
	r.Handle("/bids/{bidId}/status", md.UserExistsMiddleware(http.HandlerFunc(bHandler.GetBidStatus))).Methods(http.MethodGet)
	r.Handle("/bids/{bidId}/status", md.UserExistsMiddleware(http.HandlerFunc(bHandler.EditBidStatus))).Methods(http.MethodPut)
	r.Handle("/bids/{bidId}/edit", md.UserExistsMiddleware(http.HandlerFunc(bHandler.EditBid))).Methods(http.MethodPatch)
	r.Handle("/bids/{bidId}/submit_decision", md.UserExistsMiddleware(http.HandlerFunc(bHandler.SubmitDecision))).Methods(http.MethodPut)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		a.log.Info("Start server on ", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Error("Error in listen: ", err.Error())
		}
	}()

	sig := <-signalCh
	a.log.Info("Received signal: ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		a.log.Fatal("Server shutdown failed: ", err.Error())
	}
	return nil
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
