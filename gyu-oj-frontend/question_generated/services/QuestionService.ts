/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CreateQuestionReq } from "../models/CreateQuestionReq";
import type { CreateQuestionResp } from "../models/CreateQuestionResp";
import type { DeleteQuestionReq } from "../models/DeleteQuestionReq";
import type { DeleteQuestionResp } from "../models/DeleteQuestionResp";
import type { GetQuestionListResp } from "../models/GetQuestionListResp";
import type { GetQuestionResp } from "../models/GetQuestionResp";
import type { UpdateQuestionReq } from "../models/UpdateQuestionReq";
import type { UpdateQuestionResp } from "../models/UpdateQuestionResp";
import type { CancelablePromise } from "../core/CancelablePromise";
import { OpenAPI } from "../core/OpenAPI";
import { request as __request } from "../core/request";

export class QuestionService {
  /**
   * admin create question
   * @param authorization
   * @param body
   * @returns CreateQuestionResp A successful response.
   * @throws ApiError
   */
  public static createQuestion(
    authorization: string,
    body: CreateQuestionReq
  ): CancelablePromise<CreateQuestionResp> {
    return __request(OpenAPI, {
      method: "POST",
      url: "/gyu_oj/v1/question/add",
      headers: {
        authorization: authorization,
      },
      body: body,
    });
  }

  /**
   * admin delete question
   * @param authorization
   * @param body
   * @returns DeleteQuestionResp A successful response.
   * @throws ApiError
   */
  public static deleteQuestion(
    authorization: string,
    body: DeleteQuestionReq
  ): CancelablePromise<DeleteQuestionResp> {
    return __request(OpenAPI, {
      method: "POST",
      url: "/gyu_oj/v1/question/delete",
      headers: {
        authorization: authorization,
      },
      body: body,
    });
  }

  /**
   * query question List
   * @param current
   * @param pageSize
   * @param sortField
   * @param sortOrder
   * @param title
   * @param tags
   * @returns GetQuestionListResp A successful response.
   * @throws ApiError
   */
  public static queryQuestionList(
    current: number,
    pageSize: number,
    sortField?: string,
    sortOrder?: string,
    title?: string,
    tags?: string
  ): CancelablePromise<GetQuestionListResp> {
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
  public static queryQuestion(id: string): CancelablePromise<GetQuestionResp> {
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
    authorization: string,
    body: UpdateQuestionReq
  ): CancelablePromise<UpdateQuestionResp> {
    return __request(OpenAPI, {
      method: "POST",
      url: "/gyu_oj/v1/question/update",
      headers: {
        authorization: authorization,
      },
      body: body,
    });
  }
}