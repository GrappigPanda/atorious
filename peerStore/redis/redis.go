package redisPeerStore

import (
	"bytes"
	"fmt"
	"gopkg.in/redis.v3"
	"strings"
	"time"
)

// EXPIRETIME signifies how long a peer will live under the specified info_hash
// until the reaper removes it.
var EXPIRETIME int64 = 5 * 60

// OpenClient opens a connection to redis.
func OpenClient() (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return
}

// SetIPMember sets a key as a member of an infohash and sets a timeout.
func SetIPMember(c *redis.Client, infoHash, ipPort string) (retval int) {
	if c == nil {
		c = OpenClient()
	}

	keymember := concatenateKeyMember(infoHash, "ip")

	currTime := int64(time.Now().UTC().AddDate(0, 0, 1).Unix())

	key := fmt.Sprintf("%s:%v", ipPort, currTime)

	if err := c.SAdd(keymember, key).Err(); err != nil {
		retval = 0
		panic("Failed to add key")

	} else {
		retval = 1
	}

	return
}

// SetKeyVal Sets a key to the specified value. Used mostly with adding a
// peer into an info_hash
func SetKeyVal(c *redis.Client, keymember string, value string) {
	if c == nil {
		c = OpenClient()
	}

	// SetKeyVal sets a key:member's value to value. Returns nothing as of
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

// GetKeyVal Lookup a peer in the specified infohash at `key`
func GetKeyVal(c *redis.Client, key string) []string {
	if c == nil {
		c = OpenClient()
	}

	// GetKeyVal retrieves a value from the  store by looking up the
	// provided key. If the key does not yet exist, we Create the key in the KV
	// storage or if the value is empty, we add the current requester to the
	// list.
	keymember := concatenateKeyMember(key, "complete")

	val, err := c.SMembers(keymember).Result()
	if err != nil {
		// Fail because the key doesn't exist in the KV storage.
		CreateNewTorrentKey(c, keymember)
	}

	return val
}

// GetAllPeers fetches all peers from the info_hash at `key`
func GetAllPeers(c *redis.Client, key string) []string {
	if c == nil {
		c = OpenClient()
	}

	keymember := concatenateKeyMember(key, "complete")

	val, err := c.SRandMemberN(keymember, 30).Result()
	if err != nil {
		// Fail because the key doesn't exist in the KV storage.
		CreateNewTorrentKey(c, keymember)
	}

	if len(val) == 30 {
		return val
	}

	keymember = concatenateKeyMember(key, "incomplete")

	val2, err := c.SRandMemberN(keymember, int64(30-len(val))).Result()
	if err != nil {
		panic("Failed to get incomplete peers for")
	} else {
		val = append(val, val2...)
	}

	return val
}

// GetCount counts all of the peers at `info_hash`
func GetCount(c *redis.Client, infoHash string, member string) (retval int, err error) {
	if c == nil {
		c = OpenClient()
	}

	// A generic function which is used to retrieve either the complete count
	// or the incomplete count for a specified `info_hash`.
	keymember := concatenateKeyMember(infoHash, member)

	x, err := c.SMembers(keymember).Result()
	if err != nil {
		// TODO(ian): Add actual error checking here.
		err = fmt.Errorf("The info hash %s with member %s doesn't exist", infoHash, member)
	}

	retval = len(x)
	return
}

// GetBoolKeyVal Checks if a `key` exists
func GetBoolKeyVal(c *redis.Client, key string) bool {
	if c == nil {
		c = OpenClient()
	}

	ret, _ := c.Exists(key).Result()

	return ret
}

// SetKeyIfNotExists Set a key if it doesn't exist.
func SetKeyIfNotExists(c *redis.Client, keymember string, value string) (rv bool) {
	if c == nil {
		c = OpenClient()
	}

	rv = GetBoolKeyVal(c, keymember)
	if !rv {
		SetKeyVal(c, keymember, value)
	}
	return
}

// RemoveKeysValue Remove a `value` from `key` in the redis kv storage. `key` is typically
// a keymember of info_hash:(in)complete and the value is typically the
// ip:port concatenated.
func RemoveKeysValue(c *redis.Client, key string, value string) {
	if c == nil {
		c = OpenClient()
	}

	c.SRem(key, value)
}

// CreateNewTorrentKey Creates a new key. By default, it adds a member
// ":ip". I don't think this ought to ever be generalized, as I just want
//  to function in one specific way in notorious.
func CreateNewTorrentKey(c *redis.Client, key string) {
	if c == nil {
		c = OpenClient()
	}

	c.SAdd(key, "complete", "incomplete")

}

// concatenateKeyMember concatenates the key and the member delimited by the
// character ":"
func concatenateKeyMember(key string, member string) string {
	var buffer bytes.Buffer

	buffer.WriteString(key)
	buffer.WriteString(":")
	buffer.WriteString(member)

	return buffer.String()
}

// CreateIPPortPair Creates a string formatted ("%s:%s", value.ip,
// value.port) looking like so: "127.0.0.1:6886" and returns this value.
func createIPPortPair(ip, port string) string {
	return fmt.Sprintf("%v:%v", ip, port)
}
