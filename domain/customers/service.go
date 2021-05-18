package customers

import (
	"net"

	"github.com/nuts-foundation/go-did/did"
	nutsApi "github.com/nuts-foundation/nuts-node/vdr/api/v1"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
)

type Service struct {
	VDRClient  nutsApi.HTTPClient
	Repository Repository
}

func (s Service) ConnectCustomer(id, name, town string, serviceProviderID did.DID) (*domain.Customer, error) {
	selfControl := false
	capabilityInvocation := false
	controllers := []string{serviceProviderID.String()}

	didDoc, err := s.VDRClient.Create(nutsApi.DIDCreateRequest{
		SelfControl:          &selfControl,
		Controllers:          &controllers,
		CapabilityInvocation: &capabilityInvocation,
	})
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
