namespace API {
    // 用户信息（响应 protojson camelCase；uid/mobile 为 int64，读取处注意）
    type UserInfo = {
        uid: number;
        createAt?: number;
        deleteAt?: number;
        mobile: string;
        nickname: string;
        email?: string;
        gender?: number;
        role?: number;
        avatarUrl?: string;
    }
    // 登录表单（请求体）
    type LoginForm = {
        mobile: string;
        password: string;
    }
    // 注册表单（请求体）
    type RegisterForm = {
        mobile: string;
        password: string;
        repassword: string;
    }
    // 查询用户列表参数（请求 query）
    type QueryUserListParams = {
        page: number;
        page_size: number;
        keyword: string;
    }
}
