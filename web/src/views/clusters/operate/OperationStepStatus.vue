<template>
  <div style="flex-grow: 1; text-align: right;">
    <el-button v-if="loading" :loading="true" type="primary" plain> Loading </el-button>
    <el-tag type="danger" effect="dark" size="defualt" v-else-if="error">
      {{ error.response.status || error.status }}
    </el-tag>
    <span v-else>
      <el-tag :type="computedStatus.success == computedStatus.total ? 'success' : 'danger'">
        {{ computedStatus.success }}
      </el-tag>
      /
      <el-tag style="margin-right: 5px;">
        {{ computedStatus.total }}
      </el-tag>
      <el-button @click="visible = true" link icon="el-icon-pointer">查看状态明细</el-button>
    </span>
    <div style="text-align: left;">
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
    </div>
  </div>
</template>

<script>
import yaml from 'js-yaml';

export default {
  props: {
    cluster: { type: Object, required: true },
    currentOperation: { type: Number, required: true },
    currentStep: { type: Number, required: true },
  },
  data() {
    return {
      status: undefined,
      visible: false,
      error: undefined,
      loading: false
    }
  },
  computed: {
    computedStatus() {
      if (this.status === undefined) {
        return { loading: true };
      }
      let total = 0;
      let success = 0;
      for (let i in this.status.data) {
        let node = this.status.data[i];
        for (let j in node) {
          let action = node[j];
          total++;
          if (action.stdout == "0") {
            success++;
          }
        }
      }
      return {
        loading: false,
        success, total
      };
    },
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
  mounted() {
    this.checkStatus()
  },
  watch: {
    currentOperation: function () {
      this.checkStatus()
    },
    currentStep: function () {
      this.checkStatus()
    }
  },
  methods: {
    async checkStatus() {
      let req = {}
      this.loading = true;
      this.status = undefined;
      this.error = undefined;
      let path = "/clusters/" + this.cluster.name + "/operation/" + this.cluster.resourcePackage.operations[this.currentOperation].name;
      path += "/step/" + this.cluster.resourcePackage.operations[this.currentOperation].steps[this.currentStep].name;
      this.kuboardSprayApi.get(path).then(resp => {
        this.status = resp.data;
        this.loading = false;
      }).catch(e => {
        console.log(e);
        this.error = e;
        this.loading = false;
      })
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