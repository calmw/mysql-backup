package backup

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"path/filepath"
	"text/template"
)

type Email struct {
	Username string
	Password string
	Host     string
	Port     int
}

func NewEmail(username, password, host string, port int) *Email {
	return &Email{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

// Send
// To:      "recipient@example.com",
// From:    "sender@example.com",
// Subject: "Email with attachment",
// Body:    "This is the body of the email",
// Attach:  []string{"path/to/your/file.txt"},
func (e *Email) Send(to, from, subject, body string, attach []string) error {
	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)
	t := template.Must(template.New("email").Parse(`To: {{.to}}
From: {{.from}}
Subject: {{.subject}}
MIME-version: 1.0
Content-Type: multipart/mixed; boundary="nextpart"
 
--nextpart
Content-Type: text/plain; charset="UTF-8"
 
{{.body}}
--nextpart
`))

	var buffer bytes.Buffer
	if err := t.Execute(&buffer, e); err != nil {
		return err
	}

	msg := buffer.Bytes()
	for _, file := range attach {
		part, err := createPart(file)
		if err != nil {
			return err
		}
		msg = append(msg, part...)
	}

	msg = append(msg, []byte("--nextpart--")...)

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", e.Host, e.Port),
		auth,
		from,
		[]string{to},
		msg,
	)
}

func createPart(path string) ([]byte, error) {
	partBoundary := "--nextpart"
	filename := filepath.Base(path)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := new(bytes.Buffer)
	body.WriteString(partBoundary + "\n")
	body.WriteString("Content-Type: application/octet-stream; name=\"" + filename + "\"\n")
	body.WriteString("Content-Disposition: attachment; filename=\"" + filename + "\"\n")
	body.WriteString("Content-Transfer-Encoding: base64\n")
	body.WriteString("\n")

	// Copy the file data into the buffer, base64 encoding it
	// ...
	// body.WriteString("\n" + partBoundary + "--")

	return body.Bytes(), nil
}
