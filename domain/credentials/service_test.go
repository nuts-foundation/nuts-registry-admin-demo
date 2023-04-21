package credentials

import (
	"encoding/json"
	ssi "github.com/nuts-foundation/go-did"
	"github.com/nuts-foundation/go-did/vc"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
	"github.com/nuts-foundation/nuts-registry-admin-demo/domain"
	"testing"
)

func TestService_search(t *testing.T) {
	datra, _ := json.Marshal(SearchVCQuery{
		Type:    []ssi.URI{ssi.MustParseURI(credential.NutsOrganizationCredentialType), ssi.MustParseURI(vc.VerifiableCredentialType)},
		Context: []ssi.URI{ssi.MustParseURI(vc.VCContextV1), ssi.MustParseURI(credential.NutsV1Context)},
		CredentialSubject: domain.NutsOrganizationCredentialSubject{
			ID: "nuts:did:123",
			Organization: domain.Organization{
				Name: "*",
				City: "*",
			},
		},
	})
	println(string(datra))
}
