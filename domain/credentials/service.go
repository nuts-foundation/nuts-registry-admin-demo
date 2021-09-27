package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	didmanAPI "github.com/nuts-foundation/nuts-node/didman/api/v1"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain/sp"

	ssi "github.com/nuts-foundation/go-did"
	vcrApi "github.com/nuts-foundation/nuts-node/vcr/api/v1"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
)

type Service struct {
	NutsNodeAddr string
	SPService    sp.Service
	DIDManClient didmanAPI.HTTPClient
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
		if customer.City == nil {
			return fmt.Errorf("customer.City must be set for issuing a credential")
		}
		issue := true
		for _, credential := range credentials {
			if credential.Organization.Name == customer.Name && credential.Organization.City == *customer.City {
				// Customer already has credential with given name/ & city, so we don't need to issue one.
				issue = false
				break
			}
		}
		if issue {
			err = s.issueNutsOrgCredential(customer)
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
		return fmt.Errorf("unable to manage NutsOrgCredential for customer %d: %w", customer.Id, err)
	}
	return nil
}

func (s Service) GetCredentials(customer domain.Customer) ([]domain.OrganizationConceptCredential, error) {
	searchBody := vcrApi.SearchJSONRequestBody{
		Params: []vcrApi.KeyValuePair{
			{Key: "subject", Value: *customer.Did},
		},
	}

	return s.search(searchBody)
}

func (s Service) SearchOrganizations(name, city string) ([]domain.OrganizationConceptCredential, error) {
	searchBody := vcrApi.SearchJSONRequestBody{
		Params: []vcrApi.KeyValuePair{
			{Key: "organization.name", Value: name},
			{Key: "organization.city", Value: city},
		},
	}
	return s.search(searchBody)
}

func (s Service) search(searchBody vcrApi.SearchJSONRequestBody) ([]domain.OrganizationConceptCredential, error) {
	i := 0
	for _, param := range searchBody.Params {
		if len(strings.TrimSpace(param.Value)) > 0 {
			searchBody.Params[i] = param
			i++
		}
	}
	searchBody.Params = searchBody.Params[:i]

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	response, err := s.client().Search(
		ctx,
		"organization",
		searchBody,
	)
	if err != nil {
		return nil, domain.UnwrapAPIError(err)
	}

	body, _ := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status 200: %s", response.Status)
	}
	var credentials = make([]domain.OrganizationConceptCredential, 0)
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
			issuer, err := s.getIssuer(id)
			if err != nil {
				return result, err
			}
			issuers[i] = domain.CredentialIssuer{Trusted: true, ServiceProvider: *issuer}
		}
		for i, id := range untrustedDIDs {
			issuer, err := s.getIssuer(id)
			if err != nil {
				return result, err
			}
			issuers[len(trustedDIDs)+i] = domain.CredentialIssuer{Trusted: false, ServiceProvider: *issuer}
		}
		result.Set(credential, issuers)
	}
	return result, nil
}

func (s Service) getIssuer(id ssi.URI) (*domain.ServiceProvider, error) {
	sp := &domain.ServiceProvider{Id: id.String()}
	if id.Scheme != "did" {
		return sp, nil
	}
	contactInformation, err := s.DIDManClient.GetContactInformation(sp.Id)
	if err != nil {
		return nil, domain.UnwrapAPIError(err)
	}
	if contactInformation != nil {
		sp.Email = contactInformation.Email
		sp.Name = contactInformation.Name
		sp.Phone = contactInformation.Phone
		sp.Website = contactInformation.Website
	}
	return sp, nil
}

func (s Service) fetchCredentialIssuers(credential string, clientFn func(ctx context.Context, credentialType string, reqEditors ...vcrApi.RequestEditorFn) (*http.Response, error)) ([]ssi.URI, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	response, err := clientFn(ctx, credential)
	defer cancel()
	if err != nil {
		return nil, err
	}

	body, _ := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, err
	}
	var issuerDIDs []ssi.URI
	err = json.Unmarshal(body, &issuerDIDs)
	return issuerDIDs, err

}

func (s Service) issueNutsOrgCredential(customer domain.Customer) error {
	vendorDID, err := s.SPService.Get()
	if err != nil {
		return err
	}
	if vendorDID == nil {
		return err
	}

	logrus.Infof("Issuing NutsOrganizationCredential (did=%s,name=%s,city=%s)", *customer.Did, customer.Name, *customer.City)

	issuerURI, _ := ssi.ParseURI(vendorDID.Id)
	typeURI, _ := ssi.ParseURI("NutsOrganizationCredential")
	var credentialSubject = make([]interface{}, 0)
	credentialSubject = append(credentialSubject, domain.NutsOrganizationCredentialSubject{ID: *customer.Did, Organization: domain.Organization{
		Name: customer.Name,
		City: *customer.City,
	}})
	requestBody := vcrApi.CreateJSONRequestBody{
		Type:              []ssi.URI{*typeURI},
		Issuer:            *issuerURI,
		CredentialSubject: credentialSubject,
	}

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

func (s Service) RevokeCredentials(credentials []domain.OrganizationConceptCredential) error {
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
		response, err := s.client().Revoke(ctx, url.PathEscape(credential.ID))
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

func (s Service) ManageIssuerTrust(credentialType string, issuerID ssi.URI, trusted bool) (*domain.CredentialIssuer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		response *http.Response
		err      error
	)

	if trusted {
		requestBody := vcrApi.TrustIssuerJSONRequestBody{
			CredentialType: credentialType,
			Issuer:         issuerID.String(),
		}
		response, err = s.client().TrustIssuer(ctx, requestBody)
	} else {
		requestBody := vcrApi.UntrustIssuerJSONRequestBody{
			CredentialType: credentialType,
			Issuer:         issuerID.String(),
		}
		response, err = s.client().UntrustIssuer(ctx, requestBody)
	}
	if response.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("expected status 204: %s", response.Status)
	}

	if err != nil {
		return nil, domain.UnwrapAPIError(err)
	}

	sp, err := s.getIssuer(issuerID)
	if err != nil {
		return nil, domain.UnwrapAPIError(err)
	}
	return &domain.CredentialIssuer{
		ServiceProvider: *sp,
		Trusted:         trusted,
	}, nil
}
