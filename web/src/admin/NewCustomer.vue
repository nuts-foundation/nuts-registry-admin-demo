<template>
  <nrad-modal :cancelRoute="{name: 'admin.customers'}" :confirmFn="confirm" confirmText="Connect Customer" title="Connect existing customer" type="add">
    <p class="mb-3 text-sm">Here you can link an existing customer to the Nuts network by creating a new Nuts DID.</p>
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
      customer: {
        id: '123',
        name: 'CareOrg',
      }
    }
  },
  methods: {
    confirm() {
      this.$api.post('web/customers', this.customer)
          .then(response => this.$router.push({name: 'admin.customers'}))
          .catch(response => this.apiError = response.statusText)
    }
  }
}
</script>