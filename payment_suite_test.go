package main_test

import (
	"bytes"
	"github.com/go-pg/pg/orm"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"github.com/plusspeed/payments-api/internal/api"
	"github.com/plusspeed/payments-api/internal/model"
	"github.com/plusspeed/payments-api/internal/repository"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPaymentsAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Payment Suite")
}

const (
	// values from the docker compose
	pgUsername = "test"
	pgPassword = "example"
	pgAddress  = "127.0.0.1:5432"
	basePath   = "/v1"
)

var dbTest *repository.Repository
var router *mux.Router

var _ = Describe("Payment", func() {

	BeforeEach(func() {
		clearDB(*dbTest)
	})

	Describe("when payment api is running", func() {
		It("health endpoint should return 200", func() {
			req1, _ := http.NewRequest("Get", "/health", nil)
			response := executeRequest(*router, req1)
			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("when I create payments", func() {
		var paymentID = uuid.NewRandom().String()

		It("should return Status Bad Request if request is empty", func() {
			req, _ := http.NewRequest("POST", "/v1/payment", bytes.NewBufferString("{}"))
			response := executeRequest(*router, req)
			Expect(http.StatusBadRequest).To(Equal(response.Code))
		})

		It("should return Status Bad Request if request is invalid", func() {
			req, _ := http.NewRequest("POST", "/v1/payment", bytes.NewBuffer(createBadRequest(paymentID)))
			response := executeRequest(*router, req)
			Expect(http.StatusBadRequest).To(Equal(response.Code))
		})

		It("should return Status Created if request is valid", func() {
			req, _ := http.NewRequest("POST", "/v1/payment", bytes.NewBuffer(createRequest(paymentID)))
			response := executeRequest(*router, req)
			Expect(http.StatusCreated).To(Equal(response.Code))
		})
		It("should be idempotent", func() {
			req, _ := http.NewRequest("POST", "/v1/payment", bytes.NewBuffer(createRequest(paymentID)))
			response := executeRequest(*router, req)
			Expect(http.StatusCreated).To(Equal(response.Code))
		})

		It("should return 409 if request if different", func() {
			newVersion := "2"
			req, _ := http.NewRequest("POST", "/v1/payment", bytes.NewBufferString("{\"type\": \"Payment\","+
				"\"id\": \""+paymentID+"\","+
				"\"version\": "+newVersion+","+
				"\"organisation_id\": \"743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb\","+
				"\"attributes\": {\"amount\": \"100.21\",\"beneficiary_party\": {\"account_name\": \"W Owens\",\"account_number\": \"31926819\",\"account_number_code\": \"BBAN\",\"account_type\": 0,\"address\": \"1 The Beneficiary Localtown SE2\",\"bank_id\": \"403000\",\"bank_id_code\": \"GBDSC\",\"name\": \"Wilfred Jeremiah Owens\"},\"charges_information\": {\"bearer_code\": \"SHAR\",\"sender_charges\": [{\"amount\": \"5.00\",\"currency\": \"GBP\"},{\"amount\": \"10.00\",\"currency\": \"USD\"}],\"receiver_charges_amount\": \"1.00\",\"receiver_charges_currency\": \"USD\"},\"currency\": \"GBP\",\"debtor_party\": {\"account_name\": \"EJ Brown Black\",\"account_number\": \"GB29XABC10161234567801\",\"account_number_code\": \"IBAN\",\"address\": \"10 Debtor Crescent Sourcetown NE1\",\"bank_id\": \"203301\",\"bank_id_code\": \"GBDSC\",\"name\": \"Emelia Jane Brown\"},\"end_to_end_reference\": \"Wil piano Jan\",\"fx\": {\"contract_reference\": \"FX123\",\"exchange_rate\": \"2.00000\",\"original_amount\": \"200.42\",\"original_currency\": \"USD\"},\"numeric_reference\": \"1002001\",\"payment_id\": \"123456789012345678\",\"payment_purpose\": \"Paying for goods/services\",\"payment_scheme\": \"FPS\",\"payment_type\": \"Credit\",\"processing_date\": \"2017-01-18\",\"reference\": \"Attributes for Em's piano lessons\",\"scheme_payment_sub_type\": \"InternetBanking\",\"scheme_payment_type\": \"ImmediatePayment\",\"sponsor_party\": {\"account_number\": \"56781234\",\"bank_id\": \"123123\",\"bank_id_code\": \"GBDSC\"}}}"))
			response := executeRequest(*router, req)
			Expect(http.StatusConflict).To(Equal(response.Code))
		})
	})

	Describe("when I get payment", func() {
		var paymentID = uuid.NewRandom().String()

		It("should return 404 if does not exist", func() {
			req, _ := http.NewRequest("GET", "/v1/payment/"+paymentID, bytes.NewBufferString(""))
			response := executeRequest(*router, req)
			Expect(http.StatusNotFound).To(Equal(response.Code))
		})

		Context("after creating payment", func() {

			It("should return the payment", func() {
				reqCreate, _ := http.NewRequest("POST", "/v1/payment", bytes.NewBuffer(createRequest(paymentID)))
				executeRequest(*router, reqCreate)

				newVersion := "2"
				req, _ := http.NewRequest("POST", "/v1/payment", bytes.NewBufferString("{\"type\": \"Payment\","+
					"\"id\": \""+paymentID+"\","+
					"\"version\": "+newVersion+","+
					"\"organisation_id\": \"743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb\","+
					"\"attributes\": {\"amount\": \"100.21\",\"beneficiary_party\": {\"account_name\": \"W Owens\",\"account_number\": \"31926819\",\"account_number_code\": \"BBAN\",\"account_type\": 0,\"address\": \"1 The Beneficiary Localtown SE2\",\"bank_id\": \"403000\",\"bank_id_code\": \"GBDSC\",\"name\": \"Wilfred Jeremiah Owens\"},\"charges_information\": {\"bearer_code\": \"SHAR\",\"sender_charges\": [{\"amount\": \"5.00\",\"currency\": \"GBP\"},{\"amount\": \"10.00\",\"currency\": \"USD\"}],\"receiver_charges_amount\": \"1.00\",\"receiver_charges_currency\": \"USD\"},\"currency\": \"GBP\",\"debtor_party\": {\"account_name\": \"EJ Brown Black\",\"account_number\": \"GB29XABC10161234567801\",\"account_number_code\": \"IBAN\",\"address\": \"10 Debtor Crescent Sourcetown NE1\",\"bank_id\": \"203301\",\"bank_id_code\": \"GBDSC\",\"name\": \"Emelia Jane Brown\"},\"end_to_end_reference\": \"Wil piano Jan\",\"fx\": {\"contract_reference\": \"FX123\",\"exchange_rate\": \"2.00000\",\"original_amount\": \"200.42\",\"original_currency\": \"USD\"},\"numeric_reference\": \"1002001\",\"payment_id\": \"123456789012345678\",\"payment_purpose\": \"Paying for goods/services\",\"payment_scheme\": \"FPS\",\"payment_type\": \"Credit\",\"processing_date\": \"2017-01-18\",\"reference\": \"Attributes for Em's piano lessons\",\"scheme_payment_sub_type\": \"InternetBanking\",\"scheme_payment_type\": \"ImmediatePayment\",\"sponsor_party\": {\"account_number\": \"56781234\",\"bank_id\": \"123123\",\"bank_id_code\": \"GBDSC\"}}}"))
				response := executeRequest(*router, req)
				Expect(http.StatusConflict).To(Equal(response.Code))
			})
		})
	})

	Describe("when I update payment", func() {
		var paymentID = uuid.NewRandom().String()
		var organisationId = uuid.NewRandom().String()

		It("should return 404 if does not exist", func() {
			req, _ := http.NewRequest("GET", "/v1/payment/"+paymentID, bytes.NewBufferString(""))
			response := executeRequest(*router, req)
			Expect(http.StatusNotFound).To(Equal(response.Code))
		})

		Context("after creating payment", func() {

			It("should update the payment", func() {

				reqCreate, _ := http.NewRequest("POST", "/v1/payment", bytes.NewBuffer(createRequest(paymentID)))
				executeRequest(*router, reqCreate)

				reqUpdate, _ := http.NewRequest("PUT", "/v1/payment/"+paymentID, bytes.NewBuffer(createRequestUpdated(paymentID, organisationId)))
				executeRequest(*router, reqUpdate)

				req, _ := http.NewRequest("GET", "/v1/payment/"+paymentID, bytes.NewBufferString(""))
				response := executeRequest(*router, req)

				Expect(http.StatusOK).To(Equal(response.Code))
				Expect(response.Body.String()).To(ContainSubstring("\"id\":\"" + paymentID + "\""))
				Expect(response.Body.String()).To(ContainSubstring("\"organisation_id\":\"" + organisationId + "\""))
			})
			It("should return Status Bad Request if request is invalid", func() {
				req, _ := http.NewRequest("POST", "/v1/payment", bytes.NewBuffer(createBadRequest(paymentID)))
				response := executeRequest(*router, req)
				Expect(http.StatusBadRequest).To(Equal(response.Code))
			})
		})
	})

	Describe("when I delete payment", func() {
		var paymentId = uuid.NewRandom().String()

		It("should return 404 if does not exist", func() {
			req, _ := http.NewRequest("DELETE", "/v1/payment/"+paymentId, bytes.NewBufferString(""))
			response := executeRequest(*router, req)
			Expect(http.StatusNotFound).To(Equal(response.Code))
		})

		Context("after creating payment", func() {

			It("should return Status OK", func() {
				reqCreate, _ := http.NewRequest("POST", "/v1/payment", bytes.NewBuffer(createRequest(paymentId)))
				executeRequest(*router, reqCreate)

				req, _ := http.NewRequest("DELETE", "/v1/payment/"+paymentId, bytes.NewBufferString(""))
				response := executeRequest(*router, req)
				Expect(http.StatusNoContent).To(Equal(response.Code))
			})
		})
	})

	Describe("when I fetch payments", func() {
		It("should return Status Not Found because does not exist", func() {
			req, _ := http.NewRequest("GET", "/v1/payments", nil)
			response := executeRequest(*router, req)
			Expect(http.StatusNotFound).To(Equal(response.Code))
		})

		It("should return list of payments", func() {

			var paymentId1 = uuid.NewRandom().String()
			var paymentId2 = uuid.NewRandom().String()

			reqPayment1, _ := http.NewRequest("POST", "/v1/payment", bytes.NewBuffer(createRequest(paymentId1)))
			executeRequest(*router, reqPayment1)

			reqPaymentId2, _ := http.NewRequest("POST", "/v1/payment", bytes.NewBuffer(createRequest(paymentId2)))
			executeRequest(*router, reqPaymentId2)

			req, _ := http.NewRequest("GET", "/v1/payments?limit=100&offset=0", nil)
			response := executeRequest(*router, req)
			Expect(response.Body.String()).To(ContainSubstring("\"id\":\"" + paymentId1 + "\""))
			Expect(response.Body.String()).To(ContainSubstring("\"id\":\"" + paymentId2 + "\""))
			Expect(http.StatusOK).To(Equal(response.Code))
		})

		It("if ", func() {

			var paymentId1 = uuid.NewRandom().String()
			var paymentId2 = uuid.NewRandom().String()

			reqPayment1, _ := http.NewRequest("POST", "/v1/payment", bytes.NewBuffer(createRequest(paymentId1)))
			executeRequest(*router, reqPayment1)

			reqPaymentId2, _ := http.NewRequest("POST", "/v1/payment", bytes.NewBuffer(createRequest(paymentId2)))
			executeRequest(*router, reqPaymentId2)

			req, _ := http.NewRequest("GET", "/v1/payments?limit=100&offset=0", nil)
			response := executeRequest(*router, req)
			Expect(response.Body.String()).To(ContainSubstring("\"id\":\"" + paymentId1 + "\""))
			Expect(response.Body.String()).To(ContainSubstring("\"id\":\"" + paymentId2 + "\""))
			Expect(http.StatusOK).To(Equal(response.Code))
		})
	})
})

var _ = BeforeSuite(func() {
	dbTest = repository.New(pgAddress, "test", pgUsername, pgPassword)
	router = api.NewRouter(basePath, dbTest)
})

var _ = AfterSuite(func() {
	dbTest.Database.Close()
})

func clearDB(dbTest repository.Repository) {
	_, _ = dbTest.Database.Exec("DROP TABLE PAYMENT")

	err := dbTest.Database.CreateTable(&model.Payment{}, &orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		panic(err.Error())
	}
}

func executeRequest(router mux.Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func createRequest(paymentId string) []byte {
	return []byte("{\"type\": \"Payment\"," +
		"\"id\": \"" + paymentId + "\"," +
		"\"version\": 0," +
		"\"organisation_id\": \"743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb\"," +
		"\"attributes\": {\"amount\": \"100.21\",\"beneficiary_party\": {\"account_name\": \"W Owens\",\"account_number\": \"31926819\",\"account_number_code\": \"BBAN\",\"account_type\": 0,\"address\": \"1 The Beneficiary Localtown SE2\",\"bank_id\": \"403000\",\"bank_id_code\": \"GBDSC\",\"name\": \"Wilfred Jeremiah Owens\"},\"charges_information\": {\"bearer_code\": \"SHAR\",\"sender_charges\": [{\"amount\": \"5.00\",\"currency\": \"GBP\"},{\"amount\": \"10.00\",\"currency\": \"USD\"}],\"receiver_charges_amount\": \"1.00\",\"receiver_charges_currency\": \"USD\"},\"currency\": \"GBP\",\"debtor_party\": {\"account_name\": \"EJ Brown Black\",\"account_number\": \"GB29XABC10161234567801\",\"account_number_code\": \"IBAN\",\"address\": \"10 Debtor Crescent Sourcetown NE1\",\"bank_id\": \"203301\",\"bank_id_code\": \"GBDSC\",\"name\": \"Emelia Jane Brown\"},\"end_to_end_reference\": \"Wil piano Jan\",\"fx\": {\"contract_reference\": \"FX123\",\"exchange_rate\": \"2.00000\",\"original_amount\": \"200.42\",\"original_currency\": \"USD\"},\"numeric_reference\": \"1002001\",\"payment_id\": \"123456789012345678\",\"payment_purpose\": \"Paying for goods/services\",\"payment_scheme\": \"FPS\",\"payment_type\": \"Credit\",\"processing_date\": \"2017-01-18\",\"reference\": \"Payment for Em's piano lessons\",\"scheme_payment_sub_type\": \"InternetBanking\",\"scheme_payment_type\": \"ImmediatePayment\",\"sponsor_party\": {\"account_number\": \"56781234\",\"bank_id\": \"123123\",\"bank_id_code\": \"GBDSC\"}}}")
}

func createRequestUpdated(paymentId, organisationId string) []byte {
	return []byte("{\"type\": \"Payment\"," +
		"\"id\": \"" + paymentId + "\"," +
		"\"version\": 0," +
		"\"organisation_id\": \"" + organisationId + "\"," +
		"\"attributes\": {\"amount\": \"100.21\",\"beneficiary_party\": {\"account_name\": \"W Owens\",\"account_number\": \"31926819\",\"account_number_code\": \"BBAN\",\"account_type\": 0,\"address\": \"1 The Beneficiary Localtown SE2\",\"bank_id\": \"403000\",\"bank_id_code\": \"GBDSC\",\"name\": \"Wilfred Jeremiah Owens\"},\"charges_information\": {\"bearer_code\": \"SHAR\",\"sender_charges\": [{\"amount\": \"5.00\",\"currency\": \"GBP\"},{\"amount\": \"10.00\",\"currency\": \"USD\"}],\"receiver_charges_amount\": \"1.00\",\"receiver_charges_currency\": \"USD\"},\"currency\": \"GBP\",\"debtor_party\": {\"account_name\": \"EJ Brown Black\",\"account_number\": \"GB29XABC10161234567801\",\"account_number_code\": \"IBAN\",\"address\": \"10 Debtor Crescent Sourcetown NE1\",\"bank_id\": \"203301\",\"bank_id_code\": \"GBDSC\",\"name\": \"Emelia Jane Brown\"},\"end_to_end_reference\": \"Wil piano Jan\",\"fx\": {\"contract_reference\": \"FX123\",\"exchange_rate\": \"2.00000\",\"original_amount\": \"200.42\",\"original_currency\": \"USD\"},\"numeric_reference\": \"1002001\",\"payment_id\": \"123456789012345678\",\"payment_purpose\": \"Paying for goods/services\",\"payment_scheme\": \"FPS\",\"payment_type\": \"Credit\",\"processing_date\": \"2017-01-18\",\"reference\": \"Payment for Em's piano lessons\",\"scheme_payment_sub_type\": \"InternetBanking\",\"scheme_payment_type\": \"ImmediatePayment\",\"sponsor_party\": {\"account_number\": \"56781234\",\"bank_id\": \"123123\",\"bank_id_code\": \"GBDSC\"}}}")
}

func createBadRequest(paymentId string) []byte {
	return []byte("{\"type\": \"Payment\"," +
		"\"id\": \"" + paymentId + "\"," +
		"\"version\": 0," +
		"\"organisation_id\": \"743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb\"," +
		"\"attributes\": {\"amount\": \"100.21\",\"beneficiary_party\": {\"account_number\": \"31926819\",\"account_type\": 0,\"address\": \"1 The Beneficiary Localtown SE2\",\"bank_id\": \"403000\",\"bank_id_code\": \"GBDSC\",\"name\": \"Wilfred Jeremiah Owens\"},\"charges_information\": {\"bearer_code\": \"SHAR\",\"sender_charges\": [{\"amount\": \"5.00\",\"currency\": \"GBP\"},{\"amount\": \"10.00\",\"currency\": \"USD\"}],\"receiver_charges_amount\": \"1.00\",\"receiver_charges_currency\": \"USD\"},\"currency\": \"GBP\",\"debtor_party\": {\"account_name\": \"EJ Brown Black\",\"account_number\": \"GB29XABC10161234567801\",\"account_number_code\": \"IBAN\",\"address\": \"10 Debtor Crescent Sourcetown NE1\",\"bank_id\": \"203301\",\"bank_id_code\": \"GBDSC\",\"name\": \"Emelia Jane Brown\"},\"end_to_end_reference\": \"Wil piano Jan\",\"fx\": {\"contract_reference\": \"FX123\",\"exchange_rate\": \"2.00000\",\"original_amount\": \"200.42\",\"original_currency\": \"USD\"},\"numeric_reference\": \"1002001\",\"payment_id\": \"123456789012345678\",\"payment_purpose\": \"Paying for goods/services\",\"payment_scheme\": \"FPS\",\"payment_type\": \"Credit\",\"processing_date\": \"2017-01-18\",\"reference\": \"Payment for Em's piano lessons\",\"scheme_payment_sub_type\": \"InternetBanking\",\"scheme_payment_type\": \"ImmediatePayment\",\"sponsor_party\": {\"account_number\": \"56781234\",\"bank_id\": \"123123\",\"bank_id_code\": \"GBDSC\"}}}")
}
