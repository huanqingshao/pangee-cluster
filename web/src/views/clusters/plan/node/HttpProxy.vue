<i18n>
en:
  proxy: Proxy
  proxyDescription: Proxy (Generally speaking, you don't need this.)
  proxyUsage1: PangeeCluster can use the following proxy params to fetch content from internet that are not included in the resource package.
  proxyUsage2: "Usually speaking, all the resources (exclude docker-ce / OS mirror) are already included in the resource package, you don't need the proxy params here."
zh:
  proxy: 代理
  proxyDescription: 代理（通常不需要设置）
  proxyUsage1: PangeeCluster 可以使用如下代理参数从外网获取资源包中未找到的资源；
  proxyUsage2: 通常资源包中包含所需资源，您无需设置此处的代理参数；
</i18n>

<template>
  <ConfigSection v-model:enabled="proxyEnabled" :label="t('proxy')" :description="t('proxyDescription')" anti-freeze>
    <el-alert class="app_margin_bottom" :closable="false">
      <li>{{ t('proxyUsage1') }}</li>
      <li>{{ t('proxyUsage2') }}</li>
    </el-alert>
    <EditString v-model="vars.http_proxy" label="http_proxy" prop="all.children.target.vars" anti-freeze>
    </EditString>
    <EditString v-model="vars.https_proxy" label="https_proxy" prop="all.children.target.vars" anti-freeze>
    </EditString>
    <EditString v-model="vars.additional_no_proxy" label="no_proxy" prop="all.children.target.vars" anti-freeze>
    </EditString>
  </ConfigSection>
</template>

<script>
export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {

    }
  },
  computed: {
    inventory: {
      get() { return this.cluster.inventory },
      set() { },
    },
    vars: {
      get() {
        return this.cluster.inventory.all.children.target.vars
      },
      set() { }
    },
    proxyEnabled: {
      get() {
        if (this.vars) {
          return this.vars.http_proxy !== undefined
        }
        return false
      },
      set(v) {
        if (v) {
          this.vars = this.vars || {}
          this.vars.http_proxy = ''
        } else {
          delete this.vars.http_proxy
          delete this.vars.https_proxy
          delete this.vars.additional_no_proxy
        }
      }
    },
  },
  components: {},
  mounted() {
  },
  methods: {
  }
}
</script>

<style scoped lang="css"></style>
