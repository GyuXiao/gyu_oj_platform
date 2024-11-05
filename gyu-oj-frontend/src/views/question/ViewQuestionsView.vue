<template>
  <div id="ViewQuestionView">
    <a-row :gutter="[24, 24]">
      <a-col :md="12" :xs="24">
        <a-tabs default-active-key="question">
          <a-tab-pane key="question" title="题目">
            <a-card v-if="question" :title="question.title">
              <a-descriptions :column="{ xs: 1, md: 2, lg: 3 }">
                <a-descriptions-item label="时间限制">
                  {{ question.judgeConfig.timeLimit ?? 0 }} MS
                </a-descriptions-item>
                <a-descriptions-item label="内存限制">
                  {{ question.judgeConfig.memoryLimit ?? 0 }} KB
                </a-descriptions-item>
              </a-descriptions>
              <MdViewer :value="question.content || ''" />
              <template #extra>
                <a-space wrap>
                  <a-tag
                    v-for="(tag, index) of question.tags"
                    :key="index"
                    :color="getTagsColor(tag)"
                    >{{ tag }}
                  </a-tag>
                </a-space>
              </template>
            </a-card>
          </a-tab-pane>
          <a-tab-pane key="comment" title="评论" disabled> 评论区</a-tab-pane>
          <a-tab-pane key="answer" title="答案"> 暂时无法查看答案</a-tab-pane>
        </a-tabs>
      </a-col>
      <a-col :md="12" :xs="24">
        <a-form :model="form" layout="inline">
          <a-form-item
            field="language"
            label="编程语言"
            style="min-width: 240px"
          >
            <a-select
              v-model="form.language"
              :style="{ width: '320px' }"
              placeholder="选择编程语言"
            >
              <a-option>java</a-option>
              <a-option>cpp</a-option>
              <a-option>go</a-option>
            </a-select>
          </a-form-item>
        </a-form>
        <CodeEditor
          :value="form.submitCode"
          :language="form.language"
          :handle-change="changeCode"
        />
        <a-divider size="0" />
        <div style="text-align: right">
          <a-button
            type="dashed"
            status="success"
            style="min-width: 200px"
            @click="doSubmitCode"
          >
            提交
          </a-button>
        </div>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { defineProps, onMounted, ref, withDefaults } from "vue";
import {
  CreateQuestionSubmitReq,
  QuestionService,
  QuestionSubmitService,
  QuestionVO,
} from "../../../question_generated";
import message from "@arco-design/web-vue/es/message";
import MdViewer from "@/components/MdViewer.vue";
import CodeEditor from "@/components/CodeEditor.vue";
import { useRouter } from "vue-router";

// 从跳转过来的页面传进来 id
interface Props {
  id: string;
}

const props = withDefaults(defineProps<Props>(), {
  id: () => "",
});

const question = ref<QuestionVO>();
const loadData = async () => {
  const res = await QuestionService.queryQuestion(props.id);
  if (res.code === 200) {
    question.value = res.data.question;
    if (question.value) {
      question.value.judgeConfig = JSON.parse(res.data.question.judgeConfig);
    }
  } else {
    message.error("加载失败， " + res.msg);
  }
};

const form = ref<CreateQuestionSubmitReq>({
  language: "go",
  submitCode: "",
  questionId: "",
});

const changeCode = (value: string) => {
  form.value.submitCode = value;
};

/**
 * 提交代码
 */
const doSubmitCode = async () => {
  if (!question.value?.id) {
    return;
  }
  if (form.value.submitCode == "") {
    message.error("代码不能为空哦~");
    return;
  }
  const res = await QuestionSubmitService.createQuestionSubmit({
    ...form.value,
    questionId: question.value?.id,
  });
  if (res.code === 200) {
    message.success("提交成功");
    toQuestionSubmitViewPage();
  } else {
    message.error("提交失败," + res.msg);
  }
};

const router = useRouter();
/**
 * 跳转到题目提交页
 */
const toQuestionSubmitViewPage = () => {
  router.push({
    path: `/question_submit`,
  });
};

const tagsColorsMap = new Map([
  ["简单", "green"],
  ["中等", "gold"],
  ["困难", "orangered"],
  ["default", "blue"],
]);
const getTagsColor = (tag: string) => {
  if (tag == "" || tag == undefined || !tagsColorsMap.has(tag)) {
    return tagsColorsMap.get("default");
  }
  return tagsColorsMap.get(tag);
};

onMounted(() => {
  loadData();
});
</script>

<style scoped>
#ViewQuestionView {
  max-width: 1280px;
  margin: 0 auto;
}
</style>
