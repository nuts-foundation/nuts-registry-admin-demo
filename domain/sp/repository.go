package sp

import (
	"encoding/json"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
	"go.etcd.io/bbolt"
)

const serviceProviderBucketName = "ServiceProvider"
const defaultServiceProviderKey = "default"

type Repository interface {
	// Get returns the DID of the Service Provider.
	Get() (string, error)
	Set(did string) error
}

type bboltRepository struct {
	DB *bbolt.DB
}

func NewBBoltRepository(db *bbolt.DB) Repository {
	return &bboltRepository{DB: db}
}

func (b bboltRepository) Get() (string, error) {
	spDID := ""
	err := b.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(serviceProviderBucketName))
		if b == nil {
			return nil
		}
		spData := b.Get([]byte(defaultServiceProviderKey))
		// Backwards compatibility for when we persisted the complete SP including contact information, rather than
		// reading it from the DID document. Can be removed when everyone has updated viewed their SP once.
		result := domain.ServiceProvider{}
		err := json.Unmarshal(spData, &result)
		if err == nil {
			return err
		}
		spDID = result.Id
		return nil
	})
	return spDID, err
}

func (b bboltRepository) Set(did string) error {
	sp := domain.ServiceProvider{}
	sp.Id = did
	return b.DB.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(serviceProviderBucketName))
		if err != nil {
			return err
		}
		return b.Put([]byte(defaultServiceProviderKey), []byte(sp.Id))
	})
}