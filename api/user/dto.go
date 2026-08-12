package user

// 发送短信请求
type SendSmsReq struct {
	Phone string `json:"phone" binding:"required,len=11"`
}

// 注册请求
type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone" binding:"required,len=11"`
	Code     string `json:"code" binding:"required,len=6"`
}

// 登录请求
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录响应
type LoginResp struct {
	Token  string `json:"token"`
	UserID uint64 `json:"user_id"`
}
