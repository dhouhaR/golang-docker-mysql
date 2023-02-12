package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"golang-docker-todo/api/controllers"
	"golang-docker-todo/api/models"
)

var server = controllers.Server{}
var userInstance = models.User{}
var taskInstance = models.Task{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
	/*fmt.Printf("TestDbUser :  %s \n", os.Getenv("TestDbUser"))
	fmt.Printf("TestDbDriver :  %s \n", os.Getenv("TestDbDriver"))
	fmt.Printf("DBURL :  %s \n", DBURL)*/
	server.DB, err = gorm.Open(TestDbDriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", TestDbDriver)
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	refreshUserTable()

	user := models.User{
		Username: "Nicolas",
		Email:    "nicolas@mycompany.com",
		Password: "password",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	return user, nil
}

func seedUsers() ([]models.User, error) {

	var err error
	if err != nil {
		return nil, err
	}
	users := []models.User{
		models.User{
			Username: "Dani",
			Email:    "dani@mycompany.com",
			Password: "password",
		},
		models.User{
			Username: "Han",
			Email:    "han@mycompany.com",
			Password: "password",
		},
	}

	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}
	return users, nil
}


func refreshUserAndTaskTable() error {

	err := server.DB.DropTableIfExists(&models.User{}, &models.Task{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Task{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserAndOneTask() (models.Task, error) {

	err := refreshUserAndTaskTable()
	if err != nil {
		return models.Task{}, err
	}
	user := models.User{
		Username: "Tom",
		Email:    "tom@mycompany.com",
		Password: "password",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Task{}, err
	}
	task := models.Task{
		Content:  "This is the task of Tom",
		AuthorID: user.ID,
	}
	err = server.DB.Model(&models.Task{}).Create(&task).Error
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func seedUsersAndTasks() ([]models.User, []models.Task, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Task{}, err
	}
	var users = []models.User{
		models.User{
			Username: "Dani",
			Email:    "dani@mycompany.com",
			Password: "password",
		},
		models.User{
			Username: "Ron",
			Email:    "ron@mycompany.com",
			Password: "password",
		},
	}
	var tasks = []models.Task{
		models.Task{
			Content: "Hello world 1",
		},
		models.Task{
			Content: "Hello world 2",
		},
	}

	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		tasks[i].AuthorID = users[i].ID

		err = server.DB.Model(&models.Task{}).Create(&tasks[i]).Error
		if err != nil {
			log.Fatalf("cannot seed tasks table: %v", err)
		}
	}
	return users, tasks, nil
}