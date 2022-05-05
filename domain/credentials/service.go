package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/nuts-foundation/go-did/vc"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
	"github.com/sirupsen/logrus"

	didmanAPI "github.com/nuts-foundation/nuts-node/didman/api/v1"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain/sp"

	ssi "github.com/nuts-foundation/go-did"
	vcrApi "github.com/nuts-foundation/nuts-node/vcr/api/v2"
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

	if !shouldHaveCredential {
		if len(credentials) > 0 {
			if err := s.RevokeCredentials(credentials); err != nil {
				return fmt.Errorf("unable to revoke NutsOrgCredentials for customer %d: %w", customer.Id, err)
			}
		}

		return nil
	}

	if customer.City == nil {
		return fmt.Errorf("customer.City must be set for issuing a credential")
	}

	// Is there a single credential issued with the exact same organization / city name? Then we don't have to do anything
	if len(credentials) == 1 &&
		credentials[0].Organization.Name == customer.Name &&
		credentials[0].Organization.City == *customer.City {
		return nil
	}

	if len(credentials) > 0 {
		if err := s.RevokeCredentials(credentials); err != nil {
			return fmt.Errorf("unable to revoke NutsOrgCredentials for customer %d: %w", customer.Id, err)
		}
	}

	if err := s.issueNutsOrgCredential(customer); err != nil {
		return fmt.Errorf("unable to manage NutsOrgCredential for customer %d: %w", customer.Id, err)
	}

	return nil
}

func (s Service) GetCredentials(customer domain.Customer) ([]domain.OrganizationConceptCredential, error) {
	return s.search(vcrApi.SearchVCQuery{
		Type:    []ssi.URI{ssi.MustParseURI(credential.NutsOrganizationCredentialType), ssi.MustParseURI(vc.VerifiableCredentialType)},
		Context: []ssi.URI{ssi.MustParseURI(vc.VCContextV1), ssi.MustParseURI(credential.NutsContext)},
		CredentialSubject: []interface{}{domain.NutsOrganizationCredentialSubject{
			ID: *customer.Did,
		}},
	})
}

func (s Service) SearchOrganizations(name, city string) ([]domain.OrganizationConceptCredential, error) {
	return s.search(vcrApi.SearchVCQuery{
		Type:    []ssi.URI{ssi.MustParseURI(credential.NutsOrganizationCredentialType), ssi.MustParseURI(vc.VerifiableCredentialType)},
		Context: []ssi.URI{ssi.MustParseURI(vc.VCContextV1), ssi.MustParseURI(credential.NutsContext)},
		CredentialSubject: []interface{}{domain.NutsOrganizationCredentialSubject{
			Organization: domain.Organization{
				Name: name,
				City: city,
			},
		}},
	})
}

func (s Service) search(credential vcrApi.SearchVCQuery) ([]domain.OrganizationConceptCredential, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	queryAsjson, _ := json.MarshalIndent(credential, "", "  ")
	logrus.Warn(string(queryAsjson))
	defer cancel()

	response, err := s.client().SearchVCs(ctx, vcrApi.SearchVCsJSONRequestBody{Query: credential})

	if err != nil {
		return nil, domain.UnwrapAPIError(err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status 200: %s", response.Status)
	}
	searchResponse, err := vcrApi.ParseSearchVCsResponse(response)
	if err != nil {
		return nil, err
	}
	var results []domain.OrganizationConceptCredential
	for _, curr := range searchResponse.JSON200.VerifiableCredentials {
		var subjects []domain.NutsOrganizationCredentialSubject
		err = curr.VerifiableCredential.UnmarshalCredentialSubject(&subjects)
		if err != nil {
			return nil, err
		}
		results = append(results, domain.OrganizationConceptCredential{
			ID:           curr.VerifiableCredential.ID.String(),
			Issuer:       curr.VerifiableCredential.Issuer.String(),
			Organization: subjects[0].Organization,
			Subject:      subjects[0].ID,
		})
	}
	return results, nil
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
		// ignore so we can still see the DID
		logrus.Warnf("Unable to get contactinfo (did=%s)", id.String())
		return sp, nil
		//return nil, domain.UnwrapAPIError(err)
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

	var credentialSubject = make([]interface{}, 0)
	credentialSubject = append(credentialSubject, domain.NutsOrganizationCredentialSubject{ID: *customer.Did, Organization: domain.Organization{
		Name: customer.Name,
		City: *customer.City,
	}})

	visiblity := vcrApi.IssueVCRequestVisibilityPublic
	requestBody := vcrApi.IssueVCJSONRequestBody{
		Type:              "NutsOrganizationCredential",
		Issuer:            vendorDID.Id,
		CredentialSubject: credentialSubject,
		Visibility:        &visiblity,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := s.client().IssueVC(ctx, requestBody)
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
		response, err := s.client().RevokeVC(ctx, url.PathEscape(credential.ID))
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
