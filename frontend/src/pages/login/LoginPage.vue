<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import { userRegisterService, userLoginService } from '@/api/userController'
import { useUserStore } from '@/stores'

const isRegister = ref(false)
const form = ref()
const register_form = ref(<API.RegisterForm>{
  mobile: '',
  password: '',
  repassword: ''
})
const login_form = ref(<API.LoginForm>{
  mobile: '',
  password: ''
})

// 整个表单的校验规则
// 1. 非空校验 required: true      message消息提示，  trigger触发校验的时机 blur change
// 2. 长度校验 min:xx, max: xx
// 3. 正则校验 pattern: 正则规则    \S 非空字符
// 4. 自定义校验 => 自己写逻辑校验 (校验函数)
//    validator: (rule, value, callback)
//    (1) rule  当前校验规则相关的信息
//    (2) value 所校验的表单元素目前的表单值
//    (3) callback 无论成功还是失败，都需要 callback 回调
//        - callback() 校验成功
//        - callback(new Error(错误信息)) 校验失败
const rules = {
  mobile: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式错误', trigger: 'blur' }
  ],
  password: [ 
    { required: true, message: '请输入密码', trigger: 'blur' },
    {
      pattern: /^\S{4,15}$/,
      message: '密码必须是 4-15位 的非空字符',
      trigger: 'blur'
    }
  ],
  repassword: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    {
      pattern: /^\S{4,15}$/,
      message: '密码必须是 4-15位 的非空字符',
      trigger: 'blur'
    },
    {
      validator: (rule: any, value: string, callback: Function) => {
        // 判断 value 和 当前 form 中收集的 password 是否一致
        if (value !== register_form.value.password) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback() // 就算校验成功，也需要callback
        }
      },
      trigger: 'blur'
    }
  ]
}

// 登录和注册页面切换 -- 清空表单
watch(() => isRegister.value, (newValue) => {
  register_form.value = {
    mobile: '',
    password: '',
    repassword: ''
  }
  login_form.value = {
    mobile: '',
    password: ''
  }
})

const register = async () => {
  await form.value.validate()
  const res = await userRegisterService(register_form.value)
  console.log(res)
  ElMessage.success('注册成功')
  isRegister.value = false
}
const router = useRouter()
const userStore = useUserStore()
const login = async () => {
  await form.value.validate()
  const res = await userLoginService(login_form.value)
  console.log(res)
  ElMessage.success('登录成功')
  const user = <API.UserProfile>{
    uid: res.data.data.uid,
    mobile: res.data.data.mobile,
    nickname: res.data.data.nickname,
    gender: res.data.data.gender,
    email: res.data.data.email,
    avatar_url: res.data.data.avatar_url,
    role: res.data.data.role
  }
  userStore.setUserInfo(user)
  userStore.setToken(res.data.data.token)
  router.push('/')
}
</script>

<template>
  <!-- 
    1. 结构相关
      el-row表示一行，一行分成24份 
       el-col表示列  
       (1) :span="12"  代表在一行中，占12份 (50%)
       (2) :span="6"   表示在一行中，占6份  (25%)
       (3) :offset="3" 代表在一行中，左侧margin份数

       el-form 整个表单组件
       el-form-item 表单的一行 （一个表单域）
       el-input 表单元素（输入框）
    2. 校验相关
       (1) el-form => :model="ruleForm"      绑定的整个form的数据对象 { xxx, xxx, xxx }
       (2) el-form => :rules="rules"         绑定的整个rules规则对象  { xxx, xxx, xxx }
       (3) 表单元素 => v-model="ruleForm.xxx" 给表单元素，绑定form的子属性
       (4) el-form-item => prop配置生效的是哪个校验规则 (和rules中的字段要对应)
  -->
  <el-row class="login-page">
    <!-- 左边背景区域 -->
    <el-col :span="12" class="bg"></el-col>
    <!-- 右边表单区域 -->
    <el-col :span="6" :offset="3" class="form">
      <!-- 注册相关表单 -->
      <el-form
        :model="register_form"
        :rules="rules"
        ref="form"
        size="large"
        autocomplete="off"
        v-if="isRegister"
      >
        <el-form-item>
          <h1>注册</h1>
        </el-form-item>
        <el-form-item prop="mobile">
          <el-input
            v-model="register_form.mobile"
            :prefix-icon="User"
            placeholder="请输入手机号"
          ></el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="register_form.password"
            :prefix-icon="Lock"
            type="password"
            placeholder="请输入密码"
          ></el-input>
        </el-form-item>
        <el-form-item prop="repassword">
          <el-input
            v-model="register_form.repassword"
            :presfix-icon="Lock"
            type="password"
            placeholder="请输入再次密码"
          ></el-input>
        </el-form-item>
        <el-form-item>
          <el-button
            @click="register"
            class="button"
            type="primary"
            auto-insert-space
          >
            注册
          </el-button>
        </el-form-item>
        <el-form-item class="flex">
          <el-link type="info" :underline="false" @click="isRegister = false">
            ← 返回
          </el-link>
        </el-form-item>
      </el-form>

      <!-- 登录相关表单 -->
      <el-form
        :model="login_form"
        :rules="rules"
        ref="form"
        size="large"
        autocomplete="off"
        v-else
      >
        <el-form-item>
          <h1>登录</h1>
        </el-form-item>
        <el-form-item prop="mobile">
          <el-input
            v-model="login_form.mobile"
            :prefix-icon="User"
            placeholder="请输入手机号"
          ></el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="login_form.password"
            name="password"
            :prefix-icon="Lock"
            type="password"
            placeholder="请输入密码"
          ></el-input>
        </el-form-item>
        <el-form-item class="flex">
          <div class="flex">
            <el-checkbox>记住我</el-checkbox>
            <el-link type="primary" :underline="false">忘记密码？</el-link>
          </div>
        </el-form-item>
        <el-form-item>
          <el-button
            @click="login"
            class="button"
            type="primary"
            auto-insert-space
            >登录</el-button
          >
        </el-form-item>
        <el-form-item class="flex">
          <el-link type="info" :underline="false" @click="isRegister = true">
            注册 →
          </el-link>
        </el-form-item>
      </el-form>
    </el-col>
  </el-row> 

</template>

<style lang="scss" scoped>
.login-page {
  height: 100vh;
  background-color: #fff;
  .bg {
    background: url('@/assets/cat.png') no-repeat 60% center / 240px auto,
      url('@/assets/login_bg.jpg') no-repeat center / cover;
    border-radius: 0 20px 20px 0;
  }
  .form {
    display: flex;
    flex-direction: column;
    justify-content: center;
    user-select: none;
    .title {
      margin: 0 auto;
    }
    .button {
      width: 100%;
    }
    .flex {
      width: 100%;
      display: flex;
      justify-content: space-between;
    }
  }
}
</style>
