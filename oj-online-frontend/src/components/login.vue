<template>
  <div class="login">
    <div class="container">
      <div class="title">
        <h3>手机号登录</h3>
        <p>未注册的手机号登录后将自动注册</p>
      </div>

      <div class="form">
        <div class="form-item">
          <input
            v-model="mobile"
            class="inp"
            maxlength="11"
            placeholder="请输入手机号码"
            type="text"
          />
        </div>
        <div class="form-item">
          <input
            v-model="captchaValue"
            class="inp"
            maxlength="5"
            placeholder="请输入图形验证码"
            type="text"
          />
          <img v-if="captchaUrl" :src="captchaUrl" @click="getPicCode" alt="" />
        </div>
        <div class="form-item">
          <input
            v-model="smsCode"
            class="inp"
            placeholder="请输入短信验证码"
            type="text"
          />
          <button @click="getSmsCode">
            {{
              second === totalSecond ? '获取验证码' : second + '秒后重新发送'
            }}
          </button>
        </div>
      </div>

      <div class="login-btn" @click="login">登录</div>
    </div>
  </div>
</template>

<script>
import { getPicCode, getSmsCode, mobileLogin } from '@/api/login'
export default {
  name: 'LoginPage',
  data () {
    return {
      captchaID: '', // 图形验证码的key
      captchaUrl: '', // 存储请求渲染的图片地址

      second: 60, // 倒计时秒数
      totalSecond: 60, // 总秒数
      timer: null,

      mobile: '', // 手机号
      captchaValue: '', // 用户输入的图形验证码
      smsCode: ''
    }
  },
  created () {
    this.getPicCode()
  },
  methods: {
    async getPicCode () {
      const res = await getPicCode()
      this.captchaID = res.data.captchaID
      this.captchaUrl = res.data.captchaUrl
    },
    async getSmsCode () {
      // 获取短信验证码前要简单校验手机号和图形验证码是否合法
      if (!this.validate()) {
        return
      }
      // 请求
      if (!this.timer && this.second === this.totalSecond) {
        const res = await getSmsCode(this.mobile, this.captchaID, this.captchaValue)
        if (res.message !== 'OK') {
          this.$message({
            message: res.message,
            type: 'fail'
          })
          return
        }

        this.$message({
          message: '短信验证码发送成功，请注意查收',
          type: 'success'
        })

        // 开启倒计时
        this.timer = setInterval(() => {
          this.second--
          if (this.second <= 0) {
            clearInterval(this.timer)
            this.timer = null
            this.second = this.totalSecond
          }
        }, 1000)
      }
    },
    async login () {
      console.log('login')
      if (!this.validate()) {
        return
      }
      const res = await mobileLogin(this.mobile, this.smsCode)
      console.log(res)
      this.$message({
        message: '登录成功',
        type: 'success'
      })
      this.$store.commit('user/setUserInfo', res.data)

      // 登录成功后，通知父组件关闭对话框
      this.$emit('close')
    },
    validate () {
      // 校验 手机号 和 图形验证码 是否合法
      if (!/^1[3-9]\d{9}$/.test(this.mobile)) {
        this.$message({
          message: '请输入正确的手机号',
          type: 'warning'
        })
        return false
      }
      if (!/^\w{4}$/.test(this.captchaValue)) {
        this.$message({
          message: '请输入正确的图形验证码',
          type: 'warning'
        })
        return false
      }
      return true
    }
  }
}
</script>

<style lang="less" scoped>
.container {
  padding: 20px 10px;

  .title {
    margin-bottom: 20px;
    h3 {
      font-size: 26px;
      font-weight: normal;
    }
    p {
      line-height: 40px;
      font-size: 14px;
      color: #b8b8b8;
    }
  }

  .form-item {
    border-bottom: 1px solid #f3f1f2;
    padding: 8px;
    margin-bottom: 14px;
    display: flex;
    align-items: center;
    .inp {
      display: block;
      border: none;
      outline: none;
      height: 32px;
      font-size: 14px;
      flex: 1;
    }
    img {
      width: 100px;
      height: 50px;
    }
    button {
      height: 31px;
      border: none;
      font-size: 13px;
      color: #cea26a;
      background-color: transparent;
      padding-right: 9px;
    }
  }

  .login-btn {
    width: 100%;
    height: 42px;
    margin-top: 39px;
    background: linear-gradient(90deg, #ecb53c, #ff9211);
    color: #fff;
    border-radius: 39px;
    box-shadow: 0 10px 20px 0 rgba(0, 0, 0, 0.1);
    letter-spacing: 2px;
    display: flex;
    justify-content: center;
    align-items: center;
    cursor: pointer;
  }
}
</style>
