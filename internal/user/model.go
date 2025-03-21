// 用户数据结构
// 定义用户相关的数据结构，例如用户ID、登录状态、权限、卡牌状态等
// JSON（JavaScript Object Notation）：一种轻量级数据交换格式，格式为键值对

package user

// User 用户数据类型
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` //管理员或玩家，默认玩家
}

// RegisterInput 注册请求输入结构体
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginInput 登录请求输入结构体
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
