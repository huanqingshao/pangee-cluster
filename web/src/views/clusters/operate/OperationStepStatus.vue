<template>
  <div style="flex-grow: 1; text-align: right;">
    <el-button>xx</el-button>
    {{ computedStatus }}
  </div>
</template>

<script>

export default {
  props: {
    cluster: { type: Object, required: true },
    currentOperation: { type: Number, required: true },
    currentStep: { type: Number, required: true },
  },
  data() {
    return {
      status: undefined
    }
  },
  computed: {
    computedStatus() {
      if (this.status === undefined) {
        return "Loading";
      }
      return "";
    }

  },
  methods: {
    async checkStatus() {
      let req = {}
      let path = "/clusters/" + this.cluster.name + "/operation/" + this.cluster.resourcePackage.operations[this.currentOperation].name;
      path += "/step/" + this.cluster.resourcePackage.operations[this.currentOperation].steps[this.currentStep].name;
      this.kuboardSprayApi.get(path).then(resp => {
        this.status = resp.data
      })
    }
  }
}
</script>