<template>
  <el-dialog v-model="visible" top="5vh" width="90vw">
    <div v-if="visible" class="status_table">
      <el-table :data="rows" style="width: 100%">
        <el-table-column label="任务名称" width="150">
          <template #default="scope">
            {{ scope.row._metadata.name }}
          </template>
        </el-table-column>
        <el-table-column label="任务描述" width="270">
          <template #default="scope">
            {{ scope.row._metadata.description[locale] }}
          </template>
        </el-table-column>
        <el-table-column label="查询命令" width="120">
          <template #default="scope">
            <el-popover placement="top" trigger="hover" width="75vw" effect="dark">
              <template #reference>
                <el-button link class="app_text_mono" icon="el-icon-pointer">查看</el-button>
              </template>
              <div class="shell_response">
                <div class="app_text_mono">
                  {{ scope.row._metadata.cmd }}
                </div>
              </div>
            </el-popover>
          </template>
        </el-table-column>
        <template v-for="node in columns">
          <el-table-column :label="node">
            <template #default="scope">
              <el-tag v-if="scope.row[node].stdout == '0'" type="success" effect="dark" round size="small">
                <el-icon>
                  <ElIconCheck />
                </el-icon>
              </el-tag>
              <el-popover v-else width="75vw" effect="dark">
                <template #reference>
                  <el-tag type="danger" effect="dark" round size="small">
                    <el-icon>
                      <ElIconClose />
                    </el-icon>
                  </el-tag>
                </template>
                <div class="shell_response">
                  <div class="app_text_mono">
                    {{ scope.row[node].stdout }}
                  </div>
                </div>
              </el-popover>
            </template>
          </el-table-column>
        </template>
      </el-table>
    </div>
  </el-dialog>
</template>

<script>
import yaml from 'js-yaml';

export default {
  data() {
    return {
      status: undefined,
      visible: false
    }
  },
  computed: {
    columns() {
      let columns = [];
      for (let i in this.status.data) {
        columns.push(i);
      }
      return columns;
    },
    rows() {
      let rows = {};
      for (let i in this.status.data) {
        let node = this.status.data[i];
        for (let j in node) {
          let row = rows[j] || {};
          row._metadata = yaml.load(j);
          row._metadata.cmd = node[j].cmd;
          row[i] = node[j];
          rows[j] = row;
        }
      }
      let tmp = [];
      for (let i in rows) {
        tmp.push(rows[i]);
      }
      return tmp;
    }
  },
  methods: {
    show(status) {
      this.status = status;
      this.visible = true;
    }
  }
}
</script>

<style lang="scss" scoped>
.status_table {
  width: 100%;
  height: calc(90vh - 100px);
}

.shell_response {
  display: flex;
  max-height: 60vh;
  overflow: hidden;
  overflow-y: auto;

  div {
    flex: 0 1 auto;
    overflow-x: auto;
    white-space: pre;
    word-break: keep-all;
    overflow-wrap: normal;
  }
}
</style>