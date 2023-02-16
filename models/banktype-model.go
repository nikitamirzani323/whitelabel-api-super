package models

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/configs"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/db"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/entities"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/helpers"
	"github.com/nleeper/goment"
)

func Fetch_banktypeHome() (helpers.Response, error) {
	var obj entities.Model_banktype
	var arraobj []entities.Model_banktype
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			A.idbanktype , B.nmcatebank, A.nmbanktype,
			A.imgbanktype , A.statusbanktype,
			A.createbanktype, to_char(COALESCE(A.createdatebanktype,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			A.updatebanktype, to_char(COALESCE(A.updatedatebanktype,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + configs.DB_tbl_mst_banktype + ` as A   
			JOIN ` + configs.DB_tbl_mst_catebank + ` as B ON B.idcatebank = A.idcatebank   
			ORDER BY A.updatedatebanktype DESC   
	`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idbanktype_db, nmcatebank_db, nmbanktype_db, imgbanktype_db, statusbanktype_db     string
			createbanktype_db, createdatebanktype_db, updatebanktype_db, updatedatebanktype_db string
		)

		err = row.Scan(&idbanktype_db, &nmcatebank_db, &nmbanktype_db, &imgbanktype_db, &statusbanktype_db,
			&createbanktype_db, &createdatebanktype_db, &updatebanktype_db, &updatedatebanktype_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if createbanktype_db != "" {
			create = createbanktype_db + ", " + createdatebanktype_db
		}
		if updatebanktype_db != "" {
			update = updatebanktype_db + ", " + updatedatebanktype_db
		}

		obj.Banktype_id = idbanktype_db
		obj.Banktype_nmcatebank = nmcatebank_db
		obj.Banktype_img = imgbanktype_db
		obj.Banktype_status = statusbanktype_db
		obj.Banktype_create = create
		obj.Banktype_update = update
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
func Save_banktype(admin, idrecord, name, img, status, sData string, idcatebank int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(configs.DB_tbl_mst_banktype, "idbanktype", idrecord)
		if !flag {
			sql_insert := `
				insert into
				` + configs.DB_tbl_mst_banktype + ` (
					idbanktype , idcatebank, nmbanktype, imgbanktype, statusbanktype,  
					createbanktype, createdatebanktype 
				) values (
					$1, $2, $3, $4, $5
					$6, $7 
				)
			`
			flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_banktype, "INSERT",
				idrecord, idcatebank, name, img, status,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				msg = "Succes"
			} else {
				log.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_mst_banktype + `  
				SET idcatebank=$1, nmbanktype=$2, imgbanktype=$3 , statusbanktype=$4
				updatebanktype=$5, updatedatebanktype=$6 
				WHERE idbanktype =$7  
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_banktype, "UPDATE",
			idcatebank, name, img, status, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
			flag = true
			msg = "Succes"
		} else {
			log.Println(msg_update)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
