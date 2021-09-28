<template>
  <h1>Manage Verifiable Credentials</h1>

  <div class="mt-8 bg-white p-5 shadow-lg rounded-lg">
    <p v-if="!!fetchError" class="m-4">Could not fetch credential issuers: {{ fetchError }}</p>
    <div class="m-4" v-if="loading">Loading...</div>
    <div class="m-4" v-if="!loading && Object.keys(credentialIssuers).length === 0 && !fetchError">No credential issuers yet.</div>
    <div v-for="(issuers, type) in credentialIssuers">
      <div class="font-medium p-2">{{ type }}</div>

      <ul v-for="issuer in issuers">
        <li class="flex justify-between p-2">
          <span>{{ issuer.serviceProvider.id }} - {{ issuer.serviceProvider.name }}</span>
          <form-checkbox v-model="issuer.trusted" @update:modelValue="toggleTrust(type, issuer)">Trusted</form-checkbox>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import FormCheckbox from "../components/FormCheckbox.vue";

export default {
  components: {FormCheckbox},
  emits: ['statusUpdate'],
  data() {
    return {
      credentialIssuers: {},
      loading: true,
      fetchError: ""
    }
  },
  methods: {
    fetchIssuers() {
      this.$api.get("web/private/credentials/issuers")
          .then((response) => {
            this.credentialIssuers = response
          })
          .catch((reason) => {
            console.log("fetch failed: ", reason)
            this.fetchError = reason
          })
          .finally(() => this.loading = false)
    },
    toggleTrust(type, issuer) {
      console.log("toggle", type, issuer)
      this.$api.put(`web/private/credential/${type}/issuer/${encodeURIComponent(issuer.serviceProvider.id)}`, {trusted: issuer.trusted})
          .then((response) => {
          }).catch((reason => console.log("update status failed:", reason)))
          .finally(this.fetchIssuers)
    },
  },
  mounted() {
    this.fetchIssuers()
  }
}
</script>
