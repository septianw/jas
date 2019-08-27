package main

/**
This app will contain multilevel bootstrap event
*/
import (
	"context"
	"fmt"
	"log"
	"net/http"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var Modloc string = ""
var Libloc string = ""
var Routers *gin.Engine

func main() {
	BootstrapAll()

	// r := SetupRouter()
	srv := &http.Server{
		Addr:    ListenAddr,
		Handler: Routers,
	}
	fmt.Printf("Listening at: %s", ListenAddr)

	// srv.ListenAndServe()

	// gracefull shutdown procedure

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("wow")
	log.Println("Shutdown Server ...")
	os.Remove("/tmp/shinyRuntimeFile")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
