<template>
  <div>

    <h1 class="mb-4">Service Provider Configuration</h1>

    <p>A Service Provider offers (Nuts) services to its care organizations.</p>
    <p>Here you can update the contact information of your Service Provider.</p>

    <form @submit.stop.prevent="updateServiceProvider">
      <div class="mt-8 bg-white p-5 shadow-lg rounded-lg">
        <div class="space-y-4 w-full">
          <div v-if="feedbackMsg"
               :class="{ 'bg-green-300': responseState === 'success', 'bg-red-300': responseState === 'error'}"
               class="py-2 px-4 border rounded-md text-white">
            <svg v-if="responseState === 'success'" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 inline"
                 fill="none"
                 viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>
            <svg v-if="responseState === 'error'" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 inline" fill="none"
                 viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
            </svg>
            {{ feedbackMsg }}
          </div>

          <div v-if="serviceProvider.id">
            <label for="did-input">DID</label>
            <input id="did-input" type="text" disabled v-model="serviceProvider.id">
          </div>

          <div>
            <label for="name-input">Name of the Service Provider</label>
            <input id="name-input" v-model="serviceProvider.name" type="text">
          </div>

          <div>
            <label for="email-input">Support email address (required)</label>
            <input id="email-input" v-model="serviceProvider.email" type="email" required>
          </div>

          <div>
            <label for="website-input">Service Provider website</label>
            <input id="website-input" v-model="serviceProvider.website" type="text">
          </div>

          <div>
            <label for="endpoint-input">Nuts node endpoint of the Service Provider</label>
            <input id="endpoint-input" v-model="serviceProvider.endpoint" type="text" placeholder="grpc://nuts.nl:5555">
            <div class="text-sm">
              Public address of the Nuts Node endpoint which other nodes connect to, e.g. <pre class="inline">grpc://nuts.nl:5555</pre>.
              See <a href="https://nuts-node.readthedocs.io/en/latest/pages/getting-started/3-configure-your-node.html#configure-node" target="_blank">the documentation</a> for more information.
            </div>
          </div>
        </div>
      </div>

      <div class="mt-4">
        <button class="btn btn-primary">
          {{ !serviceProvider.id ? 'Create' : 'Update' }} Service Provider
        </button>
      </div>
    </form>

    <div class="flex justify-between mt-14" v-if="serviceProvider.id">
      <h2>Endpoints</h2>

      <p class="pl-6 w-9/12 text-left">An endpoint is a simple registration of a named URL. It can be used as a building
        block for Services.</p>

      <button
          @click="$router.push({name: 'admin.newEndpoint'})"
          class="float-right inline-flex items-center bg-nuts w-10 h-10 rounded-lg justify-center shadow-md"
      >
        <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#fff">
          <path d="M0 0h24v24H0V0z" fill="none"/>
          <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
        </svg>
      </button>
    </div>

    <div v-if="endpoints.length > 0" class="mt-4 bg-white p-5 shadow-lg rounded-lg">
      <table class="min-w-full divide-y divide-gray-200">
        <thead>
        <tr>
          <th class="thead">Type</th>
          <th class="thead">URL</th>
          <th class="thead">Delete</th>
        </tr>
        </thead>
        <tbody class="tbody">
        <tr class="hover:bg-gray-100"
            v-for="endpoint in endpoints"
            :key="endpoint.id"
        >
          <td class="tcell">{{ endpoint.type }}</td>
          <td class="tcell">{{ endpoint.url }}</td>
          <td class="tcell">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-300 hover:text-gray-500"
                 @click="deleteEndpoint(endpoint.id)" fill="none" viewBox="0 0 24 24"
                 stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
            </svg>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <div class="flex justify-between mt-14" v-if="serviceProvider.id">
      <h2>Services</h2>

      <p class="pl-8 w-9/12">A service is set of endpoints which can be offered to care organizations.</p>

      <button
          @click="$router.push({name: 'admin.newCompoundService'})"
          class="float-right inline-flex items-center bg-nuts w-10 h-10 rounded-lg justify-center shadow-md"
      >
        <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#fff">
          <path d="M0 0h24v24H0V0z" fill="none"/>
          <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
        </svg>
      </button>
    </div>

    <div v-if="services.length > 0" class="mt-4 bg-white p-5 shadow-lg rounded-lg">
      <table class="min-w-full divide-y divide-gray-200">
        <thead>
        <tr>
          <th class="thead">Service name</th>
          <th class="thead">Endpoints</th>
          <th class="thead">Actions</th>
        </tr>
        </thead>
        <tbody class="tbody">
        <tr class="hover:bg-gray-100"
            v-for="service in services"
            :key="service.id"
        >
          <td class="tcell">{{ service.name }}</td>
          <td class="tcell"><p v-for="(endpoint, name) in service.serviceEndpoint" :key="endpoint.id">{{ name }} &#8594;
            did:SP-DID?type={{ endpoint.split('=')[1] }}</p></td>
          <td class="tcell">
            <!-- Edit -->
            <svg xmlns="http://www.w3.org/2000/svg"
                 class="h-6 w-6 text-gray-300 hover:text-gray-500 cursor-pointer inline" fill="none" viewBox="0 0 24 24"
                 stroke="currentColor"
                 @click="$router.push({name: 'admin.editCompoundService', params: {serviceID: service.id}})">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
            </svg>
            <!-- Delete -->
            <svg xmlns="http://www.w3.org/2000/svg"
                 class="h-6 w-6 text-gray-300 hover:text-gray-500 cursor-pointer inline"
                 @click="deleteEndpoint(service.id)" fill="none" viewBox="0 0 24 24"
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
      responseState: '',
      feedbackMsg: '',
      serviceProvider: {
        id: '',
        name: '',
        email: '',
        phone: '',
        website: ''
      },
      endpoints: [],
      services: []
    }
  },
  emits: ['statusUpdate'],
  watch: {
    '$route.params': {
      handler (toParams, previousParams) {
        // Fetch data when the route change (e.g. from the modal back to the list)
        this.fetchServiceProvider()
      },
      immediate: true
    }
  },
  methods: {
    updateStatus (event) {
      this.$emit('statusUpdate', event)
    },
    updateServiceProvider () {
      this.feedbackMsg = ''

      this.$api.put('web/private/service-provider', this.serviceProvider)
        .then(responseData => {
          this.responseState = 'success'
          this.$emit('statusUpdate', 'Service Provider Saved')
          this.serviceProvider = responseData
          this.feedbackMsg = ''
        })
        .catch(reason => {
          console.error('failure', reason)
          this.responseState = 'error'
          this.feedbackMsg = reason
        })
    },
    fetchServiceProvider () {
      this.feedbackMsg = ''

      this.$api.get('web/private/service-provider')
        .then(responseData => {
          this.responseState = 'success'
          this.serviceProvider = responseData
          this.fetchData()
        })
        .catch(reason => {
          if (reason !== 'Not Found') {
            console.error('failure', reason)
            this.responseState = 'error'
            this.feedbackMsg = reason
          }
        })
    },
    fetchData () {
      this.feedbackMsg = ''

      this.$api.get('web/private/service-provider/endpoints')
        .then(responseData => {
          this.endpoints = responseData
        })
        .catch(reason => {
          console.log('error while fetching endpoints: ', reason)
        })

      this.$api.get('web/private/service-provider/services')
        .then(responseData => {
          this.services = responseData
        })
        .catch(reason => {
          console.error('failure', reason)
          this.responseState = 'error'
          this.feedbackMsg = reason
        })
    },
    deleteEndpoint (id) {
      if (confirm('Are you sure you want to delete this endpoint/service?')) {
        this.feedbackMsg = ''

        this.$api.delete(`web/private/service-provider/endpoints/${escape(id)}`, id)
          .then(response => {
            this.$emit('statusUpdate', 'Endpoint deleted')
          })
          .catch(reason => {
            console.error('failure', reason)
            this.responseState = 'error'
            this.feedbackMsg = reason
          })
          .finally(() => {
            this.fetchData()
          })
      }
    }
  }
}
</script>
<style>

.thead {
  @apply px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider;
}

.tcell {
  @apply px-6 py-4 text-left;
}
</style>
