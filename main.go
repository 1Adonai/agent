package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db    *sql.DB
	store = sessions.NewCookieStore([]byte("super-secret-key"))
)

type ContactForm struct {
    Name         string
    ContactType  string
    ContactInfo  string
    SelectOption string
    Message      string
}

func main() {
	initDB()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/osago", osagoHandler)
	http.HandleFunc("/kasko", kaskoHandler)
	http.HandleFunc("/dom", domHandler)
	http.HandleFunc("/house", houseHandler)
	http.HandleFunc("/contact", contactHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Server starting at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	message, ok := session.Values["message"].(string)
	if ok {
		delete(session.Values, "message")
		session.Save(r, w)
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	data := struct {
		Message string
	}{
		Message: message,
	}
	tmpl.Execute(w, data)
}

func osagoHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	message, ok := session.Values["message"].(string)
	if ok {
		delete(session.Values, "message")
		session.Save(r, w)
	}

	tmpl := template.Must(template.ParseFiles("templates/OSAGO.html"))
	data := struct {
		Message string
	}{
		Message: message,
	}
	tmpl.Execute(w, data)
}

func kaskoHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	message, ok := session.Values["message"].(string)
	if ok {
		delete(session.Values, "message")
		session.Save(r, w)
	}

	tmpl := template.Must(template.ParseFiles("templates/KASKO.html"))
	data := struct {
		Message string
	}{
		Message: message,
	}
	tmpl.Execute(w, data)
}

func houseHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	message, ok := session.Values["message"].(string)
	if ok {
		delete(session.Values, "message")
		session.Save(r, w)
	}

	tmpl := template.Must(template.ParseFiles("templates/HOUSE.html"))
	data := struct {
		Message string
	}{
		Message: message,
	}
	tmpl.Execute(w, data)
}

func domHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	message, ok := session.Values["message"].(string)
	if ok {
		delete(session.Values, "message")
		session.Save(r, w)
	}

	tmpl := template.Must(template.ParseFiles("templates/DOM.html"))
	data := struct {
		Message string
	}{
		Message: message,
	}
	tmpl.Execute(w, data)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        log.Println("Received a POST request to /contact")

        err := r.ParseForm()
        if err != nil {
            log.Printf("Error parsing form: %v", err)
            http.Error(w, "Ошибка парсинга формы", http.StatusBadRequest)
            return
        }

        contactType := r.FormValue("contactType") // Тип контакта: email или phone
        contactInfo := r.FormValue("contactInfo") // Контактные данные (email или телефон)

        contactForm := ContactForm{
            Name:         r.FormValue("name"),
            ContactType:  contactType,
            ContactInfo:  contactInfo,
            SelectOption: r.FormValue("selectOption"),
            Message:      r.FormValue("message"),
        }

        log.Printf("Parsed form: %+v", contactForm)

        insertQuery := `INSERT INTO contacts (name, contactType, contactInfo, selectOption, message) VALUES (?, ?, ?, ?, ?)`
        _, err = db.Exec(insertQuery, contactForm.Name, contactForm.ContactType, contactForm.ContactInfo, contactForm.SelectOption, contactForm.Message)
        if err != nil {
            log.Printf("Error inserting data into database: %v", err)
            http.Error(w, "Ошибка вставки данных в базу", http.StatusInternalServerError)
            return
        }

        // Установка сообщения в сессию
        session, _ := store.Get(r, "session-name")
        session.Values["message"] = "Сообщение успешно отправлено! Мы скоро свяжемся с вами."
        session.Save(r, w)

        http.Redirect(w, r, "/", http.StatusSeeOther)
    } else {
        log.Printf("Unsupported request method: %s", r.Method)
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./contacts.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS contacts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        contactType TEXT, -- Тип контакта 
        contactInfo TEXT, -- Сам контакт 
        selectOption TEXT,
        message TEXT
    );`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}
