<template>
  <div id="userLoginView">
    <a-form
      style="max-width: 480px; margin: 0 auto"
      label-align="left"
      auto-label-width
      :model="form"
      @submit="handleSubmit"
    >
      <a-form-item
        field="用户名"
        tooltip="名称长度不能少于 6 位"
        label="用户名称"
        validate-trigger="input"
        required
      >
        <a-input v-model="form.username" placeholder="请输入你的用户名称" />
      </a-form-item>
      <a-form-item
        field="密码"
        tooltip="密码长度不能少于 8 位"
        label="用户密码"
        validate-trigger="input"
        required
      >
        <a-input-password
          v-model="form.password"
          placeholder="请输入你的用户密码"
        />
      </a-form-item>
      <a-form-item field="autoLogin">
        <a-checkbox v-model="autoLoginForm.autoLogin"> 自动登录</a-checkbox>
        <div style="width: 75%; text-align: right">
          <a style="cursor: pointer; color: #165dff" @click="ToRegister"
            >没有账号？去注册</a
          >
        </div>
      </a-form-item>
      <a-form-item>
        <a-button type="primary" html-type="submit" long>登 录</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup lang="ts">
import { reactive } from "vue";
import { LoginReq, UserService } from "../../../generated";
import message from "@arco-design/web-vue/es/message";
import { useStore } from "vuex";
import { useRouter } from "vue-router";
import { setToken } from "@/plugins/token";

const store = useStore();
const router = useRouter();

const form = reactive({
  username: "",
  password: "",
} as LoginReq);

const autoLoginForm = reactive({
  autoLogin: true,
});

// 跳转到注册页面
const ToRegister = () => {
  router.push({
    path: "/user/register",
    replace: true,
  });
};

const handleSubmit = async () => {
  const res = await UserService.login(form);
  if (res.code === 200) {
    setToken(res.data.token);
    await store.dispatch("user/getLoginUser");
    router.push({
      path: "/",
      replace: true,
    });
  } else {
    message.error("登陆失败 " + res.msg);
  }
};
</script>
