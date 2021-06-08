package customers

import (
	"fmt"
	"github.com/nuts-foundation/go-did/did"
	didmanAPI "github.com/nuts-foundation/nuts-node/didman/api/v1"
	nutsApi "github.com/nuts-foundation/nuts-node/vdr/api/v1"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
)

type Service struct {
	VDRClient    nutsApi.HTTPClient
	Repository   Repository
	DIDManClient didmanAPI.HTTPClient
}

func (s Service) ConnectCustomer(id, name, city string, serviceProviderID did.DID) (*domain.Customer, error) {
	selfControl := false
	capabilityInvocation := false
	controllers := []string{serviceProviderID.String()}

	didDoc, err := s.VDRClient.Create(nutsApi.DIDCreateRequest{
		SelfControl:          &selfControl,
		Controllers:          &controllers,
		CapabilityInvocation: &capabilityInvocation,
	})
	if err != nil {
		return nil, domain.UnwrapAPIError(err)
	}

	customer := domain.Customer{
		Did:  didDoc.ID.String(),
		Id:   id,
		Name: name,
		City: &city,
	}

	return s.Repository.NewCustomer(customer)
}

func (s Service) EnableService(customerID string, spDID string, serviceType string) error {
	customer, err := s.Repository.FindByID(customerID)
	if err != nil {
		return err
	}
	ref := "ref://" + spDID + "serviceEndpoints?type=" + serviceType

	_, err = s.DIDManClient.AddEndpoint(customer.Did, serviceType, ref)
	if err != nil {
		return fmt.Errorf("unable to add new service reference to did doc: %w", err)
	}
	return nil
}

func (s Service) DisableService(customerID, serviceType string) error {
	customer, err := s.Repository.FindByID(customerID)
	if err != nil {
		return err
	}
	// Add Delete endpoint to didman un the nuts-node
	return s.DIDManClient.DeleteEndpoint(customer.Did, serviceType)
}

func (s Service) ManageServices(customerDIDStr string, serviceIDStr []string) ([]did.Service, error) {
	customerDIDDoc, _, err := s.VDRClient.Get(customerDIDStr)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch customer DID Document: %w", err)
	}

	currentServices := customerDIDDoc.Service
	// Contains IDs of compound services in the Service Providers DID Doc
	addedServices := []string{}
	// Contains IDs of services in the customers DID Doc
	removedServices := []string{}

	for _, cs := range customerDIDDoc.Service {
		serviceID, ok := cs.ServiceEndpoint.(string)
		if !ok {
			continue
		}
		found := false
		for _, newServiceID := range serviceIDStr {
			if serviceID == newServiceID {
				found = true
				break
			}
		}
		if !found {
			removedServices = append(removedServices, cs.ID.String())
		}
	}

	for _, newServiceID := range serviceIDStr {
		found := false
		for _, cs := range currentServices {
			serviceID, ok := cs.ServiceEndpoint.(string)
			if !ok {
				continue
			}
			if serviceID == newServiceID {
				found = true
				break
			}
		}
		if !found {
			addedServices = append(addedServices, newServiceID)
		}
	}

	for _, addedServiceID := range addedServices {
		serviceDIDURL, err := did.ParseDID(addedServiceID)
		if err != nil {
			return nil, fmt.Errorf("new service DID has incorrect format: %w", err)
		}
		serviceDIDURL.Fragment = ""
		doc, _, err := s.VDRClient.Get(serviceDIDURL.String())
		if err != nil {
			return nil, fmt.Errorf("unable to resolve did doc for new service %w", err)
		}
		referencedService := findService(doc.Service, addedServiceID)
		if referencedService == nil {
			fmt.Errorf("new service not found on DID")
		}
		ref := serviceDIDURL.String() + "?type=" + referencedService.Type

		_, err = s.DIDManClient.AddEndpoint(customerDIDStr, referencedService.Type, ref)
		if err != nil {
			return nil, fmt.Errorf("unable to add new service reference to did doc: %w", err)
		}
	}

	//todo: remove services

	customerDIDDoc, _, err = s.VDRClient.Get(customerDIDStr)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch customer DID Document: %w", err)
	}
	return customerDIDDoc.Service, nil
}

func findService(services []did.Service, serviceIDStr string) *did.Service {
	for _, s := range services {
		if s.ID.String() == serviceIDStr {
			return &s
		}
	}
	return nil
}

func (s Service) GetServices(customerID string) ([]did.Service, error) {
	customer, err := s.Repository.FindByID(customerID)
	if err != nil {
		return nil, err
	}

	customerDIDDoc, _, err := s.VDRClient.Get(customer.Did)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch customer DID Document: %w", err)
	}

	return customerDIDDoc.Service, nil
}
