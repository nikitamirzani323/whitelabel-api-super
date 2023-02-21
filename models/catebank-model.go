package models

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/configs"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/db"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/entities"
	"github.com/nikitamirzani323/whitelabel/whitelabel_api_super/helpers"
	"github.com/nleeper/goment"
)

func Fetch_catebankHome() (helpers.Response, error) {
	var obj entities.Model_catebank
	var arraobj []entities.Model_catebank
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			idcatebank , nmcatebank,statuscatebank,
			createcatebank, to_char(COALESCE(createdatecatebank,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			updatecatebank, to_char(COALESCE(updatedatecatebank,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + configs.DB_tbl_mst_catebank + `  
			ORDER BY updatedatecatebank DESC   
	`

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idcatebank_db                                                                      int
			nmcatebank_db, statuscatebank_db                                                   string
			createcatebank_db, createdatecatebank_db, updatecatebank_db, updatedatecatebank_db string
		)

		err = row.Scan(&idcatebank_db, &nmcatebank_db, &statuscatebank_db,
			&createcatebank_db, &createdatecatebank_db, &updatecatebank_db, &updatedatecatebank_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if createcatebank_db != "" {
			create = createcatebank_db + ", " + createdatecatebank_db
		}
		if updatecatebank_db != "" {
			update = updatecatebank_db + ", " + updatedatecatebank_db
		}

		obj.Catebank_id = idcatebank_db
		obj.Catebank_name = nmcatebank_db
		obj.Catebank_status = statuscatebank_db
		obj.Catebank_create = create
		obj.Catebank_update = update
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
func Save_catebank(admin, name, status, sData string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()

	if sData == "New" {
		sql_insert := `
				insert into
				` + configs.DB_tbl_mst_catebank + ` (
					idcatebank , nmcatebank,statuscatebank, 
					createcatebank, createdatecatebank
				) values (
					$1, $2, $3, 
					$4, $5
				)
			`
		field_column := configs.DB_tbl_mst_catebank + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, configs.DB_tbl_mst_catebank, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter), name, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
		} else {
			fmt.Println(msg_insert)
		}
	} else {
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_mst_catebank + `  
				SET nmcatebank =$1, statuscatebank=$2, 
				updatecatebank=$3, updatedatecatebank=$4 
				WHERE idcatebank =$5 
			`

		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_mst_catebank, "UPDATE",
			name, status, admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
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
