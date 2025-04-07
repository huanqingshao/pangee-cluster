<i18n>
en:
  noHistory: No History
zh:
  noHistory: 暂无历史操作
</i18n>

<template>
  <div v-if="history.length > 0">
    <template v-for="(item, idx) in history" :key="'history' + idx">
      <div style="line-height: 28px">
        <el-link>{{ formatTime(item.time) }}</el-link>
      </div>
    </template>
  </div>
  <div v-else>
    <div style="line-height: 28px">
      <el-link type="danger" disabled>{{ t("noHistory") }}</el-link>
    </div>
  </div>
</template>


<script>
import dayjs from 'dayjs';

export default {
  props: {
    cluster: { type: Object, required: true },
    currentOperation: { type: Number, required: true },
    currentStep: { type: Number, required: true },
  },
  data() {
    return {
      history: []
    }
  },
  computed: {
    apiPath() {
      let path = "/clusters/" + this.cluster.name;
      if (this.cluster.resourcePackage.operations[this.currentOperation] == undefined) {
        return "";
      }
      path += "/history/" + this.cluster.resourcePackage.operations[this.currentOperation]?.name;

      if (this.cluster.resourcePackage.operations[this.currentOperation]?.steps[this.currentStep] == undefined) {
        return "";
      }
      path += "/step/" + this.cluster.resourcePackage.operations[this.currentOperation]?.steps[this.currentStep]?.name;
      return path;
    }
  },
  components: {},
  mounted() {
    this.refresh()
  },
  watch: {
    apiPath() {
      this.refresh()
    }
  },
  methods: {
    formatTime(time) {
      let temp = time.padEnd(23, "0")
      return dayjs(temp, "YYYY-MM-DD_HH-mm-ss.SSS").format("YYYY-MM-DD HH:mm:ss")
    },
    /**
     * 这是一个异步刷新函数。如果存在apiPath，则通过kuboardSprayApi获取数据，并将返回的数据赋值给history变量，同时在控制台打印出返回的数据。
     *
     * @async
     * @function refresh
     * @return {Promise}  无返回值，但会改变this.history的值，并在控制台打印出返回的数据。
     */
    async refresh() {
      if (this.apiPath) {
        this.kuboardSprayApi
          .get(this.apiPath).then(resp => {
            this.history = resp.data.data.history
          })
      }
    }
  }
}
</script>