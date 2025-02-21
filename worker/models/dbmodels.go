package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DBJob struct {
	Id       uuid.UUID `gorm:"primaryKey,column:id"`
	TaskId   uuid.UUID `gorm:"column:task_id;type:uuid;not null"`
	Status   string    `gorm:"column:status;type:varchar(100);not null"`
	Data     string    `gorm:"column:data;type:text;not null"`
	Response string    `gorm:"column:response;type:text"`
}

func (DBJob) TableName() string {
	return "job_tbl"
}

func (*DBJob) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("id", uuid)
	return nil
}
