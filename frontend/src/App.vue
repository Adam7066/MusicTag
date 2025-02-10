<template>
  <t-config-provider :global-config="globalConfig">
    <Music />
  </t-config-provider>
</template>

<script setup lang="ts">
import type { GlobalConfigProvider } from "tdesign-vue-next";
import { merge } from "lodash-es";
import enConfig from "tdesign-vue-next/es/locale/en_US";
import Music from "./components/Music.vue";

const globalConfig: GlobalConfigProvider = merge(enConfig);
</script>
