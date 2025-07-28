package redis

import (
	"MKICS/backend/constant"
	"MKICS/backend/global"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func AcquireRedisLockWithRetry(ctx context.Context, key string, ttl time.Duration, retryInterval time.Duration, maxWait time.Duration) (lockValue string, ok bool) {
	lockValue = uuid.NewString()
	start := time.Now()

	for {
		ok = global.RDS.SetNX(ctx, key, lockValue, ttl).Val()
		if ok {
			return lockValue, true
		}
		if time.Since(start) > maxWait {
			return "", false
		}
		time.Sleep(retryInterval)
	}
}

func ReleaseRedisLock(ctx context.Context, key string, lockValue string) {
	luaScript := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	global.RDS.Eval(ctx, luaScript, []string{key}, lockValue)
}

func SetToken(uuid, ip, token string) error {
	ctx := context.Background()
	key := buildTokenKey(uuid, ip)
	jwtExpire := global.CONF.AuthConfig.JwtExpired
	return global.RDS.Set(ctx, key, token, jwtExpire).Err()
}

func GetToken(uuid, ip string) (string, error) {
	ctx := context.Background()
	key := buildTokenKey(uuid, ip)
	return global.RDS.Get(ctx, key).Result()
}

func DelToken(uuid, ip string) error {
	ctx := context.Background()
	key := buildTokenKey(uuid, ip)
	return global.RDS.Del(ctx, key).Err()
}

func ValidateToken(uuid, ip, token string) bool {
	val, err := GetToken(uuid, ip)
	return err == nil && val == token
}

func buildTokenKey(uuid, ip string) string {
	return fmt.Sprintf("%s:%s:%s", constant.KeyUserTokenPrefix, uuid, ip)
}
