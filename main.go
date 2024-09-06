package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
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

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}
}

func main() {
	initDB()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/admin", adminHandler)

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
func adminHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || username != os.Getenv("ADMIN_USERNAME") || password != os.Getenv("ADMIN_PASSWORD") {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == "GET" {
		contacts, err := getAllContacts()
		if err != nil {
			http.Error(w, "Ошибка при получении контактов", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("templates/admin.html")
		if err != nil {
			http.Error(w, "Ошибка при загрузке шаблона", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, contacts)
		if err != nil {
			http.Error(w, "Ошибка при отображении шаблона", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func getAllContacts() ([]ContactForm, error) {
	query := `SELECT name, contactType, contactInfo, selectOption, message FROM contacts ORDER BY id DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []ContactForm
	for rows.Next() {
		var contact ContactForm
		err := rows.Scan(&contact.Name, &contact.ContactType, &contact.ContactInfo, &contact.SelectOption, &contact.Message)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
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
