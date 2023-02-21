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

const Fieldcatebank_home_redis = "LISTCATEBANK_BACKEND_ISBPANEL"

func Catebankhome(c *fiber.Ctx) error {
	var obj entities.Model_catebank
	var arraobj []entities.Model_catebank
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldcatebank_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		catebank_id, _ := jsonparser.GetInt(value, "catebank_id")
		catebank_name, _ := jsonparser.GetString(value, "catebank_name")
		catebank_status, _ := jsonparser.GetString(value, "catebank_status")
		catebank_create, _ := jsonparser.GetString(value, "catebank_create")
		catebank_update, _ := jsonparser.GetString(value, "catebank_update")

		obj.Catebank_id = int(catebank_id)
		obj.Catebank_name = catebank_name
		obj.Catebank_status = catebank_status
		obj.Catebank_create = catebank_create
		obj.Catebank_update = catebank_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_catebankHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldcatebank_home_redis, result, 60*time.Minute)
		fmt.Println("CATEBANK MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("CATEBANK CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func CatebankSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_catebanksave)
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

	//admin, name, status, sData string, idrecord int
	result, err := models.Save_catebank(
		client_admin,
		client.Catebank_name, client.Catebank_status, client.Sdata, client.Catebank_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	_deleteredis_catebank()
	return c.JSON(result)
}
func _deleteredis_catebank() {
	val_master := helpers.DeleteRedis(Fieldcatebank_home_redis)
	log.Printf("Redis Delete BACKEND CATE BANK : %d", val_master)

}
