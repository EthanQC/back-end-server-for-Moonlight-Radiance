// 卡牌和地图逻辑

package game

import (
	"math/rand"
	"time"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/user"
)

// CardEngine 卡牌机制引擎
func CardEngine(player *user.User) string {
	// 假设这里有一些复杂的卡牌逻辑
	rand.Seed(time.Now().Unix())
	cards := []string{"Sun", "Moon", "Star", "Comet"}
	return cards[rand.Intn(len(cards))]
}
