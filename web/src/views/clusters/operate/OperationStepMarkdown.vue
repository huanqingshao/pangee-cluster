<template>
  <div>
    <div>{{ test }}</div>
    <div v-if="markdownContent" v-html="docsHtml" class="operation-step-markdown"></div>
  </div>
</template>

<script>
import markdown from 'markdown-it'

export default {
  props: {
    cluster: { type: Object, required: true },
    operationIndex: { type: Number, required: true },
    stepIndex: { type: Number, required: true }
  },
  data() {
    return {
      markdownContent: "",
      test: "![web](./dev.assets/iShot_2022-08-06_20.21.03.png)".replaceAll("](./", "/resource-package/")
    }
  },
  computed: {
    path() {
      let result = this.cluster.resourcePackage.metadata.version;
      result += "/content/operations/";
      if (this.cluster.resourcePackage.operations[this.operationIndex] == undefined) {
        return "";
      }
      result += this.cluster.resourcePackage.operations[this.operationIndex]?.name;
      if (this.cluster.resourcePackage.operations[this.operationIndex]?.steps[this.stepIndex] == undefined) {
        return "";
      }
      result += "/" + this.cluster.resourcePackage.operations[this.operationIndex]?.steps[this.stepIndex]?.name;
      return result;
    },
    docsHtml() {
      if (this.markdownContent) {
        let md = new markdown()
        return md.render(this.markdownContent)
      }
      return undefined
    }
  },
  watch: {
    path(newValue, oldValue) {
      this.refresh();
    }
  },
  mounted() {
    this.refresh();
  },
  methods: {
    refresh() {
      if (this.path) {
        this.kuboardSprayApi
          .get(`/filebrowser/get?filePath=data:///resource/${this.path}/README.md`).then(resp => {
            let content = resp.data.data.content;
            content = content.replaceAll("](./", "](/resource-package/" + this.path + "/");
            this.markdownContent = content;
          }).catch(e => {
            if (e.status == 404) {
              this.markdownContent = e.response.data.message
            }
          })
      } else {
        this.markdownContent = "";
      }
    }
  }
}
</script>

<style lang="scss" scoped></style>

<style lang="scss">
.operation-step-markdown {
  font-family: Consolas, Menlo, Bitstream Vera Sans Mono, Monaco, "微软雅黑", monospace;
  font-size: 14px;

  h1 {
    font-size: 18px;
  }

  h2 {
    font-size: 16px;
    background-color: var(--el-color-primary-light-9);
    padding: 10px 20px;
  }

}
</style>