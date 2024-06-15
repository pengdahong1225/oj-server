<template>
  <pre class="code-block"><code v-html="formattedCode"></code></pre>
</template>

<script>
export default {
  name: 'CodeBlock',
  props: {
    code: {
      type: String,
      required: true
    }
  },
  computed: {
    formattedCode () {
      return this.highlight(this.code)
    }
  },
  methods: {
    highlight (code) {
      // 在这里添加代码高亮的逻辑
      // 你可能需要使用一个库，比如highlight.js
      return code
    }
  }
}
</script>

<style scoped lang="less">
.code-block{
    background-color: #f0f0f0; /* 浅灰色 */
    .code{
        color: #495060;
        font-size: 100%;
        word-break: normal;
        font-family: Consolas,Menlo,Courier,monospace;
    }
}
</style>
