package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"userhub/internal/adapters/memory"
	"userhub/internal/app"
	"userhub/internal/ports/httpport"
	"userhub/internal/security"
)

func main(){
	

	secret := []byte(env("APP_SECRET", "dev-secret-samuel-go-userhub"))
	pepper := []byte(env("APP_PEPPER", "dev-pepper"))

	repo := memory.NewUserRepo()
	hasher := security.NewHasher(pepper)
	jwt := security.NewJWT(secret)

	svc := &app.Service{Repo: repo, Hash: hasher, Token: jwt}
	h := &httpport.Handlers{Svc: svc}
	router := httpport.NewRouter(h, jwt)

	srv := &http.Server{Addr: ":8080", Handler: router.Mux, ReadHeaderTimeout: 5 * time.Second}

	go func(){
		log.Println("listening on:", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed)	{
			log.Fatal(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}