import type { CurrentUserResp } from "./CurrentUserResp";
import type { LoginResp } from "./LoginResp";
import type { RegisterResp } from "./RegisterResp";

// current 请求的返回参数
export type BaseCurrentResponse = {
  code: number;
  msg: string;
  data?: CurrentUserResp;
};

export type BaseUserLoginResponse = {
  code: number;
  msg: string;
  data?: LoginResp;
};

export type BaseUserRegisterResponse = {
  code: number;
  msg: string;
  data?: RegisterResp;
};