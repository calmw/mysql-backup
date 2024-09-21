package backup

import (
	"fmt"
	"testing"
)

func TestDump_Dump(t *testing.T) {
	backup := NewBackup("127.0.0.1", "3306", "root", "root", "test2", "./backup", 2)
	err, sqlFile := backup.Dump()
	fmt.Println(err, sqlFile)
}
