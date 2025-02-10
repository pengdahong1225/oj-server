namespace API {
    // 用户信息
    type UserProfile = {
        uid: number;
        mobile?: number;
        nickname?: string;
        email?: string;
        gender?: number;
        role?: number;
        avatar_url?: string;
    }

    // 登录表单
    type LoginForm = {
        mobile: string; // 用户名
        password: string; // 密码
    }

    // 注册表单
    type RegisterForm = {
        mobile: string; // 用户名
        password: string; // 密码
        repassword: string; // 确认密码
    }
}
