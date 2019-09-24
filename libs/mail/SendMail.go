package mail

func SendMail(toEmail string, Name string, code string) bool {
	subject := "Email verify code for forgot password"
	r := NewRequest([]string{toEmail}, subject)
	return r.Send("template/send_mail_reset_password.html", map[string]string{"name": Name, "code": code})
}