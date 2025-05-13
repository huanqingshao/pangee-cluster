<template>
  <div>
    <div v-if="markdownContent" v-html="docsHtml" class="plan-markdown"></div>
    <div id="markdown-image-modal" class="plan-markdown-image-modal" onclick="hideMarkdownImageModal()">
      <img id="markdown-image-modal-image" class="plan-markdown-image-modal-content">
    </div>
  </div>
</template>

<script>
import markdown from 'markdown-it'
import axios from 'axios'

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
      if (this.cluster.resourcePackage.data.operations[this.operationIndex] == undefined) {
        return "";
      }
      result += this.cluster.resourcePackage.data.operations[this.operationIndex]?.name;
      if (this.cluster.resourcePackage.data.operations[this.operationIndex]?.steps[this.stepIndex] == undefined) {
        return "";
      }
      result += "/" + this.cluster.resourcePackage.data.operations[this.operationIndex]?.steps[this.stepIndex]?.name;
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
        axios
          .get(`/resource-package/${this.path}/README.md`).then(resp => {
            let content = resp.data;
            content = content.replaceAll("](./", "](/resource-package/" + this.path + "/");
            this.markdownContent = content;
          }).catch(e => {
            if (e.status == 404) {
              this.markdownContent = "* Failed to get `/resource-package/" + this.path + "/README.md` \n* " + e.message;
              console.log(this.markdownContent)
            }
          })
      } else {
        this.markdownContent = "";
      }
    }
  }
}
</script>

<style lang="css" scoped>
@import "../plan/markdown.scss";
</style>