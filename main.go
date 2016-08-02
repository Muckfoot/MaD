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

type torrentInfo struct {
	name string
	dir  string
	id   int
}

const pathToHD = "/Volumes/1 TB WD/"
const MB = 1000000

var (
	stat syscall.Statfs_t
	tNFO torrentInfo
)

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
		var torrentList []torrentInfo
		torrents, err := client.GetTorrents()
		checkErr(err)
		for _, torrent := range torrents {
			if torrent.PercentDone == 1 {
				tNFO.dir = torrent.DownloadDir
				tNFO.name = torrent.Name
				tNFO.id = torrent.ID
				torrentList = append(torrentList, tNFO)

			}
		}

		// if _, err := os.Stat(pathToHD); os.IsNotExist(err) {
		// _, err = logf.WriteString("HD not found")
		// continue
		// }

		for _, file := range torrentList {
			// fmt.Println(torrentList)
			// errorVal := ""

			err := syscall.Statfs(pathToHD, &stat)
			if err != nil {
				// fmt.Println(err)
				logf.WriteString("Unable to locate HD: " + pathToHD)
				time.Sleep(time.Minute * 1)
				break
			}
			avbSpace := stat.Bavail * uint64(stat.Bsize)

			fmt.Println("Available space: ", stat.Bavail*uint64(stat.Bsize)/MB, "MB")

			if _, err := os.Stat(pathToHD + file.name); err == nil {
				fmt.Printf("File %s already exists in %s... skipping\n", file.name, pathToHD)
				continue
			}

			r, err := os.Open(file.dir + "/" + file.name)
			checkErr(err)

			defer r.Close()

			fileStat, err := r.Stat()
			checkErr(err)

			fileSize := fileStat.Size()

			if avbSpace-uint64(fileSize) > 0 {

				w, err := os.Create(pathToHD + file.name)
				checkErr(err)
				defer w.Close()

				fmt.Printf("Copying %s to %s\n", file.name, pathToHD)

				n, err := io.Copy(w, r)
				checkErr(err)

				fmt.Printf("copied %v Megabytes \n", n/MB)

				/* delCmd, err := transmission.NewDelCmd(file.id, true) */
				// checkErr(err)

				// _, err = client.ExecuteCommand(delCmd)
				/* checkErr(err) */

			} else {
				_, _ = logf.WriteString("Not enough available space in the HD")
				continue
			}

		}
		time.Sleep(time.Second * 5)
	}

}
