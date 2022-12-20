package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/middleware"
	"a21hc3NpZ25tZW50/model"

	"github.com/google/uuid"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var creds model.Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(400)
		error := model.ErrorResponse{
			Error: "Internal Server Error",
		}
		jsonError, err := json.Marshal(error)
		if err != nil {
			return
		}
		w.Write([]byte(jsonError))
		return
	}

	if creds.Username == "" || creds.Password == "" {
		w.WriteHeader(400)
		error := model.ErrorResponse{
			Error: "Username or Password empty",
		}
		jsonError, err := json.Marshal(error)
		if err != nil {
			return
		}
		w.Write([]byte(jsonError))
		return
	}

	for _, j := range db.Users {
		if creds.Username == j {
			w.WriteHeader(409)
			error := model.ErrorResponse{
				Error: "Username already exist",
			}
			jsonError, err := json.Marshal(error)
			if err != nil {
				return
			}
			w.Write([]byte(jsonError))
			return
		}
	}

	db.Users["username"] = creds.Username
	db.Users["password"] = creds.Password

	w.WriteHeader(200)
	success := model.SuccessResponse{
		Username: creds.Username,
		Message:  "Register Success",
	}

	jsonSuccess, err := json.Marshal(success)
	if err != nil {
		return
	}
	w.Write([]byte(jsonSuccess))
	return
	// TODO: answer here
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds model.Credentials
	var exist = false

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(400)
		error := model.ErrorResponse{
			Error: "Internal Server Error",
		}
		jsonError, err := json.Marshal(error)
		if err != nil {
			return
		}
		w.Write([]byte(jsonError))
		return
	}

	if creds.Username == "" || creds.Password == "" {
		w.WriteHeader(400)
		error := model.ErrorResponse{
			Error: "Username or Password empty",
		}
		jsonError, err := json.Marshal(error)
		if err != nil {
			return
		}
		w.Write([]byte(jsonError))
		return
	}

	for _, j := range db.Users {
		if creds.Username == j {
			exist = true
		}
	}

	if exist == false {
		expectedPassword, ok := db.Users[creds.Username]
		if !ok || expectedPassword != creds.Password {
			w.WriteHeader(401)
			error := model.ErrorResponse{
				Error: "Wrong User or Password!",
			}
			jsonError, err := json.Marshal(error)
			if err != nil {
				return
			}
			w.Write([]byte(jsonError))
			return
		}
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(5 * time.Hour)

	db.Sessions[sessionToken] = model.Session{
		Username: creds.Username,
		Expiry:   expiresAt,
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	w.WriteHeader(200)
	success := model.SuccessResponse{
		Username: creds.Username,
		Message:  "Login Success",
	}

	jsonSuccess, err := json.Marshal(success)
	if err != nil {
		return
	}
	w.Write([]byte(jsonSuccess))
	return
	// TODO: answer here
}

func AddToDo(w http.ResponseWriter, r *http.Request) {
	var todo model.Todo
	username := fmt.Sprintf("%s", r.Context().Value("username"))

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		w.WriteHeader(400)
		error := model.ErrorResponse{
			Error: "Internal Server Error",
		}
		jsonError, err := json.Marshal(error)
		if err != nil {
			return
		}
		w.Write([]byte(jsonError))
		return
	}

	todoId := uuid.NewString()

	todoAdd := model.Todo{
		Id:   todoId,
		Task: todo.Task,
		Done: todo.Done,
	}

	if len(db.Task[username]) == 0 {
		db.Task[username] = []model.Todo{}
	}
	db.Task[username] = append(db.Task[username], todoAdd)

	massage := fmt.Sprintf("Task %s added!", todo.Task)

	w.WriteHeader(200)
	success := model.SuccessResponse{
		Username: username,
		Message:  massage,
	}

	jsonSuccess, err := json.Marshal(success)
	if err != nil {
		return
	}
	w.Write([]byte(jsonSuccess))
	return
	// TODO: answer here
}

func ListToDo(w http.ResponseWriter, r *http.Request) {
	username := fmt.Sprintf("%s", r.Context().Value("username"))
	var exist = false

	for i := range db.Task {
		if i == username {
			exist = true
		}
	}

	if exist == false {
		w.WriteHeader(404)
		error := model.ErrorResponse{
			Error: "Todolist not found!",
		}

		jsonError, err := json.Marshal(error)
		if err != nil {
			return
		}
		w.Write([]byte(jsonError))
		return
	}

	w.WriteHeader(200)
	successJson, err := json.Marshal(db.Task[username])
	if err != nil {
		return
	}
	w.Write([]byte(successJson))
	return
	// TODO: answer here
}

func ClearToDo(w http.ResponseWriter, r *http.Request) {
	username := fmt.Sprintf("%s", r.Context().Value("username"))

	db.Task[username] = []model.Todo{}

	w.WriteHeader(200)
	success := model.SuccessResponse{
		Username: username,
		Message:  "Clear ToDo Success",
	}

	successJson, err := json.Marshal(success)
	if err != nil {
		return
	}
	w.Write([]byte(successJson))
	// TODO: answer here
}

func Logout(w http.ResponseWriter, r *http.Request) {
	username := fmt.Sprintf("%s", r.Context().Value("username"))
	c, err := r.Cookie("session_token")
	if err != nil {
		return
	}

	sessionToken := c.Value
	delete(db.Sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
	})

	w.WriteHeader(200)
	success := model.SuccessResponse{
		Username: username,
		Message:  "Logout Success",
	}

	jsonSuccess, err := json.Marshal(success)
	if err != nil {
		return
	}
	w.Write([]byte(jsonSuccess))
	return
	// TODO: answer here
}

func ResetToDo(w http.ResponseWriter, r *http.Request) {
	db.Task = map[string][]model.Todo{}
	w.WriteHeader(http.StatusOK)
}

type API struct {
	mux *http.ServeMux
}

func NewAPI() API {
	mux := http.NewServeMux()
	api := API{
		mux,
	}

	mux.Handle("/user/register", middleware.Post(http.HandlerFunc(Register)))
	mux.Handle("/user/login", middleware.Post(http.HandlerFunc(Login)))
	mux.Handle("/user/logout", middleware.Get(middleware.Auth(http.HandlerFunc(Logout))))

	// TODO: answer here

	mux.Handle("/todo/reset", http.HandlerFunc(ResetToDo))
	mux.Handle("/todo/create", middleware.Post(middleware.Auth(http.HandlerFunc(AddToDo))))
	mux.Handle("/todo/read", middleware.Get(middleware.Auth(http.HandlerFunc(ListToDo))))
	mux.Handle("/todo/clear", middleware.Delete(middleware.Auth(http.HandlerFunc(ClearToDo))))

	return api
}

func (api *API) Handler() *http.ServeMux {
	return api.mux
}

func (api *API) Start() {
	fmt.Println("starting web server at http://localhost:8080")
	http.ListenAndServe(":8080", api.Handler())
}

func main() {
	mainAPI := NewAPI()
	mainAPI.Start()
}
