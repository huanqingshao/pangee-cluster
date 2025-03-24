<i18n locale="en" lang="yaml">
hi: kuboard
</i18n>

<i18n locale="zh" lang="yaml">
hi: kuboard
</i18n>

<template>
  <div>
    <ConfigGlobal :cluster="cluster" @refresh="$emit('refresh')"></ConfigGlobal>
    <ConfigK8sCluster :cluster="cluster" @switchTab="$emit('switchTab', $event)"></ConfigK8sCluster>
    <ConfigEtcd :cluster="cluster"></ConfigEtcd>
    <ConfigAddons
      :cluster="cluster"
      :currentTab="currentPropertiesTab"
      :pingpong="pingpong"
      @refresh="$emit('refresh')"
    ></ConfigAddons>
  </div>
</template>

<script>
import ConfigK8sCluster from "./k8s_cluster/ConfigK8sCluster.vue";
import ConfigGlobal from "./global/ConfigGlobal.vue";
import ConfigAddons from "./addons/ConfigAddons.vue";
import ConfigNode from "./node/ConfigNode.vue";
import CopyGapNodeToInventory from "./node/CopyGapNodeToInventory.vue";
import ConfigEtcd from "./etcd/ConfigEtcd.vue";

export default {
  props: {
    cluster: { type: Object, required: true }
  },
  data() {
    return {
      currentPropertiesTab: "global_config",
      pingpong: {},
      pingpong_loading: false
    };
  },
  components: {
    ConfigK8sCluster,
    ConfigGlobal,
    ConfigAddons,
    ConfigNode,
    ConfigEtcd,
    CopyGapNodeToInventory
  }
};
</script>

<style scoped lang="scss"></style>
