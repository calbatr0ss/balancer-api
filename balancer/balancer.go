package balancer

import (
	"balancer-api/models"
)

type RecordService interface {
	GetRecord(string) (*models.Record, error)
	GetAllRecords() (*[]models.Record, error)
	GetRecordsByType(string) (*[]models.Record, error)
	CreateRecord(*models.Record) (*uint64, error)
	UpdateRecord(*models.Record) error
	DeleteRecord(id uint64) error
}
