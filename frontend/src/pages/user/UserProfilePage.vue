<script lang="ts" setup>
import { onMounted } from 'vue'
import { User } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores'
import { getUserInfoService } from '@/api/userController'
import avatar from '@/assets/default.png'

const userStore = useUserStore()
const rules = {}

onMounted(() => {
    getUserInfo()
})
const getUserInfo = async () => {
    const res = await getUserInfoService(userStore.userInfo.uid)
    console.log(res)
    // userStore.setUserInfo(res.data.data)
}

const onUpdateProfile = () => {
    console.log('onUpdateProfile')
}
const onUpdateAvator = () => {
    console.log('onUpdateAvator')
}

</script>

<template>
    <div class="container">
        <!-- 信息区域 -->
        <div class="user-profile">
            <el-card>
                <template #header>
                    <div style="text-align: left;"><strong>Your Profile</strong></div>
                </template>
                <!-- 表单 -->
                <el-form :model="userStore.userInfo" :rules="rules" label-width="120px">
                    <el-form-item label="ID"  prop="uid" style="width: 500px;">
                        <el-input v-model="userStore.userInfo.uid" disabled />
                    </el-form-item>
                    <el-form-item label="NickName" prop="nickname" style="width: 500px;">
                        <el-input v-model="userStore.userInfo.nick_name" />
                    </el-form-item>
                    <el-form-item label="Mobile" prop="mobile" style="width: 500px;">
                        <el-input v-model="userStore.userInfo.mobile" />
                    </el-form-item>
                    <el-form-item label="Email" prop="email" style="width: 500px;">
                        <el-input v-model="userStore.userInfo.email" />
                    </el-form-item>
                    <el-form-item label="Role" prop="role" style="width: 500px;">
                        <template #default>
                            <el-tag v-if="userStore.userInfo.role === 0" type="primary">user</el-tag>
                            <el-tag v-else-if="userStore.userInfo.role === 1" type="warning">admin</el-tag>
                            <el-tag v-else type="danger">nil</el-tag>
                        </template>
                    </el-form-item>
                    <el-form-item label="Gender" prop="gender" style="width: 500px;">
                        <el-radio-group v-model="userStore.userInfo.gender">
                            <el-radio :value ="1">male</el-radio>
                            <el-radio :value ="0">female</el-radio>
                        </el-radio-group>
                    </el-form-item>
                    <el-button type="primary" @click="onUpdateProfile">
                        update
                    </el-button>
                </el-form>
            </el-card>
        </div>
        <!-- 头像区域 -->
        <div class="user-avatar" style="text-align: center">
            <el-card>
                <template #header>
                    <div style="text-align: left;"><strong>Avatar</strong></div>
                </template>
                <el-avatar shape="circle" :size="90" :src="userStore.userInfo.avatar_url || avatar"
                    style="margin: 10px;" /><br>
                <el-button type="primary" @click="onUpdateAvator">
                    update
                </el-button>
            </el-card>
        </div>
    </div>
</template>

<style lang="less" scoped>
.container {
    display: flex;
    width: 65%;
    text-align: center;
    margin: auto;

    .user-profile {
        width: 65%;
        margin-right: 20px;

        .el-card {
            &:hover {
                box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
            }
        }
    }

    .user-avatar {
        width: 35%;
        margin-left: 20px;

        .el-card {
            &:hover {
                box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
            }
        }
    }
}
</style>