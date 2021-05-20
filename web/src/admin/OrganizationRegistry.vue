<template>
  <h1 class="text-3xl">Care Organisation Registry</h1>
  <p v-if="fetchError" class="m-4">Could not fetch care organizations: {{ fetchError }}</p>
  <p>Search the Nuts Network for care organizations:</p>

  <form class="space-x-3 flex">
    <div>
      <label for="nameInput">Name:</label>
      <input type="text" v-model="query.name" id="nameInput" v-on:input="search" v-on:focusout="search">
    </div>
    <div>
      <label for="cityInput">City:</label>
      <input type="text" v-model="query.city" id="cityInput" v-on:input="search" v-on:focusout="search">
    </div>
  </form>

  <h1 class="text-xl pt-3">Search Results</h1>

  <div class="space-y-3 shadow overflow-hidden border-gray-200 rounded">
    <table class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
      <tr>
        <th>Name</th>
        <th>City</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="result in results">
        <td class="tcell">{{ result.organization.name }}</td>
        <td class="tcell">{{ result.organization.city }}</td>
      </tr>
      </tbody>
      <tfoot>
        <tr>
          <td colspan="2" class="px-6 py-2 bg-gray-50" >Found {{results.length}} result{{results.length != 1 ? 's':''}}</td>
        </tr>

      </tfoot>
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
      this.$api.post('web/private/organizations', this.query)
          .then(data => this.results = data)
          .catch(response => {
            this.fetchError = response.statusText
            this.results = []
          })
    }
  }
}
</script>

<style>

.thead {
  @apply px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider;
}

.tcell {
  @apply px-6 py-4 text-left;
}

.body {
  @apply bg-white divide-y divide-gray-200;
}
</style>