package entities

type Model_catebank struct {
	Catebank_id     int    `json:"catebank_id"`
	Catebank_name   string `json:"catebank_name"`
	Catebank_status string `json:"catebank_status"`
	Catebank_create string `json:"catebank_create"`
	Catebank_update string `json:"catebank_update"`
}
type Controller_catebanksave struct {
	Page            string `json:"page" validate:"required"`
	Sdata           string `json:"sdata" validate:"required"`
	Catebank_id     int    `json:"catebank_id" validate:"required"`
	Catebank_name   string `json:"catebank_name" validate:"required"`
	Catebank_status string `json:"catebank_status" validate:"required"`
}
