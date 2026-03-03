package main

import "elearning-api/cmd"

func main() {
	server := cmd.ApiServer{}
	server.Run()
}
