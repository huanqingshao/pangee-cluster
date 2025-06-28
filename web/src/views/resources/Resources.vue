<i18n>
en:
  title: Resource Packages
  resourceDescription1: Kuboard provides a set of pre-verified resource packages.
  resourceDescription2: You can also refer to https://github.com/opencmit/pangee-cluster-resource to create your own resource package.
  pangeecluster_version_min: PangeeCluster min version
  pangeecluster_version_max: PangeeCluster max version
  version: Version
  container_engine: Containerd Version
  supported_os: Supported OS
  kube_version: K8S Version
  kubespray_version: Kubespray Version
  import: Import
  import_status: Import status
  import_status_true: True
  import_status_false: False
  minVersionRequired: Min version required to pangeecluster
  cannot_reach_online_repository: Current browser cannot reach https://pangee-cluster.cn/support
zh:
  title: 资源包列表
  resourceDescription1: Kuboard 提供一组经过预先测试验证的资源包列表，可以帮助您快速完成集群安装
  resourceDescription2: 您也可以参考项目 https://github.com/opencmit/pangee-cluster-resource 自己创建资源包
  pangeecluster_version_min: PangeeCluster最低版本
  pangeecluster_version_max: PangeeCluster最高版本
  version: 资源包版本
  container_engine: Containerd 版本
  supported_os: 支持的操作系统
  kube_version: K8S 版本
  kubespray_version: Kubespray 版本
  import: 导 入
  import_status: 导入状态
  import_status_true: 已导入
  import_status_false: 未导入
  minVersionRequired: PangeeCluster最低版本要求
  cannot_reach_online_repository: 当前浏览器不能在线获取可选的资源包列表，请在可以访问外网的浏览器打开地址 https://pangee-cluster.cn/support
</i18n>

<template>
  <div>
    <div class="app_block_title">{{ t('title') }}</div>
    <div class="app_block_content">
      <el-alert :title="t('title')" type="info" :closable="false" class="app_white_alert">
        <div class="description">
          <li>{{ t('resourceDescription1') }}</li>
          <li>{{ t('resourceDescription2') }}</li>
        </div>
      </el-alert>
      <div style="text-align: right;">
        <PangeeClusterLink v-if="cannot_reach_online_repository" href="https://pangee-cluster.cn/support" type="danger"
          style="margin-right: 10px; color: var(--el-color-danger)">{{ t('cannot_reach_online_repository') }}
        </PangeeClusterLink>
        <ResourcesCreateOffline class="app_margin_top"></ResourcesCreateOffline>
      </div>
      <div class="contentList">
        <el-table ref="table" v-if="mergedPackageList" :data="mergedPackageList" style="width: 100%" row-key="version">
          <el-table-column prop="version" :label="t('version')" width="200px">
            <template #default="scope">
              <template v-if="hideLink">
                {{ scope.row.version }}
              </template>
              <template v-else-if="importedPackageMap">
                <router-link v-if="importedPackageMap[scope.row.version]"
                  :to="`/settings/resources/${scope.row.version}`">
                  <el-icon :size="12" style="width: 12px; height: 12px; vertical-align: middle;">
                    <el-icon-link></el-icon-link>
                  </el-icon>
                  {{ scope.row.version }}
                </router-link>
                <router-link v-else :to="`/settings/resources/${scope.row.version}/on_air`">
                  <el-icon :size="12" style="width: 12px; height: 12px; vertical-align: middle;">
                    <el-icon-link></el-icon-link>
                  </el-icon>
                  {{ scope.row.version }}
                </router-link>
              </template>
            </template>
          </el-table-column>
          <el-table-column :label="t('kube_version')" width="90px">
            <template #default="scope">
              <span v-if="packageYaml[scope.row.version]">
                {{ depVersion("kubernetes", scope.row.yaml) }}
              </span>
              <el-icon v-else :size="12" style="width: 12px; height: 12px; vertical-align: middle;" class="is-loading">
                <el-icon-loading></el-icon-loading>
              </el-icon>
            </template>
          </el-table-column>
          <el-table-column :label="t('container_engine')" width="110px">
            <template #default="scope">
              <template v-if="packageYaml[scope.row.version]">
                {{ depVersion("containerd", scope.row.yaml) }}
              </template>
              <el-icon v-else :size="12" style="width: 12px; height: 12px; vertical-align: middle;" class="is-loading">
                <el-icon-loading></el-icon-loading>
              </el-icon>
            </template>
          </el-table-column>
          <el-table-column :label="t('supported_os')">
            <template #default="scope">
              <template v-if="packageYaml[scope.row.version]">
                <span v-for="(os, key) in scope.row.yaml.metadata.supported_os" :key="`os${scope.index}_${key}`"
                  style="margin-right: 10px;">
                  <el-tag>
                    {{ os.distribution }}<span v-for="(v, i) in os.versions" :key="key + 'v' + i">_{{ v }}</span>
                  </el-tag>
                </span>
              </template>
              <el-icon v-else :size="12" style="width: 12px; height: 12px; vertical-align: middle;" class="is-loading">
                <el-icon-loading></el-icon-loading>
              </el-icon>
            </template>
          </el-table-column>
          <el-table-column :label="t('import_status')" width="200px">
            <template #default="scope">
              <template v-if="importedPackageMap">
                <template v-if="importedPackageMap[scope.row.version]">
                  <el-tag type="success" effect="dark">
                    <el-icon :size="14" style="width: 14px; height: 14px; vertical-align: bottom;">
                      <el-icon-download v-if="scope.row.isOffline"></el-icon-download>
                      <el-icon-cloudy v-else></el-icon-cloudy>
                    </el-icon>
                    {{ t('import_status_true') }}
                  </el-tag>
                </template>
                <el-tag v-else-if="scope.row.meetPangeeClusterVersion" type="warning" effect="dark">
                  <el-icon :size="14" style="width: 14px; height: 14px; vertical-align: bottom;">
                    <el-icon-circle-close></el-icon-circle-close>
                  </el-icon>
                  {{ t('import_status_false') }}
                </el-tag>
                <template v-else-if="scope.row.yaml">
                  <el-tag type="danger" effect="dark">{{ t('minVersionRequired') }}</el-tag>
                  <el-tag type="danger" class="app_text_mono">
                    {{ scope.row.yaml.metadata.kuboard_spray_version.min }}
                  </el-tag>
                </template>
              </template>
              <el-icon v-else :size="12" style="width: 12px; height: 12px; vertical-align: middle;" class="is-loading">
                <el-icon-loading></el-icon-loading>
              </el-icon>
            </template>
          </el-table-column>
          <slot name="columns">
            <el-table-column :label="t('msg.operations')" width="200px">
              <template #default="scope">
                <template v-if="importedPackageMap">
                  <template v-if="importedPackageMap[scope.row.version]">
                    <el-button type="primary" plain icon="el-icon-view"
                      @click="$router.push(`/settings/resources/${scope.row.version}`)">{{ $t('msg.view') }}</el-button>
                    <el-button type="danger" icon="el-icon-delete"
                      @click="$router.push(`/settings/resources/${scope.row.version}`)">{{ $t('msg.delete')
                      }}</el-button>
                  </template>
                  <template v-else>
                    <el-button type="primary" plain icon="el-icon-view"
                      @click="$router.push(`/settings/resources/${scope.row.version}/on_air`)">{{ $t('msg.view')
                      }}</el-button>
                    <el-button type="primary" v-if="scope.row.meetPangeeClusterVersion" icon="el-icon-download"
                      @click="$router.push(`/settings/resources/${scope.row.version}/on_air`)">{{ t('import')
                      }}</el-button>
                  </template>
                </template>
              </template>
            </el-table-column>
          </slot>
        </el-table>
      </div>
    </div>

  </div>
</template>

<script>
import axios from 'axios'
import yaml from 'js-yaml'
import ResourcesCreateOffline from './ResourcesCreateOffline.vue'
import compareVersions from 'compare-versions'
import repositoryPrefix from './repository_prefix.js'

export default {
  props: {
    hideLink: { type: Boolean, required: false, default: false }
  },
  data() {
    return {
      availablePackageList: undefined,
      importedPackageMap: {},
      packageYaml: {},
      cannot_reach_online_repository: false,
    }
  },
  computed: {
    mergedPackageList() {
      let result = []
      for (let i in this.importedPackageMap) {
        let flag = false
        for (let j in this.availablePackageList) {
          if (this.availablePackageList[j].version === i) {
            flag = true
          }
        }
        if (flag === false) {
          result.push({
            version: i,
            isOffline: true,
          })
        }
      }
      for (let i in this.availablePackageList) {
        result.push(this.availablePackageList[i])
      }
      return result
    },
  },
  components: { ResourcesCreateOffline },
  mounted() {
    this.refresh()
  },
  methods: {
    depVersion(depName, packageYaml) {
      for (let i in packageYaml?.data.dependency) {
        let dep = packageYaml.data.dependency[i];
        if (dep.name == depName) {
          return dep.version;
        }
      }
      return undefined;
    },
    async refresh() {
      this.importedPackageMap = {}
      this.availablePackageList = undefined
      this.cannot_reach_online_repository = false
      await axios.get(`${repositoryPrefix()}/package-list.yaml?nocache=${new Date().getTime()}`).then(resp => {
        this.availablePackageList = yaml.load(resp.data).items
      }).catch(e => {
        console.log(e)
        // this.$message.error('离线环境')
        this.cannot_reach_online_repository = true
      })
      await this.pangeeClusterApi.get(`/resources`).then(resp => {
        for (let i in resp.data.data) {
          this.importedPackageMap[resp.data.data[i]] = true
          for (let j in this.availablePackageList) {
            let pkg = this.availablePackageList[j]
            if (pkg.version === resp.data.data[i]) {
              pkg.imported = true
            }
          }
        }
      }).catch(e => {
        console.log(e)
      })
      for (let i in this.mergedPackageList) {
        let packageVersion = this.mergedPackageList[i]
        if (packageVersion.isOffline || packageVersion.imported) {
          this.loadPackageLocal(packageVersion)
        } else {
          this.loadPackageFromRepository(packageVersion)
        }
      }
    },
    loadPackageLocal(packageVersion) {
      this.pangeeClusterApi.get(`/resources/${packageVersion.version}`).then(resp => {
        this.packageYaml[packageVersion.version] = resp.data.data.package
        packageVersion.yaml = resp.data.data.package
        packageVersion.meetPangeeClusterVersion = compareVersions(window.PangeeCluster.version.trimed, packageVersion.yaml.metadata.kuboard_spray_version.min) >= 0
      })
    },
    loadPackageFromRepository(packageVersion) {
      axios.get(`${repositoryPrefix()}/${packageVersion.version}/package.yaml?nocache=${new Date().getTime()}`).then(resp => {
        setTimeout(() => {
          let temp = yaml.load(resp.data)
          this.packageYaml[packageVersion.version] = temp
          packageVersion.yaml = temp
          packageVersion.meetPangeeClusterVersion = compareVersions(window.PangeeCluster.version.trimed, packageVersion.yaml.metadata.kuboard_spray_version.min) >= 0
        }, 500)
      }).catch(e => {
        this.$message.error(e + '')
      })
    }
  }
}
</script>

<style scoped lang="css">
.description {
  line-height: 28px;
}

.contentList {
  margin: 10px 0;
}
</style>
