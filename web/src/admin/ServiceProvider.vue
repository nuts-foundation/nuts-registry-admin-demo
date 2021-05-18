<template>
  <h1 class="page-title">Service Provider Configuration</h1>
  <p>A Service Provider offers (Nuts) services to its customers.</p>
  <p>Here you can create or update the contact information of your Service Provider.</p>
  <form class="vertical-form">

    <div class="space-y-4">
      <div v-if="!!serviceProvider.id">
        ID: {{ serviceProvider.id }}
      </div>

      <div>
        <label for="name-input">Name of the Service Provider</label>
        <input id="name-input" v-model="serviceProvider.name" type="text">
      </div>

      <div>
        <label for="email-input">Support email address</label>
        <input id="email-input" v-model="serviceProvider.email" type="email" required>
      </div>


      <div>
        <label for="phone-input">Emergency phone number</label>
        <input id="phone-input" v-model="serviceProvider.phone" type="text">
      </div>

      <div>
        <label for="website-input">Service Provider website</label>
        <input id="website-input" v-model="serviceProvider.website" type="text">
      </div>

      <button v-if="!serviceProvider.id" class="btn-submit" @click="createServiceProvider">Create Service Provider
      </button>
      <button v-if="!!serviceProvider.id" class="btn-submit" @click="updateServiceProvider">Update Service Provider
      </button>
      <div v-if="!!feedbackMsg"
           :class="{ 'bg-green-300': responseState === 'success', 'bg-red-300': responseState === 'error'}"
           class="py-2 px-4 border rounded-md text-white">
        <svg v-if="responseState === 'success'" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 inline" fill="none"
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
    </div>
  </form>

  <div class="flex justify-between mb-6">
    <h2 class="page-subtitle">Endpoints</h2>
    <button
        class="bg-blue-400 hover:bg-blue-500 text-white font-medium rounded-md px-3 py-2"
        @click="$router.push({name: 'admin.newEndpoint'})">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 inline" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
      </svg>
      Add
    </button>
  </div>

  <table class="min-w-full divide-y divide-gray-200">
    <thead class="bg-gray-50">
    <tr>
      <th class="thead">Type</th>
      <th class="thead">URL</th>
    </tr>
    </thead>
    <tbody class="tbody">
    <tr class="hover:bg-gray-100" v-for="endpoint in serviceProvider.endpoints">
      <td class="tcell">{{ endpoint.type }}</td>
      <td class="tcell">{{ endpoint.url }}</td>
    </tr>
    </tbody>
  </table>
  <router-view name="modal" @statusUpdate="updateStatus"></router-view>
</template>

<script>
export default {
  data() {
    return {
      responseState: '',
      feedbackMsg: '',
      serviceProvider: {
        id: '',
        name: '',
        email: '',
        phone: '',
        endpoints: [],
      }
    }
  },
  emits: ["statusUpdate"],
  created() {
    this.fetchData()
  },
  methods: {
    updateStatus(event) {
      this.$emit("statusUpdate", event)
    },
    updateServiceProvider() {
      this.$api.post("web/private/service-provider", this.serviceProvider)
          .then(responseData => {
            this.responseState = 'success'
            this.$emit("statusUpdate", "Service Provider Saved")
            this.serviceProvider = responseData
            this.feedbackMsg = ''
          })
          .catch(reason => {
            console.error("failure", reason)
            this.responseState = 'error'
            this.feedbackMsg = reason
          })
    },
    createServiceProvider() {
      this.$api.put("web/private/service-provider", this.serviceProvider)
          .then(responseData => {
            this.responseState = 'success'
            this.feedbackMsg = ''
            this.$emit("statusUpdate", "Service Provider Saved")
            this.serviceProvider = responseData
          })
          .catch(reason => {
            console.error("failure", reason)
            this.responseState = 'error'
            this.feedbackMsg = reason
          })
    },
    fetchData() {
      this.$api.get("web/private/service-provider")
          .then(responseData => {
            this.responseState = 'success'
            this.feedbackMsg = ''
            this.serviceProvider = responseData
          })
          .catch(reason => {
            console.error("failure", reason)
            this.responseState = 'error'
            this.feedbackMsg = reason
          })
    }
  }
}
</script>