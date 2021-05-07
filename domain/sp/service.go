package sp

import (
	vdrAPI "github.com/nuts-foundation/nuts-node/vdr/api/v1"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
	"net"
)

type Service struct {
	Repository Repository
	VDRClient  vdrAPI.HTTPClient
}

// Get tries to find the default service provider from the database.
// Returns nil when no default service provider was found
func (svc Service) Get() (*domain.ServiceProvider, error) {
	sp := &domain.ServiceProvider{}

	spDID, err := svc.Repository.Get()
	if err != nil {
		return nil, unwrapAPIError(err)
	}
	if spDID == "" {
		return nil, nil
	}
	// TODO: Use DIDMAN!
	document, _, err := svc.VDRClient.Get(spDID)
	if err != nil {
		return nil, err
	}
	for _, service := range document.Service {
		if service.Type == "node-contact-info" {
			contactInfo := make(map[string]string, 0)
			err := service.UnmarshalServiceEndpoint(&contactInfo)
			if err == nil {
				sp.Name = contactInfo["name"]
				sp.Email = contactInfo["email"]
				sp.Phone = contactInfo["tel"]
				continue
			}
		}
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



	return &sp, nil
}

func unwrapAPIError(err error) error {
	if _, ok := err.(net.Error); ok {
		return domain.ErrNutsNodeUnreachable
	}
	return err
}
