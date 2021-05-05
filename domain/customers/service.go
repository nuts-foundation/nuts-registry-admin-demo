package customers

import (
	"time"

	nutsApi "github.com/nuts-foundation/nuts-node/vdr/api/v1"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
)

type Service struct {
	NutsNodeAddr string
	Repository Repository
}

func (s Service) ConnectCustomer(id, name, town string) (*domain.Customer, error) {
	nodeClient := nutsApi.HTTPClient{
		ServerAddress: s.NutsNodeAddr,
		Timeout:       5*time.Second,
	}

	didDoc, err := nodeClient.Create()
	if err != nil {
		return nil, err
	}

	customer := domain.Customer{
		Did:  didDoc.ID.String(),
		Id:   id,
		Name: name,
		Town: &town,
	}

	return s.Repository.NewCustomer(customer)
}


