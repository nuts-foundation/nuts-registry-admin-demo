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
    <div class="m-4" v-if="loading">Loading...</div>
    <div class="m-4" v-if="!loading && customers.length == 0 && !fetchError">No customers yet, add one!</div>
    <table v-if="customers.length > 0" class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
      <tr>
        <th class="thead">Customer ID</th>
        <th class="thead">Name</th>
        <th class="thead">Town</th>
        <th class="thead">Published</th>
      </tr>
      </thead>
      <tbody class="tbody">
      <tr class="hover:bg-gray-100 cursor-pointer"
          @click="$router.push({name: 'admin.editCustomer', params: {id: customer.id} })"
          v-for="customer in customers">
        <td class="tcell">
          {{ customer.id }}
        </td>
        <td class="tcell">{{ customer.name }}</td>
        <td class="tcell">{{ customer.town }}</td>
        <td class="tcell">{{ customer.active }}</td>
      </tr>
      </tbody>
    </table>
  </div>
  <router-view name="modal" @statusUpdate="updateStatus"></router-view>
</template>

<script>

export default {
  data() {
    return {
      fetchError: "",
      customers: [],
      loading: true,
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
  emits: ["statusUpdate"],
  methods: {
    updateStatus(event) {
      this.$emit("statusUpdate", event)
    },
    openCustomer(customer) {
      console.log("open customer", customer.name)

    },
    fetchData() {
      this.$api.get('web/private/customers')
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
          .finally(() => this.loading = false)
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