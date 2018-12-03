# go-deluge

Basic client library to interact with Deluge

```
client := New("http://awesome.com", "p4Ssw0Rd")

// Add a torrent given the URL
id, err := client.AddTorrent("http://releases.ubuntu.com/17.10/ubuntu-17.10.1-desktop-amd64.iso.torrent",
                        "/deluge/incomplete/path",
                        "/deluge/completed/path")

// Get information for a torrent by id
torrent, err := client.Torrent(id)

// Get information for all torrents
torrents, err := client.Torrents()

// Pause a torrent
err := client.PauseTorrent(id)

// Remove a torrent
removed, err := client.RemoveTorrent(id, removeData)

// Move a torrent's storage location
err := client.MoveStorage(id, "/deluge/new/completed/path")
```
