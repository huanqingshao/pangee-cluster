<template>
  <div class="bread-crumb">
    <!-- <div class="kind">{{kind}}</div> -->
    <div class="wrapper noselect" :style="hasNext ? '' : 'cursor: auto'">
      <router-link v-if="to" class="label" :to="to" style="cursor: pointer;">{{label}}</router-link>
      <span v-else :class="hasNext ? 'label' : 'label isLeaf'"
        class="app_text_mono" @click="$emit('click')">{{label}}</span>
      <slot></slot>
      <el-icon v-if="hasNext" class="icon" :size="24">
        <!-- <el-icon-arrow-right></el-icon-arrow-right> -->
        /
      </el-icon>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    kind: { type: String, required: false },
    label: { type: String, required: false },
    to: { type: String, required: false, default: undefined },
    hasNext: { type: Boolean, required: false, default: true },
  },
  data() {
    return {

    }
  },
  computed: {
  },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped lang="css">
.bread-crumb {
  height: 60px;
  line-height: 60px;
  display: inline-block;
}
.wrapper {
  height: 60px;
  line-height: 60px;
}
.kind {
  opacity: 0.6;
  font-size: 13px;
  line-height: 22px;
  color: var(--el-text-color-placeholder);
  transition: 0.5s;
}
.bread-crumb:hover .kind {
  opacity: 1;
}
.label {
  color: var(--text-3);
  font-weight: bold;
  font-size: 14px;
  /* border-bottom: dashed 1px var(--el-color-primary); */
}
.isLeaf {
  color: var(--text-1);
  font-weight: bold;
  font-size: 14px;
  border-bottom: none;
}
.icon {
  margin: 0 5px;
  color: var(--el-color-primary-light-5);
  /* font-weight: 800; */
  vertical-align: middle;
}
</style>
