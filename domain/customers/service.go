package customers

import (
	nutsApi "github.com/nuts-foundation/nuts-node/vdr/api/v1"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
	"net"
)

type Service struct {
	VDRClient  nutsApi.HTTPClient
	Repository Repository
}

func (s Service) ConnectCustomer(id, name, town string) (*domain.Customer, error) {
	didDoc, err := s.VDRClient.Create(nutsApi.DIDCreateRequest{})
	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, domain.ErrNutsNodeUnreachable
		}
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
