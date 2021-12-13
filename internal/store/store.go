//go:generate mockgen -package=store -source=store.go -destination=./store_mock.go Store
package store

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/moynur/gateway/internal/state"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Storer interface {
	Create(request *Transaction) error
	GetLatestByTransactionId(TransactionId uuid.UUID) (Transaction, error)
	Update(request *Transaction) error
}

type Store struct {
	db *gorm.DB
}

type Transaction struct {
	ID              uint `gorm:"primaryKey"`
	TransactionId   uuid.UUID
	OperationId     uuid.UUID
	Amount          int
	AmountAvailable int
	Currency        string
	OperationType   state.State
	Pan             string
}

func NewStore() (*Store, error) {
	dsn := "host=db user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	log.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Println(err)
	}
	return &Store{
		db: db,
	}, nil
}

func (s *Store) Create(storeRequest *Transaction) error {
	log.Println("store request", &storeRequest)
	err := s.db.Create(storeRequest).Error
	if err != nil {
		return fmt.Errorf("failed to create record %e", err)
	}
	return nil
}

func (s *Store) GetLatestByTransactionId(TransactionId uuid.UUID) (Transaction, error) {
	log.Println("get store request", TransactionId)
	var FindResult Transaction
	resp := s.db.Where(&Transaction{TransactionId: TransactionId}).Last(&FindResult)
	if resp.Error != nil {
		return Transaction{}, fmt.Errorf("failed to get record %e", resp.Error)
	}
	log.Println("find result", FindResult)
	return FindResult, nil
}

func (s *Store) Update(storerequest *Transaction) error {
	return errors.New("not implemented")
}
