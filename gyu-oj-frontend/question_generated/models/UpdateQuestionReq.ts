/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { JudgeCase } from "./JudgeCase";
import type { JudgeConfig } from "./JudgeConfig";

export type UpdateQuestionReq = {
  id?: string;
  title?: string;
  content?: string;
  tags?: Array<string>;
  answer?: string;
  judgeCase?: Array<JudgeCase>;
  judgeConfig?: JudgeConfig;
};

