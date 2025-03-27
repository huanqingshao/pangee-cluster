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
        <template v-for="(operation, index) in cluster.resourcePackage.operations" :key="'op_' + index">
          <el-radio-button :label="operation.title[locale]" :value="index"></el-radio-button>
        </template>
      </el-radio-group>
      <div class="operation-params">
        参数设置
      </div>
    </div>
    <div style="display: flex; gap: 10px; flex: 1; overflow-y: hidden;">
      <el-card shadow="never" class="operation-card" style="width: 240px; flex-grow: 0;">
        <div style="display: flex; flex-direction: column; gap: 10px">
          <template v-for="(step, index) in cluster.resourcePackage.operations[currentOperation].steps"
            :key="'step_' + index">
            <OperationStep :is-current="currentStep == index" :step="step" status="completed"
              @click="currentStep = index"></OperationStep>
          </template>
        </div>
      </el-card>
      <el-card shadow="never" class="operation-card" style="max-width: 50%; min-width: 50%;"
        :body-style="{ height: 'calc(100% - 40px)' }">
        <div style="height: 0px; text-align: right;">
          <el-button @click="showFileBrowser" type="primary" icon="el-icon-pointer">查看代码</el-button>
        </div>
        <div style="height: 100%; overflow: hidden; overflow-y: auto;">
          <OperationStepMarkdown :cluster="cluster" :operation-index="currentOperation" :step-index="currentStep">
          </OperationStepMarkdown>
        </div>
      </el-card>
      <el-card shadow="never" class="operation-card">
        执行历史日志
        <p>{{ currentStep }}</p>
      </el-card>
    </div>
    <FileBrowser :root="`data:///resource/${cluster.resourcePackage.metadata.version}`" ref="filebrowser"></FileBrowser>
  </div>
</template>

<script>
import FileBrowser from "./filebrowser/FileBrowser.vue";
import OperationStep from "./OperationStep.vue";
import OperationStepMarkdown from "./OperationStepMarkdown.vue"

export default {
  props: {
    cluster: { type: Object, required: true },
  },
  data() {
    return {
      currentOperation: 0,
      currentStep: 0,
    };
  },
  components: { FileBrowser, OperationStep, OperationStepMarkdown },
  methods: {
    stepClass(step) {
      if (this.currentStep == step) {
        return "step active";
      }
      return "step";
    },
    showFileBrowser() {
      this.$refs.filebrowser.show();
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
  height: calc(100% - 2px);
  overflow: hidden;
}
</style>
