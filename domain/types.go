package domain

type OrganizationCredential struct {
	Id                string              `json:"id"`
	CredentialSubject []CredentialSubject `json:"credentialSubject"`
}

type CredentialSubject struct {
	ID           string       `json:"id"`
	Organization Organization `json:"organization"`
}

type Organization struct {
	Name string `json:"name"`
	City string `json:"city"`
}
