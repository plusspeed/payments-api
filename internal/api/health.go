package api

import (
	"github.com/plusspeed/payments-api/internal/repository"
	"io"
	"net/http"
)

//HealthCheckHandler returns 200 if the app is healthy and connected to the db.
//returns 500 if can't connect to the db.
func HealthCheckHandler(repo repository.Repository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if _, err := repo.Database.ExecOne("select 'It is running'"); err != nil {
			io.WriteString(w, `{"alive": false}`)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.WriteString(w, `{"alive": true}`)
		w.WriteHeader(http.StatusOK)
		return
	})
}
