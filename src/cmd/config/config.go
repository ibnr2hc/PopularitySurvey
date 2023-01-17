package config

import "time"

// SHOW_RANKING_LIMIT ランキングに表示するフォロワー数上限
const SHOW_RANKING_LIMIT = 30

// RATE_LIMIT_WAITING_SECOND TWITTER APIレート制限時の待機秒
const RATE_LIMIT_WAITING_SECOND = time.Second * 120

// RECONNECT_HTTP_WAITING_SECOND HTTPコネクトが切断された時の再接続待機秒
const RECONNECT_HTTP_WAITING_SECOND = time.Second * 60
