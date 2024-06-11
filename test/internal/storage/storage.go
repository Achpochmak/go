package storage

import (
	"encoding/json"
	"errors"
	"os"

	"HOMEWORK-1/internal/models"
)

type Storage struct {
	fileName string
}

// NewStorage .. TODO сделать описание функции
func NewStorage(fileName string) Storage {
	return Storage{fileName: fileName}
}

func (s Storage) AddOrder(order models.Order) error {
	if _, err := os.Stat(s.fileName); errors.Is(err, os.ErrNotExist) {
		// создаем файл
		if errCreateFile := s.createFile(); errCreateFile != nil {
			return errCreateFile
		}
	}

	// прочитать
	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return err
	}

	var records []orderRecord
	if errUnmarshal := json.Unmarshal(b, &records); errUnmarshal != nil {
		return errUnmarshal
	}
	records = append(records, transform(order))

	bWrite, errMarshal := json.MarshalIndent(records, "  ", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	return os.WriteFile(s.fileName, bWrite, 0666)
}

func (s Storage) ReWrite(orders []models.Order) error {
	if _, err := os.Stat(s.fileName); errors.Is(err, os.ErrNotExist) {
		// создаем файл
		if errCreateFile := s.createFile(); errCreateFile != nil {
			return errCreateFile
		}
	}
	var records []orderRecord
	for _, order := range orders {
		records = append(records, transform(order))
	}
	bWrite, errMarshal := json.MarshalIndent(records, "  ", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	return os.WriteFile(s.fileName, bWrite, 0666)
}

func (s Storage) UpdateOrder(order models.Order) error {
	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return err
	}
	var records []orderRecord
	if errUnmarshal := json.Unmarshal(b, &records); errUnmarshal != nil {
		return errUnmarshal
	}
	for i, record := range records {
		if models.Id(record.Id) == order.Id {
			records[i] = transform(order)
			break
		}
	}
	bWrite, errMarshal := json.MarshalIndent(records, "  ", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	return os.WriteFile(s.fileName, bWrite, 0666)
}

func (s Storage) FindOrder(id models.Id) (models.Order, error) { 
	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return models.Order{}, err
	}
	var records []orderRecord
	if errUnmarshal := json.Unmarshal(b, &records); errUnmarshal != nil {
		return models.Order{}, errUnmarshal
	}

	for _, record := range records {
		if models.Id(record.Id) == id {
			return record.toDomain(), nil
		}
	}

	return models.Order{}, errors.New("name not found")
}

func (s Storage) ListOrder() ([]models.Order, error) {
	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return nil, err
	}
	var records []orderRecord
	if errUnmarshal := json.Unmarshal(b, &records); errUnmarshal != nil {
		return nil, errUnmarshal
	}

	result := make([]models.Order, 0, len(records))
	for _, record := range records {
		result = append(result, record.toDomain())
	}
	return result, nil
}

func (s Storage) createFile() error {
	f, err := os.Create(s.fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}
