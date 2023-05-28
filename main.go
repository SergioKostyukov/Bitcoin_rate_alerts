package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go-gomail/gomail"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"net/http"
)

type CoinGeckoResponse struct {
	MarketData struct {
		CurrentPrice struct {
			UAH float64 `json:"uah"`
		} `json:"current_price"`
	} `json:"market_data"`
}

type Email struct {
	Address string `json:"email"`
}

func handleRateRequest(w http.ResponseWriter) string {
	// Отримання актуального курсу
	price, err := http.Get("https://api.coingecko.com/api/v3/coins/bitcoin")
	if err != nil {
		return ""
	}
	defer price.Body.Close()

	var data CoinGeckoResponse
	err = json.NewDecoder(price.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	// Серіалізація структури у формат JSON
	jsonData, err := json.Marshal(data.MarketData.CurrentPrice.UAH)

	// Конвертація рядка JSON в рядок string
	jsonStr := string(jsonData)

	// Виводимо результат на localhost
	fmt.Fprintf(w, "The current Bitcoin price in UAH is $%s", jsonStr)

	return jsonStr
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/rate" {
		if r.Method == http.MethodGet {
			handleRateRequest(w)
		} else {
			http.Error(w, "Метод не підтримується", http.StatusMethodNotAllowed)
		}
	} else if r.URL.Path == "/subscribe" {
		if r.Method == http.MethodPost {
			handleSubscribeRequest(w, r)
		} else {
			http.Error(w, "Метод не підтримується", http.StatusMethodNotAllowed)
		}
	} else if r.URL.Path == "/sendEmails" {
		if r.Method == http.MethodPost {
			handleSendEmails(w)
		} else {
			http.Error(w, "Метод не підтримується", http.StatusMethodNotAllowed)
		}
	} else {
		http.NotFound(w, r)
	}
}

func handleSendEmails(w http.ResponseWriter) {
	fmt.Fprintf(w, "Емейли розіслано")
	// Implement your sendEmails logic here
	subject := "Bitcoin Price Update"

	var price string = handleRateRequest(w)
	fmt.Println("Current Bitcoin price in UAH:", price)

	body := fmt.Sprintf("The current Bitcoin price in UAH is %s", price)

	file, err := os.Open(subscriptionFile)
	if err != nil {
		fmt.Println("Помилка відкриття файлу:", err)
		return
	}
	defer file.Close()

	// Створення нового читача файлу
	reader := bufio.NewScanner(file)

	// Читання рядків з файлу
	for reader.Scan() {
		line := reader.Text()
		sendEmail(line, subject, body)
		//fmt.Println("Email to %s sent successfully!", line)
	}

	if err := reader.Err(); err != nil {
		fmt.Println("Помилка читання файлу:", err)
		return
	}
	fmt.Println("Email sent successfully!")
}

func sendEmail(to, subject, body string) error {
	from := "kostyukov.sergey2003@gmail.com"
	password := "asjkkfrqahkinbpc"
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	email := gomail.NewMessage()
	email.SetHeader("From", from)
	email.SetHeader("To", to)
	email.SetHeader("Subject", subject)
	email.SetBody("text/plain", body)

	dialer := gomail.NewDialer(smtpHost, smtpPort, from, password)

	if err := dialer.DialAndSend(email); err != nil {
		return err
	}

	return nil
}

const subscriptionFile = "subscriptions.txt"

func handleSubscribeRequest(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Parse the request body
	var email Email
	err = json.Unmarshal(body, &email)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check for duplicate email
	exists, err := CheckEmailExists(email.Address)
	if err != nil && exists {
		http.Error(w, "Email already subscribed", http.StatusBadRequest)
		return
	}

	// Save the email to file
	err = SaveEmailToFile(email.Address)
	if err != nil {
		http.Error(w, "Failed to save email", http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Email subscribed successfully")
}

// CheckEmailExists перевіряє, чи існує електронна адреса в файлі підписок.
func CheckEmailExists(email string) (bool, error) {
	filePath := subscriptionFile

	// Відкриття файлу
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Читання файлу рядок за рядком
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Перевірка на співпадіння електронної адреси
		if strings.TrimSpace(scanner.Text()) == strings.TrimSpace(email) {
			return true, nil
		}
	}

	// Перевірка на помилку читання файлу
	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("error while reading file: %v", err)
	}

	return false, nil
}

// SaveEmailToFile зберігає електронну адресу в файлі підписок.
func SaveEmailToFile(email string) error {
	file, err := os.OpenFile(subscriptionFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(email + "\n")
	if err != nil {
		return err
	}

	return nil
}

func main() {
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
