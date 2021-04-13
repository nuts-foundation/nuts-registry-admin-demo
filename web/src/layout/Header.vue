<template>
  <div class="header-wrapper p-4 bg-white border-b border-gray-light">
    <div class="header container flex items centered">
      <div class="menu ml-auto font-light space-x-3">
        <router-view name="menu"></router-view>
      </div>
    </div>
  </div>
</template>