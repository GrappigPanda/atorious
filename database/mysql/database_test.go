package mysql

import (
	"github.com/GrappigPanda/notorious/config"
	. "github.com/GrappigPanda/notorious/database/schemas"
	"os"
	"testing"
	"time"
)

var CONFIG = config.ConfigStruct{
	"mysql",
	"localhost",
	"3306",
	"testuser",
	"testuser",
	"testdb",
	false,
	nil,
}

var DBCONN, ERR = OpenConnectionWithConfig(&CONFIG)

func TestOpenConnPostgres(t *testing.T) {
	if ERR != nil {
		t.Fatalf("Unable to connect %v", ERR)
	}
	InitDB(DBCONN)
}

func TestAddWhitelistedTorrent(t *testing.T) {
	newTorrent := &WhiteTorrent{
		InfoHash:  "12345123451234512345",
		Name:      "Hello Kitty Island Adventure.exe",
		AddedBy:   "127.0.0.2",
		DateAdded: time.Now().Unix(),
	}

	if !newTorrent.AddWhitelistedTorrent(DBCONN) {
		t.Fatalf("Failed to Add a whitelisted torrent")
	}
}

func TestGetWhitelistedTorrents(t *testing.T) {
	newTorrent := &WhiteTorrent{
		InfoHash:  "12345123GetWhitelistedTorrents",
		Name:      "Hello Kitty Island Adventure3.exe",
		AddedBy:   "127.0.0.2",
		DateAdded: time.Now().Unix(),
	}

	newTorrent2 := &WhiteTorrent{
		InfoHash:  "FFFFFFFFFFFFhitelistedTorrents",
		Name:      "Hello Kitty Island Adventure4.exe",
		AddedBy:   "127.0.0.2",
		DateAdded: time.Now().Unix(),
	}

	newTorrent.AddWhitelistedTorrent(DBCONN)
	newTorrent2.AddWhitelistedTorrent(DBCONN)

	_, err := GetWhitelistedTorrents(DBCONN)
	if err != nil {
		t.Fatalf("Failed to get all whitelisted torrents: %v", err)
	}
}

func TestGetWhitelistedTorrent(t *testing.T) {
	newTorrent := &WhiteTorrent{
		InfoHash:  "12345123GetWhitelistedTorrent",
		Name:      "Hello Kitty Island Adventure2.exe",
		AddedBy:   "127.0.0.2",
		DateAdded: time.Now().Unix(),
	}

	newTorrent.AddWhitelistedTorrent(DBCONN)

	retval, err := GetWhitelistedTorrent(nil, newTorrent.InfoHash)
	if err != nil {
		t.Fatalf("Failed to GetWhitelistedTorrent: %v", err)
	}

	if retval.InfoHash != newTorrent.InfoHash {
		t.Fatalf("Expected %v, got %v", retval.InfoHash,
			newTorrent.InfoHash)
	}
}

func TestUpdateStats(t *testing.T) {
	expectedReturn := &TrackerStats{
		Downloaded: 6,
		Uploaded:   21,
	}

	newStats := &TrackerStats{
		Downloaded: 1,
		Uploaded:   1,
	}
	DBCONN.Save(&newStats)

	UpdateStats(nil, 20, 5)

	retval := &TrackerStats{}
	DBCONN.First(&retval)
	if retval.Downloaded != expectedReturn.Downloaded {
		t.Fatalf("Expected %v, got %v",
			expectedReturn.Downloaded,
			retval.Downloaded)
	}

	if retval.Uploaded != expectedReturn.Uploaded {
		t.Fatalf("Expected %v, got %v",
			expectedReturn.Uploaded,
			retval.Uploaded)
	}
}

func TestUpdatePeerStats(t *testing.T) {
	expectedReturn := &PeerStats{
		Downloaded: 6,
		Uploaded:   21,
		Ip:         "127.0.0.2",
	}

	newPeer := &PeerStats{
		Downloaded: 1,
		Uploaded:   1,
		Ip:         "127.0.0.2",
	}

	DBCONN.Save(&newPeer)

	UpdatePeerStats(nil, 20, 5, "127.0.0.2")

	retval := &PeerStats{}
	DBCONN.Where("Ip = ?", "127.0.0.2").Find(&retval)

	if retval.Downloaded != expectedReturn.Downloaded {
		t.Fatalf("Expected %v, got %v",
			expectedReturn.Downloaded,
			retval.Downloaded)
	}

	if retval.Uploaded != expectedReturn.Uploaded {
		t.Fatalf("Expected %v, got %v",
			expectedReturn.Uploaded,
			retval.Uploaded)
	}

	if retval.Ip != expectedReturn.Ip {
		t.Fatalf("Expected %v, got %v",
			expectedReturn.Ip,
			retval.Ip)
	}
}

func TestMain(m *testing.M) {
	dbConn, _ := OpenConnectionWithConfig(&CONFIG)
	dbConn.DropTableIfExists(
		&PeerStats{},
		&Torrent{},
		&TrackerStats{},
	)
	InitDB(dbConn)
	os.Exit(m.Run())
}
