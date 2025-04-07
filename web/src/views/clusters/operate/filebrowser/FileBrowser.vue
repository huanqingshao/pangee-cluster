<i18n></i18n>

<template>
  <el-dialog v-model="visible" top="5vh" width="80%">
    <div style="height: calc(90vh - 92px)">
      <div>
      </div>
      <div class="filebrowser">
        <div class="left" :style="maxifyTreeWidth ? 'width: 50%' : 'width: 20%'">
          <div class="collapse">
            <el-button style="margin-top: 5px;" v-if="maxifyTreeWidth" icon="el-icon-d-arrow-left"
              @click="maxifyTreeWidth = false" circle type="primary" plain></el-button>
            <el-button v-else style="margin-top: 5px;" icon="el-icon-d-arrow-right" @click="maxifyTreeWidth = true"
              circle type="primary" plain></el-button>
          </div>
          <div class="tree">
            <el-tree ref="tree" style="max-width: 600px" node-key="path" :props="treeProps" :load="loadNode"
              :current-node-key="currentNode?.data.path" highlight-current @node-click="clickNode" lazy>
              <template #default="{ node, data }">
                <span class="app_text_mono">{{ data.name }}</span>
              </template>
            </el-tree>
          </div>
        </div>
        <div class="right" :style="maxifyTreeWidth ? 'width: 50%' : 'width: 80%'">
          <div class="c-bar app_text_mono">
            <template v-for="(item, idx) in items">
              <el-tag class="path_item noselect" @click="clickNode(item.data, item)">{{ item.data.name }}</el-tag>
              <span v-if="idx < items.length - 1">></span>
            </template>
          </div>
          <div class="editor">
            <pre v-if="contentType == 'text'"><code>{{ content }}</code></pre>
            <pre v-else-if="contentType == 'image'"><img :src="content" onclick=showMarkdownImageModal(this) /></pre>
            <pre v-else-if="contentType == ''">请先选择一个文件</pre>
            <pre v-else>不能正常显示内容类型为 {{ contentType }} 的文件</pre>
          </div>

        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script>
import axios from 'axios'

function nodeToPath(node) {
  let filePath = (node.data.path || "");
  if (node.level > 0 && node.parent.path) {
    console.log(node.parent)
    filePath = node.parent.path + "/" + node.data.name
    console.log(filePath)
  }
  return filePath;
}

export default {
  props: {
    packageName: { type: String, required: true },
  },
  data() {
    return {
      visible: false,
      maxifyTreeWidth: false,
      content: "",
      contentType: "",
      currentNode: undefined,
      treeProps: {
        label: "name",
        children: "children",
        isLeaf: (item) => {
          return !item.isDir;
        }
      },
      files: []
    }
  },
  computed: {
    items() {
      let temp = [];
      let t = this.currentNode;
      if (this.currentNode) {
        temp.push(this.currentNode);
        while (t.parent) {
          t = t.parent;
          if (!(t.data instanceof Array))
            temp.unshift(t);
        }
      }
      return temp;
    }
  },
  methods: {
    show(files) {
      this.files = files;
      this.visible = true;
    },
    loadNode(node, resolve, reject) {
      console.log("loadNode: ", node, nodeToPath(node));
      if (this.files && this.files.length > 0 && node.level == 0) {
        let items = [];
        for (let i in this.files) {
          let splited = this.files[i].path.split("/");
          this.files[i].name = splited[splited.length - 1];
          items.push(this.files[i]);
        }
        resolve(items);
      } else {
        this.kuboardSprayApi
          .get(`/filebrowser?rootPath=data:///resource/${this.packageName + "/content" + nodeToPath(node)}`).then(resp => {
            let items = resp.data.data;
            let path = nodeToPath(node);
            for (let item of items) {
              item.path = path + "/" + item.name;
            }
            resolve(items);
          })
      }
    },
    clickNode(file, node) {
      this.currentNode = node;
      if (file.isDir) {
        this.contentType = "";
        this.content = "";
        return;
      }
      let path = `/resource-package/${this.packageName}/content/${nodeToPath(node)}`;
      axios.head(path).then(resp => {
        if (resp.headers['content-type']?.indexOf("text") >= 0) {
          this.contentType = "text";
          axios.get(path).then(respTxt => {
            this.content = respTxt.data;
          })
        } else if (resp.headers['content-type']?.indexOf("image") >= 0) {
          this.contentType = "image";
          this.content = path;
        } else {
          this.contentType = resp.headers['content-type'];
        }
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.filebrowser {
  height: 100%;
  display: flex;
  gap: 10px;

  .left {
    height: 100%;

    .collapse {
      height: 0px;
      text-align: right;
      overflow: visible;
      padding-right: 5px;
      position: relative;
      z-index: 800;
    }

    .tree {
      overflow: hidden;
      overflow-y: auto;
      border: 1px solid var(--el-border-color);
      border-radius: 8px;
      width: 100%;
      height: 100%;
    }
  }

  .right {
    display: flex;
    flex-grow: 1;
    flex-direction: column;
    gap: 10px;

    .c-bar {
      border: 1px solid var(--el-border-color);
      border-radius: 8px;
      height: 32px;
      width: 100%;
      display: flex;
      gap: 5px;
      align-items: center;
      padding: 0 10px;

      .path_item {
        cursor: pointer;
      }
    }

    .editor {
      border: 1px solid var(--el-border-color);
      border-radius: 8px;
      height: calc(100% - 42px);
      width: 100%;
      overflow-x: auto;
      overflow-y: auto;

      img {
        max-width: 100%;
        max-height: 100%;
      }
    }
  }
}
</style>