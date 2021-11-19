<template>
  <div>
    <div class="flex justify-between mb-6">
      <h1>Your care organizations</h1>

      <button
          @click="$router.push({name: 'admin.newCustomer'})"
          class="float-right inline-flex items-center bg-nuts w-10 h-10 rounded-lg justify-center shadow-md"
      >
        <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#fff">
          <path d="M0 0h24v24H0V0z" fill="none"/>
          <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
        </svg>
      </button>
    </div>

    <div class="mt-8 bg-white p-5 shadow-lg rounded-lg">
      <p v-if="fetchError" class="m-4">Could not fetch customers: {{ fetchError }}</p>
      <div class="m-4" v-if="loading">Loading...</div>
      <div class="m-4" v-if="!loading && customers.length === 0 && !fetchError">No customers yet, add one!</div>
      <table v-if="customers.length > 0" class="min-w-full divide-y divide-gray-200">
        <thead>
        <tr>
          <th class="thead">Customer ID</th>
          <th class="thead">Name</th>
          <th class="thead">City</th>
          <th class="thead">Published</th>
          <th class="thead"></th>
        </tr>
        </thead>
        <tbody class="tbody">
        <tr class="hover:bg-gray-100 cursor-pointer"
            :key="customer.id"
            @click="$router.push({name: 'admin.editCustomer', params: {id: customer.id} })"
            v-for="customer in customers">
          <td class="tcell">
            {{ customer.id }}
          </td>
          <td class="tcell">{{ customer.name }}</td>
          <td class="tcell">{{ customer.city }}</td>
          <td class="tcell">{{ customer.active }}</td>
          <td class="tcell">
            <!-- Delete -->
            <svg xmlns="http://www.w3.org/2000/svg"
                 class="h-6 w-6 text-gray-300 hover:text-gray-500 cursor-pointer inline"
                 @click.stop.prevent="deleteCustomer(customer)" fill="none" viewBox="0 0 24 24"
                 stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
            </svg>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
    <router-view name="modal" @statusUpdate="updateStatus"></router-view>
  </div>
</template>

<script>

export default {
  data () {
    return {
      fetchError: '',
      customers: [],
      loading: true
    }
  },
  created () {
    // watch the params of the route to fetch the data again
    this.$watch(
      () => this.$route.params,
      () => {
        this.fetchData()
      },
      // fetch the data when the view is created and the data is
      // already being observed
      { immediate: true }
    )
  },
  emits: ['statusUpdate'],
  methods: {
    updateStatus (event) {
      this.$emit('statusUpdate', event)
    },
    openCustomer (customer) {
      console.log('open customer', customer.name)
    },
    deleteCustomer (customer) {
      if (confirm(`Are you sure you want to delete '${customer.name}'?`)) {
        this.$api.delete(`web/private/customers/${customer.id}`)
          .then(() => this.fetchData())
          .catch(response => {
            console.error('failure', response)
            if (response.status === 403) {
              this.fetchError = 'Invalid credentials'
              return
            }
            this.fetchError = response
          })
          .finally(() => this.loading = false)
      }
    },
    fetchData () {
      this.$api.get('web/private/customers')
        .then(data => this.customers = data)
        .catch(response => {
          console.error('failure', response)
          if (response.status === 403) {
            this.fetchError = 'Invalid credentials'
            return
          }
          this.fetchError = response
        })
        .finally(() => this.loading = false)
    }
  }
}
</script>
