<template>
  <section class="view-process" v-test-id="'viewProcess'">
    <process-options></process-options>
    <process-list class="process-list"></process-list>
  </section>
</template>

<script>
  import ProcessOptions from './options.vue'
  import ProcessList from './process-list'
  export default {
    name: 'view-process',
    components: {
      ProcessOptions,
      ProcessList
    },
    data() {
      return {}
    }
  }
</script>
