<template>
  <compound-service-form :confirm-fn="registerService"
                         :endpoints="allEndpoints"
                         :error="apiError"
                         confirm-text="Register"
                         description="Here you can compose a new Service from endpoints. A service can than be enabled per care organization."
                         mode="new"
                         title="Register new service">
  </compound-service-form>
</template>

<script>
import CompoundServiceForm from './CompoundServiceForm.vue'

export default {
  components: { CompoundServiceForm },
  data () {
    return {
      apiError: '',
      allEndpoints: {}
    }
  },
  emits: ['statusUpdate'],
  watch: {
    '$route.params': {
      handler (toParams, previousParams) {
        // Fetch data when the route change (e.g. from the modal back to the list)
        this.fetchData()
      },
      immediate: true
    }
  },
  methods: {
    fetchData () {
      this.$api.get('web/private/service-provider/endpoints')
        .then(responseData => {
          responseData.forEach((el) => {
            this.allEndpoints[el.id] = el
          })
        })
        .catch(reason => {
          this.apiError = reason
        })
    },
    registerService (service) {
      return this.$api.post('web/private/service-provider/services', service)
        .then(() => {
          this.$emit('statusUpdate', 'Service registered')
          this.$router.push({ name: 'admin.serviceProvider' })
        })
        .catch(response => {
          this.apiError = response
        })
    }
  }
}
</script>
