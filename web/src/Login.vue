<template>
  <div class="flex justify-center">

    <div class="mt-12 bg-white shadow-sm border rounded-md w-96 p-8 flex flex-col">
      <logo class="mx-auto my-2"></logo>
      <h1>Nuts Registry Admin</h1>

      <form class="w-full mt-4" @submit.stop.prevent="login">
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
          <p v-if="!!loginError" class="p-2 text-center bg-red-100 rounded-md">{{ loginError }}</p>
          <button id="login_button" class="w-full btn btn-primary">Login</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import Logo from './components/Logo.vue'

export default {
  components: { Logo },
  data () {
    return {
      loginError: '',
      credentials: {
        username: 'demo@nuts.nl',
        password: ''
      }
    }
  },
  watch: {
    // Remove error when typing
    'credentials.username' () {
      this.loginError = ''
    },
    'credentials.password' () {
      this.loginError = ''
    }
  },
  methods: {
    redirectAfterLogin () {
      this.$router.push('/admin/')
    },
    login () {
      this.$api.post('web/auth', this.credentials)
        .then(responseData => {
          localStorage.setItem('session', responseData.token)
          this.redirectAfterLogin()
        })
        .catch(response => {
          console.error('failure', response)
          if (response === 'invalid credentials') {
            this.loginError = 'Invalid credentials'
          } else {
            this.loginError = response.statusText
          }
        })
    }
  },
  mounted () {
    // Check if session still valid, if so just redirect to application
    this.$api.get('web/private')
      .then(() => this.redirectAfterLogin())
      .catch(() => {
        // session is invalid, need to authenticate
      })
  }
}
</script>
