/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CreateQuestionSubmitReq } from "../models/CreateQuestionSubmitReq";
import type { CreateQuestionSubmitResp } from "../models/CreateQuestionSubmitResp";
import type { QueryQuestionSubmitResp } from "../models/QueryQuestionSubmitResp";
import type { CancelablePromise } from "../core/CancelablePromise";
import { OpenAPI } from "../core/OpenAPI";
import { request as __request } from "../core/request";

export class QuestionSubmitService {
  /**
   * create questionSubmit
   * @param authorization
   * @param body  已登陆用户才能提交代码
   * @returns CreateQuestionSubmitResp A successful response.
   * @throws ApiError
   */
  public static createQuestionSubmit(
    authorization: string,
    body: CreateQuestionSubmitReq
  ): CancelablePromise<CreateQuestionSubmitResp> {
    return __request(OpenAPI, {
      method: "POST",
      url: "/gyu_oj/v1/questionSubmit/create",
      headers: {
        "authorization": authorization
      },
      body: body
    });
  }

  /**
   * query questionSubmit List
   * @param authorization
   * @param current
   * @param pageSize
   * @param sortField
   * @param sortOrder
   * @param language
   * @param status
   * @param questionId
   * @param userId
   * @returns QueryQuestionSubmitResp A successful response.
   * @throws ApiError
   */
  public static queryQuestionSubmitList(
    authorization: string,
    current: number,
    pageSize: number,
    sortField?: string,
    sortOrder?: string,
    language?: string,
    status?: number,
    questionId?: string,
    userId?: number
  ): CancelablePromise<QueryQuestionSubmitResp> {
    return __request(OpenAPI, {
      method: "GET",
      url: "/gyu_oj/v1/questionSubmit/list",
      headers: {
        "authorization": authorization
      },
      query: {
        "current": current,
        "pageSize": pageSize,
        "sortField": sortField,
        "sortOrder": sortOrder,
        "language": language,
        "status": status,
        "questionId": questionId,
        "userId": userId
      }
    });
  }
}
