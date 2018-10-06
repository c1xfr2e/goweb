package mocking

import (
	"io/ioutil"
	"time"

	"github.com/bluecover/lm/models"
	"github.com/bluecover/lm/util"
	"github.com/jinzhu/gorm"
)

func createMockUser(db *gorm.DB) {
	users := []*models.User{
		&models.User{
			Email:    "zh",
			Password: util.GeneratePassword([]byte("000")),
		},
		&models.User{
			Email:    "1@zaoshu.io",
			Password: util.GeneratePassword([]byte("123456")),
		},
		&models.User{
			Email:    "2@zaoshu.io",
			Password: util.GeneratePassword([]byte("123456")),
		},
		&models.User{
			Email:    "3@zaoshu.io",
			Password: util.GeneratePassword([]byte("123456")),
		},
		&models.User{
			Email:    "4@zaoshu.io",
			Password: util.GeneratePassword([]byte("123456")),
		},
		&models.User{
			Email:    "5@zaoshu.io",
			Password: util.GeneratePassword([]byte("123456")),
		},
	}

	for _, u := range users {
		err := db.Create(u).Error
		if err != nil {
			panic(err)
		}
	}
}

func mockFingerprints(db *gorm.DB, userID uint) {
	for i := 0; i < 2; i++ {
		fpOrigin := util.GenerateRandomString(1024)
		hashAll := util.SHA256([]byte(fpOrigin))
		_, err := models.CreateFingerprint(db, userID, hashAll, fpOrigin)
		if err != nil {
			panic(err)
		}
	}
}

func MockDatasets(db *gorm.DB) {
	bytes, err := ioutil.ReadFile("./figures/yrd/dataset.json")
	if err != nil {
		panic(err)
	}

	dsYRD := models.Dataset{
		Name:        "YRD",
		Icon:        "YRD",
		IndexChange: -6.38,
		FigureSet:   string(bytes),
	}
	err = db.Create(&dsYRD).Error
	if err != nil {
		panic(err)
	}

	yrdMsg1 := models.DatasetMessage{
		DatasetID: dsYRD.ID,
		Msg:       "Hello YRD 1998 CNM",
		CreatedAt: time.Now(),
	}
	yrdMsg2 := models.DatasetMessage{
		DatasetID: dsYRD.ID,
		Msg:       "宜人贷数据集",
		CreatedAt: time.Now(),
	}
	db.Create(&yrdMsg1)
	db.Create(&yrdMsg2)
}

func CreateMockingData(db *gorm.DB) {
	/*
		err := db.DropTableIfExists(models.GetModels()...).Error
		if err != nil {
			panic(err)
		}
	*/

	err := db.AutoMigrate(models.GetModels()...).Error
	if err != nil {
		panic(err)
	}

	createMockUser(db)
}
