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

func (s Service) ConnectCustomer(reqCustomer domain.Customer, serviceProviderID did.DID) (*domain.Customer, error) {
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

	did := didDoc.ID.String()
	customer := domain.Customer{
		Did:  &did,
		Id:   reqCustomer.Id,
		Name: reqCustomer.Name,
		City: reqCustomer.City,
		Domain: reqCustomer.Domain,
	}

	return s.Repository.NewCustomer(customer)
}

const refTemplate = "%s/serviceEndpoint?type=%s"

// EnableService enables a service for a customer adding a reference by type to the compoundService
// to the customers DID document.
func (s Service) EnableService(customerID int, spDID string, serviceType string) error {
	customer, err := s.Repository.FindByID(customerID)
	if err != nil {
		return err
	}
	parsedDID, err := did.ParseDIDURL(spDID)
	if err != nil {
		return err
	}
	parsedDID.Fragment = ""

	ref := fmt.Sprintf(refTemplate, parsedDID.String(), serviceType)

	_, err = s.DIDManClient.AddEndpoint(*customer.Did, serviceType, ref)
	if err != nil {
		return fmt.Errorf("unable to add new service reference to DID Document: %w", err)
	}
	return nil
}

// DisableService disables a service for a customer by removing all references to a
// compoundService of a certain type from the customers DID document.
func (s Service) DisableService(customerID int, serviceType string) error {
	customer, err := s.Repository.FindByID(customerID)
	if err != nil {
		return err
	}
	return s.DIDManClient.DeleteEndpointsByType(*customer.Did, serviceType)
}

// GetServices returns all the enabled services for a customer.
func (s Service) GetServices(customerID int) ([]did.Service, error) {
	customer, err := s.Repository.FindByID(customerID)
	if err != nil {
		return nil, err
	}

	customerDIDDoc, _, err := s.VDRClient.Get(*customer.Did)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch customer DID Document: %w", err)
	}

	return customerDIDDoc.Service, nil
}
