<template>
  <EditCommon v-model="modelValue" :prop="prop" :rules="rules" :required="required" :label="label"
    :placeholder="placeholder" :helpString="helpString" :helpLink="helpLink">
    <template #edit="scope">
      <el-input v-if="showPassword" v-model="modelValue" show-password :disabled="disabled" :clearable="clearable"
        :placeholder="scope.placeholder"></el-input>
      <el-input v-else v-model.trim="modelValue" :disabled="disabled" :clearable="clearable"
        :placeholder="scope.placeholder"></el-input>
    </template>
    <template #view>
      <span v-if="modelValue">
        {{ showPassword && hidePassword ? modelValue.replace(/.?/g, '*') : modelValue }}
      </span>
      <el-button style="margin-left: 5px;" v-if="showPassword" type="primary" text icon="el-icon-view"
        @click="hidePassword = !hidePassword"></el-button>
    </template>
  </EditCommon>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import EditCommon from "./EditCommon.vue"

const props = withDefaults(defineProps<{
  prop?: string;
  required?: boolean;
  rules?: any[];
  label?: string;
  placeholder?: string;
  antiFreeze?: boolean;
  readOnly?: boolean;
  labelWidth?: string;
  helpString?: string;
  helpLink?: string;
  showPassword?: boolean;
  disabled?: boolean;
  clearable?: boolean;
}>(), {
  required: false,
  readOnly: false,
  showPassword: false,
  disabled: false,
  clearable: false
})


const hidePassword = ref(true)

const modelValue = defineModel<string>();


</script>

<style scoped lang="css"></style>
