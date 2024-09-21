package backup

import (
	"fmt"
	"testing"
)

func TestBackup_Dump(t *testing.T) {
	backup := NewBackup("127.0.0.1", "3306", "root", "root", "test_db", "./backup", 2)
	err, zipFile := backup.Dump()
	fmt.Println(err, zipFile)
}
