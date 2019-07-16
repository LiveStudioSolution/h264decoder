package logger

import (
	"io/ioutil"
	"log"
	"os"
)

//var Log = log.New(os.Stderr, "", log.LstdFlags)
var Log = log.New(ioutil.Discard,"",log.LstdFlags)

func init(){

}

func SetLogger(logger *log.Logger){
	Log = logger
}

func EnableStderrLog(){
	Log = log.New(os.Stderr, "", log.LstdFlags)
}


