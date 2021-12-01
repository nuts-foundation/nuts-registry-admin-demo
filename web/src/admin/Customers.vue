<template>
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
    <div class="m-4" v-if="!loading && customers.length == 0 && !fetchError">No customers yet, add one!</div>
    <table v-if="customers.length > 0" class="min-w-full divide-y divide-gray-200">
      <thead>
      <tr>
        <th class="thead">Customer ID</th>
        <th class="thead">Name</th>
        <th class="thead">City</th>
        <th class="thead">Published</th>
      </tr>
      </thead>
      <tbody class="tbody">
      <tr class="hover:bg-gray-100 cursor-pointer"
          @click="$router.push({name: 'admin.editCustomer', params: {id: customer.id} })"
          v-for="{id, name, city, active} in customers" :key="id">
        <td class="tcell">{{ id }}</td>
        <td class="tcell">{{ name }}</td>
        <td class="tcell">{{ city }}</td>
        <td class="tcell">{{ active }}</td>
      </tr>
      </tbody>
    </table>
  </div>
  <router-view name="modal" @statusUpdate="updateStatus"></router-view>
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
