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

var ErrNotFound = errors.New("not found")

type Repository interface {
	NewCustomer(customer domain.Customer) (*domain.Customer, error)
	FindByID(id int) (*domain.Customer, error)
	Update(id int, updateFn func(c domain.Customer) (*domain.Customer, error)) (*domain.Customer, error)
	All() ([]domain.Customer, error)
}

type flatFileRepo struct {
	filepath string
	mutex    sync.Mutex
	// records is a cache
	records map[int]domain.Customer
}

func NewFlatFileRepository(filepath string) *flatFileRepo {
	f, err := os.OpenFile(filepath, os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	return &flatFileRepo{
		filepath: filepath,
		mutex:    sync.Mutex{},
		records:  make(map[int]domain.Customer, 0),
	}
}

// NewCustomer creates a new customer with a valid id
func (db *flatFileRepo) NewCustomer(customer domain.Customer) (*domain.Customer, error) {
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

func (db *flatFileRepo) FindByID(id int) (*domain.Customer, error) {
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

	return nil, fmt.Errorf("could not FindCustomerByID with id: %s, reason: %w", id, ErrNotFound)
}

func (db *flatFileRepo) Update(id int, updateFn func(c domain.Customer) (*domain.Customer, error)) (*domain.Customer, error) {
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

	return nil, fmt.Errorf("could update customer with id: %s, reason: %w", id, ErrNotFound)
}

// WriteAll writes all records to the file, truncating the file if it exists
func (db *flatFileRepo) WriteAll() error {

	bytes, err := json.Marshal(db.records)
	if err != nil {
		return fmt.Errorf("unable to marshall db records to json: %w", err)
	}

	if err = os.WriteFile(db.filepath, bytes, 0666); err != nil {
		return fmt.Errorf("unable to write db to file: %w", err)
	}
	return nil
}

func (db *flatFileRepo) ReadAll() error {
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

func (db *flatFileRepo) All() ([]domain.Customer, error) {
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
