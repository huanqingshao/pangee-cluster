<template>
  <div>
    <component :is="RemoteConfigComponent" :cluster="cluster"></component>
  </div>
</template>

<script>
import { loadModule } from 'vue3-sfc-loader';
import * as Vue from 'vue';
import { defineAsyncComponent } from "vue";
import axios from "axios";
import i18n from "@/i18n/index.js";
import yaml from "js-yaml";

import compareVersions from 'compare-versions'

export default {
  props: {
    cluster: { type: Object, required: true }
  },
  data() {
    return {

    }
  },
  setup(props) {

    window.compareVersions = compareVersions;

    let options = {
      moduleCache: { vue: Vue },
      async getFile(url) {
        const res = await axios.get(url);
        return res.data;
      },
      addStyle(textContent) {
        const style = Object.assign(document.createElement('style'), { textContent });
        const ref = document.head.getElementsByTagName('style')[0] || null;
        document.head.insertBefore(style, ref);
      },
      customBlockHandler(block, filename, options) {

        if (block.type !== 'i18n')
          return

        const messages = yaml.load(block.content);
        for (let locale in messages)
          i18n.global.mergeLocaleMessage(locale, messages[locale]);
      }
    }
    let url = `/resource-package/${props.cluster.resourcePackage.metadata.version}/content/vue/index.vue`;
    let RemoteConfigComponent = defineAsyncComponent(() =>
      loadModule(url, options)
    );
    console.log("加载远程组件 " + url);
    return { RemoteConfigComponent };
  }
}



</script>