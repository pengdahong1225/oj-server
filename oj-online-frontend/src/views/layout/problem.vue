<template>
  <div class="problem">
    <div class="main-container">
      <div class="panel-title">{{ title }}</div>

      <div class="title">Description</div>
      <div class="content">{{ problemInfo.description }}</div>

      <div class="title">Input</div>
      <div class="content">{{ problemInfo.input }}</div>
      <div class="title">Output</div>
      <div class="content">{{ problemInfo.output }}</div>

      <div class="test-case">
        <div class="in-box">
          <div class="title">Sample-input</div>
          <input
            class="content"
            type="text"
            :value="problemInfo.test.input"
            readonly
          />
        </div>
        "
        <div class="out-box">
          <div class="title">Sample-output</div>
          <input
            class="content"
            type="text"
            :value="problemInfo.test.output"
            readonly
          />
        </div>
      </div>

      <div class="title">Hint</div>
      <div class="content">
        <el-collapse>
          <el-collapse-item title="C" name="C">
            <CodeBlock :code="problemInfo.code"></CodeBlock>
          </el-collapse-item>
        </el-collapse>
      </div>

      <el-divider></el-divider>

      <div class="editor-header">
        <span class="languages">languages: </span>
        <el-select v-model="activeLang" placeholder="请选择">
          <el-option
            v-for="item in languages"
            :key="item.value"
            :value="item.value"
            :disabled="item.disabled"
          >
          </el-option>
        </el-select>
      </div>

      <CodeBlock :code="problemInfo.code"></CodeBlock>
      <!-- <CodeEditor v-model="code"></CodeEditor> -->

      <div class="submit-box">
        <el-button type="primary" @click="submitHandler">submit</el-button>
      </div>
    </div>
    <div class="information-container">
      <div>Information</div>
      <br />
      <div v-for="i in 4" :key="i" class="text item">
        {{ '列表内容 ' + i }}
      </div>
    </div>
  </div>
</template>

<script>
import CodeBlock from '@/components/codeBlock.vue'
import { getProblemDetail } from '@/api/problem'

export default {
  name: 'ProblemPage',
  components: {
    CodeBlock
  },
  data () {
    return {
      problemInfo: {
        id: Number,
        create_at: '',
        title: '',
        description: '',
        input: '',
        output: '',
        test: {
          input: '',
          output: ''
        },
        code: `#include <stdio.h>    
int main(){
    int a, b;
    scanf("%d%d", &a, &b);
    printf("%d\\n", a+b);
    return 0;
}
`
      },
      activeLang: 'C',
      languages: [
        { value: 'C', label: 'C', disabled: false },
        { value: 'C++', label: 'C++', disabled: false },
        { value: 'Java', label: 'Java', disabled: true },
        { value: 'Python', label: 'Python', disabled: true },
        { value: 'Golang', label: 'Golang', disabled: false }
      ]
    }
  },
  created () {
    // 拉取题目详细数据
    this.getProblemDetail(this.$route.params.id)
  },
  methods: {
    async getProblemDetail (id) {
      const { data } = await getProblemDetail(id)
      this.problemInfo.id = data.id
      this.problemInfo.create_at = data.create_at
      this.problemInfo.title = data.title

      const des = JSON.parse(data.description)
      this.problemInfo.description = des.Description
      this.problemInfo.input = des.Input
      this.problemInfo.output = des.Output

      const testCast = JSON.parse(data.test_case)
      this.problemInfo.test.input = testCast.input[0]
        ? testCast.input[0].content
        : ''
      this.problemInfo.test.output = testCast.output[0]
        ? testCast.output[0].content
        : ''
    },
    submitHandler () {
      console.log('submitHandler')
    }
  }
}
</script>

<style lang="less" scoped>
@main-width: 85%; // 主容器宽度
@info-width: 15%; // 信息容器宽度
@gap: 20px;

.problem {
  display: flex;
  justify-content: space-between;
  width: 100%;

  .main-container {
    width: @main-width;
    box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.1); // 添加阴影
    margin-right: @gap;
    margin-left: 20px;

    .panel-title {
      font-size: 21px;
      font-weight: 500;
      line-height: 30px;
      padding: 5px 10px;
    }

    .title {
      color: #3091f2;
      font-size: 20px;
      font-weight: 400;
      padding: 15px 20px;
    }
    .content {
      margin-left: 25px;
      margin-right: 20px;
      font-size: 17px;
      line-height: 24px;
    }

    .test-case {
      display: flex;
      width: 90%;
      .in-box {
        width: 50%;
      }
      .out-box {
        width: 50%;
      }
      .content {
        width: 80%;
        border: 1px solid rgba(143, 143, 143, 0.5);
      }
    }
  }
  .information-container {
    width: @info-width;
    box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.1); // 添加阴影
    align-self: flex-start;
    margin-left: auto;
  }

  .editor-header {
    margin-left: 20px;
    margin-right: 20px;
    font-size: 17px;
    line-height: 24px;
  }

  .submit-box {
    text-align: right; // 靠右
    margin-top: 10px;
    margin-bottom: 10px;
    .el-button {
      font-size: 20px;
    }
  }
}
</style>
