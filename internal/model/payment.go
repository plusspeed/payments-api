package model

//Payment Root element
type Payment struct {
	Type           string     `json:"type" sql:",notnull" validate:"required"`
	ID             string     `json:"id" validate:"required"`
	Version        int        `json:"version" sql:",notnull"`
	OrganisationID string     `json:"organisation_id" sql:",notnull" validate:"required"`
	Attributes     Attributes `json:"attributes" sql:",notnull" validate:"required"`
}

//Attributes contains details about a payment
type Attributes struct {
	Amount           string `json:"amount" sql:",notnull" validate:"required"`
	BeneficiaryParty struct {
		AccountName       string `json:"account_name" sql:",notnull" validate:"required"`
		AccountNumber     string `json:"account_number" sql:",notnull" validate:"required"`
		AccountNumberCode string `json:"account_number_code" sql:",notnull" validate:"required"`
		AccountType       int    `json:"account_type" sql:",notnull"`
		Address           string `json:"address" sql:",notnull" validate:"required"`
		BankID            string `json:"bank_id" sql:",notnull" validate:"required"`
		BankIDCode        string `json:"bank_id_code" sql:",notnull" validate:"required"`
		Name              string `json:"name" sql:",notnull" validate:"required"`
	} `json:"beneficiary_party" sql:",notnull" validate:"required"`
	ChargesInformation struct {
		BearerCode    string `json:"bearer_code" sql:",notnull" validate:"required"`
		SenderCharges []struct {
			Amount   string `json:"amount" sql:",notnull" validate:"required"`
			Currency string `json:"currency" sql:",notnull" validate:"required"`
		} `json:"sender_charges" sql:",notnull" validate:"required"`
		ReceiverChargesAmount   string `json:"receiver_charges_amount" sql:",notnull" validate:"required"`
		ReceiverChargesCurrency string `json:"receiver_charges_currency" sql:",notnull" validate:"required"`
	} `json:"charges_information" sql:",notnull" validate:"required"`
	Currency    string `json:"currency" sql:",notnull" validate:"required"`
	DebtorParty struct {
		AccountName       string `json:"account_name" sql:",notnull" validate:"required"`
		AccountNumber     string `json:"account_number" sql:",notnull" validate:"required"`
		AccountNumberCode string `json:"account_number_code" sql:",notnull" validate:"required"`
		Address           string `json:"address" sql:",notnull" validate:"required"`
		BankID            string `json:"bank_id" sql:",notnull" validate:"required"`
		BankIDCode        string `json:"bank_id_code" sql:",notnull" validate:"required"`
		Name              string `json:"name" sql:",notnull" validate:"required"`
	} `json:"debtor_party" sql:",notnull" validate:"required"`
	EndToEndReference string `json:"end_to_end_reference" sql:",notnull" validate:"required"`
	Fx                struct {
		ContractReference string `json:"contract_reference" sql:",notnull" validate:"required"`
		ExchangeRate      string `json:"exchange_rate" sql:",notnull" validate:"required"`
		OriginalAmount    string `json:"original_amount" sql:",notnull" validate:"required"`
		OriginalCurrency  string `json:"original_currency" sql:",notnull" validate:"required"`
	} `json:"fx" sql:",notnull" validate:"required"`
	NumericReference     string `json:"numeric_reference" sql:",notnull" validate:"required"`
	ID                   string `json:"payment_id" sql:",notnull" validate:"required"`
	Purpose              string `json:"payment_purpose" sql:",notnull" validate:"required"`
	Scheme               string `json:"payment_scheme" sql:",notnull" validate:"required"`
	Type                 string `json:"payment_type" sql:",notnull" validate:"required"`
	ProcessingDate       string `json:"processing_date" sql:",notnull" validate:"required"`
	Reference            string `json:"reference" sql:",notnull" validate:"required"`
	SchemePaymentSubType string `json:"scheme_payment_sub_type" sql:",notnull" validate:"required"`
	SchemePaymentType    string `json:"scheme_payment_type" sql:",notnull" validate:"required"`
	SponsorParty         struct {
		AccountNumber string `json:"account_number" sql:",notnull" validate:"required"`
		BankID        string `json:"bank_id" sql:",notnull" validate:"required"`
		BankIDCode    string `json:"bank_id_code" sql:",notnull" validate:"required"`
	} `json:"sponsor_party" sql:",notnull" validate:"required"`
}
