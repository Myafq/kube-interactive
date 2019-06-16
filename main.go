package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//go build -ldflags "-X main.hash=$(date -j -f "%a %b %d %T %Z %Y" "`date`" "+%s")"
var hash string = "123456789"

func main() {
	configMap := flag.NewFlagSet("ConfigMap", flag.ExitOnError)
	configMapTask := configMap.String("t", "Zero", "Number of ConfigMap task")
	if len(os.Args) < 2 {
		fmt.Println("expected 'foo' or 'bar' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "configmap":
		configMap.Parse(os.Args[2:])
		ConfigMap(*configMapTask)
	case "bar":
		fmt.Println("Placeholder.")
	default:
		fmt.Println("expected 'foo' or 'bar' subcommands")
		os.Exit(1)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8084", nil))

}

func ConfigMap(t string) {
	fmt.Println("This is the", t, "task in configmap topic.")
	switch t {
	case "first":
		envVar := os.Getenv("GRIDU_CONFIGMAP_ENV")
		fmt.Println("The value of GRIDU_CONFIGMAP_ENV is", envVar)
		if envVar == "KUBERNETES_IS_VERY_FUN" {
			hasher := sha1.New()
			hasher.Write([]byte("first" + hash))
			answer := hex.EncodeToString(hasher.Sum(nil))
			fmt.Println("Everything is looks like expected. Here is the correct answer:", answer[:8])
		} else {
			fmt.Println("Env variable value is not equal to KUBERNETES_IS_VERY_FUN \nTry to change your configmap or pod specification.")
		}
	case "second":
		fdat, err := ioutil.ReadFile("/mnt/GRIDU_CONFIGMAP_ENV")
		check(err)
		fmt.Println("Content of the /mnt/GRIDU_CONFIGMAP_ENV:\n", string(fdat))

	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
