<i18n>
en:
  title: "Load resource package {name}."
  selectSource: "Select a source"
  download: "Download resource package"
  useProxy: " Use proxy"
zh:
  title: "加载资源包 {name}"
  selectSource: "选择一个源"
  download: "加载资源包"
  useProxy: " 使用代理"
</i18n>

<template>
  <ExecuteTask v-if="resource" :history="resource.history" :startTask="startTask" :label="t('download')"
    :title="t('title', { name: resource.package.metadata.version })" :loading="loading" @refresh="$emit('refresh')">
    <el-form @submit.prevent.stop :model="form" ref="form" label-position="left" label-width="120px">
      <!-- <el-form-item :label="t('selectSource')" prop="downloadFrom" :rules="sourceRules">
        <el-radio-group v-model="form.downloadFrom">
          <div style="line-height: 28px; padding-top: 5px">
            <div v-for="(source, index) in resource.package.metadata.available_at" :key="'source' + index">
              <el-radio :label="source" :value="source">
                <span class="app_text_mono">{{ source }}</span>
              </el-radio>
            </div>
          </div>
        </el-radio-group>
      </el-form-item> -->
      <el-form-item :label="t('useProxy')">
        <el-switch v-model="form.enableProxy" />
      </el-form-item>
      <el-form-item v-if="form.enableProxy" label="Http proxy">
        <el-input v-model="form.httpProxy" />
      </el-form-item>
    </el-form>
  </ExecuteTask>
</template>

<script>
import clone from "clone";
import ExecuteTask from "../common/task/ExecuteTask.vue";
import { useI18n } from "vue-i18n";

export default {
  props: {
    resource: { type: Object, required: true },
    loading: { type: Boolean, required: false },
    action: { type: String, required: true },
    source: { type: String, required: true }
  },
  data() {
    return {
      form: {
        downloadFrom: undefined,
        enableProxy: false,
        httpProxy: ''
      },
      sourceRules: [{ required: true, message: this.i18n("selectSource"), trigger: "change" }]
    };
  },
  setup() {
    const { t, locale } = useI18n({
      inheritLocale: true,
      useScope: "local"
    });
    return { i18n: t };
  },
  components: { ExecuteTask },
  emits: ["refresh"],
  mounted() {
    // this.form.downloadFrom = this.resource.package.metadata.available_at[0];
  },
  methods: {
    startTask() {
      return new Promise((resolve, reject) => {
        this.$refs.form.validate(async flag => {
          if (flag) {
            let request = {
              package: clone(this.resource.package),
              downloadFrom: this.source + '.com',
              enableProxy: this.form.enableProxy.toString(), // 转为 "true" 或 "false"
              httpProxy: this.form.httpProxy
            };
            this.pangeeClusterApi
              .post(`/resources/${request.package.metadata.version}/${this.action}`, request)
              .then(resp => {
                this.$router.replace(`/settings/resources/${request.package.metadata.version}`);
                resolve(resp.data.data.pid);
              })
              .catch(e => {
                reject(e);
              });
          } else {
            resolve();
          }
        });
      });
    }
  }
};
</script>

<style scoped lang="css">
.confirmText {
  font-size: 15px;
  color: var(--el-color-danger);
  font-weight: bold;
}
</style>
