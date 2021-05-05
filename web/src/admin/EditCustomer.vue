<template>
  <modal-window
      :cancelRoute="{name: 'admin.customers'}"
      :confirmFn="checkForm"
      confirmText="Save Customer"
      title="Edit Customer" type="edit"
  >
    <div v-if="loading && !apiError">Loading...</div>
    <div class="p-2 rounded-md w-full bg-red-100" v-if="apiError">
      Error during server communication: {{ apiError }}
    </div>
    <customer-form mode="edit" :value="customer" @input="(newCustomer)=> {customer = newCustomer}" v-if="!loading"/>
  </modal-window>
</template>
<style>
input:checked + svg {
  display: block;
}
</style>

<script>
import ModalWindow from "../components/ModalWindow.vue";
import CustomerForm from "./CustomerForm.vue";

export default {
  components: {ModalWindow, CustomerForm},
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