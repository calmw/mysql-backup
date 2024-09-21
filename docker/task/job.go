package task

import (
	backup "github.com/calmw/mysql-backup"
	"log"
	"os"
	"strconv"
)

func BackupMysql(send bool) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	backupPath := os.Getenv("BACKUP_PATH")
	backupFileNumberStr := os.Getenv("BACKUP_FILE_NUM")
	backupFileNumber, err := strconv.Atoi(backupFileNumberStr)
	if err != nil {
		backupFileNumber = 2
	}

	bp := backup.NewBackup(dbHost, dbPort, dbUser, dbPassword, dbName, backupPath, backupFileNumber)
	err, zipFile := bp.Dump()
	log.Printf("备份结果：%v, %s \n", err, zipFile)

	if !send {
		return
	}

	senderEmail := os.Getenv("SENDER_EMAIL")
	senderUsername := os.Getenv("SENDER_USERNAME")
	senderPassword := os.Getenv("SENDER_PASSWORD")
	smtpHost := os.Getenv("EMAIL_HOST")
	smtpPortStr := os.Getenv("EMAIL_PORT")
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		smtpPort = 587
	}

	to := os.Getenv("TO")
	subject := os.Getenv("SUBJECT")
	body := os.Getenv("BODY")

	email := backup.NewEmail(senderEmail, senderUsername, senderPassword, smtpHost, smtpPort)
	err = email.Send(to, subject, body, []string{zipFile})
	log.Printf("发送结果：%v \n", err)
}
