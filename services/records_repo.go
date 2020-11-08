package services

import (
	"balancer-api/models"
	"log"

	"gorm.io/gorm"
)

type RecordService struct {
	DB *gorm.DB
}

func (rs *RecordService) GetRecord(id string) (*models.Record, error) {
	var record models.Record

	if result := rs.DB.Find(&record, id); result.Error != nil {
		log.Println("Error reading entry from db")
		return nil, result.Error
	}

	return &record, nil
}

func (rs *RecordService) GetAllRecords() (*[]models.Record, error) {
	var records []models.Record

	if result := rs.DB.Find(&records); result.Error != nil {
		log.Println("Error reading entries in db")
		return nil, result.Error
	}

	return &records, nil
}

func (rs *RecordService) GetRecordsByType(typeKey string) (*[]models.Record, error) {
	var records []models.Record

	if result := rs.DB.Where("type = ?", typeKey).Find(&records); result.Error != nil {
		log.Println("Error reading entries in db")
		return nil, result.Error
	}

	return &records, nil
}

func (rs *RecordService) CreateRecord(record *models.Record) (*uint64, error) {
	if result := rs.DB.Create(&record); result.Error != nil {
		log.Println("Error creating entry in db")
		return nil, result.Error
	}

	return &record.ID, nil
}

func (rs *RecordService) UpdateRecord(record *models.Record) error {
	if result := rs.DB.Save(&record); result.Error != nil {
		log.Println("Error updating entry in db")
		return result.Error
	}

	return nil
}

func (rs *RecordService) DeleteRecord(id uint64) error {
	if result := rs.DB.Delete(&models.Record{}, id); result.Error != nil {
		log.Println("Error deleting entry in db")
		return result.Error
	}

	return nil
}
