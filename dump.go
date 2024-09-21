package dump

import (
	"io"
	"os"
	"os/exec"
	"time"
)

/**
 *
 * 备份MySql数据库
 * @param host:         数据库地址:
 * @param port:         端口:
 * @param user:         用户名:
 * @param password:     密码:
 * @param databaseName: 需要备份的数据库名:
 * @param tableName:    需要备份的表名:
 * @param sqlPath:      备份SQL存储路径:
 * @return backupPath   返回备份路径
 *
 */

const (
	BackupPath = "/备份所在文件夹"
)

// BackupMySqlDb 备份
func BackupMySqlDb(host, port, user, password, databaseName, sqlPath string) (error, string) {
	//获得一个当前的时间戳
	now := time.Now().Format("20060102150405")
	var backupPath string
	// 判断文件夹不存在时自动创建
	if !FileExists(sqlPath) {
		if err := os.MkdirAll(BackupPath, os.ModePerm); err != nil {
			return err, ""
		}
	}
	//设置备份文件的路径
	backupPath = sqlPath + databaseName + "_" + now + ".sql"
	mysqldumpCmd := `mysqldump -h ` + host + ` -P ` + port + ` -u` + user + ` -p` + password + ` --databases ` + databaseName + ` --ignore-table=` + databaseName + `.logs` + ` >` + backupPath
	//--ignore-table=库名.表名 表示备份忽略该表
	if err := ExecutiveCommand(mysqldumpCmd); err != nil {
		return err, ""
	}
	return nil, backupPath
}

// RecoverMySqlDb 恢复数据表
func RecoverMySqlDb(host, port, user, password, databaseName, backupPath string) error {
	//恢复表 mysql -h[地址] -P[端口] -u[用户名] -p[密码] [数据库名] <[备份文件]
	mysqldumpCmd := `mysql -h` + host + ` -P` + port + ` -u` + user + ` -p` + password + ` ` + databaseName + ` <` + backupPath
	if err := ExecutiveCommand(mysqldumpCmd); err != nil {
		return err
	}
	return nil
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

// FileExists 检查文件或文件夹是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
