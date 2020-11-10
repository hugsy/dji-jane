package djijane

import (
	"dji-joe"
	
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apsdehal/go-logger"
)


const PROGNAME string = "DJI-Jane"
const VERSION string = "0.1"

var Log *logger.Logger
var probes djijoe.Probes

func StartServer(listenAddress string){

	Log.Debug("Preparing server")
	s := &http.Server{
		Addr:           listenAddress,
		Handler:        NewRouter(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	Log.Debug("Placing sighandlers")
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func(){
		sig := <-sigc
		Log.InfoF("Got '%+v': stopping cleanly", sig)
		s.Close() // <-- this is why Go 1.8+ is required
	}()

	Log.InfoF("Starting REST server listening on '%s'", listenAddress)
	s.ListenAndServe()
}
