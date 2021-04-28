package customers

import (
	"time"

	nutsApi "github.com/nuts-foundation/nuts-node/vdr/api/v1"
)

type Service struct {
	NutsNodeAddr string
}

func (s Service) ConnectCustomer(id, name string) (*Customer, error) {
	nodeClient := nutsApi.HTTPClient{
		ServerAddress: s.NutsNodeAddr,
		Timeout:       5*time.Second,
	}

	didDoc, err := nodeClient.Create()
	if err != nil {
		return nil, err
	}

	// todo: save customer

	return &Customer{
		Did:  didDoc.ID.String(),
		ID:   id,
		Name: name,
	}, nil
}


