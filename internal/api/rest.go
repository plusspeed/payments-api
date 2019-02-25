package api

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/plusspeed/payments-api/internal/model"
	"github.com/plusspeed/payments-api/internal/repository"
	"gopkg.in/go-playground/validator.v8"
	"net/http"
	"strconv"
	"strings"
)

//NewRouter starts the service. In the case of a service failure, it will PANIC.
func NewRouter(basePath string, db *repository.Repository) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/health", HealthCheckHandler(*db))

	r.HandleFunc(basePath+"/payment", CreatePayment(*db)).Methods("POST")
	r.HandleFunc(basePath+"/payment/{paymentID}", GetPayment(*db)).Methods("GET")
	r.HandleFunc(basePath+"/payment/{paymentID}", WithPaymentCtx(*db, DeletePayment)).Methods("DELETE")
	r.HandleFunc(basePath+"/payment/{paymentID}", WithPaymentCtx(*db, UpdatePayment)).Methods("PUT")
	r.HandleFunc(basePath+"/payments", GetAllPayments(*db)).
		Queries("offset", "{offset}", "limit", "{limit}").
		Methods("GET")

	http.Handle("/", r)

	return r
}

//CreatePayment creates a new payment transaction resource
func CreatePayment(repo repository.Repository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var t *model.Payment
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		err := validate(t)
		if err != nil {
			SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		dup, err := repo.Get(t.ID)
		if err == nil {
			if cmp.Equal(*t, *dup) {
				w.WriteHeader(http.StatusCreated)
				return
			}
			SendErrorResponse(w, r, http.StatusConflict, errors.New("already exists"))
			return
		}

		err = repo.Create(t)
		if err != nil {
			SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
}

// WithPaymentCtx checks if there is a transaction with the paymentID.
// If the payment exists, calls next handlerFunc.
// else returns an error message
func WithPaymentCtx(repo repository.Repository, next func(repository.Repository) http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paymentID := mux.Vars(r)["paymentID"]
		if paymentID == "" {
			SendErrorResponse(w, r, http.StatusNotFound, errors.Errorf("bad request paymentID:%s", paymentID))
			return
		}
		_, err := repo.Get(paymentID)
		if err != nil {
			if err == repository.ErrNotFound {
				SendErrorResponse(w, r, http.StatusNotFound, errors.Errorf("paymentID:%s not found", paymentID))
				return
			}
			SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		next(repo).ServeHTTP(w, r)
	})
}

//GetPayment returns the payment if exist.
func GetPayment(repo repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paymentID := mux.Vars(r)["paymentID"]
		val, err := repo.Get(paymentID)
		if err != nil {
			if err == repository.ErrNotFound {
				SendErrorResponse(w, r, http.StatusNotFound, errors.Errorf("paymentID:%s not found", paymentID))
				return
			}
			SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		SendResponse(w, r, http.StatusOK, val)
		return
	}
}

//DeletePayment deletes a resource payment if exist.
func DeletePayment(repo repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paymentID := mux.Vars(r)["paymentID"]

		err := repo.Delete(paymentID)
		if err != nil && err == repository.ErrNotFound {
			SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

//UpdatePayment updates a previous transaction.
func UpdatePayment(repo repository.Repository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paymentID := mux.Vars(r)["paymentID"]

		var t *model.Payment
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}
		err := validate(t)
		if err != nil {
			SendErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}
		//to ensure that the users does not try to modify a different payment
		t.ID = paymentID
		err = repo.Update(t)
		if err != nil {
			SendErrorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	})
}

//GetAllPayments Returns all the payments.
//Query params are optional.
//The default limit is 100 and max is 100000. Offset default value is 0.
func GetAllPayments(repo repository.Repository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var offset, limit int
		var err error
		if offset, err = strconv.Atoi(mux.Vars(r)["offset"]); err != nil {
			offset = 0
		}
		if limit, err = strconv.Atoi(mux.Vars(r)["limit"]); err != nil {
			limit = 100
		}

		if offset <= 0 {
			offset = 0
		}
		if limit >= 0 || 100000 <= limit {
			limit = 100
		}

		list, err := repo.List(offset, limit)
		if err != nil {
			if err == repository.ErrNotFound {
				SendErrorResponse(w, r, http.StatusNotFound, errors.Errorf("no payments found"))
				return
			}
			SendErrorResponse(w, r, http.StatusNotFound, err)
			return
		}
		SendResponse(w, r, http.StatusOK, list)
		return

	})
}

func validate(t *model.Payment) error {
	config := &validator.Config{TagName: "validate"}
	validate := validator.New(config)

	err := validate.Struct(t)
	if err != nil {
		return err
	}
	if !strings.EqualFold(t.Type, "Payment") {
		return errors.New("type is not Payment")
	}
	return nil
}
