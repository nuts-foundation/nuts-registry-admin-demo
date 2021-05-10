<template>
  <h1>Manage Verifiable Credentials</h1>
  <div v-for="(issuers, type) in credentialIssuers">
    <span>Type: {{type}}</span>
    <div>
      <h1>Trusted Issuers:</h1>
      <ul v-for="issuer in issuers.trusted">
        <li>{{issuer.did}}{{issuer.name}}</li>
      </ul>
    </div>
    <div>
      <h1>Untrusted Issuers:</h1>
      <ul v-for="issuer in issuers.untrusted">
        <li>{{issuer.did}}{{issuer.name}}</li>
      </ul>
    </div>
  </div>
</template>
<script>
export default {
  emits: ['statusUpdate'],
  data() {
    return {
      credentialIssuers: {
        'NutsOrganizationCredential': {trusted: [{did: "did:nuts:123", name:"CareSoft NV"}, {did: "did:nuts:456", name:"Should we Care?"}], untrusted: [{did: "did:nuts:666", name: "EvilSoft"}]},
        'LRZACredential': {trusted: [{did: "did:web:lrza.nl", name: "LRZa"}], untrusted: [{did: "did:nuts:666", name: "EvilSoft"}]}
      }
    }
  },
  methods: {
    fetchTrusted(credentialType) {

    },
    fetchUntrusted(credentialType) {

    },
    trustIssuer(credentialType, issuerID) {

    },
    untrustIssuer(credentialType, issuerID) {

    }
  }
}
</script>