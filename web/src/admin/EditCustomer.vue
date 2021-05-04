<template>
  <nrad-modal :cancelRoute="{name: 'admin.customers'}" :confirmFn="checkForm" confirmText="Save Customer"
              title="Edit Customer" type="edit">
    <div v-if="loading && !apiError">Loading...</div>
    <div class="p-2 rounded-md w-full bg-red-100" v-if="apiError">Error during server communication: {{ apiError }}</div>
    <form v-if="!loading" class="space-y-3">
      <div>
        <label>Internal customer ID:</label>
        <div>
          {{ customer.id }}
        </div>
      </div>
      <div>
        <label for="customerNameInput">Customer name</label>
        <input type="text" v-model="customer.name" id="customerNameInput">
      </div>
      <div>
        <label for="customerTownInput">Town</label>
        <input type="text" v-model="customer.town" id="customerTownInput">
      </div>
      <div>
        <label for="customerActiveInput">Customer activated</label>
        <input id="customerActiveInput" type="checkbox" v-model="customer.active">
      </div>
    </form>
  </nrad-modal>
</template>

<script>
export default {
  data() {
    return {
      customer: {
        name: '',
        town: ''
      },
      formErrors: [],
      apiError: '',
      loading: true,
    }
  },
  watch: {
    "$route.params": {
      handler(toParams, previousParams) {
        if (toParams && 'id' in toParams) {
          this.fetchCustomer(toParams.id)
        }
      },
      immediate: true
    }

  },
  methods: {
    checkForm() {
      // reset the errors
      this.formErrors.length = 0
      this.apiError = ''

      if (this.customer.name) {
        return this.saveCustomer()
      }

      if (this.active && !this.customer.town) {
        this.formErrors.push("Name required")
      }

      if (!this.customer.name) {
        this.formErrors.push("Name required")
      }
    },
    fetchCustomer(id) {
      console.log("id: ", id)
      this.$api.get(`web/customers/${id}`)
          .then((customer) => {
            this.customer = customer
            this.loading = false
          })
          .catch((reason) => {
            this.apiError = reason.statusText
            console.log("failed:", reason)
          })
    },
    saveCustomer() {
      this.$api.put(`web/customers/${this.customer.id}`, this.customer)
          .then((customer) => {
            this.customer = customer
            this.$router.push({name: 'admin.customers'})
          })
          .catch((reason) => {
            this.apiError = reason.statusText
            console.log("failed:", reason)
          })
    }
  }
}
</script>