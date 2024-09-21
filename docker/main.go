package main

import (
	"docker/task"
	"github.com/jasonlvhit/gocron"
	"os"
	"strconv"
)

func main() {
	backupIntervalStr := os.Getenv("BACKUP_INTERVAL")
	sendEmailStr := os.Getenv("SEND_EMAIL")
	backupInterval, err := strconv.Atoi(backupIntervalStr)
	if err != nil {
		backupInterval = 86400
	}
	var sendEmail bool
	if sendEmailStr == "true" {
		sendEmail = true
	}

	s := gocron.NewScheduler()
	_ = s.Every(uint64(backupInterval)).Seconds().From(gocron.NextTick()).Do(task.BackupMysql, sendEmail)
	<-s.Start()
}
