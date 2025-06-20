<i18n>
en:
  cluster: Cluster
  clusterList: Clusters List
  plan: Cluster Plan
  access: Access Cluster
  node: Nodes Maintainance
  health_check: Health Check
  backup: Backup/Restore
  csi_scan: CSI Scan
  upgrade: Upgrade Cluster
  faile_validate_form: Please validate the form
  operation: Operation
  resourcePackage: Resource Package
zh:
  clusterList: 集群列表
  cluster: 集群
  access: 访问集群
  plan: 规划集群
  nodes: 节点维护
  health_check: 状态检查
  backup: 备份/恢复
  csi_scan: CIS扫描
  upgrade: 升级集群
  faile_validate_form: 请完善表单
  operation: 操作集群
  resourcePackage: 资源包
</i18n>

<template>
  <div style="height: calc(100% - 54px)">
    <ControlBar :title="name">
      <span v-show="cluster && !cluster.history.processing" style="margin-right: 10px">
        <template v-if="currentTab === 'plan'">
          <template v-if="mode === 'view'">
            <el-button type="primary" icon="el-icon-edit" @click="$router.replace(`/clusters/${name}?mode=edit`)"
              :loading="loading">
              {{ $t("msg.edit") }}
            </el-button>
          </template>
          <template v-if="mode === 'edit'">
            <el-popconfirm :confirm-button-text="$t('msg.ok')" :cancel-button-text="$t('msg.cancel')"
              placement="bottom-start" icon="el-icon-warning" icon-color="red" :title="$t('msg.confirmToCancel')"
              @confirm="cancelEdit">
              <template #reference>
                <el-button type="default" icon="el-icon-close">{{ $t("msg.cancel") }}</el-button>
              </template>
            </el-popconfirm>
            <el-button type="primary" icon="el-icon-check" :disabled="noSaveRequired" @click="save">
              {{ $t("msg.save") }}
            </el-button>
          </template>
        </template>
      </span>
      <template v-if="cluster">
        <ExecuteTask :history="cluster.history" @refresh="refresh" :prompt-on-processing="true"></ExecuteTask>
      </template>
      <template v-if="cluster && isClusterInstalled">
        <ClusterStateNodes :cluster="cluster"></ClusterStateNodes>
      </template>
    </ControlBar>
    <el-card shadow="never" v-if="loading">
      <el-skeleton animated :rows="10" style="height: calc(100vh - 190px)"></el-skeleton>
    </el-card>
    <el-tabs type="border-card" v-else v-model="currentTab" class="app_scrollable_tabs">
      <el-tab-pane name="resourcePackage">
        <template #label>
          {{ t("resourcePackage") }}
        </template>
        <div class="app_scroll_content">
          <ConfigPangeeCluster :cluster="cluster"></ConfigPangeeCluster>
        </div>
      </el-tab-pane>
      <el-tab-pane :label="t('plan')" name="plan" style="overflow: visible;">
        <Plan v-if="cluster" ref="plan" :cluster="cluster" :mode="mode" @refresh="refresh"
          @switchTab="currentTab = $event">
        </Plan>
      </el-tab-pane>
      <el-tab-pane :label="t('operation')" name="operation">
        <Operation v-if="currentTab == 'operation'" ref="operation" :cluster="cluster" @refresh="refresh" @go-to-plan-hosts="onGoToPlanHosts"></Operation>
      </el-tab-pane>
      <el-tab-pane :label="t('access')" name="access" :disabled="disableNonePlanTab || !isClusterOnline">
        <Access v-if="currentTab == 'access'" ref="access" :cluster="cluster" :loading="loading"
          @switch="currentTab = $event"></Access>
      </el-tab-pane>
      <el-tab-pane v-if="false && cluster && cluster.inventory" :disabled="disableNonePlanTab || !isClusterOnline"
        :label="t('health_check')" name="health_check">
        <ClusterHealthCheck v-if="isClusterInstalled && isClusterOnline && currentTab == 'health_check'"
          :cluster="cluster" @refresh="refresh"></ClusterHealthCheck>
        <el-skeleton v-else style="height: calc(100vh - 220px)"></el-skeleton>
      </el-tab-pane>
      <el-tab-pane v-if="false" :disabled="disableNonePlanTab || !isClusterOnline" :label="t('backup')" name="backup">
        <Backup v-if="isClusterInstalled && isClusterOnline && currentTab == 'backup'" :cluster="cluster"
          @refresh="refresh"></Backup>
        <el-skeleton v-else animated :rows="10" style="height: calc(100vh - 220px)"></el-skeleton>
      </el-tab-pane>
      <el-tab-pane v-if="false" :disabled="disableNonePlanTab || !isClusterOnline" :label="t('upgrade')" name="upgrade">
        <template v-if="currentTab == 'upgrade'">
          <el-skeleton v-if="loading"></el-skeleton>
          <Upgrade v-else :cluster="cluster" @refresh="refresh"></Upgrade>
        </template>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import mixin from "../../mixins/mixin.js";
import yaml from "js-yaml";
import Access from "./access/Access.vue";
import Backup from "./backup/Backup.vue";
import { computed } from "vue";
import ClusterStateNodes from "./ClusterStateNodes.vue";
import clone from "clone";
import Plan from "./plan/Plan.vue";
import Upgrade from "./upgrade/Upgrade.vue";
import ClusterHealthCheck from "./health_check/ClusterHealthCheck.vue";
import Operation from "./operate/Operation.vue";
import ConfigPangeeCluster from "./plan/pangeecluster/ConfigPangeeCluster.vue";
import ExecuteTask from "../common/task/ExecuteTask.vue";

export default {
  mixins: [mixin],
  props: {
    name: { type: String, required: true },
    mode: { type: String, required: false, default: "view" }
  },
  percentage() {
    return this.percentage;
  },
  breadcrumb() {
    return [{ label: this.t("clusterList"), to: "/clusters" }, { label: this.name }];
  },
  refresh() {
    this.refresh();
  },
  data() {
    return {
      percentage: 0,
      cluster: undefined,
      originalInventoryYaml: "",
    };
  },
  computed: {
    currentTab: {
      get() {
        if (this.$store.state.cluster[this.name] == undefined) {
          return "plan";
        }
        return this.$store.state.cluster[this.name].currentTab || "plan"
      },
      set(v) {
        this.$store.commit("cluster/CHANGE_TAB", { cluster: this.name, currentTab: v });
      }
    },
    loading() {
      return this.percentage < 100;
    },
    noSaveRequired() {
      if (this.cluster === undefined) {
        return true;
      }
      return this.originalInventoryYaml == yaml.dump(this.cluster.inventory);
    },
    isClusterInstalled() {
      return !(this.cluster && this.cluster.history.success_tasks.length == 0);
    },
    isClusterOnline() {
      if (this.cluster) {
        if (this.cluster.state && this.cluster.state.code === 200) {
          return true;
        }
      }
      return false;
    },
    disableNonePlanTab() {
      if (this.mode !== "view") {
        return true;
      }
      return !this.isClusterInstalled;
    }
  },
  provide() {
    return {
      isClusterInstalled: computed(() => {
        return this.isClusterInstalled;
      }),
      isClusterOnline: computed(() => {
        return this.isClusterOnline;
      }),
      pendingRemoveNodes: computed(() => {
        return this.pendingNodes("delete_node");
      }),
      pendingAddNodes: computed(() => {
        return this.pendingNodes("add_node");
      }),
      onlineNodes: computed(() => {
        let result = {};
        if (this.cluster && this.cluster.state) {
          for (let key in this.cluster.state.nodes) {
            result[key] = { k8s_node: this.cluster.state.nodes[key] };
          }
          for (let key in this.cluster.state.etcd_members) {
            for (let nodeName in this.cluster.inventory.all.children.target.children.etcd.hosts) {
              // let etcdNode = this.cluster.inventory.all.hosts[nodeName];
              let etcdState = this.cluster.state.etcd_members[key];
              if (etcdState.name.indexOf(nodeName) > 0) {
                result[nodeName] = result[nodeName] || {};
                result[nodeName].etcd_member = etcdState;
              }
            }
          }
        }
        return result;
      })
    };
  },
  components: {
    Plan,
    Access,
    ClusterStateNodes,
    ClusterHealthCheck,
    Upgrade,
    Backup,
    Operation,
    ConfigPangeeCluster,
    ExecuteTask
  },
  watch: {
    "cluster.inventory.all.hosts.localhost.pangeecluster_resource_package": function () {
      this.loadResourcePackage();
    }
  },
  mounted() {
    this.refresh();
  },
  methods: {
    onGoToPlanHosts() {
      console.log("onGoToPlanHosts")
      this.currentTab = 'plan';
      this.$store.commit("cluster/CHANGE_CLUSTER_STATE",
        {
          cluster: this.cluster.name,
          tab: "plan",
          key: "currentTab",
          value: "hosts"
        }
      );
    },
    pendingNodes(action) {
      let result = [];
      if (this.isClusterInstalled) {
        for (let key in this.cluster.inventory.all.hosts) {
          let host = this.cluster.inventory.all.hosts[key];
          if (host.pangeecluster_node_action === action) {
            let h = clone(host);
            h.name = key;
            if (this.cluster.inventory.all.children.target.children.k8s_cluster.children.kube_control_plane.hosts[key]) {
              h.isControlPlane = true;
            }
            if (this.cluster.inventory.all.children.target.children.k8s_cluster.children.kube_node.hosts[key]) {
              h.isNode = true;
            }
            if (this.cluster.inventory.all.children.target.children.etcd.hosts[key]) {
              h.isEtcd = true;
            }
            result.push(h);
          }
        }
      }
      return result;
    },
    cancelEdit() {
      this.$router.replace(`/clusters/${this.name}`);
      this.refresh();
    },
    async refresh() {
      // console.log('refresh cluster.')
      this.percentage = 10;
      let _this = this;
      await this.pangeeClusterApi
        .get(`/clusters/${this.name}`)
        .then(async function (resp) {
          _this.cluster = resp.data.data;
          _this.originalInventoryYaml = yaml.dump(_this.cluster.inventory);
          for (let key in _this.cluster.inventory.all.hosts) {
            if (_this.cluster.inventory.all.hosts[key].pangeecluster_node_action === "upgrade_node") {
              _this.currentTab = "upgrade";
              break;
            }
          }
          _this.percentage = 30;
          await _this.loadResourcePackage();
          _this.percentage = 40;
          setTimeout(() => {
            if (_this.$refs.plan) {
              _this.$refs.plan.ping();
            }
          }, 200);
          if (_this.isClusterInstalled) {
            _this.percentage = 100;
            await _this.loadStateNodes();
          } else {
            _this.percentage = 100;
          }
        })
        .catch(e => {
          _this.percentage = 100;
          console.log(e);
        });
    },
    stateNodesLoaded(state) {
      if (state.code > 0 && state.etcd_code > 0) {
        this.cluster.state = state;
      }
    },
    async loadStateNodes() {
      let temp = { nodes: {}, code: 0, etcd_members: {}, etcd_code: 0, addons: {}, addon_code: 0 };
      await this.pangeeClusterApi
        .get(`/clusters/${this.name}/state/nodes`)
        .then(resp => {
          if (resp.data.data.stdout_obj && resp.data.data.stdout_obj.items) {
            temp.code = 200;
            for (let item of resp.data.data.stdout_obj.items) {
              temp.nodes[item.metadata.name] = item;
            }
          }
          if (resp.data.data.return_code === "" && resp.data.data.stdout_obj.changed === false && resp.data.data.stdout_obj.msg) {
            temp.code = 500;
            temp.msg = resp.data.data.stdout_obj.msg;
          }
          if (resp.data.data.stderr) {
            temp.code = resp.data.data.rc;
            temp.msg = resp.data.data.stderr;
          }
        })
        .catch(e => {
          temp.code = 500;
          temp.msg = e.response.data.message;
        });
      // this.percentage += 20
      this.pangeeClusterApi
        .get(`/clusters/${this.name}/state/etcd_member_health`)
        .then(resp => {
          // console.log(resp.data.data)
          if (resp.data.data && resp.data.data.length == 2) {
            temp.etcd_code = 200;
            let count = 0;
            for (let item of resp.data.data[0].stdout_obj.members) {
              temp.etcd_members[item.name] = item;
              for (let health of resp.data.data[1].stdout_obj) {
                if (item.clientURLs.indexOf(health.endpoint) >= 0) {
                  item.health = health;
                }
              }
              count++;
            }
            temp.etcd_members_count = count;
          }
          if (resp.data.data.return_code == "" && resp.data.data.stdout_obj.changed === false && resp.data.data.stdout_obj.msg) {
            temp.etcd_code = 500;
            temp.etcd_msg = resp.data.data.stdout_obj.msg;
          }
          for (let index in resp.data.data) {
            let a = resp.data.data[index];
            if (a.stderr) {
              temp.etcd_code = a.rc;
              temp.etcd_msg = a.stderr;
            }
          }
          this.stateNodesLoaded(temp);
        })
        .catch(e => {
          console.log(e);
          temp.etcd_code = 500;
          temp.etcd_msg = e.response.data.message;
          this.stateNodesLoaded(temp);
        });
      // this.percentage += 20
      // this.pangeeClusterApi
      //   .get(`/clusters/${this.name}/state/addons`)
      //   .then(resp => {
      //     temp.addon_code = 200;
      //     temp.addons = resp.data.data;
      //     this.stateNodesLoaded(temp);
      //   })
      //   .catch(e => {
      //     temp.addon_code = 500;
      //     temp.addons = {};
      //     console.error(e);
      //     this.stateNodesLoaded(temp);
      //   });
      // this.percentage += 20
      // this.cluster.state = temp
    },
    async loadResourcePackage() {
      this.cluster.resourcePackage = undefined;
      let newValue = this.cluster.inventory.all.hosts.localhost.pangeecluster_resource_package;
      if (newValue) {
        await this.pangeeClusterApi
          .get(`/resources/${newValue}`)
          .then(resp => {
            this.cluster.resourcePackage = resp.data.data.package;
          })
          .catch(e => {
            console.log(e);
          });
      }
    },
    save() {
      this.$refs.plan.validate(flag => {
        if (flag) {
          this.pangeeClusterApi
            .put(`/clusters/${this.name}`, this.cluster.inventory)
            .then(() => {
              this.$message.success(this.$t("msg.save_succeeded"));
              this.refresh();
              this.$router.replace(`/clusters/${this.name}`);
            })
            .catch(e => {
              this.$message.error(this.$t("msg.save_failed", e.response.data.message));
            });
        } else {
          this.$message.error(this.t("faile_validate_form"));
        }
      });
    }
  }
};
</script>

<style scoped lang="scss"></style>
