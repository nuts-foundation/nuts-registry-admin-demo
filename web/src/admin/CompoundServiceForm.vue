<template>
  <modal-window :cancelRoute="{name: 'admin.serviceProvider'}" :confirmFn="checkForm" :confirmText="confirmText"
                :title="title" :type="mode">

    <p class="mb-3 text-sm">
      {{ description }}
    </p>

    <p v-if="error" class="p-3 bg-red-100 rounded-md">{{ error }}</p>

    <div class="p-3 bg-red-100 rounded-md" v-if="formErrors.length">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="error in formErrors">* {{ error }}</li>
      </ul>
    </div>

    <form class="space-y-3" @submit.prevent novalidate>
      <div>
        <label for="endpointTypeInput">Name</label>
        <input type="text" v-model="service.name" id="endpointTypeInput" required>
      </div>
      <label>Endpoints</label>
      <table>
        <tr>
          <th>Type</th>
          <th colspan="2">Endpoint name</th>
        </tr>
        <tr v-for="(endpointRef, name) in service.serviceEndpoint">
          <td>{{ name }}</td>
          <td>{{ endpointRef.split('=')[1] }}</td>
          <td>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-300 hover:text-gray-500 cursor-pointer"
                 @click="deleteEndpoint(name, endpointRef)" fill="none" viewBox="0 0 24 24"
                 stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
            </svg>
          </td>
        </tr>
        <tr v-if="Object.keys(availableEndpoints).length">
          <td>
            <input type="text" placeholder="type" v-model="newEndpointType">
          </td>
          <td>
            <select :disabled="Object.keys(availableEndpoints).length === 0" name="" id="" class="form-select"
                    v-model="selectedEndpoint">
              <option selected value="">Select an endpoint</option>
              <option v-for="(endpoint, endpointId) in availableEndpoints" :value="endpointId">
                {{ endpoint.type }}
              </option>
            </select>
          </td>
          <td>
            <button class="btn btn-primary btn-sm" @click="addEndpoint(newEndpointType, selectedEndpoint)">Add</button>
          </td>
        </tr>
      </table>
    </form>

  </modal-window>
</template>

<script>
import ModalWindow from "../components/ModalWindow.vue";

export default {
  components: {ModalWindow},
  props: {
    confirmText: String,
    confirmFn: Function,
    description: String,
    title: String,
    error: String,
    allEndpoints: Object,
    availableEndpoints: Object,
    existingService: Object,
    mode: String, // add | edit
  },
  data() {
    return {
      formErrors: [],
      selectedEndpoint: "",
      newEndpointType: "",
      service: this.existingService ?? {
        name: '',
        serviceEndpoint: {}
      },
    }
  },
  watch: {
    existingService: function(service) { this.service = service },
  },
  emits: ['input'],
  methods: {
    addEndpoint(type, endpointID) {
      if (!type) {
        return
      }
      let endpoint = this.allEndpoints[endpointID]
      if (!endpoint) {
        return
      }
      let did = endpointID.split("#")[0]

      this.service.serviceEndpoint[type] = `${did}/serviceEndpoint?type=${endpoint.type}`
      this.selectedEndpoint = ""
      this.newEndpointType = ""
      delete this.availableEndpoints[endpointID]
    },
    deleteEndpoint(name, ref) {
      let type = ref.split("=")[1]
      let endpoint = Object.values(this.allEndpoints).find((el) => el.type == type)
      if (!endpoint) {
        return
      }
      delete this.service.serviceEndpoint[name]
      this.availableEndpoints[endpoint.id] = endpoint
    },
    checkForm(e) {
      this.formErrors.length = 0

      if (!this.service.name) {
        this.formErrors.push("Name required")
      }
      if (!Object.keys(this.service.serviceEndpoint).length > 0) {
        this.formErrors.push("At least one endpoint required")
      }

      if (this.formErrors.length === 0) {
        return this.confirmFn(this.service)
      } else {
        e.preventDefault()
      }
    },
  }
}
</script>
