<script lang="ts" setup>
import { SuccessFilled } from '@element-plus/icons-vue'
</script>

<template>
    <div class="container">
        <div class="result-container">
            <el-descriptions size="large">
                <template #title>
                    <div class="icon-box">
                        <el-icon color="#19be6b">
                            <SuccessFilled />
                        </el-icon>
                        Accepted
                    </div>
                </template>
                <el-descriptions-item label="Username">kooriookami</el-descriptions-item>
                <el-descriptions-item label="Telephone">18100000000</el-descriptions-item>
                <el-descriptions-item label="Place">Suzhou</el-descriptions-item>
                <el-descriptions-item label="Remarks">
                    <el-tag size="small">School</el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="Address">
                    No.1188, Wuzhong Avenue, Wuzhong District, Suzhou, Jiangsu Province
                </el-descriptions-item>
            </el-descriptions>
        </div>

        <div class="code-container">
            source code
        </div>
    </div>
</template>

<style lang="scss" scoped>
.container {
    display: block;
    justify-content: space-between;
    width: 70%;
    margin: auto;

    .result-container {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-top: 30px;
        background-color: #e8f9f0;
        border: 1px solid #d1f2e1;

        .icon-box {
            justify-content: space-between;
            font-size: 35px;
            display: flex;
            align-items: center;
            margin: auto;
        }
    }

    .code-container {
        margin-top: 30px;
        border: 1px solid #00ffea;
    }
}
</style>