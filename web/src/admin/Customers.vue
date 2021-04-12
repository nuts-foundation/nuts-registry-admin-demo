<template>
  <h1>customer</h1>
</template>

<script>
import { onMounted, reactive } from "vue";

export default {
  data() {
    return {
      customers: []
    }
  },
  created() {
    // watch the params of the route to fetch the data again
    this.$watch(
        () => this.$route.params,
        () => {
          this.fetchData()
        },
        // fetch the data when the view is created and the data is
        // already being observed
        { immediate: true }
    )
  },
  methods: {
    fetchData() {
      fetch("api/customers")
          .then(response => response.json())
          .then(data => this.customers = data)
    }
  },
}
</script>