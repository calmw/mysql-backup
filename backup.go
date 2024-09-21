package backup

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

/// 备份mysql数据库

type Backup struct {
	Host             string // 数据库Host地址
	Port             string // 端口
	User             string // 用户名
	Password         string // 密码
	DatabaseName     string // 需要备份的数据库名
	BackupPath       string // 备份路径
	BackupFileNumber int    // 保留备份文件的数量。根据时间，保留最近的文件，删除其他历史备份文件。 0 全部保留
}

func NewBackup(host, port, user, password, databaseName, backupPath string, backupFileNumber int) *Backup {
	return &Backup{
		Host:             host,
		Port:             port,
		User:             user,
		Password:         password,
		DatabaseName:     databaseName,
		BackupPath:       backupPath,
		BackupFileNumber: backupFileNumber,
	}
}

// Dump 备份
func (d Backup) Dump() (error, string) {
	//获得一个当前的时间戳
	now := time.Now().Format("20060102150405")
	// 判断文件夹不存在时自动创建
	if !fileExists(d.BackupPath) {
		if err := os.MkdirAll(d.BackupPath, os.ModePerm); err != nil {
			return err, ""
		}
	}
	//设置备份文件
	fileName := d.DatabaseName + "_" + now
	backupPath := strings.TrimRight(d.BackupPath, "/")
	backupFile := backupPath + "/" + fileName + ".sql"
	//mysqldumpCmd := `mysqldump -h ` + d.Host + ` -P ` + d.Port + ` -u` + d.User + ` -p` + d.Password + ` --databases ` + d.DatabaseName + ` --ignore-table=` + d.DatabaseName + `.logs` + ` >` + backupFile
	//--ignore-table=库名.表名 表示备份忽略该表
	mysqldumpCmd := `mysqldump -h ` + d.Host + ` -P ` + d.Port + ` -u` + d.User + ` -p` + d.Password + ` --databases ` + d.DatabaseName + ` >` + backupFile
	if err := ExecutiveCommand(mysqldumpCmd); err != nil {
		return err, ""
	}
	zipCmd := fmt.Sprintf(`cd %s && zip %s.zip %s.sql`, backupPath, fileName, fileName)
	if err := ExecutiveCommand(zipCmd); err != nil {
		return err, ""
	}
	err := os.Remove(backupFile)
	if err != nil {
		return err, ""
	}
	if err = deleteExpiredFiles(backupPath, d.BackupFileNumber); err != nil {
		return err, ""
	}

	return nil, backupFile
}

func ExecutiveCommand(command string) error {
	//需要执行命令:command
	cmd := exec.Command("/bin/bash", "-c", command)
	// 获取管道输入
	output, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	_, err = io.ReadAll(output)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

// fileExists 检查文件或文件夹是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// 删除过期文件
func deleteExpiredFiles(directory string, maxNum int) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Printf("Error reading directory: %s\n", err)
		return err
	}

	var fileMap = make(map[int64]os.DirEntry)
	var fileModTimeSli Int64Slice

	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			log.Printf("Error getting file info: %s\n", err)
			continue
		}

		fileMap[fileInfo.ModTime().Unix()] = file
		fileModTimeSli = append(fileModTimeSli, fileInfo.ModTime().Unix())
	}

	sort.Sort(fileModTimeSli)

	count := fileModTimeSli.Len()
	if count <= maxNum {
		return nil
	}
	directory = strings.TrimRight(directory, "/")

	for i, k := range fileModTimeSli {
		if i < (count - maxNum) {
			err = os.Remove(directory + "/" + fileMap[k].Name())
			if err != nil {
				log.Printf("delete file error: %s\n", err)
				return err
			}
		}
	}
	return nil
}
