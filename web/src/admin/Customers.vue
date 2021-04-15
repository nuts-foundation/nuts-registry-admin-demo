<template>
  <h1 class="py-4 font-medium text-xl">Customers</h1>
  <div class="customer-container">
    <p v-if="fetchError" class="m-4">Could not fetch customers: {{ fetchError }}</p>
    <table v-if="customers.length > 0" class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
      <tr>
        <th class="thead">Name</th>
        <th class="thead">DID</th>
      </tr>
      </thead>
      <tbody class="tbody">
      <tr v-for="customer in customers">
        <td>
          <div class="m-4">
            {{ customer.name }}
          </div>
        </td>
        <td>{{ customer.did }}</td>
      </tr>
      </tbody>
    </table>
  </div>
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
      fetch("web/customers", {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem("session")}`
        }
      }).then(response => {
        if (!response.ok) {
          if (response.status == 403) {
            throw "Invalid credentials"
          }
          throw response.statusText
        }
        return response.json()
      }).then(data => this.customers = data)
          .catch(reason => {
            console.log(reason)
            this.fetchError = reason
          })
    }
  }
}
</script>

<style scoped>

.thead {
  @apply px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider;
}

.body {
  @apply bg-white divide-y divide-gray-200;
}

.customer-container {
  @apply shadow overflow-hidden border-b border-gray-200 sm:rounded-lg;
}

</style>