package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"golang-docker-todo/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateTask(t *testing.T) {

	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}
	token, err := server.SignIn(user.Email, "password")
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		content      string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"content": "the task content", "author_id": 1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			content:      "the task content",
			author_id:    user.ID,
			errorMessage: "",
		},
		{
			// When no token is passed
			inputJSON:    `{"content": "the task content passed with no token", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"content": "the task content passed with incorrect token", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"content": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Content",
		},
		{
			inputJSON:    `{"content": "the task content without author"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Author",
		},
		{
			// When user with id = 2 uses user with id = 1 token
			inputJSON:    `{"content": "the task content", "author_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/tasks", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateTask)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id))
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetTasks(t *testing.T) {

	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedUsersAndTasks()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetTasks)
	handler.ServeHTTP(rr, req)

	var tasks []models.Task
	err = json.Unmarshal([]byte(rr.Body.String()), &tasks)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(tasks), 2)
}
func TestGetTaskByID(t *testing.T) {

	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatal(err)
	}
	task, err := seedOneUserAndOneTask()
	if err != nil {
		log.Fatal(err)
	}
	taskSample := []struct {
		id           string
		statusCode   int
		content      string
		author_id    uint32
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(task.ID)),
			statusCode: 200,
			content:    task.Content,
			author_id:  task.AuthorID,
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}
	for _, v := range taskSample {

		req, err := http.NewRequest("GET", "/tasks", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetTask)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, task.Content, responseMap["content"])
			assert.Equal(t, float64(task.AuthorID), responseMap["author_id"])
		}
	}
}

func TestUpdateTask(t *testing.T) {

	var TaskUserEmail, TaskUserPassword string
	var AuthTaskAuthorID uint32
	var AuthTaskID uint64

	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatal(err)
	}
	users, tasks, err := seedUsersAndTasks()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		TaskUserEmail = user.Email
		TaskUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(TaskUserEmail, TaskUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the first task
	for _, task := range tasks {
		if task.ID == 2 {
			continue
		}
		AuthTaskID = task.ID
		AuthTaskAuthorID = task.AuthorID
	}
	// fmt.Printf("this is the auth task: %v\n", AuthTaskID)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		content      string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"content": "This is the updated content", "author_id": 1}`,
			statusCode:   200,
			content:      "This is the updated content",
			author_id:    AuthTaskAuthorID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token is passed
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"content": "This is the updated content with no token", "author_id": 1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"content": "This is the updated content with incorrect token", "author_id": 1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"content": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Content",
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(AuthTaskID)),
			updateJSON:   `{"content": "This is the updated content with other user", "author_id": 2}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/tasks", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateTask)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id))
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteTask(t *testing.T) {

	var TaskUserEmail, TaskUserPassword string
	var TaskUserID uint32
	var AuthTaskID uint64

	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatal(err)
	}
	users, tasks, err := seedUsersAndTasks()
	if err != nil {
		log.Fatal(err)
	}
	// Get only the Second user (id = 2)
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		TaskUserEmail = user.Email
		TaskUserPassword = "password"
	}
	// Login the user with id = 2 and get the authentication token
	token, err := server.SignIn(TaskUserEmail, TaskUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the second task (id = 2)
	for _, task := range tasks {
		if task.ID == 1 {
			continue
		}
		AuthTaskID = task.ID
		TaskUserID = task.AuthorID
	}
	taskSample := []struct {
		id           string
		author_id    uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			id:           strconv.Itoa(int(AuthTaskID)),
			author_id:    TaskUserID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// Empty token is passed
			id:           strconv.Itoa(int(AuthTaskID)),
			author_id:    TaskUserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// Incorrect token is passed
			id:           strconv.Itoa(int(AuthTaskID)),
			author_id:    TaskUserID,
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(1)),
			author_id:    1,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range taskSample {

		req, _ := http.NewRequest("GET", "/tasks", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteTask)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}