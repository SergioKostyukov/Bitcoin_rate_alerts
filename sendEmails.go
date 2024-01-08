package main

import (
	"bufio"
	"fmt"
	"go-gomail/gomail"
	"net/http"
	"os"
)

// HandleSendEmails оброблює запит /sendEmail,
// виконує оновлення курсу та викликає функцію SendEmail для розсилки на всі підписані emails з файлу subscriptionFile
func HandleSendEmails(w http.ResponseWriter) {
	subject := "Bitcoin Price Update"

	// Отримання актуальної інформації
	var price string = HandleRateRequest(w)
	body := fmt.Sprintf("The current Bitcoin price in UAH is %s", price)
	//fmt.Println("Current Bitcoin price in UAH:", price) // дюблювання інформації в консоль

	file, err := os.Open(subscriptionFile)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	// Читання файлу рядок за рядком
	reader := bufio.NewScanner(file)

	// Читання рядків з файлу
	for reader.Scan() {
		line := reader.Text()
		SendEmail(line, subject, body)
		//fmt.Println("Email to %s sent successfully!", line) // вивід в консоль інформації про успішне надсилання повідомлення на пошту
	}

	if err := reader.Err(); err != nil {
		fmt.Println("Error while reading file:", err)
		return
	}
	//fmt.Println("Email sent successfully!") // розсилка листів є успішною
}

// SendEmail власне реалізує надсилання листа
// (to - отримувач, subject - тума листа, body - тіло листа з інформацією)
func SendEmail(to, subject, body string) error {
	from := "kostyukov.sergey2003@gmail.com" // пошта відправлення
	password := "asjkkfrqahkinbpc"           // ключ локального пристрою розробника
	smtpHost := "smtp.gmail.com"             // безкоштовний хостинг для надсилання листів
	smtpPort := 587

	// створення email
	email := gomail.NewMessage()
	email.SetHeader("From", from)
	email.SetHeader("To", to)
	email.SetHeader("Subject", subject)
	email.SetBody("text/plain", body)

	// надислання email
	dialer := gomail.NewDialer(smtpHost, smtpPort, from, password)

	if err := dialer.DialAndSend(email); err != nil {
		return err
	}

	return nil
}
