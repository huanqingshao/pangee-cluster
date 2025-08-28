<i18n>
en:
  title: "Load resource package {name}."
  selectSource: "Select a source"
  download: "Download resource package"
  upload: "Upload Resource Package"
  useProxy: " Use proxy"
zh:
  title: "加载资源包 {name}"
  selectSource: "选择一个源"
  download: "加载资源包"
  upload: "加载离线资源包"
  useProxy: " 使用代理"
</i18n>

<template>
  <ExecuteTask 
    :history="resource.history"
    :startTask="action !== 'upload' ? startTask : localPack"
    :label="action !== 'upload' ? t('download') : t('upload')"
    :title="action !== 'upload' ? t('title', { name: resource.package.metadata.version }) : t('upload')"
    :loading="loading"
    @refresh="$emit('refresh')"
  >
    <el-form @submit.prevent.stop :model="form" ref="form" label-position="left" label-width="120px">
      <el-form-item :label="t('useProxy')">
        <el-switch v-model="form.enableProxyOnDownload" />
      </el-form-item>
      <el-form-item v-if="form.enableProxyOnDownload" label="Http proxy">
        <el-input v-model="form.httpProxy" />
      </el-form-item>
    </el-form>
    <input 
      ref="fileInput" 
      type="file" 
      accept=".zip" 
      class="hidden" 
      @change="handleFileSelect"
    >
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
        enableProxyOnDownload: false,
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
        let request = {
          package: clone(this.resource.package),
          downloadFrom: this.source + '.com',
          tagName: this.resource.tag_name,
          fileName: this.resource.file_name,
          enableProxyOnDownload: this.form.enableProxyOnDownload.toString(), // 转为 "true" 或 "false"
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
      });
    },
    localPack() {
      this.$refs.fileInput.click();
      // 返回一个Promise供ExecuteTask使用，实际完成在handleFileSelect中
      return new Promise((resolve) => {
        this.uploadResolve = resolve;
      });
    },
    
    async handleFileSelect(e) {
      const file = e.target.files[0];
      if (!file) {
        this.uploadResolve?.();
        return;
      }
      
      // 验证文件类型
      if (!file.name.toLowerCase().endsWith('.zip')) {
        e.target.value = '';
        this.uploadResolve?.();
        return;
      }

      this.resource.history.task_name = file.name.replace(/\.zip$/i, '')

      try {
        const request = {
          enableProxyOnDownload: this.form.enableProxyOnDownload.toString(),
          httpProxy: this.form.httpProxy
        };
        const formData = new FormData();
        formData.append('file', file);

        Object.keys(request).forEach(key => {
          formData.append(key, request[key]);
        });

        this.pangeeClusterApi
          .post('/resources/upload', formData, { headers: { 'Content-Type': 'multipart/form-data' }})
          .then(resp => {
            this.$router.replace(`/settings/resources/${this.resource.history.task_name}`);
            this.uploadResolve(resp.data.data.pid);
          })
          .catch(e => {
            this.uploadResolve();
          });
      } catch (error) {
        console.error('Upload error:', error);
        this.uploadResolve();
      } finally {
        e.target.value = '';
      }
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

.hidden {
  display: none;
}
</style>
