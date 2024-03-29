<template>
  <modal-window
      :cancelRoute="{name: 'admin.customers'}"
      :confirmFn="checkForm"
      confirmText="Save Customer"
      title="Edit Customer" type="edit"
  >
    <div v-if="loading && !apiError">Loading...</div>
    <div class="p-2 rounded-md w-full bg-red-100" v-if="apiError">
      Error during server communication: {{ apiError }}
    </div>
    <customer-form mode="edit"
                   :value="customer"
                   @input="(newCustomer)=> {customer = newCustomer}"
                   v-if="!loading"/>
    <div class="pt-3 space-y-1">
      <p>Service configuration:</p>
      <p class="text-sm" v-if="!this.availableServices.length">No services provided by the Service Provider.</p>
      <label class="flex justify-start items-start" v-for="service in availableServices" :key="service.id">
        <div
            class="bg-white border rounded-md border-gray-300 w-6 h-6 flex flex-shrink-0 justify-center items-center mr-2 focus-within:border-blue-500">
          <input :data-service-name="service.name"
                 class="opacity-0 absolute" type="checkbox"
                 :checked="enabledServices.map(v=>v.type).includes(service.name)"
                 @input="toggleService(service.id, $event)">
          <svg class="fill-current hidden w-4 h-4 text-green-500 pointer-events-none" viewBox="0 0 20 20">
            <path d="M0 11l2-2 5 5L18 3l2 2L7 18z"/>
          </svg>
        </div>
        <div>{{ service.name }}</div>
      </label>
    </div>
  </modal-window>
</template>
<style>
input:checked + svg {
  display: block;
}
</style>

<script>
import ModalWindow from '../components/ModalWindow.vue'
import CustomerForm from './CustomerForm.vue'

export default {
  components: { ModalWindow, CustomerForm },
  data () {
    return {
      customer: {
        name: '',
        city: '',
        domain: '',
        active: false
      },
      availableServices: [],
      enabledServices: [],
      formErrors: [],
      apiError: '',
      loading: true,
      servicesToAdd: [],
      servicesToRemove: []
    }
  },
  emits: ['statusUpdate'],
  watch: {
    '$route.params': {
      handler (toParams, previousParams) {
        if (toParams && 'id' in toParams) {
          this.fetchCustomer(toParams.id)
          this.fetchServices()
        }
      },
      immediate: true
    }
  },
  methods: {
    checkForm () {
      // reset the errors
      this.formErrors.length = 0
      this.apiError = ''

      if (this.customer.name) {
        return this.saveCustomer()
      }

      if (this.active && !this.customer.city) {
        this.formErrors.push('To be active a city is required')
      }

      if (!this.customer.name) {
        this.formErrors.push('Name required')
      }
    },
    toggleService (id, event) {
      const newState = event.target.checked
      const service = this.availableServices.find(v => v.id === id)
      if (!service) {
        return
      }
      if (newState) {
        const idx = this.servicesToRemove.map(v => v.id).indexOf(service.id)
        if (idx > -1) {
          this.servicesToRemove.splice(idx, 1)
        } else {
          this.servicesToAdd.push(service)
        }
      } else {
        const idx = this.servicesToAdd.map(v => v.id).indexOf(service.id)
        if (idx > -1) {
          this.servicesToAdd.splice(idx, 1)
        } else {
          this.servicesToRemove.push(service)
        }
      }
    },
    fetchCustomer (id) {
      this.$api.get(`web/private/customers/${id}`)
        .then((customer) => {
          this.customer = {
            ...customer
          }
          this.loading = false
        })
        .catch((reason) => {
          this.apiError = reason.statusText
          console.log('failed:', reason)
        })
      this.$api.get(`web/private/customers/${id}/services`)
        .then((services) => {
          this.enabledServices = services
        })
        .catch((reason) => {
          this.apiError = reason.statusText
          console.log('failed:', reason)
        })
    },
    fetchServices () {
      this.$api.get('web/private/service-provider/services')
        .then(responseData => {
          this.availableServices = responseData
        })
        .catch(reason => {
          this.apiError = reason.statusText
          console.log('error while fetching services: ', reason)
        })
    },
    saveCustomer () {
      this.$api.put(`web/private/customers/${this.customer.id}`, this.customer)
        .then((customer) => {
          this.customer = customer
          this.$emit('statusUpdate', 'Customer saved')
          this.saveServices().then(() => {
            this.$router.push({ name: 'admin.customers' })
          })
        })
        .catch((reason) => {
          this.apiError = reason
          console.log('failed:', reason)
        })
    },
    saveServices () {
      const removePromises = this.servicesToRemove.map(v => {
        return this.$api.delete(`web/private/customers/${this.customer.id}/services/${v.name}`)
      })

      const addPromises = this.servicesToAdd.map(v => {
        return this.$api.post(`web/private/customers/${this.customer.id}/services`, { did: v.id, type: v.name })
      })

      return Promise.all(addPromises.concat(removePromises))
    }
  }
}
</script>
