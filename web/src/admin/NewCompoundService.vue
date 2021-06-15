<template>
  <modal-window :cancelRoute="{name: 'admin.serviceProvider'}" :confirmFn="checkForm" confirmText="Register"
                title="Register new Service" type="add">

    <p class="mb-3 text-sm">
      Here you can compose a new Service from endpoints. A service can than be enabled per customer.
    </p>

    <p v-if="apiError" class="p-3 bg-red-100 rounded-md">Could not register service: {{ apiError }}</p>

    <div class="p-3 bg-red-100 rounded-md" v-if="formErrors.length">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="error in formErrors">* {{ error }}</li>
      </ul>
    </div>

    <form class="space-y-3" novalidate>
      <div>
        <label for="endpointTypeInput">Name:</label>
        <input type="text" v-model="service.name" id="endpointTypeInput" required>
      </div>
      <label>Endpoints:
        <table>
          <tr>
            <th>Type</th>
            <th colspan="2">Endpoint name</th>
          </tr>
          <tr v-for="(endpointRef, name) in service.serviceEndpoint">
            <td>{{ name }}</td>
            <td>{{ endpointRef.split('=')[1] }}</td>
            <td>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-gray-300 hover:text-gray-500"
                   @click="deleteEndpoint(name, endpointRef)" fill="none" viewBox="0 0 24 24"
                   stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </td>
          </tr>
          <tr v-if="Object.keys(availableEndpoints).length" >
            <td>
              <input type="text" placeholder="type" v-model="newEndpointType">
            </td>
            <td>
              <select :disabled="Object.keys(availableEndpoints).length === 0" name="" id="" class="form-select"
                      v-model="selectedEndpoint">
                <option selected value="">Select an endpoint</option>
                <option v-for="(endpoint, endpointId) in availableEndpoints" :value="endpointId">{{
                    endpoint.type
                  }}
                </option>
              </select>
            </td>
            <td>
              <button class="btn-primary" @click="addEndpoint(newEndpointType, selectedEndpoint)">Add</button>
            </td>
          </tr>
        </table>
      </label>
    </form>

  </modal-window>
</template>

<script>
import ModalWindow from "../components/ModalWindow.vue";

export default {
  components: {ModalWindow},
  data() {
    return {
      apiError: '',
      formErrors: [],
      service: {
        name: '',
        serviceEndpoint: {}
      },
      newEndpointType: "",
      allEndpoints: {},
      availableEndpoints: {},
      selectedEndpoint: ""
    }
  },
  emits: ["statusUpdate"],
  watch: {
    "$route.params": {
      handler(toParams, previousParams) {
        // Fetch data when the route change (e.g. from the modal back to the list)
        this.fetchData()
      },
      immediate: true
    }
  },
  methods: {
    fetchData() {
      this.$api.get("web/private/service-provider/endpoints")
          .then(responseData => {
            responseData.forEach((el) => {
              this.allEndpoints[el.id] = el
              this.availableEndpoints[el.id] = el
            })
          })
          .catch(reason => {
            console.log("error while fetching endpoints: ", reason)
          })
    },
    addEndpoint(type, endpointID) {
      if (!type) {
        return
      }
      let endpoint = this.allEndpoints[endpointID]
      if (!endpoint) {
        return
      }
      let did = endpointID.split("#")[0]

      this.service.serviceEndpoint[type] = `${did}?type=${endpoint.type}`
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
      // reset the errors
      this.formErrors.length = 0
      this.apiError = ''

      if (this.service.name && Object.keys(this.service.serviceEndpoint).length > 0) {
        return this.confirm()
      }

      if (!this.service.name) {
        this.formErrors.push("Name required")
      }

      if (!Object.keys(this.service.serviceEndpoint).length > 0) {
        this.formErrors.push("At least one endpoint required")
      }

      e.preventDefault()
    },
    confirm() {
      this.$api.post('web/private/service-provider/services', this.service)
          .then(() => {
            this.$emit("statusUpdate", "Endpoint registered")
            this.$router.push({name: 'admin.serviceProvider'})
          })
          .catch(response => {
            this.apiError = response
          })
    }
  }
}
</script>