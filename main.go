package main

/*
	Нагальні задачі:
	 - Можливо заюзати struct Email в файлі subscribe.go ?
	 - Запустити через Dockerfile
	 - Знайти варіант як замінити використання ключа пристрою адміна

	Варіанти апгрейду програми:
	 - Робота з БД, а не файлами
	 - Можливість відписки email
*/

import (
	"log"
	"net/http"
)

// HandleRequest приймає і оброблює запити типу /rate, /subscribe, /sendEmails
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/rate" {
		if r.Method == http.MethodGet {
			HandleRateRequest(w)
		} else {
			http.Error(w, "Incorrect method", http.StatusMethodNotAllowed)
		}
	} else if r.URL.Path == "/subscribe" {
		if r.Method == http.MethodPost {
			HandleSubscribeRequest(w, r)
		} else {
			http.Error(w, "Incorrect method", http.StatusMethodNotAllowed)
		}
	} else if r.URL.Path == "/sendEmails" {
		if r.Method == http.MethodPost {
			HandleSendEmails(w)
		} else {
			http.Error(w, "Incorrect method", http.StatusMethodNotAllowed)
		}
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	http.HandleFunc("/", HandleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
