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

const Fieldcurrency_home_redis = "LISTCURR_BACKEND_WHITELABEL"
const Fieldcurrency_frontend_redis = "LISTCURR_FRONTEND_WHITELABEL"

func Currencyhome(c *fiber.Ctx) error {
	var obj entities.Model_currency
	var arraobj []entities.Model_currency
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldcurrency_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		currency_id, _ := jsonparser.GetString(value, "currency_id")
		currency_name, _ := jsonparser.GetString(value, "currency_name")
		currency_create, _ := jsonparser.GetString(value, "currency_create")
		currency_update, _ := jsonparser.GetString(value, "currency_update")

		obj.Currency_id = currency_id
		obj.Currency_name = currency_name
		obj.Currency_create = currency_create
		obj.Currency_update = currency_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_currHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldcurrency_home_redis, result, 60*time.Minute)
		fmt.Println("CURR MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("CURR CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func CurrSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_currencysave)
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

	result, err := models.Save_currency(
		client_admin,
		client.Currency_id, client.Currency_name, client.Sdata)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_curr()
	return c.JSON(result)
}
func _deleteredis_curr() {
	val_master := helpers.DeleteRedis(Fieldcurrency_home_redis)
	log.Printf("Redis Delete BACKEND CURR : %d", val_master)

}
