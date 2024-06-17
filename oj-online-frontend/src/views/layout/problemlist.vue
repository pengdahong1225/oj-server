<template>
  <div class="problemlist">
    <div class="search-box">
      <!-- 搜索框 -->
      <el-autocomplete
        :fetch-suggestions="querySearchAsync"
        placeholder="搜索题目"
        @select="handleSelect"
        suffix-icon="el-icon-search"
        size="medium"
        class="search-pro"
      ></el-autocomplete>

      <!-- tag选择器 -->
      <el-select filterable placeholder="标签" class="select-tag">
        <el-option> </el-option>
      </el-select>
    </div>

    <!-- 题目列表 -->
    <el-table :data="problemList" style="width: 100%" stripe>
      <el-table-column width="60" align="center">
        <template slot-scope="scope">
          <el-icon-check
            v-if="isSolved(scope.row.id)"
            class="el-icon-check"
          ></el-icon-check>
        </template>
      </el-table-column>
      <el-table-column label="#" prop="id" width="90" sortable align="center">
      </el-table-column>
      <el-table-column label="题目" prop="title" width="400">
        <template slot-scope="scope">
          <el-link
            type="primary"
            :underline="false"
            @click="
              $router.push({
                path: `problem/${scope.row.id}`
              })
            "
            >{{ scope.row.title }}</el-link
          >
        </template>
      </el-table-column>
      <el-table-column label="难度" prop="level" :formatter="formatterLevel">
      </el-table-column>
      <el-table-column prop="tag" label="标签" width="300">
        <template slot-scope="scope">
          <el-tag
            v-for="item in scope.row.tags"
            :key="item"
            type="success"
            disable-transitions
          >
            {{ item }}</el-tag
          >
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { getProblemList } from '@/api/problem'
import { getUserSolvedList } from '@/api/user'
export default {
  name: 'ProblemListPage',
  data () {
    return {}
  },
  computed: {
    problemList () {
      return this.$store.getters.problemList
    },
    isLogin () {
      return this.$store.getters.token
    }
  },
  async created () {
    let res = await getProblemList()
    this.$store.commit('problem/setProblemInfo', res.data) // 缓存

    // 拉取用户解题信息
    if (this.isLogin) {
      res = await getUserSolvedList()
      this.$store.commit('user/setUserSolvedList', res.data)
    }
  },
  methods: {
    formatterLevel (row, column) {
      const level = row.level
      if (level === 1) {
        return '简单'
      } else if (level === 2) {
        return '中等'
      } else {
        return '困难'
      }
    },
    querySearchAsync (queryString, cb) {},
    handleSelect (item) {
      console.log(item)
    },
    isSolved (id) {
      // 判断用户是否已经AC了该题目
      const userSolvedList = this.$store.getters.userSolvedList
      if (userSolvedList.includes(id)) {
        return true
      } else {
        return false
      }
    }
  }
}
</script>

<style scoped>
.search-box .search-pro {
  float: left;
  margin-left: 20px;
  margin-top: 10px;
}
.search-box .select-tag {
  float: right;
  margin-right: 20px;
  margin-top: 10px;
}
.el-icon-check {
  color: #67c23a;
  font-size: 20px;
  margin-left: 30px;
}
</style>
