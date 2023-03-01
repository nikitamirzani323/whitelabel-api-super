package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/entities"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/helpers"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/models"
)

const Fieldcompany_home_redis = "LISTCOMPANY_BACKEND_WHITELABEL"

func Companyhome(c *fiber.Ctx) error {
	var obj entities.Model_company
	var arraobj []entities.Model_company
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldcompany_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		company_id, _ := jsonparser.GetString(value, "company_id")
		company_idcurr, _ := jsonparser.GetString(value, "company_idcurr")
		company_start, _ := jsonparser.GetString(value, "company_start")
		company_end, _ := jsonparser.GetString(value, "company_end")
		company_name, _ := jsonparser.GetString(value, "company_name")
		company_owner, _ := jsonparser.GetString(value, "company_owner")
		company_phone, _ := jsonparser.GetString(value, "company_phone")
		company_email, _ := jsonparser.GetString(value, "company_email")
		company_companyurl, _ := jsonparser.GetString(value, "company_companyurl")
		company_status, _ := jsonparser.GetString(value, "company_status")
		company_create, _ := jsonparser.GetString(value, "company_create")
		company_update, _ := jsonparser.GetString(value, "company_update")

		obj.Company_id = company_id
		obj.Company_idcurr = company_idcurr
		obj.Company_start = company_start
		obj.Company_end = company_end
		obj.Company_name = company_name
		obj.Company_owner = company_owner
		obj.Company_phone = company_phone
		obj.Company_email = company_email
		obj.Company_companyurl = company_companyurl
		obj.Company_status = company_status
		obj.Company_create = company_create
		obj.Company_update = company_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_companyHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldcompany_home_redis, result, 60*time.Minute)
		fmt.Println("COMPANY MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("COMPANY CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func CompanySave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_companysave)
	validate := validator.New()
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	err := validate.Struct(client)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element helpers.ErrorResponse
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "validation",
			"record":  errors,
		})
	}
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	client_admin, _ := helpers.Parsing_Decry(temp_decp, "==")

	result, err := models.Save_company(
		client_admin,
		client.Company_id, client.Company_idcurr, client.Company_name, client.Company_owner, client.Company_phone, client.Company_email, client.Company_companyurl, client.Company_status, client.Sdata)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_company()
	return c.JSON(result)
}
func _deleteredis_company() {
	val_master := helpers.DeleteRedis(Fieldcompany_home_redis)
	log.Printf("Redis Delete BACKEND COMPANY : %d", val_master)

}
