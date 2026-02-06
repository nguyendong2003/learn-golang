package main

import (
	"bulk-upload-csv/cmd"
)

func main() {
	server := cmd.ApiServer{}
	server.Run()
}
