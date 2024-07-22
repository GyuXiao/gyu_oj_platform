<template>
  <a-row id="globalHeader" align="center" :wrap="false">
    <a-col flex="auto">
      <a-menu
        mode="horizontal"
        :selected-keys="selectedKey"
        @menu-item-click="doMenuClick"
      >
        <a-menu-item
          key="0"
          :style="{ padding: 0, marginRight: '38px' }"
          disabled
        >
          <div class="title-bar">
            <img class="logo" src="../assets/oj_logo.jpg" />
            <div class="title">TopCoder</div>
          </div>
        </a-menu-item>
        <a-menu-item v-for="item in visibleRoutes" :key="item.path">
          {{ item.name }}
        </a-menu-item>
      </a-menu>
    </a-col>
    <a-col flex="100px">
      <div>
        <template
          v-if="
            store.state.user &&
            store.state.user.loginUser &&
            store.state.user.loginUser.username &&
            store.state.user.loginUser.username !== '未登录'
          "
        >
          <a-dropdown @select="handleSelect">
            <a-button type="text"
              >{{ store.state.user.loginUser.username }}
            </a-button>
            <template #content>
              <a-doption :value="{ value: 'logout' }">退出登录</a-doption>
            </template>
          </a-dropdown>
        </template>
        <template v-else>
          <a-button type="text" @click="toLoginPage">未登录</a-button>
        </template>
      </div>
    </a-col>
  </a-row>
</template>

<script setup lang="ts">
import { routes } from "@/router/routers";
import { useRouter } from "vue-router";
import { computed, ref } from "vue";
import { useStore } from "vuex";
import checkAccess from "@/access/checkAccess";
import { UserService } from "../../generated";
import accessEnum from "@/access/accessEnum";

const router = useRouter();
const store = useStore();

const handleSelect = (v) => {
  console.log(v);
  // 退出登录
  if (v.value === "logout") {
    UserService.logout();
    localStorage.removeItem("loginUser");
    toLoginPage();
  }
};

/**
 * 跳转到登录界面
 */
const toLoginPage = () => {
  router.push({
    path: `/user/login`,
  });
};

// 默认主页
const selectedKey = ref(["/"]);

router.afterEach((to, from, failure) => {
  selectedKey.value = [to.path];
});

const doMenuClick = (key: string) => {
  router.push({
    path: key,
  });
};

// 展示在菜单的路由数组
const visibleRoutes = computed(() => {
  return routes.filter((item, index) => {
    if (item.meta?.hideInMenu) {
      return false;
    }
    // 根据权限过滤菜单项的显隐
    if (
      !checkAccess(store.state.user.loginUser, item?.meta?.access as string)
    ) {
      return false;
    }
    return true;
  });
});
</script>

<style scoped>
.title-bar {
  display: flex;
  align-items: center;
}

.title {
  color: hotpink;
  margin-left: 8px;
}

.logo {
  height: 48px;
}
</style>
