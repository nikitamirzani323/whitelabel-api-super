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

const Fieldbanktype_home_redis = "LISTBANKTYPE_BACKEND_WHITELABEL"

func Banktypehome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_banktypehome)
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

	var obj entities.Model_banktype
	var arraobj []entities.Model_banktype
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldbanktype_home_redis)
	jsonredis := []byte(resultredis)
	message_RD, _ := jsonparser.GetString(jsonredis, "message")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		banktype_id, _ := jsonparser.GetString(value, "banktype_id")
		banktype_nmcatebank, _ := jsonparser.GetString(value, "banktype_nmcatebank")
		banktype_name, _ := jsonparser.GetString(value, "banktype_name")
		banktype_img, _ := jsonparser.GetString(value, "banktype_img")
		banktype_status, _ := jsonparser.GetString(value, "banktype_status")
		banktype_create, _ := jsonparser.GetString(value, "banktype_create")
		banktype_update, _ := jsonparser.GetString(value, "banktype_update")

		obj.Banktype_id = banktype_id
		obj.Banktype_name = banktype_name
		obj.Banktype_nmcatebank = banktype_nmcatebank
		obj.Banktype_img = banktype_img
		obj.Banktype_status = banktype_status
		obj.Banktype_create = banktype_create
		obj.Banktype_update = banktype_update
		arraobj = append(arraobj, obj)
	})
	if !flag {
		result, err := models.Fetch_banktypeHome(client.Banktype_idcatebank)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldbanktype_home_redis, result, 60*time.Minute)
		fmt.Println("BANK TYPE MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("BANK TYPE MYSQL")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": message_RD,
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func Banktypesave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_banktypesave)
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
	//admin, idrecord, name, img, status, sData string, idcatebank int
	result, err := models.Save_banktype(client_admin, client.Banktype_id,
		client.Banktype_name, client.Banktype_img, client.Banktype_status,
		client.Sdata, client.Banktype_idcatebank)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	_deleteredis_banktype()
	return c.JSON(result)
}
func _deleteredis_banktype() {
	val_master := helpers.DeleteRedis(Fieldbanktype_home_redis)
	log.Printf("Redis Delete BACKEND BANK TYPE : %d", val_master)
}
