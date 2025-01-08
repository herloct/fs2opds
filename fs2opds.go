package main

import (
	"context"
	"crypto/subtle"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/herloct/fs2opds/api/v1d2"
	"github.com/herloct/fs2opds/configs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	// Setup
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	catalogConfig := configs.GetCatalogConfig()
	if !catalogConfig.IsValid() {
		e.Logger.Panic("Invalid catalog config")
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	})

	groupOpds := e.Group("/opds")
	if catalogConfig.NeedsAuth() {
		groupOpds.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			// Be careful to use constant time comparison to prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(username), []byte(catalogConfig.Username)) == 1 &&
				subtle.ConstantTimeCompare([]byte(password), []byte(catalogConfig.Password)) == 1 {
				return true, nil
			}

			return false, nil
		}))
	}

	groupV1d2 := groupOpds.Group("/v1.2")
	v1d2.AppendRoute(groupV1d2)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
