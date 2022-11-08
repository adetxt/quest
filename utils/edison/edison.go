package edison

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	RestContext struct {
		EchoContext echo.Context
	}

	RestHandler func(ctx context.Context, clientCtx RestContext) error
)

type Edison struct {
	ec *echo.Echo
}

func NewEdison() *Edison {
	return &Edison{
		ec: echo.New(),
	}
}

func (ed *Edison) RestRouter(method, path string, h RestHandler) {
	ed.ec.Router().Add(method, path, func(c echo.Context) error {
		return h(context.Background(), RestContext{
			EchoContext: c,
		})
	})
}

func (ed *Edison) StartRestServer(port string) {
	ed.RestRouter("GET", "/__health", func(ctx context.Context, c RestContext) error {
		return c.EchoContext.String(http.StatusOK, "ok")
	})

	// Start server
	go func() {
		if err := ed.ec.Start(fmt.Sprintf(":%s", port)); err != nil && err != http.ErrServerClosed {
			ed.ec.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := ed.ec.Shutdown(ctx); err != nil {
		ed.ec.Logger.Fatal(err)
	}
}

func (c *RestContext) Bind(i interface{}) error {
	return c.EchoContext.Bind(i)
}

func (c *RestContext) JSON(code int, i interface{}, message string) error {
	isOK := code < 400

	res := map[string]interface{}{}

	if !isOK {
		res["status"] = "error"
		res["error"] = strings.ToUpper(http.StatusText(code))
		res["message"] = message
	} else {
		res["status"] = "success"
		res["message"] = strings.ToUpper(http.StatusText(code))
		res["data"] = i
	}

	return c.EchoContext.JSON(code, res)
}

func (c *RestContext) Success(i interface{}) error {
	return c.JSON(200, i, "")
}

func (c *RestContext) Error(e error, code int) error {
	return c.JSON(code, nil, e.Error())
}

func (c *RestContext) ErrorWithCustomMessage(code int, message string) error {
	return c.JSON(code, nil, message)
}
