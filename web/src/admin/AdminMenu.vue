<template>
  <router-link to="service-provider" class="page-link">My Config</router-link><el class="text-gray-400">|</el>
  <router-link to="customers" class="page-link">My Customers</router-link><el class="text-gray-400">|</el>
  <router-link to="/logout" class="page-link">Logout</router-link>
</template>

<style scoped>
.page-link {
  @apply no-underline hover:text-gray-900 text-gray-400;
}
</style>