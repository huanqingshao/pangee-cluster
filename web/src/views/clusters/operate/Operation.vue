<i18n>
en:
  singleNode: Specific Node
  global_config: Global Config
  addons: Addons
  enabledBation: Enabled
  disabledBation: Disabled
  selectANode: Please select a node from the diagram to the left.
  resourcePackage: Resource Package
zh:
  singleNode: 单个节点
  global_config: 全局设置
  addons: 可选组件
  enabledBation: 使用跳板机
  disabledBation: 不使用跳板机
  selectANode: 请从左侧图中选择一个节点
  resourcePackage: 资源包
</i18n>

<template>
  <div style="height: calc(100% - 2px); display: flex; flex-direction: column">
    <div style="margin-bottom: 10px">
      <el-radio-group v-model="currentOperation" size="default">
        <template v-for="(operation, index) in cluster.resourcePackage.data.operations" :key="'op_' + index">
          <el-radio-button :label="operation.title[locale]" :value="index"></el-radio-button>
        </template>
      </el-radio-group>
      <div class="operation-params">
        参数设置
      </div>
    </div>
    <div style="display: flex; gap: 10px; flex: 1; overflow-y: hidden;">
      <div class="operation-card" style="min-width: 240px; flex-grow: 0;">
        <div class="noselect operation-steps">
          <div>
            <el-steps direction="vertical" :active="currentStep">
              <template v-for="(step, index) in cluster.resourcePackage.data.operations[currentOperation].steps"
                :key="'step_' + index">
                <el-step :title="step.name" :description="step.title[locale]" @click="currentStep = index"
                  :status="currentStep >= index || true ? 'success' : ''"
                  :class="currentStep == index ? 'is-selected-step' : ''" />
              </template>
            </el-steps>
          </div>
        </div>
      </div>
      <div class="operation-card" style="max-width: 50%; min-width: 50%;" :body-style="{ height: 'calc(100% - 40px)' }">
        <div class="markdown-title">
          <div style="flex-grow: 1; font-weight: bolder;">{{
            cluster.resourcePackage.data.operations[currentOperation].steps[currentStep].title[locale] }}</div>
          <el-button style="float: right" @click="showFileBrowser" type="primary"
            icon="el-icon-pointer">查看代码</el-button>
        </div>
        <div class="markdown-content">
          <OperationStepMarkdown :cluster="cluster" :operation-index="currentOperation" :step-index="currentStep">
          </OperationStepMarkdown>
        </div>
      </div>
      <div class="operation-card">
        <div class="markdown-title">
          <OperationStepExecute :cluster="cluster" :currentOperation="currentOperation" :currentStep="currentStep"
            @refresh="$emit('refresh')">
          </OperationStepExecute>
          <OperationStepStatus ref="stepStatus" :cluster="cluster" :currentOperation="currentOperation"
            :currentStep="currentStep" @refresh="$refs.history.refresh()"></OperationStepStatus>
        </div>
        <div class="operation-history">
          <OperationStepHistory ref="history" :cluster="cluster" :currentOperation="currentOperation"
            :currentStep="currentStep">
          </OperationStepHistory>
        </div>
      </div>
    </div>
    <FileBrowser :package-name="`${cluster.resourcePackage.metadata.version}`" ref="filebrowser"></FileBrowser>
  </div>
</template>

<script>
import FileBrowser from "./filebrowser/FileBrowser.vue";
import OperationStepMarkdown from "./OperationStepMarkdown.vue"
import OperationStepExecute from "./OperationStepExecute.vue";
import OperationStepHistory from "./OperationStepHistory.vue";
import OperationStepStatus from "./OperationStepStatus.vue";

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
    };
  },
  computed: {
    currentOperation: {
      get() {
        if (this.$store.state.cluster[this.cluster.name] == undefined) {
          return 0;
        }
        if (this.$store.state.cluster[this.cluster.name].operation == undefined) {
          return 0;
        }
        return this.$store.state.cluster[this.cluster.name].operation.currentOperation || 0
      },
      set(v) {
        this.$store.commit("cluster/CHANGE_CLUSTER_STATE",
          {
            cluster: this.cluster.name,
            tab: "operation",
            key: "currentOperation",
            value: v
          }
        );
      }
    },
    currentStep: {
      get() {
        if (this.$store.state.cluster[this.cluster.name] == undefined) {
          return 0;
        }
        if (this.$store.state.cluster[this.cluster.name].operation == undefined) {
          return 0;
        }
        return this.$store.state.cluster[this.cluster.name].operation.currentStep || 0
      },
      set(v) {
        this.$store.commit("cluster/CHANGE_CLUSTER_STATE",
          {
            cluster: this.cluster.name,
            tab: "operation",
            key: "currentStep",
            value: v
          }
        );
      }
    }
  },
  components: {
    FileBrowser,
    OperationStepExecute,
    OperationStepMarkdown,
    OperationStepHistory,
    OperationStepStatus
  },
  methods: {
    stepClass(step) {
      if (this.currentStep == step) {
        return "step active";
      }
      return "step";
    },
    showFileBrowser() {
      let path = "/operations/" + this.cluster.resourcePackage.data.operations[this.currentOperation].name;
      path += "/" + this.cluster.resourcePackage.data.operations[this.currentOperation].steps[this.currentStep].name;
      this.$refs.filebrowser.show([{
        isDir: true,
        path: path
      }]);
    }
  },
};
</script>

<style lang="scss" scoped>
.operation-params {
  padding: 10px;
  border: 1px solid var(--el-color-info-light-9);
  border-radius: 5px;
  border-top-left-radius: 0;
  background-color: var(--el-color-primary-light-9);
}

.operation-card {
  flex-grow: 1;
  height: calc(100% - 42px);
  overflow: hidden;
  padding: 20px;
  border: 1px solid var(--el-border-color);
  border-radius: 8px;

  .operation-steps {
    height: 100%;
    overflow: hidden;
    overflow-y: auto;
  }

  .markdown-title {
    height: 36px;
    margin-bottom: 10px;
    background-color: var(--el-color-primary-light-9);
    display: flex;
    align-items: center;
    justify-content: stretch;
    padding: 0 10px;
  }

  .markdown-content {
    height: calc(100% - 46px);
    overflow: hidden;
    overflow-y: auto;
  }

  .operation-history {
    height: calc(100% - 46px);
    overflow: hidden;
    overflow-y: auto;
  }
}
</style>

<style lang="scss">
.operation-card .el-step__main {
  margin-bottom: 20px;
  cursor: pointer;

  .el-step__title {
    font-size: 13px;
    font-family: Consolas, Menlo, "Bitstream Vera Sans Mono", Monaco, "微软雅黑", monospace;
    font-weight: normal;
  }

  .el-step__description {
    border-bottom: 2px solid var(--el-color-primary-light-7);
    font-size: 14px;
    font-weight: bolder;
  }

}

.operation-card .is-selected-step {
  .el-step__icon {
    border: 2px solid var(--el-color-primary);

    svg {
      color: var(--el-color-primary);
    }
  }

  .el-step__description {
    border-bottom: 3px solid var(--el-color-primary);
  }
}
</style>