<template>
  <router-view/>
</template>

<script>

import Header from "./layout/Header.vue"
import Footer from "./layout/Footer.vue"

export default {
  components: {
    "nrad-header": Header,
    "nrad-footer": Footer
  }
}
</script>

<style>
</style>