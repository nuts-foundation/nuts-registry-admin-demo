<template>
  <modal-window :cancelRoute="{name: 'admin.serviceProvider'}" :confirmFn="checkForm" confirmText="Register"
              title="Register endpoint URL" type="add">

    <p class="mb-4 text-sm">
      Here you can register an endpoint URL of your XIS on your Service Provider's DID,
      which can be used when enabling services for your customers.
    </p>

    <p v-if="apiError" class="p-3 bg-red-100 rounded-md">Could not register endpoint: {{ apiError }}</p>

    <div class="p-3 bg-red-100 rounded-md" v-if="formErrors.length">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="(error, idx) in formErrors" :key="`err-${idx}`">* {{ error }}</li>
      </ul>
    </div>

    <form class="space-y-3">
      <div>
        <label for="endpointTypeInput">Type</label>
        <input type="text" v-model="endpoint.type" id="endpointTypeInput" required>
      </div>
      <div>
        <label for="endpointURLInput">URL</label>
        <input type="text" v-model="endpoint.url" id="endpointURLInput" required
               placeholder="https://example.com">
      </div>
    </form>

  </modal-window>
</template>

<script>
import ModalWindow from '../components/ModalWindow.vue'

export default {
  components: { ModalWindow },
  data () {
    return {
      apiError: '',
      formErrors: [],
      endpoint: {
        id: '',
        type: '',
        url: ''
      }
    }
  },
  emits: ['statusUpdate'],
  methods: {
    checkForm (e) {
      // reset the errors
      this.formErrors.length = 0
      this.apiError = ''

      if (this.endpoint.type && this.endpoint.url) {
        return this.confirm()
      }

      if (!this.endpoint.type) {
        this.formErrors.push('Type required')
      }

      if (!this.endpoint.url) {
        this.formErrors.push('URL required')
      }

      e.preventDefault()
    },
    confirm () {
      this.$api.post('web/private/service-provider/endpoints', this.endpoint)
        .then(() => {
          this.$emit('statusUpdate', 'Endpoint registered')
          this.$router.push({ name: 'admin.serviceProvider' })
        })
        .catch(response => {
          this.apiError = response
        })
    }
  }
}
</script>
