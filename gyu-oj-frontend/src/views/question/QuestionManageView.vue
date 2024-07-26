<template>
  <div id="QuestionManageView">
    <a-table
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
      <template #optional="{ record }">
        <a-space>
          <a-button type="primary" @click="doUpdate(record)"> 修改</a-button>
          <a-button status="danger" @click="doDelete(record)"> 删除</a-button>
        </a-space>
      </template>
    </a-table>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watchEffect } from "vue";
import { QuestionService, QuestionVO } from "../../../question_generated";
import message from "@arco-design/web-vue/es/message";
import { useRouter } from "vue-router";

const doDelete = async (question: QuestionVO) => {
  const res = await QuestionService.deleteQuestion({
    id: question.id,
  });
  if (res.code === 200) {
    message.success("删除成功");
    await loadData();
  } else {
    message.error("删除失败");
  }
};

const router = useRouter();
const doUpdate = (question: QuestionVO) => {
  router.push({
    path: "/question/update",
    query: {
      id: question.id,
    },
  });
};

const searchParams = ref({
  pageSize: 10,
  current: 1,
});
const dataList = ref([]);
const total = ref(0);

const loadData = async () => {
  const res = await QuestionService.queryQuestionList(
    searchParams.value.current,
    searchParams.value.pageSize
  );
  if (res.code === 200) {
    dataList.value = res.data.questionList;
    total.value = res.data.total;
  } else {
    message.error("分页获取题目列表错误， " + res.msg);
  }
};

onMounted(() => {
  loadData();
});

/**
 * 监听列表参数的变化，然后加载最新的数据
 */
watchEffect(() => {
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
    title: "id",
    dataIndex: "id",
    ellipsis: true,
    tooltip: true,
  },
  {
    title: "题目标题",
    dataIndex: "title",
    ellipsis: true,
    tooltip: true,
  },
  {
    title: "题目内容",
    dataIndex: "content",
    ellipsis: true,
    tooltip: true,
  },
  {
    title: "标签",
    dataIndex: "tags",
  },
  {
    title: "题目答案",
    dataIndex: "answer",
    ellipsis: true,
    tooltip: true,
  },
  {
    title: "提交数",
    dataIndex: "submitNum",
  },
  {
    title: "通过数",
    dataIndex: "acceptedNum",
  },
  {
    title: "判题配置",
    dataIndex: "judgeConfig",
  },
  {
    title: "判题样例",
    dataIndex: "judgeCase",
    ellipsis: true,
    tooltip: true,
  },
  {
    title: "用户 id",
    dataIndex: "userId",
  },
  {
    title: "创建时间",
    dataIndex: "createTime",
  },
  {
    title: "更新时间",
    dataIndex: "updateTime",
  },
  {
    title: "操作",
    slotName: "optional",
  },
];
</script>

<style scoped>
#QuestionManageView {
}
</style>
