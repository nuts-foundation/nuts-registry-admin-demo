package domain

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

const serviceProviderBucketName = "ServiceProvider"
const defaultServiceProviderKey = "default"

type ServiceProviderRepository struct {
	DB *bbolt.DB
}

// Get tries to find the default service provider from the database.
// Returns nil when no default service provider was found
func (repo ServiceProviderRepository) Get() (*ServiceProvider, error) {
	sp := &ServiceProvider{}
	if err := repo.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(serviceProviderBucketName))
		if b == nil {
			return nil
		}
		spJSON := b.Get([]byte(defaultServiceProviderKey))
		_ = json.Unmarshal(spJSON, sp)
		return nil
	}); err != nil {
		return nil, err
	}
	return sp, nil
}

func (repo ServiceProviderRepository) CreateOrUpdate(sp ServiceProvider) error {
	spJson, _ := json.Marshal(sp)
	return repo.DB.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(serviceProviderBucketName))
		if err != nil {
			return err
		}
		return b.Put([]byte(defaultServiceProviderKey), spJson)
	})
}
