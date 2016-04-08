package server

import (
	"bytes"
	"fmt"
	"gopkg.in/redis.v3"
	"strings"
	"time"
)

var EXPIRETIME int64 = 5 * 60

func OpenClient() (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return
}

func RedisSetIPMember(data *announceData) (retval int) {
	c := data.redisClient

	keymember := concatenateKeyMember(data.info_hash, "ip")

	currTime := int64(time.Now().UTC().AddDate(0, 0, 1).Unix())

	ipPort := fmt.Sprintf("%s:%v", createIPPortPair(data), currTime)

	if err := c.SAdd(keymember, ipPort).Err(); err != nil {
		retval = 0
		panic("Failed to add key")

	} else {
		retval = 1
	}

	return
}

func RedisSetKeyVal(c *redis.Client, keymember string, value string) {
	// RedisSetKeyVal sets a key:member's value to value. Returns nothing as of
	// yet.
	currTime := int64(time.Now().UTC().Unix())
	currTime += EXPIRETIME
	value = fmt.Sprintf("%v:%v", value, currTime)

	if sz := strings.Split(value, ":"); len(sz) >= 1 {
		// If the value being added can be converted to an int, it is a ip:port key
		// and we can set an expiration on it.
		c.SAdd(keymember, value)
	}
}

func RedisGetKeyVal(client *redis.Client, key string, value *announceData) []string {
	// RedisGetKeyVal retrieves a value from the Redis store by looking up the
	// provided key. If the key does not yet exist, we create the key in the KV
	// storage or if the value is empty, we add the current requester to the
	// list.
	keymember := concatenateKeyMember(key, "complete")

	val, err := client.SMembers(keymember).Result()
	if err != nil {
		// Fail because the key doesn't exist in the KV storage.
		CreateNewTorrentKey(client, keymember)
	}

	return val
}

func RedisGetAllPeers(c *redis.Client, key string, data *announceData) []string {
	keymember := concatenateKeyMember(key, "complete")

	val, err := c.SMembers(keymember).Result()
	if err != nil {
		// Fail because the key doesn't exist in the KV storage.
		CreateNewTorrentKey(c, keymember)
	}

	keymember = concatenateKeyMember(key, "incomplete")

	val2, err := c.SMembers(keymember).Result()
	if err != nil {
		panic("Failed to get incomplete peers for")
	} else {
		val = append(val, val2...)
	}

	return val
}

func RedisGetCount(c *redis.Client, info_hash string, member string) (retval []string, err error) {
	// A generic function which is used to retrieve either the complete count
	// or the incomplete count for a specified `info_hash`.
	keymember := concatenateKeyMember(info_hash, member)

	retval, err = c.SMembers(keymember).Result()
	if err != nil {
		// TODO(ian): Add actual error checking here.
		err = fmt.Errorf("The info hash %s with member %s doesn't exist", info_hash, member)
	}

	return
}

func RedisGetBoolKeyVal(client *redis.Client, key string, value interface{}) bool {
	_, err := client.Get(key).Result()

	return err != nil
}

func RedisSetKeyIfNotExists(c *redis.Client, keymember string, value string) (rv bool) {
	rv = RedisGetBoolKeyVal(c, keymember, value)
	if !rv {
		RedisSetKeyVal(c, keymember, value)
	}
	return
}

func RedisRemoveKeysValue(c *redis.Client, key string, value string) {
	// Remove a `value` from `key` in the redis kv storage. `key` is typically
	// a keymember of info_hash:(in)complete and the value is typically the
	// ip:port concatenated.
	c.SRem(key, value)
}

func CreateNewTorrentKey(client *redis.Client, key string) {
	// CreateNewTorrentKey creates a new key. By default, it adds a member
	// ":ip". I don't think this ought to ever be generalized, as I just want
	// Redis to function in one specific way in notorious.
	client.SAdd(key, "complete", "incomplete")

}

func concatenateKeyMember(key string, member string) string {
	var buffer bytes.Buffer

	buffer.WriteString(key)
	buffer.WriteString(":")
	buffer.WriteString(member)

	return buffer.String()
}

func createIPPortPair(value *announceData) string {
	// createIPPortPair creates a string formatted ("%s:%s", value.ip,
	// value.port) looking like so: "127.0.0.1:6886" and returns this value.
	return fmt.Sprintf("%s:%s", value.ip, value.port)
}
