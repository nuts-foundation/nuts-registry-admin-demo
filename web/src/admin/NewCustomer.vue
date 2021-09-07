<template>
  <modal-window :cancelRoute="{name: 'admin.customers'}" :confirmFn="checkForm" confirmText="Connect Customer"
              title="Connect existing customer" type="add">

    <p class="mb-3 text-sm">Here you can link an existing customer to the Nuts network by creating a new Nuts DID.</p>

    <p v-if="apiError" class="p-3 bg-red-100 rounded-md">Could not connect customer: {{ apiError }}</p>

    <div class="p-3 bg-red-100 rounded-md" v-if="formErrors.length">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="error in formErrors">* {{ error }}</li>
      </ul>
    </div>
    <customer-form mode="new" :value="customer" @input="(newCustomer)=> {customer = newCustomer}"/>

  </modal-window>
</template>

<script>
import ModalWindow from "../components/ModalWindow.vue";
import CustomerForm from "./CustomerForm.vue";

export default {
  components: {ModalWindow, CustomerForm},
  data() {
    return {
      apiError: '',
      formErrors: [],
      customer: {
        id: '',
        name: '',
        city: '',
        domain: '',
        active: false
      }
    }
  },
  emits: ["statusUpdate"],
  methods: {
    checkForm(e) {
      // reset the errors
      this.formErrors.length = 0
      this.apiError = ''

      if (this.customer.id && this.customer.name) {
        return this.confirm()
      }

      if (!this.customer.id) {
        this.formErrors.push("ID required")
      }

      if (!this.customer.name) {
        this.formErrors.push("Name required")
      }
      e.preventDefault()
    },
    confirm() {
      this.$api.post('web/private/customers', this.customer)
          .then(response => {
            this.$emit("statusUpdate", "Customer connected")
            this.$router.push({name: 'admin.customers'})
          })
          .catch(response => this.apiError = response.statusText)
    }
  }
}
</script>
