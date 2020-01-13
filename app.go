package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	w.Write(response)
}

func respondWithText(w http.ResponseWriter, code int, content string) {
	w.Header().Set("Content-Type", "text/plain")

	w.WriteHeader(code)
	w.Write([]byte(content))
}

func sendEmail(subject string, plainTextContent string) {
	sendGridApiKey := os.Getenv("SENDGRID_API_KEY")
	email := os.Args[2]

	if sendGridApiKey != "" {
		from := mail.NewEmail("Todo User", "test@todo.com")
		to := mail.NewEmail("Example User", email)
		htmlContent := "<strong>" + plainTextContent + "</strong>"
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
		_, err := client.Send(message)
		if err != nil {
			log.Println(err)
		}
	}
}

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) getInfo(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "Todos API")
}

func (a *App) getTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	t := todo{ID: id}

	if err := t.getTodo(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Todo not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

func (a *App) getTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := getTodos(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, todos)
}

func (a *App) createTodo(w http.ResponseWriter, r *http.Request) {
	var t todo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := t.createTodo(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go sendEmail("new todo just created", t.Content)
	respondWithJSON(w, http.StatusCreated, t)
}

func (a *App) updateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	var t todo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	t.ID = id

	if err := t.updateTodo(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var priorityMap map[int]string
	priorityMap = make(map[int]string)
	priorityMap[0] = "Normal"
	priorityMap[1] = "Low"
	priorityMap[2] = "High"

	var completeMap map[int]string
	completeMap = make(map[int]string)
	completeMap[0] = "Not completed"
	completeMap[1] = "Completed"

	go sendEmail("todo just updated", t.Content+
		" with "+priorityMap[t.Priority]+
		" priority as "+completeMap[t.Completed])

	respondWithJSON(w, http.StatusOK, t)
}

func (a *App) deleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Todo ID")
		return
	}

	t := todo{ID: id}
	if err := t.deleteTodo(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go sendEmail("Todo just deleted", "todo #"+vars["id"]+" just deleted. Please, check your list!")

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) Initialize(dbname string) {
	var err error
	if _, err := os.Stat(dbname); os.IsNotExist(err) {
		log.Fatal(err)
		return
	}

	a.DB, err = sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal(err)
		return
	}

	a.Router = mux.NewRouter()
	a.initializeRouters()
}

func (a *App) initializeRouters() {
	a.Router.HandleFunc("/", a.getInfo).Methods("GET")
	a.Router.HandleFunc("/todos", a.getTodos).Methods("GET")
	a.Router.HandleFunc("/todo", a.createTodo).Methods("POST")
	a.Router.HandleFunc("/todo/{id:[0-9]+}", a.getTodo).Methods("GET")
	a.Router.HandleFunc("/todo/{id:[0-9]+}", a.updateTodo).Methods("PUT")
	a.Router.HandleFunc("/todo/{id:[0-9]+}", a.deleteTodo).Methods("DELETE")
}

func (a *App) Run() {
	port, ok := os.LookupEnv("PORT")

	if ok == false {
		port = "3001"
	}

	allowedOrigins := os.Args[2]

	log.Fatal(http.
		ListenAndServe(":"+port,
			handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
				handlers.AllowedOrigins([]string{allowedOrigins}))(a.Router)))
}
