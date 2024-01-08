package main

import (
	"bufio"
	"fmt"
	"io"

	"net/http"
	"os"
	"strings"
)

// файл для зберігання emails для розсилки
const subscriptionFile = "subscriptions.txt"

// HandleSubscribeRequest оброблює запит /subscribe,
func HandleSubscribeRequest(w http.ResponseWriter, r *http.Request) {
	// Зчитування запиту на додавання email
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Використання struct Email тут (функція HandleSubscribeRequest_) - upgrade variant

	// Перевірка на повторення email
	exists, err := CheckEmailExists(string(body))
	if err != nil || exists {
		fmt.Println("Exists")
		http.Error(w, "Email already subscribed", http.StatusBadRequest) // some changes here
		return
	}

	// Збереження email в файл
	err = SaveEmailToFile(string(body))
	if err != nil {
		http.Error(w, "Failed to save email", http.StatusInternalServerError)
		return
	}

	// Надсилання повідомлення про успішну підписку
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Email subscribed successfully")
}

// CheckEmailExists перевіряє, чи існує електронна адреса в файлі підписок.
func CheckEmailExists(email string) (bool, error) {
	file, err := os.Open(subscriptionFile)
	if err != nil {
		return false, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Читання файлу рядок за рядком
	reader := bufio.NewScanner(file)
	for reader.Scan() {
		// Перевірка на співпадіння електронної адреси
		if strings.TrimSpace(reader.Text()) == strings.TrimSpace(email) {
			return true, nil
		}
	}

	// Перевірка на помилку читання файлу
	if err := reader.Err(); err != nil {
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

	// запис нового рядка
	_, err = file.WriteString(email + "\n")
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	return nil
}

/*
type Email struct {
	Address string `json:"email"`
}

func HandleSubscribeRequest_(w http.ResponseWriter, r *http.Request) {
	// Парсинг тіла запиту
	var email Email
	err = json.Unmarshal(body, &email)
	//fmt.Println(email.Address) // test
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Перевірка на повторення email
	exists, err := CheckEmailExists(email.Address)

	// Збереження email в файл
	err = SaveEmailToFile(email.Address)
}
*/
