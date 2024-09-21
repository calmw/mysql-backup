package backup

import (
	"log"
	"testing"
)

func TestEmail_Send(t *testing.T) {
	email := NewEmail("calm.fei@gmail.com", "dev server", "xxxxx", "smtp.gmail.com", 587)
	err := email.Send("calm.wang@hotmail.com", "db backup", "This is an email body.\nHere's more data.\n", []string{"./backup/test2_20240921145023.zip"})
	log.Println(err)
}
