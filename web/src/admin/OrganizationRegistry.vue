<template>
  <h1 class="text-3xl">Care Organisation Registry</h1>
  <p v-if="fetchError" class="m-4">Could not fetch care organizations: {{ fetchError }}</p>
  <p>Search for registered care organizations with trusted credentials:</p>
  <form class="space-y-3">
    <div>
      <label for="nameInput">Name:</label>
      <input type="text" v-model="query.name" id="nameInput" v-on:input="search" v-on:focusout="search">
    </div>
    <div>
      <label for="townInput">City:</label>
      <input type="text" v-model="query.city" id="townInput" v-on:input="search" v-on:focusout="search">
    </div>
  </form>
  <div class="space-y-3">
    <h1 class="text-xl">Search Results</h1>
    <p v-if="results.length === 0">
      Nothing found (yet)
    </p>
    <table v-if="results.length > 0" class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
      <tr>
        <th>Name</th>
        <th>City</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="result in results">
        <td>{{ result.organization.name }}</td>
        <td>{{ result.organization.city }}</td>
      </tr>
      </tbody>
    </table>
  </div>
</template>

<script>

export default {
  data() {
    return {
      fetchError: "",
      results: [],
      query: {
        name: "",
        city: ""
      },
    }
  },
  emits: ["statusUpdate"],
  methods: {
    search() {
      if (this.query.name === "" && this.query.city === "") {
        this.results = []
        return
      }
      this.$api.post('web/private/credentials/organizations', this.query)
          .then(data => this.results = data)
          .catch(response => {
            this.fetchError = response.statusText
            this.results = []
          })
    }
  }
}
</script>