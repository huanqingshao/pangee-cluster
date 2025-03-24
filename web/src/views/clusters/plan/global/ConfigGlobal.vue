<i18n>
</i18n>

<template>
  <div>
    <ContainerManager :cluster="cluster" @refresh="$emit('refresh')"></ContainerManager>
    <OsMirror :cluster="cluster"></OsMirror>
  </div>
</template>

<script>
import OsMirror from "./OsMirror.vue";
import ContainerManager from "./ContainerManager.vue";

export default {
  props: {
    cluster: { type: Object, required: true }
  },
  data() {
    return {};
  },
  computed: {
    inventory: {
      get() {
        return this.cluster.inventory;
      },
      set() {}
    }
  },
  components: { OsMirror, ContainerManager },
  mounted() {},
  methods: {}
};
</script>

<style scoped lang="css"></style>
