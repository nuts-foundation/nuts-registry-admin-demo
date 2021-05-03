package customers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
)

type Repository interface {
	NewCustomer(customer domain.Customer) (*domain.Customer, error)
	FindByID(id string) (*domain.Customer, error)
	Update(id string, updateFn func(c domain.Customer) (*domain.Customer, error)) (*domain.Customer, error)
	All() ([]domain.Customer, error)
}

type FlatFileDB struct {
	filepath string
	mutex    sync.Mutex
	// records is a cache
	records map[string]domain.Customer
}

func NewDB(filepath string) *FlatFileDB {
	f, err := os.OpenFile(filepath, os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	return &FlatFileDB{
		filepath: filepath,
		mutex:    sync.Mutex{},
		records:  make(map[string]domain.Customer, 0),
	}
}

// NewCustomer creates a new customer with a valid id
func (db *FlatFileDB) NewCustomer(customer domain.Customer) (*domain.Customer, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if err := db.ReadAll(); err != nil {
		return nil, err
	}
	if _, ok := db.records[customer.Id]; ok {
		return nil, errors.New("customer already exists")
	}

	db.records[customer.Id] = customer

	return &customer, db.WriteAll()
}

func (db *FlatFileDB) FindByID(id string) (*domain.Customer, error) {
	if len(db.records) == 0 {
		if err := db.ReadAll(); err != nil {
			return nil, err
		}
	}

	for _, r := range db.records {
		if r.Id == id {
			// Hazardous to return a pointer, but this is a demo.
			return &r, nil
		}
	}

	return nil, errors.New("not found")
}

func (db *FlatFileDB) Update(id string, updateFn func(c domain.Customer) (*domain.Customer, error)) (*domain.Customer, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if len(db.records) == 0 {
		if err := db.ReadAll(); err != nil {
			return nil, err
		}
	}

	for i, r := range db.records {
		if r.Id == id {
			updatedCustomer, err := updateFn(db.records[i])
			if err != nil {
				return nil, err
			}

			db.records[i] = *updatedCustomer
			err = db.WriteAll()
			if err != nil {
				return nil, err
			}
			return updatedCustomer, nil
		}
	}

	return nil, errors.New("not found")
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
	//log.Debug("Reading full customer list from file")
	bytes, err := os.ReadFile(db.filepath)
	if err != nil {
		return fmt.Errorf("unable to read db from file: %w", err)
	}

	if len(bytes) == 0 {
		return nil
	}

	if err = json.Unmarshal(bytes, &db.records); err != nil {
		return fmt.Errorf("unable to unmarshall db from file: %w", err)
	}
	return nil
}

func (db *FlatFileDB) All() ([]domain.Customer, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if err := db.ReadAll(); err != nil {
		return nil, err
	}

	v := make([]domain.Customer, len(db.records))

	idx := 0
	for _, value := range db.records {
		v[idx] = value
		idx = idx + 1
	}
	sort.SliceStable(v, func(i, j int) bool {
		return v[i].Id < v[j].Id
	})
	return v, nil
}
