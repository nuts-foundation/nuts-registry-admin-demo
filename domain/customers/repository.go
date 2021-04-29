package customers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

// Customer record in DB
type Customer struct {
	ID   string `json:"id"`
	DID  string `json:"did"`
	Name string `json:"name"`
	// Domain is also used as login constraint
	Domain string `json:"domain"`
}

type Repository interface {
	NewCustomer(customer Customer) (*Customer, error)
	FindByID(id string) (*Customer, error)
	Update(id string, updateFn func(c Customer) (*Customer, error)) error
	All() ([]Customer, error)
}

type FlatFileDB struct {
	filepath string
	mutex    sync.Mutex
	// records is a cache
	records map[string]Customer
}

func NewDB(filepath string) *FlatFileDB {
	f, err := os.OpenFile(filepath, os.O_CREATE,0666)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	return &FlatFileDB{
		filepath: filepath,
		mutex:    sync.Mutex{},
		records:  make(map[string]Customer, 0),
	}
}

// NewCustomer creates a new customer with a valid id
func (db *FlatFileDB) NewCustomer(customer Customer) (*Customer, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if err := db.ReadAll(); err != nil {
		return nil, err
	}
	if _, ok := db.records[customer.ID]; ok {
		return nil, errors.New("customer already exists")
	}

	db.records[customer.ID] = customer

	return &customer, db.WriteAll()
}

func (db *FlatFileDB) FindByID(id string) (*Customer, error) {
	if len(db.records) == 0 {
		if err := db.ReadAll(); err != nil {
			return nil, err
		}
	}

	for _, r := range db.records {
		if r.ID == id {
			// Hazardous to return a pointer, but this is a demo.
			return &r, nil
		}
	}

	return nil, errors.New("not found")
}

func (db *FlatFileDB) Update(id string, updateFn func(c Customer) (*Customer, error)) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if len(db.records) == 0 {
		if err := db.ReadAll(); err != nil {
			return err
		}
	}

	for i, r := range db.records {
		if r.ID == id {
			updatedCustomer, err := updateFn(db.records[i])
			if err != nil {
				return err
			}

			db.records[i] = *updatedCustomer
			return db.WriteAll()
		}
	}

	return errors.New("not found")
}

// WriteAll writes all records to the file, truncating the file if it exists
func (db *FlatFileDB) WriteAll() error {

	bytes, err := json.Marshal(db.records)
	if err != nil {
		return fmt.Errorf("unable to marshall db records to json: %w", err)
	}

	if err = os.WriteFile(db.filepath, bytes, 0666); err != nil {
		return fmt.Errorf("unable to write db to file: %w", err)
	}
	return nil
}

func (db *FlatFileDB) ReadAll() error {
	bytes, err := os.ReadFile(db.filepath)
	if err != nil {
		return fmt.Errorf("unable to read db from file: %w", err)
	}

	if len(bytes) == 0 {
		return  nil
	}

	if err = json.Unmarshal(bytes, &db.records); err != nil {
		return fmt.Errorf("unable to unmarshall db from file: %w", err)
	}
	return nil
}

func (db *FlatFileDB) All() ([]Customer, error) {
	if err := db.ReadAll(); err != nil {
		return nil, err
	}

	v := make([]Customer, 0, len(db.records))

	for _, value := range db.records {
		v = append(v, value)
	}
	return v, nil
}
