<template>
  <div :class="stepClass">
    <div class="app_text_mono op_name">
      {{ step.name }}
    </div>
    <div class="op_title">
      {{ step.title[locale] }}
    </div>
  </div>
</template>

<script>
export default {
  props: {
    isCurrent: { type: Boolean, required: true },
    status: { type: String, required: true },
    step: { type: Object, required: true }
  },
  data() {
    return {}
  },
  computed: {
    stepClass() {
      if (this.isCurrent) {
        return "noselect op_step active " + this.status;
      }
      return "noselect op_step " + this.status;
    }
  },
  watch: {

  },
  methods: {
    refresh() {

    }
  }
}
</script>

<style lang="scss" scoped>
.op_step {
  padding: 10px;
  border: 2px solid var(--el-fill-color-light);
  border-radius: 5px;
  background-color: var(--el-fill-color-light);
  cursor: pointer;

  .op_name {
    font-size: 13px;
    margin-bottom: 3px;
    color: var(--el-text-color-secondary);
  }

  .op_title {
    font-size: 14px;
    font-weight: bold;
  }
}

.op_step:hover {
  background-color: var(--el-color-primary-light-6);
}

.op_step.active {
  border: 2px solid var(--el-color-primary);
}
</style>