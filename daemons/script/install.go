package script

import (
	"fmt"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/util"
)
import "database/sql"

func Install() {
	// create database
	func() {
		dbHost := util.GetEnv("MYSQL_HOST", "127.0.0.1")
		dbUser := util.GetEnv("MYSQL_USER", "root")
		dbPass := util.GetEnv("MYSQL_PASS", "")
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, "")
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE `%s` DEFAULT CHARACTER SET = `utf8mb4`", util.NetworkNode))
		if err != nil {
			panic(err)
		}
		fmt.Println("Create database ", util.NetworkNode, "success!!!")
	}()

	// migration
	d := dao.New()
	d.Migration()
	defer d.Close()

	// nginx
}
