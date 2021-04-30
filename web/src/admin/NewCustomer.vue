<template>
  <nrad-modal :cancelRoute="{name: 'admin.customers'}" :confirmFn="checkForm" confirmText="Connect Customer"
              title="Connect existing customer" type="add">

    <p class="mb-3 text-sm">Here you can link an existing customer to the Nuts network by creating a new Nuts DID.</p>

    <p v-if="apiError" class="p-3 bg-red-100 rounded-md">Could not connect customer: {{ apiError }}</p>

    <div class="p-3 bg-red-100 rounded-md" v-if="formErrors.length">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="error in formErrors">* {{ error }}</li>
      </ul>
    </div>


    <form class="space-y-3">
      <div>
        <label for="newCustomerIdInput">Internal customer ID</label>
        <input type="text" v-model="customer.id" id="newCustomerIdInput">
      </div>
      <div>
        <label for="newCustomerNameInput">Customer name</label>
        <input type="text" v-model="customer.name" id="newCustomerNameInput">
      </div>
    </form>
  </nrad-modal>
</template>

<script>
export default {
  data() {
    return {
      apiError: '',
      formErrors: [],
      customer: {
        id: '',
        name: '',
      }
    }
  },
  methods: {
    checkForm(e) {
      // reset the errors
      this.formErrors.length = 0
      this.apiError = ''

      if (this.customer.id && this.customer.name) {
        return this.confirm()
      }

      if (!this.customer.name) {
        this.formErrors.push("Name required")
      }

      if (!this.customer.id) {
        this.formErrors.push("Id required")
      }

      e.preventDefault()
    },
    confirm() {
      this.$api.post('web/customers', this.customer)
          .then(response => this.$router.push({name: 'admin.customers'}))
          .catch(response => this.apiError = response.statusText)
    }
  }
}
</script>