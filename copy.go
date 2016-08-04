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

const MB = 1000000.0

func copyPaste(client transmission.TransmissionClient, cfg Configuration, logf *os.File) {

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

				if i, err := os.Stat(cfg.Paths.PathToHD + torrent.Name); err == nil && i.Size() == fileSize {
					fmt.Printf("File %s already exists in %s... skipping\n",
						torrent.Name, cfg.Paths.PathToHD)

					//removeTorrent(torrent, client)
					continue
				}

				if avbSpace-uint64(fileSize) > 0 {

					w, err := os.Create(cfg.Paths.PathToHD + torrent.Name)
					checkErr(err)
					defer w.Close()

					fmt.Printf("Copying %s to %s\n", torrent.Name, cfg.Paths.PathToHD)

					n, err := io.Copy(w, r)
					checkErr(err)

					if i, _ := os.Stat(cfg.Paths.PathToHD + torrent.Name); i.Size() == fileSize {
						fmt.Printf("copied %v Megabytes \n", n/MB)

						time.Sleep(time.Second * 2)

						//	removeTorrent(torrent, client)

					} else {
						log.Printf("Failed to copy torrent %s\n, size missmatch", torrent.Name)
					}

				} else {
					log.Printf("Not enough available space in the HD")
				}

			}
		}

		time.Sleep(time.Second * 5)
	}
}

func removeTorrent(torrent transmission.Torrent, client transmission.TransmissionClient) {
	delCmd, err := transmission.NewDelCmd(torrent.ID, true)
	if err != nil {
		checkErr(err)
		log.Printf("Failed to remove torrent %s", torrent.Name)
	}
	_, err = client.ExecuteCommand(delCmd)
	if err != nil {
		checkErr(err)
		log.Printf("Failed to remove torrent %s", torrent.Name)
	}
}
