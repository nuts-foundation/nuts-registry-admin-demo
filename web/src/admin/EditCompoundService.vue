<template>
  <compound-service-form v-if="service"
                         :confirm-fn="updateService"
                         confirm-text="Edit"
                         description="Here you can update an existing service."
                         mode="edit"
                         title="Update service"
                         :existingService="service"
                         :error="apiError"
                         :endpoints="allEndpoints">
  </compound-service-form>
</template>

<script>
import CompoundServiceForm from './CompoundServiceForm.vue'

export default {
  components: { CompoundServiceForm },
  data () {
    return {
      apiError: '',
      allEndpoints: {},
      service: null,
      serviceID: null
    }
  },
  emits: ['statusUpdate'],
  watch: {
    '$route.params': {
      handler (toParams, previousParams) {
        this.serviceID = toParams.serviceID
        // Fetch data when the route change (e.g. from the modal back to the list)
        this.fetchData()
      },
      immediate: true
    }
  },
  methods: {
    fetchData () {
      this.$api.get('web/private/service-provider/services')
        .then(services => {
          this.service = services.filter(svc => svc.id === this.serviceID)[0]
          return this.$api.get('web/private/service-provider/endpoints')
        })
        .then(endpoints =>
          endpoints.forEach((el) => {
            this.allEndpoints[el.id] = el
          })
        )
        .catch(reason => this.apiError = reason)
    },
    updateService (service) {
      // To the reader: since the delete-then-add below is not transactional, dataloss might occur when delete succeeds but register fails.
      // But since this is a demo application, it's probably good enough (just make sure to do it properly in your production implementation).
      this.$api.delete(`web/private/service-provider/endpoints/${escape(this.serviceID)}`)
        .then(() => this.$api.post('web/private/service-provider/services', service))
        .then(() => {
          this.$emit('statusUpdate', 'Service updated')
          this.$router.push({ name: 'admin.serviceProvider' })
        })
        .catch(response => this.apiError = response)
    }
  }
}
</script>
