/*
 * @Date: 2022-10-24 18:11:08
 */
package constants

import "fmt"

var GoCachePrefix = "go-cache-prefix"

const (
	OrderFmtCacheKey = "order-fmt-cache-key"
)

func generateBaseKey(key string) string {
	return fmt.Sprintf("%s-%s", GoCachePrefix, key)
}

func GetOrderFmtCacheKey(orderId int64) string {
	return fmt.Sprintf("%s-%d", generateBaseKey(OrderFmtCacheKey), orderId)
}
