package domain

import "errors"

var ErrNutsNodeUnreachable = errors.New("nuts node unreachable")

// OrganizationConceptCredential models an organization as concept, as returned by the Nuts Node's Verifiable Credential Search API.
type OrganizationConceptCredential struct {
	// ID contains the identifier that uniquely identifies this Verifiable Credential
	ID string `json:"id"`
	// Issuer contains the identifier of the entity that issued the Verifiable Credential
	Issuer string `json:"issuer"`
	// Organization contains the properties of the organization concept that was mapped from the Verifiable Credential.
	Organization Organization `json:"organization"`
	// Subject contains the identifier that uniquely identifies the holder of the Verifiable Credential,
	// the entity to which the credential was issued.
	Subject string `json:"subject"`
}

// NutsOrganizationCredentialSubject models the subject for a Verifiable Credential of type NutsOrganizationCredential
type NutsOrganizationCredentialSubject struct {
	ID           string       `json:"id"`
	Organization Organization `json:"organization"`
}

// Organization models the properties for a legally registered organization.
type Organization struct {
	Name string `json:"name"`
	City string `json:"city"`
}

const NutsCommService = "NutsComm"
