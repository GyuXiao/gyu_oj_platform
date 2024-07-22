import router from "@/router";
import store from "@/store";
import ACCESS_ENUM from "@/access/accessEnum";
import checkAccess from "@/access/checkAccess";
import accessEnum from "@/access/accessEnum";

router.beforeEach(async (to, from, next) => {
  console.log("登陆用户的信息", store.state.user.loginUser);
  let loginUser = store.state.user.loginUser;
  // todo：后续还要再优化这里的逻辑
  const userStr = localStorage.getItem("loginUser");
  if (userStr) {
    const user = JSON.parse(userStr);
    // 1,如果之前登陆过，则在页面刷新时可以请求 getLoginUser 刷新登陆
    if (user.userRole !== accessEnum.NOT_LOGIN) {
      await store.dispatch("user/getLoginUser");
      loginUser = store.state.user.loginUser;
    }
  }
  // if (loginUser.userRole !== accessEnum.NOT_LOGIN) {
  //   // 加 await 是为了等用户登录成功之后，再执行后续的代码
  //   await store.dispatch("user/getLoginUser");
  //   loginUser = store.state.user.loginUser;
  // }
  // 获取页面所需要的权限
  const needAccess = (to.meta?.access as number) ?? ACCESS_ENUM.NOT_LOGIN;
  // 2,如果待跳转页面要求用户必须登陆
  if (needAccess !== ACCESS_ENUM.NOT_LOGIN) {
    // 如果用户未登陆，则跳转到登录页面
    if (
      !loginUser ||
      !loginUser.userRole ||
      loginUser.userRole === accessEnum.NOT_LOGIN
    ) {
      next(`/user/login?redirect=${to.fullPath}`);
      return;
    }
    // 如果用户已经登陆了，但是权限不足，则跳转到无权限页面
    if (!checkAccess(loginUser, needAccess)) {
      next("/noAuth");
      return;
    }
  }
  next();
});
