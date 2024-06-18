<template>
  <div class="problem">
    <div class="main-container">
      <div class="panel-title">{{ problemInfo.problem_title }}</div>

      <div class="title-box">Description</div>
      <div class="content-box">{{ problemInfo.description }}</div>

      <div class="title-box">Input</div>
      <div class="content-box">{{ problemInfo.input }}</div>
      <div class="title-box">Output</div>
      <div class="content-box">{{ problemInfo.output }}</div>

      <div class="test-case">
        <div class="in-box">
          <div class="title-box">Sample-input</div>
          <input
            class="content-box"
            type="text"
            :value="problemInfo.test.input"
            readonly
          />
        </div>

        <div class="out-box">
          <div class="title-box">Sample-output</div>
          <input
            class="content-box"
            type="text"
            :value="problemInfo.test.output"
            readonly
          />
        </div>
      </div>

      <div class="Hint-box">
        <el-collapse>
          <el-collapse-item>
            <template slot="title">
              <span class="el-collapse-item-title-box"> Hint </span>
            </template>

            <CodeBlock
              :code="problemInfo.code"
            ></CodeBlock>
          </el-collapse-item>
        </el-collapse>
      </div>

      <el-divider></el-divider>

      <CodeEditor :fontSize="20" ref="codeEditor"></CodeEditor>

      <div class="submit-box">
        <el-button type="primary" @click="submitHandler">submit</el-button>
      </div>
    </div>

    <div class="information-container">
      <el-descriptions title="Information" :column=1>
        <el-descriptions-item label="ID">{{
          problemInfo.id
        }}</el-descriptions-item>
        <el-descriptions-item label="Time Limit"
          >{{ problemInfo.time_limit }} MS</el-descriptions-item
        >
        <el-descriptions-item label="Memory Limit"
          >{{ problemInfo.memory_limit }} MB</el-descriptions-item
        >
        <el-descriptions-item label="IO Mode">{{
          problemInfo.io_mode
        }}</el-descriptions-item>
        <el-descriptions-item label="Create By">{{
          problemInfo.create_by
        }}</el-descriptions-item>
        <el-descriptions-item label="Level">
          {{ problemInfo.level }}
        </el-descriptions-item>
      </el-descriptions>
    </div>
  </div>
</template>

<script>
import CodeEditor from '@/components/codeEditor.vue'
import CodeBlock from '@/components/codeBlock.vue'
import { getProblemDetail, submitCode } from '@/api/problem'

export default {
  name: 'ProblemPage',
  components: {
    CodeEditor,
    CodeBlock
  },
  data () {
    return {
      problemInfo: {
        id: Number,
        create_at: '',
        create_by: '',
        io_mode: '',
        memory_limit: Number,
        time_limit: Number,
        problem_title: '',
        level: '',
        description: '',
        input: '',
        output: '',
        test: {
          input: '',
          output: ''
        },
        code: 'printf("%d\\n", a+b);'
      },
      activeLang: 'c_cpp',
      languages: [
        { value: 'c_cpp', label: 'c_cpp', disabled: false },
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
      this.problemInfo.create_by = data.create_by
      this.problemInfo.problem_title = data.title
      this.problemInfo.io_mode = data.io_mode
      this.problemInfo.time_limit = data.time_limit
      this.problemInfo.memory_limit = data.memory_limit
      switch (data.level) {
        case 1:
          this.problemInfo.level = '简单'
          break
        case 2:
          this.problemInfo.level = '中等'
          break
        case 3:
          this.problemInfo.level = '困难'
          break
        default:
          this.problemInfo.level = '困难'
          break
      }

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
    async submitHandler () {
      const submitData = this.$refs.codeEditor.getValue()
      if (submitData === '') {
        this.$message.error('请写入代码')
        return
      }
      console.log(submitData)
      const lang = this.$refs.codeEditor.getLang()

      const res = await submitCode({
        problem_id: this.problemInfo.id,
        title: this.problemInfo.problem_title,
        lang: lang,
        code: submitData
      }).catch(err => {
        console.log(err)
      })

      console.log(res)
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
    // box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.1); // 添加阴影
    margin-right: @gap;
    margin-left: 20px;

    .panel-title {
      font-size: 21px;
      font-weight: 500;
      line-height: 30px;
      padding: 5px 10px;
    }

    .title-box {
      color: #3091f2;
      font-size: 20px;
      font-weight: 400;
      padding: 15px 20px;
    }
    .content-box {
      margin-left: 25px;
      margin-right: 20px;
    }
    .Hint-box {
      padding: 15px 20px;
    }
    .el-collapse-item-title-box {
      color: #3091f2;
      font-size: 20px;
      font-weight: 400;
    }

    .test-case {
      display: flex;
      width: 90%;
      .in-box {
        width: 50%;
        .content-box{
          height: 25px;
          font-size: 15px;
        }
      }
      .out-box {
        width: 50%;
        .content-box{
          height: 25px;
          font-size: 15px;
        }
      }
      .content-box {
        width: 80%;
        border: 1px solid rgba(143, 143, 143, 0.5);
      }
    }
  }
  .information-container {
    width: @info-width;
    box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.1); // 添加阴影
    align-self: flex-start;
    margin-left: 10px;
    margin-right: 10px;
    margin-top: 10px;
    .el-descriptions {
      margin-top: 10px;
      margin-bottom: 10px;
      margin-left: 10px;
    }
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
