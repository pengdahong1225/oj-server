<template>
  <div>
    <codemirror
      ref="cm"
      v-model="code"
      :options="cmOptions"
      @input="inputChange"
    ></codemirror>
  </div>
</template>

<script>
import { codemirror } from 'vue-codemirror'
import 'vue-codemirror/node_modules/codemirror/lib/codemirror.css'
// mode
import 'vue-codemirror/node_modules/codemirror/mode/sql/sql'
import 'vue-codemirror/node_modules/codemirror/mode/javascript/javascript'

export default {
  name: 'codeEditor',
  components: { codemirror },
  data () {
    return {
      code: 'select a from table1 where b = 1',
      cmOptions: {
        // 语言及语法模式
        mode: 'sql',
        // 主题
        theme: 'idea',
        // 显示函数
        line: true,
        lineNumbers: true,
        // 软换行
        lineWrapping: true,
        // tab宽度
        tabSize: 4
      }
    }
  },
  methods: {
    inputChange (content) {
      this.$nextTick(() => {
        console.log('code:' + this.code)
        console.log('content:' + content)
      })
    }
  }
}
</script>
