<template>
  <div id="userRegisterView">
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
      <a-form-item
        field="确认密码"
        tooltip="密码长度不能少于 8 位"
        label="密码确认"
        validate-trigger="input"
        required
      >
        <a-input-password
          v-model="form.confirm_password"
          placeholder="请再次输入你的用户密码"
        />
      </a-form-item>
      <a-form-item>
        <a-button type="primary" html-type="submit" long>注 册</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup lang="ts">
import { reactive } from "vue";
import { RegisterReq, UserService } from "../../../generated";
import message from "@arco-design/web-vue/es/message";
import { useRouter } from "vue-router";
import { Notification } from "@arco-design/web-vue";

const router = useRouter();

const RegisterSuccessNotification = () => {
  Notification.info({
    title: "注册成功",
    content: "再去登陆一下，就能开启 OJ 之旅了~",
    duration: 5000,
  });
};

const form = reactive({
  username: "",
  password: "",
  confirm_password: "",
} as RegisterReq);

const handleSubmit = async () => {
  const res = await UserService.register(form);
  if (res.code === 200) {
    RegisterSuccessNotification();
    router.push({
      path: "/user/login",
      replace: true,
    });
  } else {
    message.error("登陆失败" + res.msg);
  }
};
</script>
