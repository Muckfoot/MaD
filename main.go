package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/tubbebubbe/transmission"
)

const pathToHD = "/Volumes/1 TB WD/"
const MB = 1000000.0

func main() {
	logf := initLogs()
	cfg := getCfg()

	client := transmission.New("http://127.0.0.1:9091", cfg.Login.Username,
		cfg.Login.Password)

	for {
		torrents, err := client.GetTorrents()
		if err != nil {
			log.Fatal(err)
		}

		for _, torrent := range torrents {
			if torrent.PercentDone == 1 {

				var stat syscall.Statfs_t

				err := syscall.Statfs(cfg.Paths.PathToHD, &stat)
				if err != nil {
					log.Printf("Unable to locate HD: %s\n", cfg.Paths.PathToHD)
					time.Sleep(time.Minute * 1)
					break
				}

				avbSpace := stat.Bavail * uint64(stat.Bsize)

				fmt.Println("Available space: ", avbSpace/MB, "MB")

				r, err := os.Open(torrent.DownloadDir + "/" + torrent.Name)
				checkErr(err)

				defer r.Close()

				fileStat, err := r.Stat()
				checkErr(err)

				fileSize := fileStat.Size()

				if i, err := os.Stat(cfg.Paths.PathToHD,
					+torrent.Name); err == nil && i.Size() == fileSize {
					fmt.Printf("File %s already exists in %s... skipping\n",
						torrent.Name, cfg.Paths.PathToHD)
					// ADD TORRENT REMOVAL
					continue
				}

				if avbSpace-uint64(fileSize) > 0 {

					w, err := os.Create(cfg.Paths.PathToHD, +torrent.Name)
					checkErr(err)
					defer w.Close()

					fmt.Printf("Copying %s to %s\n", torrent.Name, cfg.Paths.PathToHD)

					n, err := io.Copy(w, r)
					checkErr(err)

					fmt.Printf("copied %v Megabytes \n", n/MB)

					removeTorrent(torrent.ID, client)

				} else {
					_, _ = logf.WriteString("Not enough available space in the HD")
					continue
				}

			}
		}

		time.Sleep(time.Second * 5)
	}
}

func removeTorrent(id int, client transmission.TransmissionClient) bool {
	delCmd, err := transmission.NewDelCmd(id, true)
	if err != nil {
		checkErr(err)
		return false

	}
	_, err = client.ExecuteCommand(delCmd)
	if err != nil {
		checkErr(err)
		return false

	}
	return true

}
