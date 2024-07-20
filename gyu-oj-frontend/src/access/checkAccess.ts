import ACCESS_ENUM from "@/access/accessEnum";

/**
 * 检查权限（判断当前登录用户是否具有某个权限）
 * @param loginUser 当前登录用户
 * @param needAccess 需要有的权限
 * @return boolean 有无权限
 */
const checkAccess = (loginUser: any, needAccess = ACCESS_ENUM.NOT_LOGIN) => {
  // 获取当前登录用户具有的权限（如果没有 loginUser，则表示未登录）
  const loginUserAccess = loginUser?.userRole ?? ACCESS_ENUM.NOT_LOGIN;
  if (needAccess === ACCESS_ENUM.NOT_LOGIN) {
    return true;
  }
  // 如果当前页面需要用户登录才能访问但用户没登录
  if (
    needAccess === ACCESS_ENUM.USER &&
    loginUserAccess === ACCESS_ENUM.NOT_LOGIN
  ) {
    return false;
  }
  // 如果需要管理员权限但用户不为管理员
  if (
    needAccess === ACCESS_ENUM.ADMIN &&
    loginUserAccess !== ACCESS_ENUM.ADMIN
  ) {
    return false;
  }
  return true;
};

export default checkAccess;
