package bencode

import (
    "bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

// EncodeInt encodes an int to a bencoded integer: "i<x>e"
func EncodeInt(x int) string {
	return fmt.Sprintf("i%de", x)
}

// EncodeList encodes a list of items (`items`) to a bencoded list:
// "l<item1><item2><...><itemN>e"
func EncodeList(items []string) string {
	tmp := "l"
	for i := range items {
		if items[i][0] == 'l' || items[i][0] == 'd' {
			tmp = writeStringData(tmp, items[i])
		} else {
			tmp = writeStringData(tmp, EncodeByteString(items[i]))
		}
	}
	tmp = writeStringData(tmp, "e")
	return tmp
}

// EncodeDictionary Takes a list of bencoded KVpairs and return a bencoded dictionary.
func EncodeDictionary(kvpairs []string) (retdict string) {

	retdict = "d"
	for i := range kvpairs {
		retdict += kvpairs[i]
	}
	retdict += "e"

	return
}

// EncodeByteString Encodes a string to <key length>:<key>
func EncodeByteString(key string) string {
	return fmt.Sprintf("%d:%s", utf8.RuneCountInString(key), key)
}

// EncodePeerList Handles peer list creation for non-compact responses. Mostly deprecated
// for most torrent clients nowadays as compact is the default. Returns a
// bencoded list of bencoded dictionaries containing "peer id", "ip",
// "port": "ld7:peer id20:<peer id>2:ip9:<127.0.0.1>4:port4:7878ee"
// peers contains a ip:port
func EncodePeerList(peers []string) (retlist string) {
	var tmpDict []string

	for i := range peers {
		var tmp []string
		peerSplit := strings.Split(peers[i], ":")

		// TODO(ian): Figure out an actual way to do peer id.
		tmp = append(tmp, EncodeKV("peer id", "11111111111111111111"))
		tmp = append(tmp, EncodeKV("ip", peerSplit[0]))
		tmp = append(tmp, EncodeKV("port", peerSplit[1]))

		tmpDict = append(tmpDict, EncodeDictionary(tmp))
	}

	peerList := EncodeList(tmpDict)
	peerList = EncodeKV("peers", peerList)
	retlist = fmt.Sprintf("d%se", peerList)

	return
}

// EncodeKV Encodes a KV pair into a string
func EncodeKV(key string, value string) string {
	key = EncodeByteString(key)
	if value[0] == 'i' || value[0] == 'l' || value[0] == 'd' {
		value = value
	} else {
		value = EncodeByteString(value)
	}
	return writeStringData(key, value)
}

// writeStringData is used to concatenate two strings. This is a heavily used
// function throughout the bencode section of the codebase. It's inherently
// naive and I just want it to combine two strings. You'll se some places where
// we use Sprintf still, but that's because I don't feel the need to adding
// padding to this function.
func writeStringData(val1 string, val2 string) string {
    var buffer bytes.Buffer

    buffer.WriteString(val1)
    buffer.WriteString(val2)

    return buffer.String()
}
