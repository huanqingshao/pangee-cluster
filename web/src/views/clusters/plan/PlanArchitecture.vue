<template>
  <div>
    <div v-if="markdownContent" v-html="docsHtml" class="architecture-markdown"></div>
    <div id="architecture-markdown-image-modal" class="architecture-markdown-image-modal"
      onclick="hideArchMarkdownImageModal()">
      <img id="architecture-markdown-image-modal-image" class="architecture-markdown-image-modal-content">
    </div>
  </div>
</template>

<script>
import markdown from 'markdown-it'
import axios from 'axios'

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      markdownContent: ""
    }
  },
  computed: {
    path() {
      let result = this.cluster.resourcePackage.metadata.version;
      result += "/content";
      return result;
    },
    docsHtml() {
      if (this.markdownContent) {
        let md = new markdown();
        let temp = md.render(this.markdownContent);
        temp = temp.replaceAll("<img ", `<img onclick=showArchMarkdownImageModal(this) `);
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
    window.showArchMarkdownImageModal = (img) => {
      document.getElementById('architecture-markdown-image-modal').style.display = 'block';
      document.getElementById('architecture-markdown-image-modal-image').src = img.src;
    };
    window.hideArchMarkdownImageModal = () => {
      document.getElementById('architecture-markdown-image-modal').style.display = 'none';
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

<style lang="scss" scoped>
.architecture-markdown-image-modal {
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

.architecture-markdown-image-modal-content {
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
.architecture-markdown {
  counter-reset: h1counter h2counter h3counter h4counter h5counter;
}

.architecture-markdown {
  font-family: Consolas, Menlo, Bitstream Vera Sans Mono, Monaco, "微软雅黑", monospace;
  font-size: 14px;


  h1 {
    font-size: 18px;
    counter-reset: h2counter 1;
  }

  h1::before {
    counter-increment: h1counter;
    content: counter(h1counter) ". ";
  }

  h2 {
    font-size: 16px;
    counter-reset: h3counter;
  }

  h2::before {
    counter-increment: h2counter;
    content: counter(h1counter)"." counter(h2counter) " ";
  }

  h3 {
    font-size: 15px;
    counter-reset: h4counter;
  }

  h3::before {
    counter-increment: h3counter;
    content: counter(h1counter)"." counter(h2counter)"." counter(h3counter) " ";
  }

  h4 {
    counter-reset: h5counter;
  }

  h4::before {
    counter-increment: h4counter;
    content: counter(h1counter)"." counter(h2counter)"." counter(h3counter) "." counter(h4counter) " ";
  }

  h5::before {
    counter-increment: h5counter;
    content: counter(h1counter)"." counter(h2counter)"." counter(h3counter) "." counter(h4counter) "." counter(h5counter) " ";
  }

  p img {
    max-width: 900px;
    cursor: zoom-in;
  }

  p code {
    color: #c7254e !important;
    padding: .25rem .5rem;
    margin: 0;
    font-size: .9em;
    background-color: rgba(27, 31, 35, .05);
    border-radius: 3px;
    font-family: Consolas, Menlo, "Bitstream Vera Sans Mono", Monaco, "微软雅黑", monospace;
  }

  li code {
    color: #c7254e !important;
    padding: 1px .5rem;
    margin: 0;
    font-size: .9em;
    border-radius: 3px;
    background-color: rgba(27, 31, 35, .05);
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