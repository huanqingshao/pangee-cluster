<template>
  <div>
    <div v-if="markdownContent" v-html="docsHtml" class="operation-step-markdown"></div>
    <div id="markdown-image-modal" class="markdown-image-modal" onclick="hideMarkdownImageModal()">
      <img id="markdown-image-modal-image" class="markdown-image-modal-content">
    </div>
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
      markdownContent: ""
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
        let md = new markdown();
        let temp = md.render(this.markdownContent);
        temp = temp.replaceAll("<img ", `<img onclick=showMarkdownImageModal(this) `);
        return temp;
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
    window.showMarkdownImageModal = (img) => {
      document.getElementById('markdown-image-modal').style.display = 'block';
      document.getElementById('markdown-image-modal-image').src = img.src;
    };
    window.hideMarkdownImageModal = () => {
      document.getElementById('markdown-image-modal').style.display = 'none';
    }
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

<style lang="scss" scoped>
.markdown-image-modal {
  z-index: 900000;
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.8);
  cursor: pointer;
}

.markdown-image-modal-content {
  max-width: 90%;
  max-height: 90%;
  margin: auto;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  position: absolute;
}
</style>

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

  p img {
    max-width: 100%;
    cursor: zoom-in;
  }

  p code {
    background-color: var(--el-color-black);
    color: #c7254e !important;
    padding: .25rem .5rem;
    margin: 0;
    font-size: .9em;
    background-color: rgba(27, 31, 35, .05);
    border-radius: 3px;
    font-family: Consolas, Menlo, "Bitstream Vera Sans Mono", Monaco, "微软雅黑", monospace;
  }

  pre {
    padding-left: .5rem;
    vertical-align: middle;
    line-height: 1.4;
    padding: 0.8rem 1rem;
    margin: .85rem 0;
    background-color: #606266;
    color: #fff;
    border-radius: 6px;
    overflow: auto;

    code {
      font-family: Consolas, Menlo, "Bitstream Vera Sans Mono", Monaco, "微软雅黑", monospace;
    }
  }

}
</style>