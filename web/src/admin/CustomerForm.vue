<template>
  <form class="space-y-3">
    <div v-if="mode === 'edit'">
      <label for="customer-did-input">Nuts Network DID</label>
      <input type="text" disabled v-model="localValue.did" id="customer-did-input">
    </div>
    <div>
      <label for="customer-id-input">Internal ID:</label>
      <input type="number" :disabled="mode==='edit'" v-model.number="localValue.id" id="customer-id-input" min="1">
    </div>
    <div>
      <label for="customer-name-input">Name</label>
      <input type="text" v-model="localValue.name" id="customer-name-input">
    </div>
    <div>
      <label for="customer-city-input">City</label>
      <input type="text" v-model="localValue.city" id="customer-city-input">
    </div>
    <div>
      <label for="customer-domain-input">Domain</label>
      <input type="text" v-model="localValue.domain" id="customer-domain-input">
    </div>
    <div v-if="mode==='edit'">
      <label class="flex justify-start items-start" for="customer-publish-checkbox">
        <div
            class="bg-white border rounded-md border-gray-300 w-6 h-6 flex flex-shrink-0 justify-center items-center mr-2 focus-within:border-blue-500">
          <input class="opacity-0 absolute" id="customer-publish-checkbox" type="checkbox" v-model="localValue.active">
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
  data () {
    return {
      localValue: this.value
    }
  },
  emits: ['input'],
  watch: {
    localValue () {
      this.$emit('input', this.localValue)
    }
  }
}
</script>
