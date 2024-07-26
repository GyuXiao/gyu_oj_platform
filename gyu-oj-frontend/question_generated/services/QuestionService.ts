/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CreateQuestionReq } from "../models/CreateQuestionReq";
import type { DeleteQuestionReq } from "../models/DeleteQuestionReq";
import type { GetQuestionResp } from "../models/GetQuestionResp";
import type { UpdateQuestionReq } from "../models/UpdateQuestionReq";
import type { UpdateQuestionResp } from "../models/UpdateQuestionResp";
import type { CancelablePromise } from "../core/CancelablePromise";
import { OpenAPI } from "../core/OpenAPI";
import { request as __request } from "../core/request";
import {
  BaseCreateQuestionResponse,
  BaseDeleteQuestionResponse,
  BaseQueryQuestionListResponse, BaseQueryQuestionResponse, BaseUpdateQuestionResponse
} from "../models/BaseQuestionResp";

export class QuestionService {
  /**
   * admin create question
   * @param body
   * @returns CreateQuestionResp A successful response.
   * @throws ApiError
   */
  public static createQuestion(
    body: CreateQuestionReq
  ): CancelablePromise<BaseCreateQuestionResponse | any> {
    return __request(OpenAPI, {
      method: "POST",
      url: "/gyu_oj/v1/question/add",
      body: body,
    });
  }

  /**
   * admin delete question
   * @param body
   * @returns DeleteQuestionResp A successful response.
   * @throws ApiError
   */
  public static deleteQuestion(
    body: DeleteQuestionReq
  ): CancelablePromise<BaseDeleteQuestionResponse | any> {
    return __request(OpenAPI, {
      method: "POST",
      url: "/gyu_oj/v1/question/delete",
      body: body,
    });
  }

  /**
   * query question List
   * @param current
   * @param pageSize
   * @param title
   * @param tags
   * @param sortField
   * @param sortOrder
   * @returns GetQuestionListResp A successful response.
   * @throws ApiError
   */
  public static queryQuestionList(
    current: number,
    pageSize: number,
    title?: string,
    tags?: Array<string>,
    sortField?: string,
    sortOrder?: string,
  ): CancelablePromise<BaseQueryQuestionListResponse | any> {
    return __request(OpenAPI, {
      method: "GET",
      url: "/gyu_oj/v1/question/list",
      query: {
        current: current,
        pageSize: pageSize,
        sortField: sortField,
        sortOrder: sortOrder,
        title: title,
        tags: tags,
      },
    });
  }

  /**
   * query question
   * @param id
   * @returns GetQuestionResp A successful response.
   * @throws ApiError
   */
  public static queryQuestion(id: string): CancelablePromise<BaseQueryQuestionResponse | any> {
    return __request(OpenAPI, {
      method: "GET",
      url: "/gyu_oj/v1/question/query",
      query: {
        id: id,
      },
    });
  }

  /**
   * admin update question
   * @param authorization
   * @param body
   * @returns UpdateQuestionResp A successful response.
   * @throws ApiError
   */
  public static updateQuestion(
    body: UpdateQuestionReq
  ): CancelablePromise<BaseUpdateQuestionResponse | any> {
    return __request(OpenAPI, {
      method: "POST",
      url: "/gyu_oj/v1/question/update",
      body: body,
    });
  }
}