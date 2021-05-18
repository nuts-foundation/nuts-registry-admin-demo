package sp

import (
	"fmt"
	"net"

	didmanAPI "github.com/nuts-foundation/nuts-node/didman/api/v1"
	vdrAPI "github.com/nuts-foundation/nuts-node/vdr/api/v1"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
)

type Service struct {
	Repository   Repository
	VDRClient    vdrAPI.HTTPClient
	DIDManClient didmanAPI.HTTPClient
}

// Get tries to find the default service provider from the database.
// Returns nil when no default service provider was found
func (svc Service) Get() (*domain.ServiceProvider, error) {
	spDID, err := svc.Repository.Get()
	if err != nil {
		return nil, err
	}
	if spDID == nil{
		return nil, nil
	}
	sp := &domain.ServiceProvider{Id: spDID.String()}
	contactInformation, err := svc.DIDManClient.GetContactInformation(sp.Id)
	if err != nil {
		return nil, unwrapAPIError(err)
	}
	if contactInformation != nil {
		sp.Email = contactInformation.Email
		sp.Name = contactInformation.Name
		sp.Phone = contactInformation.Phone
		sp.Website = contactInformation.Website
	}
	return sp, nil
}

func (svc Service) CreateOrUpdate(sp domain.ServiceProvider) (*domain.ServiceProvider, error) {
	if len(sp.Id) == 0 {
		// Service Provider not registered yet, so create a DID
		didDoc, err := svc.VDRClient.Create(vdrAPI.DIDCreateRequest{})
		if err != nil {
			return nil, unwrapAPIError(err)
		}
		sp.Id = didDoc.ID.String()
		if err := svc.Repository.Set(sp.Id); err != nil {
			return nil, err
		}
	}
	err := svc.DIDManClient.UpdateContactInformation(sp.Id, didmanAPI.ContactInformation{
		Name:    sp.Name,
		Email:   sp.Email,
		Website: sp.Website,
		Phone:   sp.Phone,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to update DID contact info: %w", unwrapAPIError(err))
	}
	return &sp, nil
}

func unwrapAPIError(err error) error {
	if _, ok := err.(net.Error); ok {
		return domain.ErrNutsNodeUnreachable
	}
	return err
}
