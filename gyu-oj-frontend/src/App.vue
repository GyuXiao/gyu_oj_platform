<template>
  <div id="app">
    <template v-if="route.path.startsWith('/user')">
      <router-view />
    </template>
    <template v-else>
      <BasicLayout />
    </template>
  </div>
</template>

<style>
#app {
}
</style>
<script setup lang="ts">
import BasicLayout from "@/layouts/BasicLayout.vue";
import { useRoute, useRouter } from "vue-router";
import { useStore } from "vuex";
import { onMounted } from "vue";
import accessEnum from "@/access/accessEnum";

const router = useRouter();
const store = useStore();
const route = useRoute();

/**
 * 全局初始化函数，有全局单次调用的代码，都可以写到这里
 */
const doInit = () => {
  console.log("hello OJ");
};

onMounted(() => {
  doInit();
});

router.beforeEach((to, from, next) => {
  if (to.meta?.access == "canAdmin") {
    if (store.state.user.loginUser?.userRole !== accessEnum.ADMIN) {
      next("/noAuth");
      return;
    }
  }
  next();
});
</script>
