/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { JudgeConfig } from "./JudgeConfig";

export type QuestionVO = {
  id: string;
  title: string;
  content: string;
  tags: Array<string>;
  answer: string;
  submitNum: number;
  acceptedNum: number;
  judgeConfig: JudgeConfig;
  userId: number;
  createTime: number;
  updateTime: number;
};

