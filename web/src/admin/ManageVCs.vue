<template>
  <div class="mb-6">
    <h1 class="text-3xl">Manage Verifiable Credentials</h1>
  </div>
  <div class="border rounded-md">
    <div v-for="(issuers, type) in credentialIssuers">
      <div class="bg-gray-200 p-2">{{ type }}</div>

      <ul v-for="issuer in issuers">
        <li class="flex justify-between p-2">
          <span>{{ issuer.serviceProvider.id }} - {{ issuer.serviceProvider.name }}</span>
          <form-checkbox  v-model="issuer.trusted" @update:modelValue="toggleTrust(type, issuer)">Trusted</form-checkbox>
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
      credentialIssuers: {}
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
          })
    },
    toggleTrust(type, issuer) {
      console.log("toggle", type, issuer)
      this.fetchIssuers()
    },
  },
  mounted() {
    this.fetchIssuers()
  }
}
</script>