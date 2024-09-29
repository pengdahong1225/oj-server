<template>
  <div class="test-case">
    <div class="in-box">
      <div class="title-box">Sample-input</div>
      <input class="content-box" type="text" :value="input" readonly />
    </div>

    <div class="out-box">
      <div class="title-box">Sample-output</div>
      <input class="content-box" type="text" :value="output" readonly />
    </div>
  </div>
</template>

<script>
export default {
  name: 'TestCase',
  props: {
    input: {
      type: String,
      default: '',
    },
    output: {
      type: String,
      default: '',
    },
  },
}
</script>

<style lang="less" scoped>
.test-case {
  display: flex;
  width: 90%;
  .in-box {
    width: 50%;
    .content-box {
      height: 25px;
      font-size: 15px;
    }
  }
  .out-box {
    width: 50%;
    .content-box {
      height: 25px;
      font-size: 15px;
    }
  }
  .content-box {
    width: 80%;
    border: 1px solid rgba(143, 143, 143, 0.5);
  }
}
</style>