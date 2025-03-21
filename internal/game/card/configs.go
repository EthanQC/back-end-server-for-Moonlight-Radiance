package card

import (
	"os"      // 操作系统层面，这里用于设置环境变量
	"testing" // 标准测试包，提供单元测试功能

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
)

// 返回类型为 func() 的函数，用于清理测试数据
func setupTestDB(t *testing.T) func() {
	// 设置测试数据库环境变量
	// 使用前请先创建数据库 moonlight，并修改密码
	os.Setenv("DB_DSN", "root:yourpassword@tcp(localhost:3306)/moonlight?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai")

	// 初始化数据库连接
	if err := common.InitDB(os.Getenv("DB_DSN")); err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 创建测试数据库表
	common.DB.Exec(`
	    CREATE TABLE IF NOT EXISTS cards (
	        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	        name VARCHAR(50) NOT NULL,
	        type TINYINT NOT NULL,
	        cost INT NOT NULL,
	        description VARCHAR(500) NOT NULL
	    )
	`)

	common.DB.Exec(`
	    CREATE TABLE IF NOT EXISTS PlayerCardState (
	        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	        game_id INT UNSIGNED NOT NULL,
	        player_id INT UNSIGNED NOT NULL,
			stage TINYINT NOT NULL DEFAULT 0,
	        hand_card_ids JSON,
	        deck_card_ids JSON,
	        discard_card_ids JSON,
	        hand_basic_count INT NOT NULL DEFAULT 0,
	        hand_skill_count INT NOT NULL DEFAULT 0,
	        deck_basic_count INT NOT NULL DEFAULT 0,
	        deck_skill_count INT NOT NULL DEFAULT 0,
	        basic_card_played BOOLEAN NOT NULL DEFAULT FALSE
	    )
	`)

	// 插入测试卡牌数据
	common.DB.Exec(`
	    INSERT INTO cards (name, type, cost, description) VALUES
	    ('新月', 1, 0, '基础月相牌'),
	    ('蛾眉月', 1, 0, '基础月相牌'),
	    ('上弦月', 1, 0, '基础月相牌'),
	    ('盈凸月', 1, 0, '基础月相牌'),
	    ('满月', 1, 0, '基础月相牌'),
	    ('亏凸月', 1, 0, '基础月相牌'),
	    ('下弦月', 1, 0, '基础月相牌'),
	    ('残月', 1, 0, '基础月相牌'),
	    ('测试功能牌1', 2, 2, '测试用功能牌'),
	    ('测试功能牌2', 2, 5, '测试用功能牌'),
		('测试功能牌3', 2, 4, '测试用功能牌'),
        ('测试功能牌4', 2, 3, '测试用功能牌')
	`)

	// 返回清理函数
	return func() {
		// 清理测试数据
		common.DB.Exec("TRUNCATE TABLE PlayerCardState")
		common.DB.Exec("TRUNCATE TABLE cards")
	}
}
