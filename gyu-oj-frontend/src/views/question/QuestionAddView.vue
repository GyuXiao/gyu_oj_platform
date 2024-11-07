<template>
  <div id="QuestionAddView">
    <h2>创建题目</h2>
    <a-form
      style="max-width: 1000px; margin: 0 auto"
      label-align="left"
      auto-label-width
      :model="form"
    >
      <a-form-item field="title" label="题目">
        <a-input v-model="form.title" placeholder="请输入题目标题" />
      </a-form-item>
      <a-form-item field="tags" label="标签">
        <a-input-tag v-model="form.tags" placeholder="请选择标签" allow-clear />
      </a-form-item>
      <a-form-item field="content" label="题目内容" flex="auto">
        <div class="questionContent">
          <MdEditor :value="form.content" :handle-change="onContentChange" />
        </div>
      </a-form-item>
      <a-form-item field="answer" label="题目答案">
        <div class="questionAnswer">
          <MdEditor :value="form.answer" :handle-change="onAnswerChange" />
        </div>
      </a-form-item>
      <a-form-item label="判题配置" :content-flex="false" :merge-props="false">
        <a-space direction="vertical" style="min-width: 640px">
          <a-form-item field="judgeConfig.timeLimit" label="时间限制（正整数）">
            <a-input-number
              v-model="form.judgeConfig.timeLimit"
              placeholder="请输入时间限制"
              mode="button"
              min="1"
              size="large"
            />
          </a-form-item>
          <a-form-item
            field="judgeConfig.memoryLimit"
            label="内存限制（正整数）"
          >
            <a-input-number
              v-model="form.judgeConfig.memoryLimit"
              placeholder="请输入内存限制"
              mode="button"
              min="1"
              size="large"
            />
          </a-form-item>
        </a-space>
      </a-form-item>
      <a-form-item
        label="测试样例配置"
        :content-flex="false"
        :merge-props="false"
      >
        <a-form-item
          v-for="(judgeCaseItem, index) of form.judgeCase"
          :key="index"
          no-style
        >
          <a-space
            direction="vertical"
            style="min-width: 640px; margin-bottom: 20px"
          >
            <a-form-item
              :field="`form.judgeCases[${index}].input`"
              :label="`输入样例-${index + 1}`"
              :key="index"
            >
              <a-input
                v-model="judgeCaseItem.input"
                placeholder="请输入测试输入样例"
              />
            </a-form-item>
            <a-form-item
              :field="`form.judgeCases[${index}].output`"
              :label="`输出样例-${index + 1}`"
              :key="index"
            >
              <a-input
                v-model="judgeCaseItem.output"
                placeholder="请输入测试输出样例"
              />
            </a-form-item>
            <a-button status="danger" @click="handleDelete(index)">
              删除测试样例
            </a-button>
          </a-space>
        </a-form-item>
        <div style="margin-top: 32px">
          <a-button @click="handleAdd" type="outline" status="success"
            >新增测试样例
          </a-button>
        </div>
      </a-form-item>
      <div style="margin-top: 16px" />
      <a-form-item>
        <div class="submitQuestionButton">
          <a-button type="primary" style="min-width: 200px" @click="doSubmit"
            >提交
          </a-button>
        </div>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import MdEditor from "@/components/MdEditor.vue";
import { QuestionService } from "../../../question_generated";
import message from "@arco-design/web-vue/es/message";
import { useRoute } from "vue-router";

let form = ref({
  title: "",
  content: "",
  tags: [],
  answer: "",
  judgeConfig: {
    timeLimit: 1000,
    memoryLimit: 1000,
  },
  judgeCase: [
    {
      input: "",
      output: "",
    },
  ],
});

const route = useRoute();
// 如果页面地址包含 update，视为更新页面
const updatePage = route.path.includes("update");
const loadData = async () => {
  const id = route.query.id;
  if (!id) {
    return;
  }
  const res = await QuestionService.queryQuestion(id as string);
  if (res.code === 200) {
    form.value = res.data.question;
    console.log("此时 form.value 的值为", form.value);
    // json 转为 js 对象
    if (!form.value.judgeCase) {
      form.value.judgeCase = [
        {
          input: "",
          output: "",
        },
      ];
    } else {
      form.value.judgeCase = JSON.parse(form.value.judgeCase as never);
    }
    if (!form.value.judgeConfig) {
      form.value.judgeConfig = {
        memoryLimit: 1000,
        timeLimit: 1000,
      };
    } else {
      form.value.judgeConfig = JSON.parse(form.value.judgeConfig as never);
    }
    if (!form.value.tags) {
      form.value.tags = [];
    }
  } else {
    message.error("加载失败，" + res.msg);
  }
};

onMounted(() => {
  loadData();
});

const doSubmit = async () => {
  console.log("表单信息包括：", form);
  if (updatePage) {
    const res = await QuestionService.updateQuestion(form.value);
    if (res.code === 200) {
      message.success("更新成功");
    } else {
      message.error("更新失败，" + res.msg);
    }
  } else {
    const res = await QuestionService.createQuestion(form.value);
    if (res.code === 200) {
      message.success("题目创建成功");
    } else {
      message.error("题目创建失败，" + res.msg);
    }
  }
};

/**
 * 新增判题样例
 */
const handleAdd = () => {
  form.value.judgeCase.push({
    input: "",
    output: "",
  });
};

/**
 * 删除判题样例
 */
const handleDelete = (index: number) => {
  form.value.judgeCase.splice(index, 1);
};

const onContentChange = (value: string) => {
  form.value.content = value;
};

const onAnswerChange = (value: string) => {
  form.value.answer = value;
};
</script>

<style scoped>
#QuestionAddView {
}

.questionContent {
  width: 100%;
}

.questionAnswer {
  width: 100%;
}

.submitQuestionButton {
  margin-left: auto;
  margin-right: auto;
}
</style>
