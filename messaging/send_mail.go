package messaging

import (
	"bytes"
	"html/template"
	"net/smtp"
)

var Auth smtp.Auth

/*func main() {
	auth = smtp.PlainAuth("", "dhanush@geektrust.in", "password", "smtp.gmail.com")
	templateData := struct {
		Name string
		URL  string
	}{
		Name: "Dhanush",
		URL:  "http://geektrust.in",
	}
	r := NewRequest([]string{"junk@junk.com"}, "Hello Junk!", "Hello, World!")
	err := r.ParseTemplate("template.html", templateData)
	if err != nil {
		ok, _ := r.SendEmail()
		fmt.Println(ok)
	}
}*/

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	Body    string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		Body:    body,
	}
}

func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.Body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, Auth, r.from, r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.Body = buf.String()
	return nil
}
