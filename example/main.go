package main

import (
	backup "github.com/calmw/mysql-backup"
	"log"
)

func main() {
	bp := backup.NewBackup("127.0.0.1", "3306", "root", "root", "test_db", "./backup", 2)
	err, zipFile := bp.Dump()
	log.Println(zipFile, err)

	email := backup.NewEmail("calm.fei@gmail.com", "dev server", "xx", "smtp.gmail.com", 587)
	err = email.Send("calm.wang@hotmail.com", "db backup", "This is an email body.\nHere's more data.\n", []string{zipFile})
	log.Println(err)
}
