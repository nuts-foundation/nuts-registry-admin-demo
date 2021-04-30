<template>
  <div class="flex justify-between mb-6">
    <h1 class="text-3xl">Your Care Organisations</h1>
    <button
        class="bg-blue-400 hover:bg-blue-500 text-white font-medium rounded-md px-3 py-2"
        @click="$router.push({name: 'admin.newCustomer'})"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 inline" fill="none" viewBox="0 0 24 24"
           stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
      </svg>
      Add
    </button>
  </div>
  <div class="customer-container">
    <p v-if="fetchError" class="m-4">Could not fetch customers: {{ fetchError }}</p>
    <table v-if="customers.length > 0" class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
      <tr>
        <th class="thead">Customer ID</th>
        <th class="thead">Name</th>
        <th class="thead">DID</th>
      </tr>
      </thead>
      <tbody class="tbody">
      <tr v-for="customer in customers">
        <td class="tcell">
          {{ customer.id }}
        </td>
        <td class="tcell">{{ customer.name }}</td>
        <td class="tcell">{{ customer.did }}</td>
      </tr>
      </tbody>
    </table>
  </div>
  <router-view name="modal"></router-view>
</template>

<script>

export default {
  data() {
    return {
      fetchError: "",
      customers: []
    }
  },
  created() {
    // watch the params of the route to fetch the data again
    this.$watch(
        () => this.$route.params,
        () => {
          this.fetchData()
        },
        // fetch the data when the view is created and the data is
        // already being observed
        {immediate: true}
    )
  },
  methods: {
    fetchData() {
      this.$api.get('web/customers')
          .then(data => this.customers = data)
          .catch(response => {
            console.error("failure", response)
            if (response.status === 403) {
              this.fetchError = "Invalid credentials"
              return
            }
            console.log(response)
            this.fetchError = response.statusText
          })
    }
  }
}
</script>

<style scoped>

.thead {
  @apply px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider;
}

.tcell {
  @apply px-6 py-4 text-left;
}

.body {
  @apply bg-white divide-y divide-gray-200;
}

.customer-container {
  @apply shadow overflow-hidden border-gray-200 rounded;
}

</style>