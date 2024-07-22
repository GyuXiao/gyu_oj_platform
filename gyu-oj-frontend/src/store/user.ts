import { StoreOptions } from "vuex";
import ACCESS_ENUM from "@/access/accessEnum";
import { UserService } from "../../generated";
import accessEnum from "@/access/accessEnum";

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
      userRole: accessEnum.NOT_LOGIN,
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
        // message.error("获取当前登陆用户错误: " + res.msg);
      }
    },
  },
  // 定义更新用户的方法
  mutations: {
    updateUser(state, payload) {
      state.loginUser = payload;
      console.log("此时的 state.loginUser", state.loginUser);
      localStorage.setItem("loginUser", JSON.stringify(payload));
    },
  },
} as StoreOptions<any>;
