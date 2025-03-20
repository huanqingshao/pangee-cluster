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
  <div
    style="height: calc(100vh - 220px); display: flex; flex-direction: column"
  >
    <el-card style="margin-bottom: 10px">
      <el-button-group>
        <el-button>安装集群</el-button>
        <el-button>添加节点</el-button>
        <el-button>删除节点</el-button>
      </el-button-group>
      <div>参数设置</div>
    </el-card>
    <div style="display: flex; gap: 10px; flex-grow: 1">
      <el-card class="operation-card" style="width: 180px; flex-grow: 0">
        <div style="display: flex; flex-direction: column; gap: 10px">
          <div :class="stepClass('step1')" @click="currentStep = 'step1'">
            Step1
          </div>
          <div :class="stepClass('step2')" @click="currentStep = 'step2'">
            Step2
          </div>
          <div :class="stepClass('step3')" @click="currentStep = 'step3'">
            Step3
          </div>
          <div :class="stepClass('step4')" @click="currentStep = 'step4'">
            Step4
          </div>
          <div :class="stepClass('step5')" @click="currentStep = 'step5'">
            Step5
          </div>
        </div>
      </el-card>
      <el-card class="operation-card">
        Markdown 文档
        <p>{{ currentStep }}</p>
      </el-card>
      <el-card class="operation-card">
        执行历史日志
        <p>{{ currentStep }}</p>
      </el-card>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      currentOperation: "install",
      currentStep: "step1",
    };
  },
  methods: {
    stepClass(step) {
      if (this.currentStep == step) {
        return "step active";
      }
      return "step";
    },
  },
};
</script>

<style lang="scss" scoped>
.operation-card {
  flex-grow: 1;
  height: 100%;

  .step {
    padding: 10px;
    border: 1px solid var(--el-color-info-light-9);
    border-radius: 3px;
    background-color: var(--el-color-primary-light-9);
    cursor: pointer;
  }

  .step:hover {
    background-color: var(--el-color-primary-light-6);
  }

  .step.active {
    background-color: var(--el-color-primary);
    color: var(--el-color-white);
  }
}
</style>
