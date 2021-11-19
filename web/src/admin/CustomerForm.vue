<template>
  <form class="space-y-3">
    <div v-if="mode === 'edit'">
      <label for="customerDIDInput">Nuts Network DID</label>
      <input type="text" disabled v-model="value.did" id="customerDIDInput">
    </div>
    <div>
      <label for="customerIDInput">Internal ID:</label>
      <input type="number" :disabled="mode==='edit'" v-model.number="value.id" id="customerIDInput" min="1">
    </div>
    <div>
      <label for="customerNameInput">Name</label>
      <input type="text" v-model="value.name" id="customerNameInput">
    </div>
    <div>
      <label for="customerCityInput">City</label>
      <input type="text" v-model="value.city" id="customerCityInput">
    </div>
    <div>
      <label for="customerDomainInput">Domain</label>
      <input type="text" v-model="value.domain" id="customerDomainInput">
    </div>
    <div v-if="mode==='edit'">
      <label class="flex justify-start items-start" for="customerActiveInput">
        <div
            class="bg-white border rounded-md border-gray-300 w-6 h-6 flex flex-shrink-0 justify-center items-center mr-2 focus-within:border-blue-500">
          <input class="opacity-0 absolute" id="customerActiveInput" type="checkbox" v-model="value.active">
          <svg class="fill-current hidden w-4 h-4 text-green-500 pointer-events-none" viewBox="0 0 20 20">
            <path d="M0 11l2-2 5 5L18 3l2 2L7 18z"/>
          </svg>
        </div>
        <div>Publish by this name in Nuts network</div>
      </label>
    </div>
  </form>
</template>
<script>
export default {
  props: {
    value: Object,
    mode: String
  },
  emits: ['input'],
  watch: {
    value () {
      this.$emit('input', this.value)
    }
  }
}
</script>
