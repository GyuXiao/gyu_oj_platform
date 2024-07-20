<template>
  <div id="exampleView">
    <MdEditor :value="mdValue" :handle-change="onMdChange" />
    <CodeEditor :value="codeValue" :handle-change="onCodeChange" />
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import MdEditor from "@/components/MdEditor.vue";
import CodeEditor from "@/components/CodeEditor.vue"; // @ is an alias to /src

const codeValue = ref();
const mdValue = ref();
const onCodeChange = (v: string) => {
  codeValue.value = v;
  console.log("父组件传递给 Code 的值是：" + v);
};

const onMdChange = (v: string) => {
  mdValue.value = v;
  console.log("父组件传递给 Md 的值是：" + v);
};

// export default defineComponent({
//   name: "HomeView",
//   components: {
//     MdEditor,
//     HelloWorld,
//   },
// });
</script>
