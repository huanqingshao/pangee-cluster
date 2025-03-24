<i18n locale="en" lang="yaml">
hi: kuboard
</i18n>

<i18n locale="zh" lang="yaml">
hi: kuboard
</i18n>

<template>
  <div>
    显示架构规划的文档
    <div>
      <br />
      <KuboardSprayLink href="https://kuboard-spray.cn/guide/install-k8s.html" :size="12">安装 K8S 集群？ </KuboardSprayLink>
      <br />
      <KuboardSprayLink href="https://kuboard-spray.cn/guide/maintain/ha-mode.html" :size="12"> 实现高可用？ </KuboardSprayLink>
      <br />
      <KuboardSprayLink href="https://kuboard-spray.cn/guide/maintain/upgrade.html" :size="12">升级集群？ </KuboardSprayLink>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {};
  }
};
</script>

<style scoped lang="scss"></style>
