<template>
  <EditCommon v-model="modelValue" :prop="prop" :rules="rules" :required="required" :label="label"
    :placeholder="placeholder">
    <template #edit>
      <el-switch v-model="modelValue" :disabled="disabled"></el-switch>
      <div class="desc">
        <slot name="edit_desc"></slot>
      </div>
    </template>
    <template #view>
      <el-switch v-model="modelValue" disabled></el-switch>
      <div class="desc">
        <slot name="view_desc"></slot>
      </div>
    </template>
  </EditCommon>
</template>

<script lang="ts" setup>
import EditCommon from "./EditCommon.vue"

const props = withDefaults(defineProps<{
  prop?: string;
  required?: boolean;
  rules?: any[];
  label?: string;
  placeholder?: string;
  antiFreeze?: boolean;
  readOnly?: boolean;
  disabled?: boolean;
}>(), {
  required: false,
  readOnly: false,
  disabled: false,
})

const modelValue = defineModel<boolean>();

</script>

<style scoped lang="css">
.desc {
  color: var(--el-text-color-secondary);
  display: inline-block;
  vertical-align: top;
  margin-left: 10px;
}
</style>
