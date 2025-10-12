package main

import "flag"

var runDaemon = flag.Bool("daemon", false, "Run the daemon")
var port = flag.Int("port", 3719, "The port to use")
var rootURI = flag.String("root-uri", "http://localhost:3719/", "The Root URI to use for minified URLs")

func main() {
	flag.Parse()
	if *runDaemon {
		RunDaemon(*port)
	} else {
		RunFilter(*rootURI)
	}
}
