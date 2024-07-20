/**
 * 本地存储 JWT token
 * @param token
 */
export function setToken(token: string) {
  localStorage.setItem("token", token);
}

/**
 * 获取本地存储的 JWT token
 */
export function getToken(): string {
  const jwtToken = localStorage.getItem("token");
  if (jwtToken) {
    return jwtToken as string;
  }
  return "";
}
