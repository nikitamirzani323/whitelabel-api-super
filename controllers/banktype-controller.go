package controllers

import (
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

const Fieldbanktype_home_redis = "LISTBANKTYPE_BACKEND_ISBPANEL"

func Banktypehome(c *fiber.Ctx) error {
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
		result, err := models.Fetch_banktypeHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldbanktype_home_redis, result, 60*time.Minute)
		log.Println("BANK TYPE MYSQL")
		return c.JSON(result)
	} else {
		log.Println("BANK TYPE CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": message_RD,
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func Bannersave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_bannersave)
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
	//admin, sdata, nmbanner, urlbanner, urlwebsite, devicebanner, posisibanner, status string, idrecord, display int
	result, err := models.Save_banner(client_admin, client.Sdata, client.Banner_name,
		client.Banner_url, client.Banner_urlwebsite, client.Banner_device, client.Banner_posisi, client.Banner_status,
		client.Banner_id, client.Banner_display)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	_deleteredis_banner()
	return c.JSON(result)
}
func _deleteredis_banner() {
	val_master := helpers.DeleteRedis(Fieldbanner_home_redis)
	log.Printf("Redis Delete BACKEND BANNER : %d", val_master)

	//CLIENT
	val_client_banner := helpers.DeleteRedis(Fieldbanner_frontend_redis)
	log.Printf("Redis Delete CLIENT BANNER : %d", val_client_banner)
}
