package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"golang-docker-todo/api/models"
)

var users = []models.User{
	models.User{
		Username: "Riadh Rzig",
		Email:    "riadh@maisonduweb.com",
		Password: "riadh123",
	},
	models.User{
		Username: "Sofien Slimi",
		Email:    "sofien@maisonduweb.com",
		Password: "sofien123",
	},
	models.User{
		Username: "Dhouha Rhaiem",
		Email:    "dhouha@maisonduweb.com",
		Password: "dhouha123",
	},
}

var tasks = []models.Task{
	models.Task{
		Content: "Support client 09h",
	},
	models.Task{
		Content: "Workshop 10h",
	},
	models.Task{
		Content: "Meeting 10h",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Task{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Task{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Task{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		tasks[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Task{}).Create(&tasks[i]).Error
		if err != nil {
			log.Fatalf("cannot seed tasks table: %v", err)
		}
	}
}