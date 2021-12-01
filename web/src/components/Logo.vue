<template>
  <img v-if="show === true" src="branding/logo" class="logo">
</template>
<style>
img.logo {
  display: block;
  max-width: 150px;
  max-height: 100px;
}
</style>
<script>
export default {
  data () {
    return {
      show: null
    }
  },
  methods: {
    load () {
      fetch('branding/logo').then((response) => {
        this.show = response.ok
      })
    }
  },
  mounted () {
    this.load()
  }
}
</script>
