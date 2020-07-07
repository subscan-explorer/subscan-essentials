package script

import (
	"database/sql"
	"fmt"
	"github.com/itering/subscan/util"
	"io"
	"os"
)

func Install(conf string) {
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
		defer func() {
			_ = db.Close()
		}()
		q := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET = `utf8mb4`", util.NetworkNode)
		_, err = db.Exec(q)
		if err != nil {
			panic(err)
		}
		fmt.Println("Create database", util.NetworkNode, "Success!!!")

	}()

	// conf
	_ = fileCopy(fmt.Sprintf("%s/http.toml.example", conf), fmt.Sprintf("%s/http.toml", conf))
	_ = fileCopy(fmt.Sprintf("%s/mysql.toml.example", conf), fmt.Sprintf("%s/mysql.toml", conf))
	_ = fileCopy(fmt.Sprintf("%s/redis.toml.example", conf), fmt.Sprintf("%s/redis.toml", conf))

}

func fileCopy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
