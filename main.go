package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jrcartee/router"
	"github.com/jrcartee/scheduling/database"
	"github.com/jrcartee/scheduling/httputil"

	"github.com/jrcartee/scheduling/domain/dependency"
	"github.com/jrcartee/scheduling/domain/schedule"
	"github.com/jrcartee/scheduling/domain/task"
)

type Adapter func(http.Handler) http.Handler

func makeHandler(routes *router.Router, mw []Adapter) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		h, err := routes.GetEndpoint(r)
		if err == router.ErrNoURLMatch {
			http.Error(rw, "Not Found", http.StatusNotFound)
			return
		} else if err == router.ErrNoMethodMatch {
			http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		} else if err != nil {
			httputil.InternalError(rw, err)
			return
		}

		for _, a := range mw {
			h = a(h)
		}

		h.ServeHTTP(rw, r)
	}
}

func main() {
	db, cleanup := database.Setup(database.DefaultConfig)
	defer cleanup()

	reqLog := log.New(os.Stdout, "", log.Ltime)
	middleware := []Adapter{
		RequestLogger(reqLog),
		PanicRecovery(),
	}

	routes := router.New()
	schedule.RegisterRoutes(routes, db)
	task.RegisterRoutes(routes, db)
	dependency.RegisterRoutes(routes, db)
	routes.Print()

	s := http.Server{
		Addr:    ":8080",
		Handler: makeHandler(routes, middleware),
	}
	s.ListenAndServe()
}
