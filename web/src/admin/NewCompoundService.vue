<template>
  <modal-window :cancelRoute="{name: 'admin.serviceProvider'}" :confirmFn="checkForm" confirmText="Register"
                title="Register new Service" type="add">

    <p class="mb-3 text-sm">
      Here you can compose a new Service from endpoints. A service can be than be enabled per customer.
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
      <div v-for="(endpointRef, type) in service.endpoints">
        {{ type }}:{{ endpointRef }}
      </div>
      <input type="text" placeholder="type" v-model="newEndpointType">:
      <select name="" id="" class="form-select" v-model="selectedEndpoint">
        <option defaul selected value="">Select an endpoint</option>
        <option v-for="(endpoint, endpointId) in availableEndpoints" :value="endpointId">{{ endpoint.type }}</option>
      </select>
      <button class="btn-primary" @click="addEndpoint(newEndpointType, selectedEndpoint)">Add endpoint</button>
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
        endpoints: {}
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
      let endpoint = this.allEndpoints[endpointID]
      if (!endpoint) {
        return
      }
      let did = endpointID.split("#")[0]

      this.service.endpoints[type] = `${did}?type=${endpoint.type}`
      this.selectedEndpoint = ""
      this.newEndpointType = ""
      delete this.availableEndpoints[endpointID]
    },
    checkForm(e) {
      // reset the errors
      this.formErrors.length = 0
      this.apiError = ''

      if (this.service.name && Object.keys(this.service.endpoints).length > 0) {
        return this.confirm()
      }

      if (!this.service.name) {
        this.formErrors.push("Name required")
      }

      if (!Object.keys(this.service.endpoints).length > 0) {
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