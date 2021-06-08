package sp

import (
	"fmt"

	ssi "github.com/nuts-foundation/go-did"
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
	if spDID == nil {
		return nil, nil
	}
	sp := &domain.ServiceProvider{Id: spDID.String()}
	if err = svc.enrichWithContactInfo(sp); err != nil {
		return nil, err
	}
	return sp, nil
}

func (svc Service) CreateOrUpdate(sp domain.ServiceProvider) (*domain.ServiceProvider, error) {
	if len(sp.Id) == 0 {
		// Service Provider not registered yet, so create a DID
		didDoc, err := svc.VDRClient.Create(vdrAPI.DIDCreateRequest{})
		if err != nil {
			return nil, domain.UnwrapAPIError(err)
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
		return nil, fmt.Errorf("unable to update DID contact info: %w", domain.UnwrapAPIError(err))
	}
	return &sp, nil
}

func (svc Service) RegisterEndpoint(endpoint domain.EndpointProperties) error {
	spDID, err := svc.Repository.Get()
	if err != nil {
		return err
	}
	_, err = svc.DIDManClient.AddEndpoint(spDID.String(), endpoint.Type, endpoint.Url)
	return err
}

func (svc Service) DeleteEndpoint(id ssi.URI) error {
	return svc.DIDManClient.DeleteService(id)
}

func (svc Service) enrichWithContactInfo(sp *domain.ServiceProvider) error {
	contactInformation, err := svc.DIDManClient.GetContactInformation(sp.Id)
	if err != nil {
		return domain.UnwrapAPIError(err)
	}
	if contactInformation != nil {
		sp.Email = contactInformation.Email
		sp.Name = contactInformation.Name
		sp.Phone = contactInformation.Phone
		sp.Website = contactInformation.Website
	}
	return nil
}

func (svc Service) Endpoints(sp domain.ServiceProvider) (domain.Endpoints, error) {
	document, _, err := svc.VDRClient.Get(sp.Id)
	if err != nil {
		return nil, domain.UnwrapAPIError(err)
	}
	endpoints := domain.Endpoints{}
	for _, svc := range document.Service {
		var endpoint string
		_ = svc.UnmarshalServiceEndpoint(&endpoint)
		if endpoint != "" {
			id := svc.ID.String()
			endpoints = append(endpoints, domain.Endpoint{
				EndpointID: domain.EndpointID{Id: id},
				EndpointProperties: domain.EndpointProperties{
					Type: svc.Type,
					Url:  endpoint,
				},
			})
		}
	}
	return endpoints, nil
}
func (svc Service) GetServices() (domain.Services, error) {
	spDID, err := svc.Repository.Get()
	if err != nil {
		return nil, err
	}
	services, err := svc.DIDManClient.GetCompoundServices(spDID.String())
	if err != nil {
		return nil, err
	}
	compoundServices := domain.Services{}
	for _, service := range services {
		compoundServices = append(compoundServices,  domain.Service{
			ServiceID:         domain.ServiceID{Id: service.Id},
			ServiceProperties: domain.ServiceProperties{
				ServiceEndpoint: service.ServiceEndpoint,
				Name:      service.Type,
			},
		})
	}
	return compoundServices, nil
}

func (svc Service) AddService(service domain.ServiceProperties) (*domain.Service, error) {
	spDID, err := svc.Repository.Get()
	if err != nil {
		return nil, err
	}
	endpoints := make(map[string]string, len(service.ServiceEndpoint))
	for key, val := range service.ServiceEndpoint {
		endpoints[key] = val.(string)
	}
	cs, err := svc.DIDManClient.AddCompoundService(spDID.String(), service.Name, endpoints)
	if err != nil {
		return nil, err
	}

	return &domain.Service{
		ServiceID:         domain.ServiceID{Id: cs.Id},
		ServiceProperties: domain.ServiceProperties{
			ServiceEndpoint: cs.ServiceEndpoint,
			Name:      cs.Type,
		},
	}, nil
}

