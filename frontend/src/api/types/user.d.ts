namespace API {
    // 用户信息（响应 protojson camelCase；uid/mobile 为 int64，读取处注意）
    type UserInfo = {
        uid: number;
        createAt?: number;
        deleteAt?: number;
        mobile: string;
        username: string;
        email?: string;
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
    // 创建用户表单（请求体）
    type CreateUserForm = {
        mobile: string;
        password: string;
        username: string;
        email?: string;
        role?: number;
    }
}
