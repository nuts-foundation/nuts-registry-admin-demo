<template>
  <h1 class="page-title">Login</h1>
  <form class="my-4 flex justify-center">
    <div class="space-y-4">

      <div>
        <label for="username_input" class="block text-sm font-medium text-gray-700">Username</label>
        <input id="username_input"
               v-model="credentials.username"
               type="text"
               placeholder="Username"
               class="flex-1 py-2 px-4 block border border-gray-300 rounded-md"
        />
      </div>

      <div>
        <label for="password_input" class="block text-sm font-medium text-gray-700">Password</label>
        <input
            id="password_input"
            v-model="credentials.password"
            type="password"
            placeholder="Password"
            class="flex-1 py-2 px-4 block border border-gray-300 rounded-md"
        />
      </div>
      <p v-if="!!loginError" class="p-3 bg-red-100 rounded-md">{{ loginError }}</p>
      <button
          @click="onSubmit"
          class="w-full btn-submit"
      >Login
      </button>
    </div>
  </form>
</template>

<script>
export default {
  data() {
    return {
      loginError: "",
      credentials: {
        username: 'demo@nuts.nl',
        password: ''
      }
    }
  },
  watch: {
    // Remove error when typing
    'credentials.username'() {
      this.loginError = ""
    },
    'credentials.password'() {
      this.loginError = ""
    }
  },
  methods: {
    onSubmit() {
      this.$api.post('web/auth', this.credentials)
          .then(responseData => {
            console.log("success!")
            localStorage.setItem("session", responseData.token)
            this.$router.push("/admin/")
          })
          .catch(response => {
            console.error("failure", response)
            if (response === "invalid credentials") {
              this.loginError = "Invalid credentials"
            } else {
              this.loginError = response.statusText
            }
          })
    }
  }
}
</script>