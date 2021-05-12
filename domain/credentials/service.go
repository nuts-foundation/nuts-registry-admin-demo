package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/nuts-foundation/go-did/did"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain/sp"

	ssi "github.com/nuts-foundation/go-did"
	vcrApi "github.com/nuts-foundation/nuts-node/vcr/api/v1"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
)

type Service struct {
	NutsNodeAddr string
	SPService    sp.Service
}

func (s Service) client() vcrApi.ClientInterface {
	url := s.NutsNodeAddr

	response, err := vcrApi.NewClientWithResponses(url)
	if err != nil {
		panic(err)
	}
	return response
}

func (s Service) ManageNutsOrgCredential(customer domain.Customer, shouldHaveCredential bool) error {
	credentials, err := s.GetCredentials(customer)
	if err != nil {
		return err
	}

	if shouldHaveCredential {
		if len(credentials) > 0 {
			// already has credential
			return nil
		} else {
			err = s.IssueNutsOrgCredential(customer)
		}
	} else {
		if len(credentials) == 0 {
			// no credential to revoke
			return nil
		} else {
			err = s.RevokeCredentials(credentials)
		}
	}
	if err != nil {
		return fmt.Errorf("unable to manage NutsOrgCredential for customer %s: %w", customer.Id, err)
	}
	return nil
}

func (s Service) GetCredentials(customer domain.Customer) ([]domain.OrganizationCredential, error) {
	// not enough info, return
	if customer.Town == nil {
		return []domain.OrganizationCredential{}, nil
	}

	searchBody := vcrApi.SearchJSONRequestBody{
		Params: []vcrApi.KeyValuePair{
			{Key: "organization.name", Value: customer.Name},
			{Key: "organization.city", Value: *customer.Town},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	response, err := s.client().Search(
		ctx,
		"organization",
		searchBody,
	)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, domain.ErrNutsNodeUnreachable
		}
		return nil, err
	}

	body, _ := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status 200: %s", response.Status)
	}

	var credentials = make([]domain.OrganizationCredential, 0)
	if err = json.Unmarshal(body, &credentials); err != nil {
		return nil, err
	}
	return credentials, nil
}

func (s Service) GetCredentialIssuers(credentials []string) (domain.CredentialIssuers, error) {
	result := domain.CredentialIssuers{}
	for _, credential := range credentials {

		trustedDIDs, err := s.fetchCredentialIssuers(credential, s.client().ListTrusted)
		if err != nil {
			return result, err
		}
		untrustedDIDs, err := s.fetchCredentialIssuers(credential, s.client().ListUntrusted)
		if err != nil {
			return result, err
		}
		issuers := make([]domain.CredentialIssuer, len(trustedDIDs)+len(untrustedDIDs))
		for i, id := range trustedDIDs {
			s.SPService.Get()
			issuer, err := s.GetIssuer(id)
			if err != nil {
				return result, nil
			}
			issuers[i] = domain.CredentialIssuer{Trusted: true, ServiceProvider: *issuer}
		}
		for i, id := range untrustedDIDs {
			issuer, err := s.GetIssuer(id)
			if err != nil {
				return result, nil
			}
			issuers[len(trustedDIDs)+i] = domain.CredentialIssuer{Trusted: false, ServiceProvider: *issuer}
		}
		result.Set(credential, issuers)
	}
	return result, nil
}

func (s Service) GetIssuer(id did.DID) (*domain.ServiceProvider, error) {
	sp := &domain.ServiceProvider{Id: id.String()}
	contactInformation, err := s.SPService.DIDManClient.GetContactInformation(sp.Id)
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

func (s Service) fetchCredentialIssuers(credential string, clientFn func(ctx context.Context, credentialType string) (*http.Response, error)) ([]did.DID, error) {
	var issuerDIDs []did.DID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	response, err := clientFn(ctx, credential)
	defer cancel()
	if err != nil {
		if _, ok := err.(net.Error); ok {
			return issuerDIDs, domain.ErrNutsNodeUnreachable
		}
		return issuerDIDs, err
	}

	body, _ := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		return issuerDIDs, err
	}
	err = json.Unmarshal(body, &issuerDIDs)
	return issuerDIDs, err

}

func (s Service) IssueNutsOrgCredential(customer domain.Customer) error {
	vendorDID, err := s.SPService.Get()
	if err != nil {
		return err
	}
	if vendorDID == nil {
		return err
	}

	if customer.Town == nil {
		return fmt.Errorf("customer.Town must be set for issuing a credential")
	}

	issuerURI, _ := ssi.ParseURI(vendorDID.Id)
	typeURI, _ := ssi.ParseURI("NutsOrganizationCredential")
	var credentialSubject = make([]interface{}, 0)
	credentialSubject = append(credentialSubject, domain.CredentialSubject{ID: customer.Did, Organization: domain.Organization{
		Name: customer.Name,
		City: *customer.Town,
	}})
	requestBody := vcrApi.CreateJSONRequestBody{
		Type:              []ssi.URI{*typeURI},
		Issuer:            *issuerURI,
		CredentialSubject: credentialSubject,
	}
	log.Printf("issue NutsOrgCredential requestBody :%+v", requestBody)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := s.client().Create(ctx, requestBody)
	if err != nil {
		return err
	}

	responseBody, _ := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		return errors.New(string(responseBody))
	}
	return nil
}

func (s Service) RevokeCredentials(credentials []domain.OrganizationCredential) error {
	vendorDID, err := s.SPService.Get()
	if err != nil {
		return err
	}
	if vendorDID == nil {
		return errors.New("no vendor DID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, credential := range credentials {
		response, err := s.client().Revoke(ctx, url.PathEscape(credential.Id))
		if err != nil {
			return err
		}

		responseBody, _ := io.ReadAll(response.Body)
		if response.StatusCode != http.StatusOK {
			return errors.New(string(responseBody))
		}
	}

	return nil
}

func (s Service) ManageIssuerTrust(credentialType string, issuerID did.DID, trusted bool) (*domain.CredentialIssuer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		response *http.Response
		err      error
	)

	if trusted {
		requestBody := vcrApi.TrustIssuerJSONRequestBody{
			CredentialType: credentialType,
			Issuer:         issuerID.ID,
		}
		response, err = s.client().TrustIssuer(ctx, requestBody)
	} else {
		requestBody := vcrApi.UntrustIssuerJSONRequestBody{
			CredentialType: credentialType,
			Issuer:         issuerID.ID,
		}
		response, err = s.client().UntrustIssuer(ctx, requestBody)
	}
	if response.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("expected status 204: %s", response.Status)
	}

	if err != nil {
		return nil, unwrapAPIError(err)
	}

	sp, err := s.GetIssuer(issuerID)
	if err != nil {
		return nil, unwrapAPIError(err)
	}
	return &domain.CredentialIssuer{
		ServiceProvider: *sp,
		Trusted:         trusted,
	}, nil
}

func unwrapAPIError(err error) error {
	if _, ok := err.(net.Error); ok {
		return domain.ErrNutsNodeUnreachable
	}
	return err
}
