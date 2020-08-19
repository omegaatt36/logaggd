package logaggd

import (
	"context"
	"fmt"
	"time"

	"github.com/omegaatt36/logaggd/cache"

	"github.com/go-redis/redis/v7"
)

const messageListKey = "logagg_list"
const messageCountKey = "logagg_count"
const messageKey = "logagg_message"

func record(ctx context.Context, b []byte) error {
	return cache.Redis().WatchContext(ctx, func(tx *redis.Tx) (err error) {
		var count int64
		if count, err = tx.Incr(messageCountKey).Result(); err != nil {
			return err
		}
		key := fmt.Sprintf("%s:%d", messageKey, count)
		if err = tx.Set(key, string(b), time.Second*1).Err(); err != nil {
			return err
		}
		return tx.RPush(messageListKey, key).Err()
	})
}

func deleteExpiredMessages(ctx context.Context) error {
	r := cache.Redis()
	done := false
	for !done {
		var keys []string
		err := r.LRange(messageListKey, 0, 10).ScanSlice(&keys)
		if err != nil {
			return err
		}
		if len(keys) == 0 {
			return nil
		}
		sliceCount := len(keys)
		for _, key := range keys {
			exist, err := r.Exists(key).Result()
			if err != nil {
				return err
			}
			sliceCount -= int(exist)
		}
		if sliceCount != len(keys) {
			done = true
		}
		for i := 0; i < sliceCount; i++ {
			err = r.LPop(messageListKey).Err()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
