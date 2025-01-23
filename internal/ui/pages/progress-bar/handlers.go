package progressbar

import (
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"github.com/a-h/templ"
	"github.com/gofs-cli/template/internal/ui"
	"github.com/gofs-cli/template/internal/ui/components/header"
)

type Job struct {
	mu       sync.Mutex
	progress int
}

// Global variable to store the current job state
var currentJob = &Job{}

func Index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(ui.IndexPage(layout(header.Header(), body()))).ServeHTTP(w, r)
	})
}

func StartProgressBar() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentJob.mu.Lock()
		currentJob.progress = 0
		currentJob.mu.Unlock()

		templ.Handler(startProgressBar()).ServeHTTP(w, r)
	})
}

func RunProgressBar() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentJob.mu.Lock()
		defer currentJob.mu.Unlock()

		if currentJob.progress < 100 {
			currentJob.progress += rand.Intn(16)
		}

		if currentJob.progress >= 100 {
			w.Header().Set("HX-Trigger", "done")
			templ.Handler(inProgress(strconv.Itoa(currentJob.progress))).ServeHTTP(w, r)
		} else {
			templ.Handler(inProgress(strconv.Itoa(currentJob.progress))).ServeHTTP(w, r)
		}
	})
}

func CompleteProgressBar() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(completeProgressBar()).ServeHTTP(w, r)
	})
}
