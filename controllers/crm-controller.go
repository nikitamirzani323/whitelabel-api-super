package controllers

import (
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

const Fieldcrm_home_redis = "LISTCRM_BACKEND_ISBPANEL"
const Fieldcrmsales_home_redis = "LISTCRMSALES_BACKEND_ISBPANEL"
const Fieldcrmdeposit_home_redis = "LISTCRMDEPOSIT_BACKEND_ISBPANEL"
const Fieldcrmisbtv_home_redis = "LISTCRMISBTV_BACKEND_ISBPANEL"
const Fieldcrmduniafilm_home_redis = "LISTCRMDUNIAFILM_BACKEND_ISBPANEL"
const Fieldcrm_sales_redis = "LISTCRM_SALES_ISBPANEL"

func Crmhome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_crm)
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
	if client.Crm_search != "" {
		val_crm := helpers.DeleteRedis(Fieldcrm_home_redis + "_" + client.Crm_status + "_" + strconv.Itoa(client.Crm_page) + "_" + client.Crm_search)
		log.Printf("Redis Delete BACKEND CRM : %d", val_crm)
	}
	var obj entities.Model_crm
	var arraobj []entities.Model_crm
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldcrm_home_redis + "_" + client.Crm_status + "_" + strconv.Itoa(client.Crm_page) + "_" + client.Crm_search)
	jsonredis := []byte(resultredis)
	perpage_RD, _ := jsonparser.GetInt(jsonredis, "perpage")
	totalrecord_RD, _ := jsonparser.GetInt(jsonredis, "totalrecord")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		crm_id, _ := jsonparser.GetInt(value, "crm_id")
		crm_phone, _ := jsonparser.GetString(value, "crm_phone")
		crm_name, _ := jsonparser.GetString(value, "crm_name")
		crm_source, _ := jsonparser.GetString(value, "crm_source")
		crm_totalpic, _ := jsonparser.GetInt(value, "crm_totalpic")
		crm_status, _ := jsonparser.GetString(value, "crm_status")
		crm_statuscss, _ := jsonparser.GetString(value, "crm_statuscss")
		crm_create, _ := jsonparser.GetString(value, "crm_create")
		crm_update, _ := jsonparser.GetString(value, "crm_update")

		var obj_crmsales entities.Model_crmsales_simple
		var arraobj_crmsales []entities.Model_crmsales_simple
		crm_pic, _, _, _ := jsonparser.Get(value, "crm_pic")
		jsonparser.ArrayEach(crm_pic, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			crmsales_idcrmsales, _ := jsonparser.GetInt(value, "crmsales_idcrmsales")
			crmsales_username, _ := jsonparser.GetString(value, "crmsales_username")
			crmsales_nameemployee, _ := jsonparser.GetString(value, "crmsales_nameemployee")
			crmsales_status_utama, _ := jsonparser.GetString(value, "crmsales_status_utama")
			crmsales_status, _ := jsonparser.GetString(value, "crmsales_status")
			crmsales_note, _ := jsonparser.GetString(value, "crmsales_note")
			crmsales_nmwebagen, _ := jsonparser.GetString(value, "crmsales_nmwebagen")
			crmsales_idwebagen, _ := jsonparser.GetString(value, "crmsales_idwebagen")
			crmsales_deposit, _ := jsonparser.GetFloat(value, "crmsales_deposit")

			obj_crmsales.Crmsales_idcrmsales = int(crmsales_idcrmsales)
			obj_crmsales.Crmsales_username = crmsales_username
			obj_crmsales.Crmsales_nameemployee = crmsales_nameemployee
			obj_crmsales.Crmsales_status_utama = crmsales_status_utama
			obj_crmsales.Crmsales_status = crmsales_status
			obj_crmsales.Crmsales_note = crmsales_note
			obj_crmsales.Crmsales_nmwebagen = crmsales_nmwebagen
			obj_crmsales.Crmsales_idwebagen = crmsales_idwebagen
			obj_crmsales.Crmsales_deposit = float32(crmsales_deposit)
			arraobj_crmsales = append(arraobj_crmsales, obj_crmsales)
		})

		obj.Crm_id = int(crm_id)
		obj.Crm_phone = crm_phone
		obj.Crm_name = crm_name
		obj.Crm_source = crm_source
		obj.Crm_totalpic = int(crm_totalpic)
		obj.Crm_pic = arraobj_crmsales
		obj.Crm_status = crm_status
		obj.Crm_statuscss = crm_statuscss
		obj.Crm_create = crm_create
		obj.Crm_update = crm_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_crm(client.Crm_search, client.Crm_status, client.Crm_page)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldcrm_home_redis+"_"+client.Crm_status+"_"+strconv.Itoa(client.Crm_page)+"_"+client.Crm_search, result, 60*time.Minute)
		log.Println("CRM  MYSQL")
		return c.JSON(result)
	} else {
		log.Println("CRM  CACHE")
		return c.JSON(fiber.Map{
			"status":      fiber.StatusOK,
			"message":     "Success",
			"record":      arraobj,
			"perpage":     perpage_RD,
			"totalrecord": totalrecord_RD,
			"time":        time.Since(render_page).String(),
		})
	}
}
func Crmsales(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_crmsales)
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

	var obj entities.Model_crmsales
	var arraobj []entities.Model_crmsales
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldcrmsales_home_redis + "_" + client.Crmsales_phone + "_" + client.Crmsales_status)
	jsonredis := []byte(resultredis)
	perpage_RD, _ := jsonparser.GetInt(jsonredis, "perpage")
	totalrecord_RD, _ := jsonparser.GetInt(jsonredis, "totalrecord")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		crmsales_id, _ := jsonparser.GetInt(value, "crmsales_id")
		crmsales_phone, _ := jsonparser.GetString(value, "crmsales_phone")
		crmsales_namamember, _ := jsonparser.GetString(value, "crmsales_namamember")
		crmsales_username, _ := jsonparser.GetString(value, "crmsales_username")
		crmsales_nameemployee, _ := jsonparser.GetString(value, "crmsales_nameemployee")
		crmsales_create, _ := jsonparser.GetString(value, "crmsales_create")
		crmsales_update, _ := jsonparser.GetString(value, "crmsales_update")

		obj.Crmsales_id = int(crmsales_id)
		obj.Crmsales_phone = crmsales_phone
		obj.Crmsales_namamember = crmsales_namamember
		obj.Crmsales_username = crmsales_username
		obj.Crmsales_nameemployee = crmsales_nameemployee
		obj.Crmsales_create = crmsales_create
		obj.Crmsales_update = crmsales_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_crmsales(client.Crmsales_phone, client.Crmsales_status)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldcrmsales_home_redis+"_"+client.Crmsales_phone+"_"+client.Crmsales_status, result, 60*time.Minute)
		log.Println("CRM SALES  MYSQL")
		return c.JSON(result)
	} else {
		log.Println("CRM SALES  CACHE")
		return c.JSON(fiber.Map{
			"status":      fiber.StatusOK,
			"message":     "Success",
			"record":      arraobj,
			"perpage":     perpage_RD,
			"totalrecord": totalrecord_RD,
			"time":        time.Since(render_page).String(),
		})
	}
}
func Crmdeposit(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_crmdeposit)
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

	var obj entities.Model_crmdeposit
	var arraobj []entities.Model_crmdeposit
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldcrmdeposit_home_redis + "_" + strconv.Itoa(client.Crmsales_idcrmsales))
	jsonredis := []byte(resultredis)
	perpage_RD, _ := jsonparser.GetInt(jsonredis, "perpage")
	totalrecord_RD, _ := jsonparser.GetInt(jsonredis, "totalrecord")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		crmdeposit_nmwebagen, _ := jsonparser.GetString(value, "crmdeposit_nmwebagen")
		crmdeposit_deposit, _ := jsonparser.GetFloat(value, "crmdeposit_deposit")
		crmdeposit_iduseragen, _ := jsonparser.GetString(value, "crmdeposit_iduseragen")
		crmdeposit_create, _ := jsonparser.GetString(value, "crmdeposit_create")

		obj.Crmsdeposit_nmwebagen = crmdeposit_nmwebagen
		obj.Crmsdeposit_deposit = float32(crmdeposit_deposit)
		obj.Crmsdeposit_iduseragen = crmdeposit_iduseragen
		obj.Crmsdeposit_create = crmdeposit_create
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_crmdeposit(client.Crmsales_idcrmsales)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldcrmdeposit_home_redis+"_"+strconv.Itoa(client.Crmsales_idcrmsales), result, 60*time.Minute)
		log.Println("CRM DEPOSIT  MYSQL")
		return c.JSON(result)
	} else {
		log.Println("CRM DEPOSIT  CACHE")
		return c.JSON(fiber.Map{
			"status":      fiber.StatusOK,
			"message":     "Success",
			"record":      arraobj,
			"perpage":     perpage_RD,
			"totalrecord": totalrecord_RD,
			"time":        time.Since(render_page).String(),
		})
	}
}
func Crmisbtvhome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_crmisbtv)
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
	if client.Crmisbtv_search != "" {
		val_news := helpers.DeleteRedis(Fieldcrmisbtv_home_redis + "_" + strconv.Itoa(client.Crmisbtv_page) + "_" + client.Crmisbtv_search)
		log.Printf("Redis Delete BACKEND NEWS : %d", val_news)
	}
	var obj entities.Model_crmisbtv
	var arraobj []entities.Model_crmisbtv
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldcrmisbtv_home_redis + "_" + strconv.Itoa(client.Crmisbtv_page) + "_" + client.Crmisbtv_search)
	jsonredis := []byte(resultredis)
	perpage_RD, _ := jsonparser.GetInt(jsonredis, "perpage")
	totalrecord_RD, _ := jsonparser.GetInt(jsonredis, "totalrecord")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		crmisbtv_username, _ := jsonparser.GetString(value, "crmisbtv_username")
		crmisbtv_name, _ := jsonparser.GetString(value, "crmisbtv_name")
		crmisbtv_coderef, _ := jsonparser.GetString(value, "crmisbtv_coderef")
		crmisbtv_point, _ := jsonparser.GetInt(value, "crmisbtv_point")
		crmisbtv_status, _ := jsonparser.GetString(value, "crmisbtv_status")
		crmisbtv_lastlogin, _ := jsonparser.GetString(value, "crmisbtv_lastlogin")
		crmisbtv_create, _ := jsonparser.GetString(value, "crmisbtv_create")
		crmisbtv_update, _ := jsonparser.GetString(value, "crmisbtv_update")

		obj.Crmisbtv_username = crmisbtv_username
		obj.Crmisbtv_name = crmisbtv_name
		obj.Crmisbtv_coderef = crmisbtv_coderef
		obj.Crmisbtv_point = int(crmisbtv_point)
		obj.Crmisbtv_status = crmisbtv_status
		obj.Crmisbtv_lastlogin = crmisbtv_lastlogin
		obj.Crmisbtv_create = crmisbtv_create
		obj.Crmisbtv_update = crmisbtv_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_crmisbtv(client.Crmisbtv_search, client.Crmisbtv_page)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldcrmisbtv_home_redis+"_"+strconv.Itoa(client.Crmisbtv_page)+"_"+client.Crmisbtv_search, result, 60*time.Minute)
		log.Println("CRM ISBTV MYSQL")
		return c.JSON(result)
	} else {
		log.Println("CRM ISBTV CACHE")
		return c.JSON(fiber.Map{
			"status":      fiber.StatusOK,
			"message":     "Success",
			"record":      arraobj,
			"perpage":     perpage_RD,
			"totalrecord": totalrecord_RD,
			"time":        time.Since(render_page).String(),
		})
	}
}
func Crmduniafilm(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_crmisbtv)
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
	if client.Crmisbtv_search != "" {
		val_news := helpers.DeleteRedis(Fieldcrmduniafilm_home_redis + "_" + strconv.Itoa(client.Crmisbtv_page) + "_" + client.Crmisbtv_search)
		log.Printf("Redis Delete BACKEND NEWS : %d", val_news)
	}
	var obj entities.Model_crmduniafilm
	var arraobj []entities.Model_crmduniafilm
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldcrmduniafilm_home_redis + "_" + strconv.Itoa(client.Crmisbtv_page) + "_" + client.Crmisbtv_search)
	jsonredis := []byte(resultredis)
	perpage_RD, _ := jsonparser.GetInt(jsonredis, "perpage")
	totalrecord_RD, _ := jsonparser.GetInt(jsonredis, "totalrecord")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		crmduniafilm_username, _ := jsonparser.GetString(value, "crmduniafilm_username")
		crmduniafilm_name, _ := jsonparser.GetString(value, "crmduniafilm_name")

		obj.Crmduniafilm_username = crmduniafilm_username
		obj.Crmduniafilm_name = crmduniafilm_name
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_crmduniafilm(client.Crmisbtv_search, client.Crmisbtv_page)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldcrmduniafilm_home_redis+"_"+strconv.Itoa(client.Crmisbtv_page)+"_"+client.Crmisbtv_search, result, 60*time.Minute)
		log.Println("CRM DUNIA FILM MYSQL")
		return c.JSON(result)
	} else {
		log.Println("CRM DUNIA FILM CACHE")
		return c.JSON(fiber.Map{
			"status":      fiber.StatusOK,
			"message":     "Success",
			"record":      arraobj,
			"perpage":     perpage_RD,
			"totalrecord": totalrecord_RD,
			"time":        time.Since(render_page).String(),
		})
	}
}

func CrmSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_crmsave)
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

	result, err := models.Save_crm(
		client_admin,
		client.Crm_phone, client.Crm_name, client.Crm_status, client.Sdata, client.Crm_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_crm(client.Crm_page, client.Crm_status, "", "", "")
	return c.JSON(result)
}
func CrmSavestatus(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_crmstatussave)
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

	result, err := models.Save_crmstatus(
		client_admin,
		client.Crm_status, client.Crm_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_crm(client.Crm_page, client.Crm_status, "", "", "")
	return c.JSON(result)
}
func CrmSalesSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_crmsalessave)
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

	//admin, phone, username string)
	result, err := models.Save_crmsales(
		client_admin,
		client.Crmsales_phone, client.Crmsales_username)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	_deleteredis_crm(client.Crm_page, "", client.Crmsales_phone, client.Crmsales_username, client.Search)
	return c.JSON(result)
}
func CrmSalesdelete(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_crmsalesdelete)
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
	_, idruleadmin := helpers.Parsing_Decry(temp_decp, "==")
	log.Println("RULE :" + client.Page)
	ruleadmin := models.Get_AdminRule("ruleadmingroup", idruleadmin)
	flag := models.Get_listitemsearch(ruleadmin, ",", client.Page)

	if !flag {
		c.Status(fiber.StatusForbidden)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Anda tidak bisa akses halaman ini",
			"record":  nil,
		})
	} else {
		//phone string, idrecord int
		result, err := models.Delete_crmsales(client.Crmsales_phone, client.Crmsales_id)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		_deleteredis_crm(client.Crm_page, "", client.Crmsales_phone, "", client.Search)
		return c.JSON(result)
	}
}
func CrmSavesource(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_crmsavesource)
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

	result, err := models.Save_crmsource(
		client_admin,
		string(client.Crm_data), client.Crm_source, client.Sdata)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_crm(client.Crm_page, "", "", "", "")
	return c.JSON(result)
}
func CrmSavedatabase(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_crmsavesource)
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

	result, err := models.Save_crmdatabase(
		client_admin,
		string(client.Crm_data), client.Crm_source, client.Sdata)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_crm(client.Crm_page, "", "", "", "")
	return c.JSON(result)
}
func CrmSavemaintenance(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_crmsavemaintenance)
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

	result, err := models.Save_crmmaintenance(
		client_admin,
		string(client.Crm_data), client.Sdata)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_crm(client.Crm_page, "", "", "", "")
	return c.JSON(result)
}
func _deleteredis_crm(page int, status, phone, username, search string) {
	for i := 0; i <= 5000; i = i + 250 {
		val_pusat_1 := helpers.DeleteRedis(Fieldcrm_home_redis + "_NEW_" + strconv.Itoa(i) + "_" + search)
		val_pusat_2 := helpers.DeleteRedis(Fieldcrm_home_redis + "_PROCESS_" + strconv.Itoa(i) + "_" + search)
		val_pusat_3 := helpers.DeleteRedis(Fieldcrm_home_redis + "_VALID_" + strconv.Itoa(i) + "_" + search)
		val_pusat_4 := helpers.DeleteRedis(Fieldcrm_home_redis + "_INVALID_" + strconv.Itoa(i) + "_" + search)
		val_pusat_5 := helpers.DeleteRedis(Fieldcrm_home_redis + "_MAINTENANCE_" + strconv.Itoa(i) + "_" + search)
		val_pusat_6 := helpers.DeleteRedis(Fieldcrm_home_redis + "_FOLLOWUP_" + strconv.Itoa(i) + "_" + search)
		val_pusat_7 := helpers.DeleteRedis(Fieldcrm_home_redis + "_VALIDMAINTENANCE_" + strconv.Itoa(i) + "_" + search)

		log.Printf("Redis Delete BACKEND CRM NEW  : %d", val_pusat_1)
		log.Printf("Redis Delete BACKEND CRM PROCESS  : %d", val_pusat_2)
		log.Printf("Redis Delete BACKEND CRM VALID  : %d", val_pusat_3)
		log.Printf("Redis Delete BACKEND CRM INVALID  : %d", val_pusat_4)
		log.Printf("Redis Delete BACKEND CRM MAINTENANCE  : %d", val_pusat_5)
		log.Printf("Redis Delete BACKEND CRM FOLLOWUP  : %d", val_pusat_6)
		log.Printf("Redis Delete BACKEND CRM VALIDMAINTENANCE  : %d", val_pusat_7)
	}

	val_crmsales := helpers.DeleteRedis(Fieldcrmsales_home_redis + "_" + phone)
	log.Printf("Redis Delete BACKEND CRM SALES : %d", val_crmsales)

	val_crmsales2 := helpers.DeleteRedis(Fieldcrmsales_home_redis + "_" + phone + "_")
	log.Printf("Redis Delete BACKEND CRM SALES : %d", val_crmsales2)

	val_crmsales3 := helpers.DeleteRedis(Fieldcrmsales_home_redis + "_" + phone + "_MAINTENANCE")
	log.Printf("Redis Delete BACKEND CRM SALES MAINTENANCE : %d", val_crmsales3)

	val_client_sales := helpers.DeleteRedis(Fieldcrm_sales_redis + "_" + username + "_PROCESS")
	log.Printf("Redis Delete SALES CRM SALES : %d", val_client_sales)

	val_client_sales_maintenance := helpers.DeleteRedis(Fieldcrm_sales_redis + "_" + username + "_MAINTENANCE")
	log.Printf("Redis Delete SALES CRM MAINTENANCE : %d", val_client_sales_maintenance)
}
