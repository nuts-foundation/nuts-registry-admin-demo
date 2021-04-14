<template>
  <form class="my-4 flex justify-center">
    <div class="space-y-4">

      <div>
        <label for="username_input" class="block text-sm font-medium text-gray-700">Username</label>
        <input id="username_input"
               v-model="username"
               type="text"
               placeholder="Username"
               class="flex-1 py-2 px-4 block border border-gray-300 rounded-md"
        />
      </div>

      <div>
        <label for="password_input" class="block text-sm font-medium text-gray-700">Password</label>
        <input
            id="password_input"
            v-model="password"
            type="password"
            placeholder="Password"
            class="flex-1 py-2 px-4 block border border-gray-300 rounded-md"
        />
      </div>
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
      username: 'demo@nuts.nl',
      password: ''
    }
  },
  methods: {
    onSubmit() {
      localStorage.setItem("session", "true")
      this.$router.push("/admin/")
    }
  }
}
</script>