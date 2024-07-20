import { StoreOptions } from "vuex";
import ACCESS_ENUM from "@/access/accessEnum";
import { UserService } from "../../generated";
import accessEnum from "@/access/accessEnum";
import message from "@arco-design/web-vue/es/message";

/**
 * 用户模块
 */
export default {
  namespaced: true,
  // 定义已登陆用户信息
  state: () => ({
    loginUser: {
      id: 0,
      username: "未登录",
      avatarUrl: "",
      userRole: accessEnum.NOT_LOGIN,
      token: "",
      tokenExpire: 0,
    },
  }),
  // 定义远程获取用户信息的方法
  actions: {
    async getLoginUser({ commit, state }, payload) {
      const res = await UserService.current();
      console.log(res);
      if (res.code === 200) {
        commit("updateUser", res.data);
      } else {
        commit("updateUser", {
          ...state.loginUser,
          userRole: ACCESS_ENUM.NOT_LOGIN,
        });
        message.error("获取当前登陆用户错误: " + res.msg);
      }
    },
  },
  // 定义更新用户的方法
  mutations: {
    updateUser(state, payload) {
      state.loginUser = payload;
    },
  },
} as StoreOptions<any>;
