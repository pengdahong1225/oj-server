<template>
  <div class="prolist">
    <!-- 搜索框 -->
    <el-autocomplete
      v-model="state"
      :fetch-suggestions="querySearchAsync"
      placeholder="搜索题目"
      @select="handleSelect"
      suffix-icon="el-icon-search"
      size="medium"
      class="search-pro"
    ></el-autocomplete>

    <!-- tag选择器 -->
    <el-select v-model="value" filterable placeholder="标签" class="select-tag">
      <el-option
        v-for="item in options"
        :key="item.value"
        :label="item.label"
        :value="item.value"
      >
      </el-option>
    </el-select>

    <!-- 题目列表 -->
    <el-table :data="problemList" style="width: 100%">
      <el-table-column label="#" prop="id" width="180" sortable align="center">
      </el-table-column>
      <el-table-column label="题目" prop="title" width="400"> </el-table-column>
      <el-table-column label="难度" prop="level" :formatter="formatterFunc">
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
export default {
  name: 'ProListPage',
  data () {
    return {
      state: ''
    }
  },
  computed: {
    problemList () {
      return this.$store.getters.problemList
    }
  },
  async created () {
    const res = await getProblemList()
    const obj = JSON.parse(res.data)
    this.$store.commit('problem/setProblemInfo', obj) // 缓存
  },
  methods: {
    formatterFunc (row, column) {
      const level = row.level
      if (level === 0) {
        return '简单'
      } else if (level === 1) {
        return '中等'
      } else {
        return '困难'
      }
    },
    querySearchAsync (queryString, cb) {},
    createStateFilter (queryString) {
      return (state) => {
        return (
          state.value.toLowerCase().indexOf(queryString.toLowerCase()) === 0
        )
      }
    },
    handleSelect (item) {
      console.log(item)
    }
  }
}
</script>

<style>
.search-pro {
  float: left;
  margin-left: 20px;
}
.select-tag {
  float: right;
  margin-right: 20px;
}
</style>
