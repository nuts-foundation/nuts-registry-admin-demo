package credentials

import (
	ssi "github.com/nuts-foundation/go-did"
	v2 "github.com/nuts-foundation/nuts-node/vcr/api/vcr/v2"
)

// SearchVCRequest is the request body for searching VCs
type SearchVCRequest struct {
	// A partial VerifiableCredential in JSON-LD format. Each field will be used to match credentials against. All fields MUST be present.
	Query         SearchVCQuery     `json:"query"`
	SearchOptions *v2.SearchOptions `json:"searchOptions,omitempty"`
}

// SearchVCQuery defines a helper struct to search for VerifiableCredentials.
type SearchVCQuery struct {
	// Context defines the json-ld context to dereference the URIs
	Context []ssi.URI `json:"@context"`
	// Type holds multiple types for a credential. A credential must always have the 'VerifiableCredential' type.
	Type []ssi.URI `json:"type,omitempty"`
	// CredentialSubject holds the actual data for the credential. It must be extracted using the UnmarshalCredentialSubject method and a custom type.
	CredentialSubject interface{} `json:"credentialSubject,omitempty"`
}
