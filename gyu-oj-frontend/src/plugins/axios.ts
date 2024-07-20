import axios from "axios";
import { getToken } from "@/plugins/token";

axios.interceptors.request.use(
  function (config) {
    // 如果 JWT token 存在，发出请求前需要先在 Header 带上 {Authorization: Bearer token}
    const jwtToken = getToken();
    if (jwtToken) {
      config.headers.Authorization = "Bearer " + getToken();
    }
    return config;
  },
  function (error) {
    return Promise.reject(error);
  }
);

axios.interceptors.response.use(
  function (response) {
    return response;
  },
  function (error) {
    return Promise.reject(error);
  }
);
