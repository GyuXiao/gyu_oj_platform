<template>
  <div id="questionSubmitView">
    <a-form :model="searchParams" layout="inline">
      <a-form-item field="questionId" label="题号" style="min-width: 240px">
        <a-input v-model="searchParams.questionId" placeholder="请输入题号" />
      </a-form-item>
      <a-form-item field="language" label="编程语言" style="min-width: 240px">
        <a-select
          v-model="searchParams.language"
          :style="{ width: '320px' }"
          placeholder="选择编程语言"
        >
          <a-option>go</a-option>
          <a-option>java</a-option>
          <a-option>cpp</a-option>
        </a-select>
      </a-form-item>
      <a-form-item>
        <a-button type="text" @click="doSearch">搜索</a-button>
      </a-form-item>
    </a-form>
    <a-divider size="0" />
    <a-table
      :ref="tableRef"
      :columns="columns"
      :data="dataList"
      :pagination="{
        showTotal: true,
        pageSize: searchParams.pageSize,
        current: searchParams.current,
        total,
      }"
      @page-change="onPageChange"
    >
      <template #judgeInfoMessage="{ record }">
        <template
          v-if="
            record.judgeInfo.message === undefined ||
            record.judgeInfo.message === null ||
            record.judgeInfo.message === ''
          "
        >
          <a-tag loading> {{ "loading" }}</a-tag>
        </template>
        <template v-else>
          <a-tag
            :color="getJudgeMessageColor(record.judgeInfo.message)"
            bordered
            >{{ record.judgeInfo.message }}
          </a-tag>
        </template>
      </template>
      <template #status="{ record }">
        <template v-if="record.status < 2">
          <a-tag loading :color="getJudgeStatusStyle(record.status).color">
            {{ getJudgeStatusStyle(record.status).text }}
          </a-tag>
        </template>
        <template v-else-if="record.status === 2">
          <a-tag :color="getJudgeStatusStyle(record.status).color" bordered>
            {{ getJudgeStatusStyle(record.status).text }}
          </a-tag>
        </template>
        <template v-else>
          <a-tag :color="getJudgeStatusStyle(record.status).color" bordered>
            {{ getJudgeStatusStyle(record.status).text }}
          </a-tag>
        </template>
      </template>
      <template #questionId="{ record }">
        <a-button type="text" @click="toQuestionPage(record.questionId)">
          {{ record.questionId }}
        </a-button>
      </template>
      <template #createTime="{ record }">
        {{ moment.unix(record.createTime).format("YYYY年MM月DD日") }}
      </template>
    </a-table>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, onBeforeUnmount, watchEffect } from "vue";
import {
  QueryQuestionSubmitReq,
  QuestionSubmitService,
} from "../../../question_generated";
import message from "@arco-design/web-vue/es/message";
import { useRouter } from "vue-router";
import moment from "moment";

const tableRef = ref();
const router = useRouter();

const judgeStatusObjList = [
  { text: "等待中", color: "grey" },
  { text: "判题中", color: "blue" },
  { text: "成功", color: "green" },
  { text: "失败", color: "red" },
];

const getJudgeStatusStyle = (status: number) => {
  if (status == null) {
    return judgeStatusObjList[judgeStatusObjList.length - 1];
  }
  return judgeStatusObjList[status];
};

const judgeMessageMap = new Map([
  ["Accepted", "green"],
  ["Wrong Answer", "red"],
  ["Compile Error", "red"],
  ["Runtime Error", "red"],
  ["System Error", "red"],
  ["Memory Limit Exceeded", "orange"],
  ["Time Limit Exceeded", "orange"],
  ["Waiting", "gold"],
  ["default", "grey"],
]);

const getJudgeMessageColor = (message: string) => {
  if (message == "" || message == undefined || !judgeMessageMap.has(message)) {
    return judgeMessageMap.get("default");
  }
  return judgeMessageMap.get(message);
};

const toQuestionPage = (questionId: string) => {
  router.push({
    path: `/question/view/id=${questionId}`,
  });
};

const searchParams = ref<QueryQuestionSubmitReq>({
  questionId: undefined,
  language: undefined,
  pageSize: 20,
  current: 1,
});
const dataList = ref([]);
const total = ref(0);

const loadData = async () => {
  const res = await QuestionSubmitService.queryQuestionSubmitList(
    searchParams.value.current,
    searchParams.value.pageSize,
    searchParams.value.language,
    searchParams.value.questionId
  );
  if (res.code === 200) {
    dataList.value = res.data.questionSubmitList;
    total.value = res.data.totalNum;
  } else {
    message.error("分页获取题目列表错误， " + res.msg);
  }
};

const doSearch = () => {
  searchParams.value = {
    ...searchParams.value,
    current: 1, // 细节：这里需要重置页面为 1
  };
};

const refreshFlag = ref(false);
// 设置定时器，每隔 3 秒执行一次 loadData
const timer = setInterval(() => {
  if (refreshFlag.value) {
    loadData();
    refreshFlag.value = false;
  }
}, 3000);

onMounted(() => {
  loadData();
});

// 组件销毁时清除定时器
onBeforeUnmount(() => {
  clearInterval(timer);
});

/**
 * 监听列表参数的变化，然后加载最新的数据
 */
watchEffect(() => {
  // 更新 refreshFlag 的值
  if (!refreshFlag.value) {
    refreshFlag.value = true;
  }
  loadData();
});

const onPageChange = (page: number) => {
  searchParams.value = {
    ...searchParams.value,
    current: page,
  };
};

const columns = [
  {
    title: "提交号",
    dataIndex: "id",
    ellipsis: true,
    tooltip: true,
  },
  {
    title: "提交语言",
    dataIndex: "language",
  },
  {
    title: "运行结果",
    slotName: "judgeInfoMessage",
  },
  {
    title: "运行时间 (MS)",
    dataIndex: "judgeInfo.time",
  },
  {
    title: "运行内存 (MB)",
    dataIndex: "judgeInfo.memory",
  },
  {
    title: "提交状态",
    slotName: "status",
  },
  {
    title: "题号",
    slotName: "questionId",
    ellipsis: true,
    tooltip: true,
  },
  {
    title: "提交者",
    dataIndex: "userId",
  },
  {
    title: "提交时间",
    slotName: "createTime",
  },
];
</script>

<style scoped>
#questionSubmitView {
  max-width: 1280px;
  margin: 0 auto;
}
</style>
