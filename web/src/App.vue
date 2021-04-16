<template>
  <div class="App min-w-screen min-h-screen flex flex-col bg-gray-100 font-sans">
    <nrad-header></nrad-header>

    <div class="content-wrapper pt-4 flex-1 bg-grey-lighter">
      <div class="content py-4 container max-w-6xl mx-auto bg-white m-4 p-6 rounded-md shadow-md">
        <router-view/>
      </div>
    </div>

    <nrad-footer></nrad-footer>
  </div>
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