package api

import (
	"bytes"
	"fmt"
	"github.com/go-pg/pg/orm"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"github.com/plusspeed/payments-api/internal/model"
	"github.com/plusspeed/payments-api/internal/repository"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	// values from the docker compose
	pgUsername = "test"
	pgPassword = "example"
	pgAddress  = "127.0.0.1:5432"
	basePath   = "/v1/payment"
)

func TestAllPaymentCall(t *testing.T) {
	dbTest := repository.New(pgAddress, "", pgUsername, pgPassword)
	defer dbTest.Database.Close()
	router := NewRouter("/v1", dbTest)
	clearDB(*dbTest)

	var paymentID = uuid.NewRandom().String()

	req1, _ := http.NewRequest("POST", basePath, bytes.NewBuffer(createRequest(paymentID)))
	response := executeRequest(*router, req1)

	checkResponseCode(t, http.StatusCreated, response.Code)

	req2, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", basePath, paymentID), bytes.NewBuffer(createRequest(paymentID)))
	response2 := executeRequest(*router, req2)

	checkResponseCode(t, http.StatusOK, response2.Code)

	req3, _ := http.NewRequest("PUT", fmt.Sprintf("%s/%s", basePath, paymentID), bytes.NewBuffer(createRequest(paymentID)))
	response3 := executeRequest(*router, req3)

	checkResponseCode(t, http.StatusNoContent, response3.Code)

	req5, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s?limit=100&offset=1", basePath, paymentID), bytes.NewBuffer(createRequest(paymentID)))
	response5 := executeRequest(*router, req5)

	checkResponseCode(t, http.StatusOK, response5.Code)

	req4, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", basePath, paymentID), bytes.NewBuffer(createRequest(paymentID)))
	response4 := executeRequest(*router, req4)

	checkResponseCode(t, http.StatusNoContent, response4.Code)

}

func executeRequest(router mux.Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func createRequest(paymentID string) []byte {
	return []byte("{\"type\": \"Payment\"," +
		"\"id\": \"" + paymentID + "\"," +
		"\"version\": 0," +
		"\"organisation_id\": \"743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb\"," +
		"\"attributes\": {\"amount\": \"100.21\",\"beneficiary_party\": {\"account_name\": \"W Owens\",\"account_number\": \"31926819\",\"account_number_code\": \"BBAN\",\"account_type\": 0,\"address\": \"1 The Beneficiary Localtown SE2\",\"bank_id\": \"403000\",\"bank_id_code\": \"GBDSC\",\"name\": \"Wilfred Jeremiah Owens\"},\"charges_information\": {\"bearer_code\": \"SHAR\",\"sender_charges\": [{\"amount\": \"5.00\",\"currency\": \"GBP\"},{\"amount\": \"10.00\",\"currency\": \"USD\"}],\"receiver_charges_amount\": \"1.00\",\"receiver_charges_currency\": \"USD\"},\"currency\": \"GBP\",\"debtor_party\": {\"account_name\": \"EJ Brown Black\",\"account_number\": \"GB29XABC10161234567801\",\"account_number_code\": \"IBAN\",\"address\": \"10 Debtor Crescent Sourcetown NE1\",\"bank_id\": \"203301\",\"bank_id_code\": \"GBDSC\",\"name\": \"Emelia Jane Brown\"},\"end_to_end_reference\": \"Wil piano Jan\",\"fx\": {\"contract_reference\": \"FX123\",\"exchange_rate\": \"2.00000\",\"original_amount\": \"200.42\",\"original_currency\": \"USD\"},\"numeric_reference\": \"1002001\",\"payment_id\": \"123456789012345678\",\"payment_purpose\": \"Paying for goods/services\",\"payment_scheme\": \"FPS\",\"payment_type\": \"Credit\",\"processing_date\": \"2017-01-18\",\"reference\": \"Payment for Em's piano lessons\",\"scheme_payment_sub_type\": \"InternetBanking\",\"scheme_payment_type\": \"ImmediatePayment\",\"sponsor_party\": {\"account_number\": \"56781234\",\"bank_id\": \"123123\",\"bank_id_code\": \"GBDSC\"}}}")
}

func clearDB(dbTest repository.Repository) {
	dbTest.Database.Delete(&model.Payment{})
	err := dbTest.Database.CreateTable(&model.Payment{}, &orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		panic(err.Error())
	}
}
