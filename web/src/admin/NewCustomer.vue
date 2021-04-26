<template>
  <nrad-modal :cancelRoute="{name: 'admin.customers'}" :confirmFn="confirm">
    <form>
      <input type="text">
      <input type="text">
    </form>
  </nrad-modal>
</template>

<script>
export default {
  methods: {
    confirm() {
      this.$router.push({name: 'admin.customers'})
    }
  }
}
</script>