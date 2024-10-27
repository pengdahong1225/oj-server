
<script lang="ts" setup>
import {
  Iphone,
  Location,
  OfficeBuilding,
  Tickets,
  User,
} from '@element-plus/icons-vue'
</script>

<template>
    <div class="container">
        <el-card class="user-profile" shadow="hover">
            <!-- 头像 -->
            <template #header>
                <div style="text-align: center">
                    <el-avatar shape="circle" :size="90" :src="squareUrl" />
                </div>
            </template>

            <!-- 信息 -->
            <el-descriptions class="margin-top" :column="2" border>
                <el-descriptions-item>
                    <template #label>
                        <div class="cell-item">
                            <el-icon :style="iconStyle">
                                <user />
                            </el-icon>
                            Username
                        </div>
                    </template>
                    kooriookami
                </el-descriptions-item>
                <el-descriptions-item>
                    <template #label>
                        <div class="cell-item">
                            <el-icon :style="iconStyle">
                                <iphone />
                            </el-icon>
                            Telephone
                        </div>
                    </template>
                    18100000000
                </el-descriptions-item>
                <el-descriptions-item>
                    <template #label>
                        <div class="cell-item">
                            <el-icon :style="iconStyle">
                                <location />
                            </el-icon>
                            Place
                        </div>
                    </template>
                    Suzhou
                </el-descriptions-item>
                <el-descriptions-item>
                    <template #label>
                        <div class="cell-item">
                            <el-icon :style="iconStyle">
                                <tickets />
                            </el-icon>
                            Remarks
                        </div>
                    </template>
                    <el-tag size="small">School</el-tag>
                </el-descriptions-item>
                <el-descriptions-item>
                    <template #label>
                        <div class="cell-item">
                            <el-icon :style="iconStyle">
                                <office-building />
                            </el-icon>
                            Address
                        </div>
                    </template>
                    No.1188, Wuzhong Avenue, Wuzhong District, Suzhou, Jiangsu Province
                </el-descriptions-item>
            </el-descriptions>
        </el-card>
    </div>
</template>

<style lang="scss" scoped>
.container {
    display: flex;
    width: 100%;

    .user-profile {
        width: 100%;
        margin: 20px
    }
}
</style>