package dal

import (
	"log"
	"task-queue/grpc-server/models"

	"gorm.io/gorm"
)

type DalImpl struct{}

type Dal interface {
	CreateTask(value *models.DBJob, trans *gorm.DB) error
	GetTask(id string, trans *gorm.DB) (*models.DBJob, error)
}

func NewDalRequest() (Dal, error) {
	return &DalImpl{}, nil
}

func (d *DalImpl) CreateTask(value *models.DBJob, trans *gorm.DB) error {
	transaction := trans.Begin()
	if transaction.Error != nil {
		log.Println("error while starting the transaction", transaction.Error)
		return transaction.Error
	}
	defer transaction.Rollback()
	addJob := transaction.Create(value)
	if addJob.Error != nil {
		log.Println("error while adding the job: ", addJob.Error)
		return addJob.Error
	}
	return nil
}

func (d *DalImpl) GetTask(id string, trans *gorm.DB) (*models.DBJob, error) {
	transaction := trans.Begin()
	if transaction.Error != nil {
		log.Println("error while starting the transaction", transaction.Error)
		return nil, transaction.Error
	}
	var task models.DBJob
	getJob := transaction.Where("id = ?", id).First(&task)
	if getJob.Error != nil {
		log.Println("error while getting the job: ", getJob.Error)
		return nil, getJob.Error
	}
	defer transaction.Rollback()
	return nil, nil
}
