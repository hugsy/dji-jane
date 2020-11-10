package main

import (
	"dji-jane"
	"dji-joe"

	"flag"
)


var listenAddress = flag.String("listen", "0.0.0.0:8000", "Specify the address/port to listen to")

func main(){
	flag.Parse()
	djijane.Log = djijoe.InitLogger(djijane.PROGNAME)
	djijane.Log.InfoF("Starting %s [%s]", djijane.PROGNAME, djijane.VERSION)
	djijane.StartServer(*listenAddress)
	djijane.Log.InfoF("Leaving %s", djijane.PROGNAME)
}
