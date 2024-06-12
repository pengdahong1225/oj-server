<template>
  <div class="user-profile">
    <div class="avatar-container">
      <img :src="squareUrl" alt="Avatar" class="avatar" />
    </div>

    <h3 class="username">{{ nickname }}</h3>

    <p class="status-message">The guy is so lazy that has not any profile.</p>

    <el-divider> </el-divider>

    <div class="statistics-container">
      <div class="echarts-box" id="e-box"></div>
      <div class="submit-record">提交记录</div>
    </div>
  </div>
</template>

<script>
import * as echarts from 'echarts'
import { getUserProfile, getUserSolvedList } from '@/api/user'
export default {
  name: 'UserPage',
  data () {
    return {
      size: 80,
      squareUrl:
        'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'
    }
  },
  computed: {
    nickname () {
      return this.$store.state.user.userInfo.nickname
    }
  },
  created () {
    // 拉取用户信息
    this.getUserProfile()
  },
  mounted () {
    // 初始化饼图
    this.myChart = echarts.init(document.querySelector('#e-box'))
    this.myChart.setOption({
      // 提示框
      tooltip: {
        trigger: 'item'
      },
      // 图例
      legend: {
        top: '10%',
        left: 'center'
      },
      // 数据
      series: [
        {
          name: '题目AC 占比',
          type: 'pie',
          radius: ['40%', '70%'], // 半径
          center: ['50%', '60%'],
          startAngle: 180,
          endAngle: 360,
          data: [
            { value: 1048, name: 'Easy' },
            { value: 735, name: 'Medium' },
            { value: 580, name: 'Hard' }
          ],
          itemStyle: {
            // 自定义每个扇区的颜色
            color: function (params) {
              // params是每个扇区的相关信息，包括数据、数据索引和名称
              const colorList = ['#00BFFF', '#FFFF00', '#CD5C5C']
              return colorList[params.dataIndex]
            }
          }
        }
      ]
    })
  },
  methods: {
    async getUserProfile () {
      const uid = this.$store.getters.uid
      if (!uid) {
        this.$message({
          message: '请先登录',
          type: 'warning'
        })
        return
      }

      let res = await getUserProfile(uid)
      // 更新到store
      this.$store.commit('user/setUserInfo', res.data)

      // 继续拉取用户解题信息
      res = await getUserSolvedList()
      // 更新到store
      console.log(res)
      this.$store.commit('user/setUserSolvedList', res.data)
    }
  }
}
</script>

<style scoped>
.user-profile {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.avatar-container {
  margin-bottom: 15px;
}

.avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
}

.username {
  font-size: 20px;
  margin-top: 10px;
  text-align: center;
}

.status-message {
  color: #99a9bf;
  font-size: 14px;
  margin-top: 8px;
  text-align: center;
}

.stats-container {
  display: flex;
  justify-content: space-around;
  margin-top: 70px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  font-size: 16px;
  color: #409eff;
}

.statistics-container {
  /* 弹性盒子 */
  display: flex;
}

.statistics-container .echarts-box {
  width: 400px;
  height: 300px;
  /* margin: 0 auto;
  border: 1px solid #ccc; */
}
</style>
