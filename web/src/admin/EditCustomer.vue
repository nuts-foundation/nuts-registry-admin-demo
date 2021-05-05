<template>
  <modal-window :cancelRoute="{name: 'admin.customers'}" :confirmFn="checkForm" confirmText="Save Customer"
              title="Edit Customer" type="edit">
    <div v-if="loading && !apiError">Loading...</div>
    <div class="p-2 rounded-md w-full bg-red-100" v-if="apiError">Error during server communication: {{
        apiError
      }}
    </div>
    <form v-if="!loading" class="space-y-3">
      <div>
        <label for="customerIDInput">Internal customer ID:</label>
        <input type="text" disabled v-model="customer.id" id="customerIDInput">
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
        <label class="flex justify-start items-start" for="customerActiveInput">
          <div
              class="bg-white border-2 rounded border-gray-400 w-6 h-6 flex flex-shrink-0 justify-center items-center mr-2 focus-within:border-blue-500">
            <input class="opacity-0 absolute" id="customerActiveInput" type="checkbox" v-model="customer.active">
            <svg class="fill-current hidden w-4 h-4 text-green-500 pointer-events-none" viewBox="0 0 20 20">
              <path d="M0 11l2-2 5 5L18 3l2 2L7 18z"/>
            </svg>
          </div>
          <div>Customer activated</div>
        </label>
      </div>
    </form>
  </modal-window>
</template>
<style>
input:checked + svg {
  display: block;
}
</style>

<script>
import ModalWindow from "../components/ModalWindow.vue";
export default {
  components: {ModalWindow},
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
  emits: ["statusUpdate"],
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
        this.formErrors.push("To be active a town is required")
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
            this.$emit("statusUpdate", "Customer saved")
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