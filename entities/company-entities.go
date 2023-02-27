package entities

type Model_company struct {
	Company_id         string `json:"company_id"`
	Company_idcurr     string `json:"company_idcurr"`
	Company_start      string `json:"company_start"`
	Company_end        string `json:"company_end"`
	Company_name       string `json:"company_name"`
	Company_owner      string `json:"company_owner"`
	Company_phone      string `json:"company_phone"`
	Company_email      string `json:"company_email"`
	Company_companyurl string `json:"company_companyurl"`
	Company_status     string `json:"company_status"`
	Company_create     string `json:"company_create"`
	Company_update     string `json:"company_update"`
}

type Controller_companysave struct {
	Page               string `json:"page" validate:"required"`
	Sdata              string `json:"sdata" validate:"required"`
	Company_id         string `json:"company_id"`
	Company_idcurr     string `json:"company_idcurr" validate:"required"`
	Company_start      string `json:"company_start" validate:"required"`
	Company_end        string `json:"company_end"`
	Company_name       string `json:"company_name" validate:"required"`
	Company_owner      string `json:"company_owner" validate:"required"`
	Company_phone      string `json:"company_phone"`
	Company_email      string `json:"company_email"`
	Company_companyurl string `json:"company_companyurl" validate:"required"`
	Company_status     string `json:"company_status" validate:"required"`
}
