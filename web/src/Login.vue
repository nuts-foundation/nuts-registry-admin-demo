<template>
  <form ref="form" :model="form" label-width="120px">
    <input placeholder="Username" />
    <input placeholder="Password" />
    <button>Cancel</button>
  </form>
</template>

<script>
export default {
  data() {
    return {
      form: {
        username: '',
        password: ''

      }
    }
  },
  methods: {
    onSubmit() {
      this.$router.push("/admin/")
    }
  }
}
</script>