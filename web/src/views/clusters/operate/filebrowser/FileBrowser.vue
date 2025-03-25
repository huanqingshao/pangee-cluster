<i18n></i18n>

<template>
  <el-dialog v-model="visible">
    <div>
      hello
    </div>
    <div>{{ root }}</div>
    <div>{{ tree }}</div>
  </el-dialog>
</template>

<script>
export default {
  props: {
    root: { type: String, required: true }
  },
  data() {
    return {
      visible: false,
      tree: []
    }
  },
  methods: {
    show(files) {
      this.kuboardSprayApi
        .get(`/filebrowser?rootPath=${this.root}`).then(resp => {
          console.log(resp)
          this.tree = resp.data.data
        })
      this.visible = true;
    }
  }
}
</script>

<style lang="scss" scoped></style>