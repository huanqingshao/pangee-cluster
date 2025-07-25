<i18n>
en:
  getKubeconfig: Fetch kubeconfig
  accessFromControlPlane: Use kubectl on control_plane
  accessMethods: You can use differenct ways to access the cluster
  controlPlanes: SSH connect to any of the following nodes, and use kubectl command to administrate the cluster.
  proposeKuboard: Kuboard is a popular K8S cluster management UI, you can refer to its website to learn how to install and use it.
  switchToPlan: Swith to cluster plan view.
  etcdAccess: SSH connect to any of the following nodes, and use etcdctl command to administrate the etcd cluster.
  requiredToSyncEtcd: You removed a etcd node, and it's a must to sync etcd config to all kube_control_plane and etcd nodes.
  yourcommand: Execute your command starting from here.

  terminal: Open ssh terminal
zh:
  getKubeconfig: 获取 kubeconfig 文件
  accessFromControlPlane: 在主节点上使用 kubectl
  accessMethods: 您可以使用多种方式对集群进行管理控制
  controlPlanes: 您可以 ssh 到如下节点中的任意一个，直接执行 kubectl 命令可以管理集群。
  proposeKuboard: Kuboard 是一款非常值得推荐的 K8S 集群管理界面，请参考 Kuboard 网站，安装该管理界面。
  switchToPlan: 切换到集群规划页
  etcdAccess: 您可以 ssh 到如下节点中的任意一个，执行以下指令后，可以通过 etcdctl 操作 etcd 集群。通常您并不需要直接操作 etcd。
  requiredToSyncEtcd: 您删除了 ETCD 节点，必须将 ETCD 配置同步到所有控制节点和 ETCD 节点
  yourcommand: 此处开始，执行您想要执行的 etcdctl 命令

  terminal: 打开 SSH 终端
</i18n>


<template>
  <el-skeleton v-if="cluster.state === undefined || cluster.state.code === 0" style="height: calc(100vh - 220px);"
    animated></el-skeleton>
  <el-scrollbar height="calc(100vh - 220px)" v-else-if="cluster.state.code === 200">
    <!-- <el-alert type="error" effect="dark" class="app_margin_bottom" :title="t('requiredToSyncEtcd')" :closable="false" show-icon></el-alert> -->
    <el-alert type="info" :title="t('accessMethods')" :closable="false"></el-alert>
    <div class="app_block_title">{{ t('accessFromControlPlane') }}</div>
    <div class="access_details" v-if="cluster">
      <el-alert :title="t('controlPlanes')" :closable="false" type="success"></el-alert>
      <div class="details" v-if="cluster.state && cluster.state.nodes">
        <template v-for="(item, key) in cluster.state.nodes" :key="key">
          <div v-if="isControlePlane(key, item)" class="app_margin_top">
            <el-tag class="node_text" size="default">
              <span class="app_text_mono">{{ key }}</span>
            </el-tag>
            <el-tag class="node_text" effect="light" size="default">
              <span class="app_text_mono">
                {{ cluster.inventory.all.hosts[controlePlaneIpToNodename[key]].ansible_host }}
              </span>
            </el-tag>
            <el-button @click="openUrlInBlank(`#/ssh/cluster/${cluster.name}/${controlePlaneIpToNodename[key]}`)"
              style="margin-left: 10px;" icon="el-icon-monitor" type="primary">{{ t('terminal') }}</el-button>
          </div>
        </template>
      </div>
    </div>
    <div class="app_block_title">kubeconfig</div>
    <div class="access_details">
      <el-button type="primary" icon="el-icon-files" @click="fetchKubeconfig" :loading="kubeconfigLoading">{{
        t('getKubeconfig') }}</el-button>
      <CopyToClipBoard v-if="kubeconfig" :value="kubeconfig"></CopyToClipBoard>
      <el-skeleton class="app_margin_top" v-if="kubeconfigLoading" animated></el-skeleton>
      <div v-if="kubeconfig && !kubeconfigLoading" class="app_margin_top app_codemirror_auto_height">
        <Codemirror v-model:value="kubeconfig" :options="cmOptions"></Codemirror>
      </div>
    </div>
    <div class="app_block_title">etcd</div>
    <div class="access_details" v-if="cluster.state">
      <el-alert :closable="false" type="success" effect="dark" :title="t('etcdAccess')"></el-alert>
      <div class="details">
        <template v-for="(item, key) in cluster.state.etcd_members" :key="'etcd' + key">
          <div style="margin-top: 10px;">
            <el-tag class="node_text" type="primary" size="default"> {{ etcdIp(item) }} </el-tag>
            <el-tag class="node_text" type="primary" effect="light" size="default"> {{ etcdClientUrl(item) }} </el-tag>
            <template v-for="(etcd, name) in cluster.inventory.all.children.target.children.etcd.hosts"
              :key="'eb' + name + key">
              <el-button v-if="cluster.inventory.all.hosts[name].ip === etcdIp(item)"
                @click="openUrlInBlank(`#/ssh/cluster/${cluster.name}/${name}`)" icon="el-icon-monitor" type="primary">
                {{ t('terminal') }} </el-button>
            </template>
          </div>
        </template>
        <div class="app_margin_bottom"></div>
        <CopyToClipBoard :value="etcdSsh"></CopyToClipBoard>
        <Codemirror v-if="showEtcdSsh" class="app_margin_top" v-model:value="etcdSsh" :options="etcdCmOptions">
        </Codemirror>
      </div>
    </div>
  </el-scrollbar>
  <el-alert v-else-if="cluster.state.code === 500" type="error" :closable="false" effect="dark" show-icon>
    <span v-if="cluster.state.msg" class="app_text_mono"
      v-html="cluster.state.msg.replaceAll('\n', '<br>').replaceAll('    ', '<span style=margin-right:20px;></span>')"></span>
    <span v-else>{{ cluster.state }}</span>
    <div style="margin-top: 20px;">
      <el-button type="primary" round icon="el-icon-arrow-left" @click="$emit('switch', 'plan')">
        {{ t('switchToPlan') }}
      </el-button>
    </div>
  </el-alert>
</template>

<script>
import Codemirror from "codemirror-editor-vue3"
import "codemirror/theme/darcula.css"
import "codemirror/mode/yaml/yaml.js"
import "codemirror/mode/shell/shell.js"

export default {
  props: {
    cluster: { type: Object, required: true },
    loading: { type: Boolean, required: true },
  },
  data() {
    return {
      kubeconfig: undefined,
      kubeconfigLoading: false,
      cmOptions: {
        mode: "yaml", // Language mode
        theme: "darcula", // Theme
        lineNumbers: true, // Show line number
        lineWrapping: true,
        smartIndent: true, // Smart indent
        indentUnit: 2, // The smart indent unit is 2 spaces in length
        foldGutter: true, // Code folding
        readOnly: true,
        styleActiveLine: true, // Display the style of the selected row
      },
      etcdCmOptions: {
        mode: "shell",
        theme: "darcula",
        lineNumbers: true,
        lineWrapping: true,
        smartIndent: true,
        indentUnit: 2,
        foldGutter: true,
        readOnly: true,
        styleActiveLine: true,
      },
      showEtcdSsh: false,
    }
  },
  computed: {
    etcdSsh: {
      get() {
        return `export ETCDCTL_API=3
export ETCDCTL_CERT=/etc/kubernetes/ssl/etcd_server.crt
export ETCDCTL_KEY=/etc/kubernetes/ssl/etcd_server.key
export ETCDCTL_CACERT=/etc/kubernetes/ssl/ca.crt
export ETCDCTL_ENDPOINTS=https://127.0.0.1:${this.cluster.inventory.all.children.target.children.etcd.vars.etcd_client_port || 2379}
etcdctl member list
# ${this.t('yourcommand')}
`
      },
      set() { }
    },
    kuboard_url() {
      let result = 'http://'
      for (let key in this.cluster.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts) {
        result += this.cluster.inventory.all.hosts[key].ip || this.cluster.inventory.all.hosts[key].ansible_host
        break
      }
      let addon = undefined
      for (let key in this.cluster.resourcePackage.data.addon) {
        addon = this.cluster.resourcePackage.data.addon[key]
        if (addon.name === 'kuboard')
          break
      }
      let kuboard_port = this.cluster.inventory.all.children.target.children.k8s_cluster.vars.kuboard_port || addon.params_default.kuboard_port
      if (kuboard_port != 80) {
        result += ':' + kuboard_port
      }
      return result
    },
    controlePlaneIpToNodename() {
      let result = {};
      for (let key in this.cluster.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts) {
        let t = this.cluster.inventory.all.hosts[key];
        result[t.ip] = key;
      }
      return result;
    }
  },
  components: { Codemirror },
  mounted() {
    setTimeout(() => {
      this.showEtcdSsh = true
    }, 300)
  },
  watch: {
    loading(newValue) {
      if (newValue) {
        this.kubeconfig = undefined
      }
    }
  },
  methods: {
    isControlePlane(key, node) {
      if (this.controlePlaneIpToNodename[key] != undefined) {
        return true;
      }
      return false;
    },
    fetchKubeconfig() {
      this.kubeconfigLoading = true
      this.kubeconfig = undefined
      this.pangeeClusterApi.get(`/clusters/${this.cluster.name}/access/kubeconfig`).then(resp => {
        console.log(resp.data)
        let out = resp.data.data
        this.kubeconfig = out.stdout.replace("127.0.0.1", this.cluster.inventory.all.hosts[out.ansible_inventory_hostname].ansible_host)
        this.kubeconfigLoading = false
      }).catch(e => {
        console.log(e)
        this.kubeconfigLoading = false
        this.$message.error('failed to get kubeconfig: ' + e.response.data.msg)
      })
    },
    etcdIp(item) {
      for (let i in item.clientURLs) {
        let temp = item.clientURLs[i]
        if (temp.indexOf("127.0.0.1") >= 0) {
          continue;
        }
        temp = temp.split(':')
        return temp[1].slice(2)
      }
      return ''
    },
    etcdClientUrl(item) {
      for (let i in item.clientURLs) {
        let temp = item.clientURLs[i]
        if (temp.indexOf("127.0.0.1") >= 0) {
          continue;
        }
        return temp
      }
      return ''
    }
  }
}
</script>

<style scoped lang="css">
.access_details {
  padding-left: 40px;
  margin-bottom: 20px;
}

.details {
  background-color: var(--el-color-info-light-9);
  padding: 10px 20px 20px 20px;
}

.node_text {
  margin-right: 10px;
  min-width: 150px;
}
</style>
