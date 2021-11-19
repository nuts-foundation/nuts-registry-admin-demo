<template>
  <label class="flex justify-start items-start">
    <div class="bg-white border rounded-md border-gray-300 w-6 h-6 flex flex-shrink-0 justify-center items-center mr-2 focus-within:border-blue-500">
      <input class="opacity-0 absolute" type="checkbox" :checked="modelValue" @change="$emit('update:modelValue', $event.target.checked)">
      <svg class="fill-current hidden w-4 h-4 text-green-500 pointer-events-none" viewBox="0 0 20 20">
        <path d="M0 11l2-2 5 5L18 3l2 2L7 18z"/>
      </svg>
    </div>
    <div>
      <slot></slot>
    </div>
  </label>
</template>
<script>
export default {
  name: 'formCheckbox',
  props: {
    modelValue: Boolean
  },
  emits: ['update:modelValue']
}
</script>
