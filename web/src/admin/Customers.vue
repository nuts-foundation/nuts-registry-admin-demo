<template>
  <h1 class="py-4 font-medium text-xl">Customers</h1>
  <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
    <table class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
      <tr>
        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
        <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">DID</th>
      </tr>
      </thead>
      <tbody class="bg-white divide-y divide-gray-200">
      <tr v-for="customer in customers">
        <td><div class="m-4">
          {{ customer.name }}
        </div></td>
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
      fetch("api/customers")
          .then(response => response.json())
          .then(data => this.customers = data)
    }
  },
}
</script>