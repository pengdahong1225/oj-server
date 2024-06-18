<template>
  <div>
    <el-card title="代码编辑器">
      <div
        style="
          display: flex;
          width: 100%;
          margin-right: 20px;
          margin-bottom: 10px;
        "
      >
        <div class="editor-header">
          <span class="languages">Language: </span>
          <el-select v-model="activeLang" placeholder="请选择" size="small">
            <el-option
              v-for="item in language_list"
              :key="item.value"
              :value="item.value"
              :disabled="item.disabled"
            >
            </el-option>
          </el-select>
        </div>
      </div>

      <editor
        ref="aceEditor"
        v-model="content"
        @init="editorInit"
        width="100%"
        height="300px"
        :lang="lang"
        :theme="theme"
        :options="{
          enableBasicAutocompletion: true,
          enableSnippets: true,
          enableLiveAutocompletion: true,
          tabSize: 4,
          fontSize: fontSize,
          readOnly: readOnly, //设置是否只读
          showPrintMargin: true, //去除编辑器里的竖线
        }"
      ></editor>
    </el-card>
  </div>
</template>

<script>
import Editor from 'vue2-ace-editor'
export default {
  components: {
    Editor
  },
  props: {
    // 是否只读
    readOnly: {
      type: Boolean,
      default: false
    },
    // 要展示的代码
    codeData: {
      type: String,
      default: ''
    },
    // 默认的主题
    valueTheme: {
      type: String,
      default: 'xcode'
    },
    // 默认的语言
    valueCodeLang: {
      type: String,
      default: 'c_cpp'
    },
    fontSize: {
      type: Number,
      default: 15
    }
  },
  data () {
    return {
      // 下拉框列表数据
      language_list: [
        { value: 'C', disabled: false },
        { value: 'C++', disabled: true },
        { value: 'java', disabled: true },
        { value: 'python', disabled: true },
        { value: 'golang', disabled: false }
      ],
      activeLang: 'C',

      // 主题数据
      listTheme: ['dracula', 'xcode', 'monokai'],
      listCodeLang: ['c_cpp', 'golang', 'java', 'python'],
      theme: '',
      lang: '',

      // 编辑器内容
      content: ''
    }
  },
  watch: {
    // 下拉框改变选项语言 -> 改变主题语言
    activeLang (newValue, oldValue) {
      if (newValue === 'C' || newValue === 'C++') {
        this.lang = 'c_cpp'
      } else {
        this.lang = newValue.toLowerCase()
      }
    }
  },
  mounted () {
    // 加载编辑器资源
    this.editorInit()
    // 初始化主题、语言
    this.theme = this.valueTheme
    this.lang = this.valueCodeLang
    // 若传输代码，则展示代码
    if (this.codeData) {
      this.content = this.codeData
    }
  },
  methods: {
    editorInit () {
      // 初始化
      require('brace/ext/language_tools')
      require('brace/ext/beautify')
      require('brace/ext/error_marker')
      require('brace/ext/searchbox')
      require('brace/ext/split')

      // 循坏加载语言
      for (let s = 0; s < this.listCodeLang.length; s++) {
        require('brace/snippets/' + this.listCodeLang[s])
      }
      for (let j = 0; j < this.listCodeLang.length; j++) {
        require('brace/mode/' + this.listCodeLang[j])
      }

      // 循坏加载主题
      for (let i = 0; i < this.listTheme.length; i++) {
        require('brace/theme/' + this.listTheme[i])
      }
    },

    formatCode () {
      const string = JSON.stringify(
        JSON.parse(this.$refs.aceEditor.editor.getValue()),
        null,
        2
      )
      this.$refs.aceEditor.editor.setValue(string)
    },
    // 父组件通过调用子组件方法获取数据
    getValue () {
      return this.$refs.aceEditor.editor.getValue()
    },
    getLang () {
      return this.activeLang
    }
  }
}
</script>

<style>
.editor-header {
  margin-right: 20px;
  font-size: 17px;
  line-height: 24px;
}
</style>
