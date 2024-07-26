import { CreateQuestionResp } from "./CreateQuestionResp";
import { GetQuestionListResp } from "./GetQuestionListResp";
import { DeleteQuestionResp } from "./DeleteQuestionResp";
import { GetQuestionResp } from "./GetQuestionResp";
import { UpdateQuestionResp } from "./UpdateQuestionResp";
import { CreateQuestionSubmitResp } from "./CreateQuestionSubmitResp";

export type BaseCreateQuestionResponse = {
  code: number;
  msg: string;
  data?: CreateQuestionResp;
};

export type BaseQueryQuestionListResponse = {
  code: number;
  msg: string;
  data?: GetQuestionListResp;
};

export type BaseDeleteQuestionResponse = {
  code: number;
  msg: string;
  data?: DeleteQuestionResp;
};

export type BaseQueryQuestionResponse = {
  code: number;
  msg: string;
  data?: GetQuestionResp;
};

export type BaseUpdateQuestionResponse = {
  code: number;
  msg: string;
  data?: UpdateQuestionResp;
};

export type BaseCreateQuestionSubmitResponse = {
  code: number;
  msg: string;
  data?: CreateQuestionSubmitResp;
};
