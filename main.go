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
const MB = 1000000

func main() {
	logf, err := os.OpenFile("errors.log",
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer logf.Close()
	log.SetOutput(logf)

	client := transmission.New("http://127.0.0.1:9091", "admin", "admin")
	for {
		torrents, err := client.GetTorrents()
		if err != nil {
			log.Fatal(err)
		}

		for _, torrent := range torrents {
			if torrent.PercentDone == 1 {

				// tNFO.dir = torrent.DownloadDir
				// tNFO.name = torrent.Name
				// tNFO.id = torrent.ID
				// torrentList = append(torrentList, tNFO)V

				var stat syscall.Statfs_t

				err := syscall.Statfs(pathToHD, &stat)
				if err != nil {
					// fmt.Println(err)
					log.Printf("Unable to locate HD: %s\n", pathToHD)
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

				if i, err := os.Stat(pathToHD + torrent.Name); err == nil && i.Size() == fileSize {
					fmt.Printf("File %s already exists in %s... skipping\n", torrent.Name, pathToHD)
					// ADD TORRENT REMOVAL
					continue
				}

				if avbSpace-uint64(fileSize) > 0 {

					w, err := os.Create(pathToHD + torrent.Name)
					checkErr(err)
					defer w.Close()

					fmt.Printf("Copying %s to %s\n", torrent.Name, pathToHD)

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
		// for _, file := range torrentList {
		// fmt.Println(torrentList)
		// errorVal := ""

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
