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

const Fieldemployee_home_redis = "LISTEMPLOYEE_BACKEND_ISBPANEL"
const Fieldemployee_frontend_redis = "LISTEMPLOYEE_FRONTEND_ISBPANEL"

func Employeehome(c *fiber.Ctx) error {
	var obj entities.Model_employee
	var arraobj []entities.Model_employee
	var objdepart entities.Model_listdepartement
	var arraobjdepart []entities.Model_listdepartement
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldemployee_home_redis)
	jsonredis := []byte(resultredis)
	listdepartement_RD, _, _, _ := jsonparser.Get(jsonredis, "listdepartement")
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		employee_username, _ := jsonparser.GetString(value, "employee_username")
		employee_iddepart, _ := jsonparser.GetString(value, "employee_iddepart")
		employee_nmdepart, _ := jsonparser.GetString(value, "employee_nmdepart")
		employee_name, _ := jsonparser.GetString(value, "employee_name")
		employee_phone, _ := jsonparser.GetString(value, "employee_phone")
		employee_status, _ := jsonparser.GetString(value, "employee_status")
		employee_create, _ := jsonparser.GetString(value, "employee_create")
		employee_update, _ := jsonparser.GetString(value, "employee_update")

		obj.Employee_username = employee_username
		obj.Employee_iddepart = employee_iddepart
		obj.Employee_nmdepart = employee_nmdepart
		obj.Employee_name = employee_name
		obj.Employee_phone = employee_phone
		obj.Employee_status = employee_status
		obj.Employee_create = employee_create
		obj.Employee_update = employee_update
		arraobj = append(arraobj, obj)
	})
	jsonparser.ArrayEach(listdepartement_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		departement_id, _ := jsonparser.GetString(value, "departement_id")
		departement_name, _ := jsonparser.GetString(value, "departement_name")

		objdepart.Departement_id = departement_id
		objdepart.Departement_name = departement_name
		arraobjdepart = append(arraobjdepart, objdepart)
	})
	if !flag {
		result, err := models.Fetch_employeeHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldemployee_home_redis, result, 60*time.Minute)
		log.Println("EMPLOYEE MYSQL")
		return c.JSON(result)
	} else {
		log.Println("EMPLOYEE CACHE")
		return c.JSON(fiber.Map{
			"status":          fiber.StatusOK,
			"message":         "Success",
			"listdepartement": arraobjdepart,
			"record":          arraobj,
			"time":            time.Since(render_page).String(),
		})
	}
}
func EmployeeByDepart(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_employeebydepart)
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

	var obj entities.Model_employeebydepart
	var arraobj []entities.Model_employeebydepart
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fieldemployee_home_redis + "_" + client.Employee_iddepart)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		employee_username, _ := jsonparser.GetString(value, "employee_username")
		employee_name, _ := jsonparser.GetString(value, "employee_name")
		employee_deposit, _ := jsonparser.GetInt(value, "employee_deposit")
		employee_noanswer, _ := jsonparser.GetInt(value, "employee_noanswer")
		employee_reject, _ := jsonparser.GetInt(value, "employee_reject")
		employee_invalid, _ := jsonparser.GetInt(value, "employee_invalid")

		obj.Employee_username = employee_username
		obj.Employee_name = employee_name
		obj.Employee_deposit = int(employee_deposit)
		obj.Employee_noanswer = int(employee_noanswer)
		obj.Employee_reject = int(employee_reject)
		obj.Employee_invalid = int(employee_invalid)
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_employeeByDepartement(client.Employee_iddepart)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldemployee_home_redis+"_"+client.Employee_iddepart, result, 60*time.Minute)
		log.Println("EMPLOYEE BY DEPARTEMENT MYSQL")
		return c.JSON(result)
	} else {
		log.Println("EMPLOYEE BY DEPARTEMENT CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func EmployeeBySalesPerformance(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_employeebysalesperform)
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

	var obj entities.Model_employeebysalesperform
	var arraobj []entities.Model_employeebysalesperform
	render_page := time.Now()
	string_redis := ""
	if client.Employee_startdate != "" {
		string_redis = Fieldemployee_home_redis + "_" + client.Employee_iddepart + "_" + client.Employee_username + "_" + client.Employee_startdate + "_" + client.Employee_enddate
	} else {
		string_redis = Fieldemployee_home_redis + "_" + client.Employee_iddepart + "_" + client.Employee_username
	}
	resultredis, flag := helpers.GetRedis(string_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		sales_deposit, _ := jsonparser.GetInt(value, "sales_deposit")
		sales_depositsum, _ := jsonparser.GetFloat(value, "sales_depositsum")
		sales_noanswer, _ := jsonparser.GetInt(value, "sales_noanswer")
		sales_reject, _ := jsonparser.GetInt(value, "sales_reject")
		sales_invalid, _ := jsonparser.GetInt(value, "sales_invalid")
		listdeposit_RD, _, _, _ := jsonparser.Get(value, "sales_listdeposit")
		listnoanswer_RD, _, _, _ := jsonparser.Get(value, "sales_listnoanswer")
		listinvalid_RD, _, _, _ := jsonparser.Get(value, "sales_listinvalid")

		var obj_listdeposit entities.Model_crmmemberlistdeposit
		var arraobj_listdeposit []entities.Model_crmmemberlistdeposit
		jsonparser.ArrayEach(listdeposit_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			crmdeposit_phone, _ := jsonparser.GetString(value, "crmdeposit_phone")
			crmdeposit_nama, _ := jsonparser.GetString(value, "crmdeposit_nama")
			crmdeposit_source, _ := jsonparser.GetString(value, "crmdeposit_source")
			crmdeposit_nmwebagen, _ := jsonparser.GetString(value, "crmdeposit_nmwebagen")
			crmdeposit_deposit, _ := jsonparser.GetFloat(value, "crmdeposit_deposit")
			crmdeposit_iduseragen, _ := jsonparser.GetString(value, "crmdeposit_iduseragen")
			crmdeposit_update, _ := jsonparser.GetString(value, "crmdeposit_update")

			obj_listdeposit.Crmsdeposit_phone = crmdeposit_phone
			obj_listdeposit.Crmsdeposit_nama = crmdeposit_nama
			obj_listdeposit.Crmsdeposit_source = crmdeposit_source
			obj_listdeposit.Crmsdeposit_nmwebagen = crmdeposit_nmwebagen
			obj_listdeposit.Crmsdeposit_deposit = float32(crmdeposit_deposit)
			obj_listdeposit.Crmsdeposit_iduseragen = crmdeposit_iduseragen
			obj_listdeposit.Crmsdeposit_update = crmdeposit_update
			arraobj_listdeposit = append(arraobj_listdeposit, obj_listdeposit)
		})

		var obj_listnoanswer entities.Model_crmmemberlistnoanswer
		var arraobj_listnoanswer []entities.Model_crmmemberlistnoanswer
		jsonparser.ArrayEach(listnoanswer_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			crmnoanswer_phone, _ := jsonparser.GetString(value, "crmnoanswer_phone")
			crmnoanswer_nama, _ := jsonparser.GetString(value, "crmnoanswer_nama")
			crmnoanswer_source, _ := jsonparser.GetString(value, "crmnoanswer_source")
			crmnoanswer_tipe, _ := jsonparser.GetString(value, "crmnoanswer_tipe")
			crmnoanswer_note, _ := jsonparser.GetString(value, "crmnoanswer_note")
			crmnoanswer_update, _ := jsonparser.GetString(value, "crmnoanswer_update")

			obj_listnoanswer.Crmnoanswer_phone = crmnoanswer_phone
			obj_listnoanswer.Crmnoanswer_nama = crmnoanswer_nama
			obj_listnoanswer.Crmnoanswer_source = crmnoanswer_source
			obj_listnoanswer.Crmnoanswer_tipe = crmnoanswer_tipe
			obj_listnoanswer.Crmnoanswer_note = crmnoanswer_note
			obj_listnoanswer.Crmnoanswer_update = crmnoanswer_update
			arraobj_listnoanswer = append(arraobj_listnoanswer, obj_listnoanswer)
		})

		var obj_listinvalid entities.Model_crmmemberlistinvalid
		var arraobj_listinvalid []entities.Model_crmmemberlistinvalid
		jsonparser.ArrayEach(listinvalid_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			crminvalid_phone, _ := jsonparser.GetString(value, "crminvalid_phone")
			crminvalid_nama, _ := jsonparser.GetString(value, "crminvalid_nama")
			crminvalid_source, _ := jsonparser.GetString(value, "crminvalid_source")
			crminvalid_update, _ := jsonparser.GetString(value, "crminvalid_update")

			obj_listinvalid.Crminvalid_phone = crminvalid_phone
			obj_listinvalid.Crminvalid_nama = crminvalid_nama
			obj_listinvalid.Crminvalid_source = crminvalid_source
			obj_listinvalid.Crminvalid_update = crminvalid_update
			arraobj_listinvalid = append(arraobj_listinvalid, obj_listinvalid)
		})

		obj.Sales_deposit = int(sales_deposit)
		obj.Sales_depositsum = float32(sales_depositsum)
		obj.Sales_noanswer = int(sales_noanswer)
		obj.Sales_reject = int(sales_reject)
		obj.Sales_invalid = int(sales_invalid)
		obj.Sales_listdeposit = arraobj_listdeposit
		obj.Sales_listnoanswer = arraobj_listnoanswer
		obj.Sales_listinvalid = arraobj_listinvalid
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_employeeBySalesPerformance(client.Employee_iddepart, client.Employee_username, client.Employee_startdate, client.Employee_enddate)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(string_redis, result, 30*time.Minute)
		log.Println("EMPLOYEE BY SALES MYSQL")
		return c.JSON(result)
	} else {
		log.Println("EMPLOYEE BY SALES CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func EmployeeSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_employeesave)
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

	// admin, username, password, iddepart, name, phone, status, sData, idrecord string
	result, err := models.Save_employee(
		client_admin,
		client.Employee_username, client.Employee_password, client.Employee_iddepart, client.Employee_name,
		client.Employee_phone, client.Employee_status, client.Sdata, client.Employee_username)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	_deleteredis_employee()
	return c.JSON(result)
}
func _deleteredis_employee() {
	val_master := helpers.DeleteRedis(Fieldemployee_home_redis)
	log.Printf("Redis Delete BACKEND EMPLOYEE : %d", val_master)

	//CLIENT
	val_client := helpers.DeleteRedis(Fieldemployee_frontend_redis)
	log.Printf("Redis Delete FRONTEND EMPLOYEE : %d", val_client)

	val_master_departement := helpers.DeleteRedis(Fielddepartement_home_redis)
	log.Printf("Redis Delete BACKEND DEPARTEMENT : %d", val_master_departement)

	//CLIENT
	val_client_departement := helpers.DeleteRedis(Fielddepartement_frontend_redis)
	log.Printf("Redis Delete FRONTEND DEPARTEMENT : %d", val_client_departement)
}
