/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { LoginReq } from "../models/LoginReq";
import type { RegisterReq } from "../models/RegisterReq";
import type { CancelablePromise } from "../core/CancelablePromise";
import { OpenAPI } from "../core/OpenAPI";
import { request as __request } from "../core/request";
import type {
  BaseCurrentResponse,
  BaseUserLoginResponse,
  BaseUserLogoutResponse,
  BaseUserRegisterResponse,
} from "../models/BaseResponse";

export class UserService {
  /**
   * get current user
   * @returns CurrentUserResp A successful response.
   * @throws ApiError
   */
  public static current(): CancelablePromise<BaseCurrentResponse | any> {
    return __request(OpenAPI, {
      method: "GET",
      url: "/gyu_oj/v1/user/current",
    });
  }

  /**
   * userLogin
   * @param body
   * @returns LoginResp A successful response.
   * @throws ApiError
   */
  public static login(
    body: LoginReq
  ): CancelablePromise<BaseUserLoginResponse | any> {
    return __request(OpenAPI, {
      method: "POST",
      url: "/gyu_oj/v1/user/login",
      body: body,
    });
  }

  /**
   * userLogout
   * @param authorization
   * @param body
   * @returns LogoutResp A successful response.
   * @throws ApiError
   */
  public static logout(): CancelablePromise<BaseUserLogoutResponse | any> {
    return __request(OpenAPI, {
      method: "POST",
      url: "/gyu_oj/v1/user/logout",
    });
  }

  /**
   * userRegister
   * @param body
   * @returns RegisterResp A successful response.
   * @throws ApiError
   */
  public static register(
    body: RegisterReq
  ): CancelablePromise<BaseUserRegisterResponse | any> {
    return __request(OpenAPI, {
      method: "POST",
      url: "/gyu_oj/v1/user/register",
      body: body,
    });
  }
}
