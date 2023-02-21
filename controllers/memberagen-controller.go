package controllers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/entities"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/helpers"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/models"
)

const Fieldmember_home_redis = "LISTMEMBER_BACKEND_ISBPANEL"
const Fieldmemberselect_home_redis = "LISTMEMBERSELECT_BACKEND_ISBPANEL"

func Memberhome(c *fiber.Ctx) error {
	var obj entities.Model_member
	var arraobj []entities.Model_member
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldmember_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		member_phone, _ := jsonparser.GetString(value, "member_phone")
		member_name, _ := jsonparser.GetString(value, "member_name")
		member_create, _ := jsonparser.GetString(value, "member_create")
		member_update, _ := jsonparser.GetString(value, "member_update")

		var objwebsiteagen entities.Model_memberagen
		var arraobjwebsiteagen []entities.Model_memberagen
		record_memberagen_RD, _, _, _ := jsonparser.Get(value, "member_agen")
		jsonparser.ArrayEach(record_memberagen_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			memberagen_idwebagen, _ := jsonparser.GetInt(value, "memberagen_idwebagen")
			memberagen_username, _ := jsonparser.GetString(value, "memberagen_username")
			memberagen_website, _ := jsonparser.GetString(value, "memberagen_website")

			objwebsiteagen.Memberagen_idwebagen = int(memberagen_idwebagen)
			objwebsiteagen.Memberagen_username = memberagen_username
			objwebsiteagen.Memberagen_website = memberagen_website
			arraobjwebsiteagen = append(arraobjwebsiteagen, objwebsiteagen)
		})

		obj.Member_phone = member_phone
		obj.Member_name = member_name
		obj.Member_agen = arraobjwebsiteagen
		obj.Member_create = member_create
		obj.Member_update = member_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_member()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldmember_home_redis, result, 60*time.Minute)
		fmt.Println("MEMBER  MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("MEMBER CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func Memberhomeselect(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_memberagenselect)
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
	var obj entities.Model_memberagenselect
	var arraobj []entities.Model_memberagenselect
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldmemberselect_home_redis + "_" + strconv.Itoa(client.Memberagen_idwebagen))
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		memberagen_id, _ := jsonparser.GetInt(value, "memberagen_id")
		memberagen_website, _ := jsonparser.GetString(value, "memberagen_website")
		memberagen_username, _ := jsonparser.GetString(value, "memberagen_username")
		memberagen_phone, _ := jsonparser.GetString(value, "memberagen_phone")
		memberagen_name, _ := jsonparser.GetString(value, "memberagen_name")

		obj.Memberagen_id = int(memberagen_id)
		obj.Memberagen_website = memberagen_website
		obj.Memberagen_username = memberagen_username
		obj.Memberagen_phone = memberagen_phone
		obj.Memberagen_name = memberagen_name
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_memberSelect(client.Memberagen_idwebagen)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldmemberselect_home_redis+"_"+strconv.Itoa(client.Memberagen_idwebagen), result, 60*time.Minute)
		fmt.Println("MEMBER AGEN SELECT  MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("MEMBER AGEN SELECT CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func MemberSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_membersave)
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

	// admin, phone, nama, sData string,
	// idrecord int
	result, err := models.Save_member(
		client_admin,
		client.Member_phone, client.Member_name, string(client.Member_listagen), client.Sdata, client.Member_phone)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_memberagen(client.Member_phone)
	return c.JSON(result)
}

func _deleteredis_memberagen(phone string) {
	val_member := helpers.DeleteRedis(Fieldmember_home_redis)
	log.Printf("Redis Delete BACKEND MEMBER  : %d", val_member)

	val_memberselect := helpers.DeleteRedis(Fieldmemberselect_home_redis)
	log.Printf("Redis Delete BACKEND MEMBER SELECT  : %d", val_memberselect)

}
