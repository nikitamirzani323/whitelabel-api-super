package entities

type Model_currency struct {
	Currency_id     string `json:"currency_id"`
	Currency_name   string `json:"currency_name"`
	Currency_create string `json:"currency_create"`
	Currency_update string `json:"currency_update"`
}

type Controller_currencysave struct {
	Page          string `json:"page" validate:"required"`
	Sdata         string `json:"sdata" validate:"required"`
	Currency_id   string `json:"currency_id"`
	Currency_name string `json:"currency_name" validate:"required"`
}
