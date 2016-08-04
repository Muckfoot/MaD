package main

import "github.com/tubbebubbe/transmission"

func main() {

	logf := initLogs()
	cfg := getCfg()

	client := transmission.New("http://127.0.0.1:9091", cfg.Login.Username,
		cfg.Login.Password)
	copyPaste(client, cfg, logf)

}
