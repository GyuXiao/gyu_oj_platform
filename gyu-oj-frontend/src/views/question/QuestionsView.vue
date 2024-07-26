<template>
  <div id="QuestionsView">
    <a-form :model="searchParams" layout="inline">
      <a-form-item field="title" label="名称" style="min-width: 240px">
        <a-input v-model="searchParams.title" placeholder="请输入题目名称" />
      </a-form-item>
      <!--      <a-form-item field="tags" label="标签" style="min-width: 240px">-->
      <!--        <a-input-tag v-model="searchParams.tags" placeholder="请输入标签" />-->
      <!--      </a-form-item>-->
      <a-form-item>
        <a-button type="primary" @click="doSearchByTagsOrTitle">搜索</a-button>
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
      <template #tags="{ record }">
        <a-space wrap>
          <a-tag
            v-for="(tag, index) of record.tags"
            :key="index"
            :color="getTagsColor(tag)"
            >{{ tag }}
          </a-tag>
        </a-space>
      </template>
      <template #acRate="{ record }">
        {{
          `${record.submitNum ? record.acceptedNum / record.submitNum : "0"}%(${
            record.acceptedNum
          }/${record.submitNum})`
        }}
      </template>
      <template #createTime="{ record }">
        {{ moment.unix(record.createTime).format("YYYY年MM月DD日") }}
      </template>
      <template #optional="{ record }">
        <a-space>
          <a-button type="primary" @click="toQuestionPage(record)">
            去做题~
          </a-button>
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
import moment from "moment";

const tableRef = ref();
const router = useRouter();

const tagsColorsMap = new Map([
  ["simple", "#7bc616"],
  ["medium", "#ffb400"],
  ["hard", "#ff5722"],
  ["default", "#165dff"],
]);

const getTagsColor = (tag: string) => {
  if (tag == "" || tag == undefined || !tagsColorsMap.has(tag)) {
    return tagsColorsMap.get("default");
  }
  return tagsColorsMap.get(tag);
};

const toQuestionPage = (question: QuestionVO) => {
  router.push({
    path: `/question/view/id=${question.id}`,
  });
};

const searchParams = ref({
  title: "",
  tags: [],
  pageSize: 10,
  current: 1,
});
const dataList = ref([]);
const total = ref(0);

const loadData = async () => {
  const res = await QuestionService.queryQuestionList(
    searchParams.value.current,
    searchParams.value.pageSize,
    searchParams.value.title,
    searchParams.value.tags
  );
  if (res.code === 200) {
    dataList.value = res.data.questionList;
    total.value = res.data.total;
  } else {
    message.error("分页获取题目列表错误， " + res.msg);
  }
};

const doSearchByTagsOrTitle = () => {
  searchParams.value = {
    ...searchParams.value,
    current: 1, // 细节：这里需要重置页面为 1
  };
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
    title: "题号",
    dataIndex: "id",
    ellipsis: true,
    tooltip: true,
  },
  {
    title: "题目标题",
    dataIndex: "title",
  },
  {
    title: "标签",
    slotName: "tags",
  },
  {
    title: "通过率",
    slotName: "acRate",
  },
  {
    title: "创建时间",
    slotName: "createTime",
  },
  {
    slotName: "optional",
  },
];
</script>

<style scoped>
#QuestionsView {
  max-width: 1280px;
  margin: 0 auto;
}
</style>
