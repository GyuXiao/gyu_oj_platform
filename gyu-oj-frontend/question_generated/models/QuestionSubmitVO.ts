/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { JudgeInfo } from "./JudgeInfo";

export type QuestionSubmitVO = {
  id: string;
  language: string;
  submitCode: string;
  judgeInfo: JudgeInfo;
  status: number;
  questionId: string;
  userId: number;
  createTime: number;
  updateTime: number;
};

