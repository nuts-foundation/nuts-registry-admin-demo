<template>
  <h1 class="page-header">Service Provider Configuration</h1>
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
        <input id="email-input" v-model="serviceProvider.email" type="email">
      </div>


      <div>
        <label for="phone-input">Emergency phone number</label>
        <input id="phone-input" v-model="serviceProvider.phone" type="text">
      </div>

      <button v-if="!serviceProvider.id" class="btn-submit" @click="createServiceProvider">Create Service Provider
      </button>
      <button v-if="!!serviceProvider.id" class="btn-submit" @click="updateServiceProvider">Update Service Provider
      </button>
      <div v-if="!!feedbackMsg" :class="{ 'bg-green-300': responseState === 'success', 'bg-red-300': responseState === 'error'}" class="py-2 px-4 border rounded-md text-white">
        <svg v-if="responseState === 'success'" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        <svg v-if="responseState === 'error'" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        {{ feedbackMsg }}
      </div>
    </div>
  </form>

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
        phone: ''
      }
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    updateServiceProvider() {
      fetch("web/service-provider", {
        method: "PUT",
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(this.serviceProvider)
      }).then(response => {
        if (!response.ok) {
          if (response.status == 403) {
            throw "Invalid credentials"
          }
          throw response.statusText
        }
        return response.json()
      }).then(responseData => {
        console.log("success!")
        this.responseState = 'success'
        this.feedbackMsg = "Service Provider Updated"
        this.serviceProvider = responseData
      }).catch(reason => {
        console.error("failure", reason)
        this.responseState = 'error'
        this.feedbackMsg = reason
      })
    },
    createServiceProvider() {
      fetch("web/service-provider", {
        method: "POST",
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(this.serviceProvider)
      }).then(response => {
        if (!response.ok) {
          if (response.status == 403) {
            throw "Invalid credentials"
          }
          throw response.statusText
        }
        return response.json()
      }).then(responseData => {
        this.responseState = 'success'
        this.feedbackMsg = "Service Provider Created"
        console.log("success!")
        this.serviceProvider = responseData
      }).catch(reason => {
        console.error("failure", reason)
        this.responseState = 'error'
        this.feedbackMsg = reason
      })
    },
    fetchData() {
      fetch("web/service-provider")
          .then(response => {
            if (!response.ok) {
              if (response.status == 403) {
                throw "Invalid credentials"
              }
              throw response.statusText
            }
            return response.json()
          }).then(responseData => {
        console.log("success!")
        this.serviceProvider = responseData
      }).catch(reason => {
        console.error("failure", reason)
        this.feedbackMsg = reason
      })
    }
  }
}
</script>