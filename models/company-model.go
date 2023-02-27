package models

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/configs"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/db"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/entities"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/helpers"
	"github.com/nleeper/goment"
)

func Fetch_companyHome(idcatebank int) (helpers.Response, error) {
	var obj entities.Model_company
	var arraobj []entities.Model_company
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idcompany , idcurr, 
			to_char(COALESCE(startjoincompany,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			to_char(COALESCE(endjoincompany,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			nmcompany, nmowner, phoneowner, emailowner, companyurl, statuscompany, 
			createcompany, to_char(COALESCE(createdatecompany,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			updatecompany, to_char(COALESCE(updatedatecompany,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + configs.DB_tbl_mst_company + ` 
			ORDER BY updatedatecompany DESC   
	`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idcompany_db, idcurr_db, startjoincompany_db, endjoincompany_db                         string
			nmcompany_db, nmowner_db, phoneowner_db, emailowner_db, companyurl_db, statuscompany_db string
			createcompany_db, createdatecompany_db, updatecompany_db, updatedatecompany_db          string
		)

		err = row.Scan(&idcompany_db, &idcurr_db, &startjoincompany_db, &endjoincompany_db, &nmcompany_db,
			&nmowner_db, &phoneowner_db, &emailowner_db, &companyurl_db, &statuscompany_db,
			&createcompany_db, &createdatecompany_db, &updatecompany_db, &updatedatecompany_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if createcompany_db != "" {
			create = createcompany_db + ", " + createdatecompany_db
		}
		if updatecompany_db != "" {
			update = updatecompany_db + ", " + updatedatecompany_db
		}

		obj.Company_id = idcompany_db
		obj.Company_idcurr = idcurr_db
		obj.Company_start = startjoincompany_db
		obj.Company_end = endjoincompany_db
		obj.Company_name = nmcompany_db
		obj.Company_phone = phoneowner_db
		obj.Company_owner = nmowner_db
		obj.Company_email = emailowner_db
		obj.Company_companyurl = companyurl_db
		obj.Company_create = create
		obj.Company_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_company(admin, idrecord, name, img, status, sData string, idcatebank int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_mst_company, "idbanktype", idrecord)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_mst_company + ` (
					idbanktype , idcatebank, nmbanktype, imgbanktype, statusbanktype,  
					createbanktype, createdatebanktype 
				) values (
					$1, $2, $3, $4, $5, 
					$6, $7 
				)
			`
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_company, "INSERT",
				idrecord, idcatebank, name, img, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				msg = "Succes"
			} else {
				fmt.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_mst_company + `  
				SET nmbanktype=$1, imgbanktype=$2, statusbanktype=$3, 
				updatebanktype=$4, updatedatebanktype=$5  
				WHERE idbanktype =$6   
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_company, "UPDATE",
			name, img, status, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			flag = true
			msg = "Succes"
		} else {
			fmt.Println(msg_update)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
