<template>
  <div id="addQuestionView">
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
          <a-form-item field="judgeConfig.timeLimit" label="时间限制">
            <a-input-number
              v-model="form.judgeConfig.timeLimit"
              placeholder="请输入时间限制"
              mode="button"
              min="0"
              size="large"
            />
          </a-form-item>
          <a-form-item field="judgeConfig.memoryLimit" label="内存限制">
            <a-input-number
              v-model="form.judgeConfig.memoryLimit"
              placeholder="请输入内存限制"
              mode="button"
              min="0"
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
          v-for="(judgeCaseItem, index) of form.judgeCases"
          :key="index"
          no-style
        >
          <a-space direction="vertical" style="min-width: 640px">
            <a-form-item
              :field="`form.judgeCases[${index}].input`"
              :label="`输入样例-${index}`"
              :key="index"
            >
              <a-input
                v-model="judgeCaseItem.input"
                placeholder="请输入测试输入样例"
              />
            </a-form-item>
            <a-form-item
              :field="`form.judgeCases[${index}].output`"
              :label="`输出样例-${index}`"
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
import { ref } from "vue";
import MdEditor from "@/components/MdEditor.vue";

let form = ref({
  title: "",
  content: "",
  tags: [],
  answer: "",
  judgeConfig: {
    timeLimit: 1000,
    memoryLimit: 1000,
  },
  judgeCases: [
    {
      input: "",
      output: "",
    },
  ],
});

const doSubmit = () => {
  console.log("表单信息包括：", form);
};

/**
 * 新增判题样例
 */
const handleAdd = () => {
  form.value.judgeCases.push({
    input: "",
    output: "",
  });
};

/**
 * 删除判题样例
 */
const handleDelete = (index: number) => {
  form.value.judgeCases.splice(index, 1);
};

const onContentChange = (value: string) => {
  form.value.content = value;
};

const onAnswerChange = (value: string) => {
  form.value.answer = value;
};
</script>

<style scoped>
#addQuestionView {
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
