export interface User {
    id: Number;
    username: String;
    nickname: String;
    email: String;
    user_pic: String;
}

// 注册表单参数
export interface RegisterParams {
    username: string; // 用户名
    password: string; // 密码
    repassword: string; // 确认密码
}
