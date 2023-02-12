package tests

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang-docker-todo/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllTasks(t *testing.T) {

	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatalf("Error refreshing user and task table %v\n", err)
	}
	_, _, err = seedUsersAndTasks()
	if err != nil {
		log.Fatalf("Error seeding user and task  table %v\n", err)
	}
	tasks, err := taskInstance.FindAllTasks(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the tasks: %v\n", err)
		return
	}
	assert.Equal(t, len(*tasks), 2)
}

func TestSaveTask(t *testing.T) {

	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatalf("Error user and task refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newTask := models.Task{
		ID:       1,
		Content:  "This is the task content",
		AuthorID: user.ID,
	}
	savedTask, err := newTask.SaveTask(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the task: %v\n", err)
		return
	}
	assert.Equal(t, newTask.ID, savedTask.ID)
	assert.Equal(t, newTask.Content, savedTask.Content)
	assert.Equal(t, newTask.AuthorID, savedTask.AuthorID)

}

func TestFindTaskByID(t *testing.T) {

	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatalf("Error refreshing user and task table: %v\n", err)
	}
	task, err := seedOneUserAndOneTask()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundTask, err := taskInstance.FindTaskByID(server.DB, task.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundTask.ID, task.ID)
	assert.Equal(t, foundTask.Content, task.Content)
}

func TestUpdateATask(t *testing.T) {

	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatalf("Error refreshing user and task table: %v\n", err)
	}
	task, err := seedOneUserAndOneTask()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	taskUpdate := models.Task{
		ID:       1,
		Content:  "updateme@mycompany.com",
		AuthorID: task.AuthorID,
	}
	updatedTask, err := taskUpdate.UpdateATask(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedTask.ID, taskUpdate.ID)
	assert.Equal(t, updatedTask.Content, taskUpdate.Content)
	assert.Equal(t, updatedTask.AuthorID, taskUpdate.AuthorID)
}

func TestDeleteATask(t *testing.T) {

	err := refreshUserAndTaskTable()
	if err != nil {
		log.Fatalf("Error refreshing user and task table: %v\n", err)
	}
	task, err := seedOneUserAndOneTask()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := taskInstance.DeleteATask(server.DB, task.ID, task.AuthorID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}

	assert.Equal(t, isDeleted, int64(1))
}