package main

import "os"
import "fmt"
import "flag"

import "io/ioutil"

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

}

func ConfigMap(t string) {
	fmt.Println("This is the", t, "task in configmap topic.")
	switch t {
	case "first":
		envVar := os.Getenv("GRIDU_CONFIGMAP_ENV")
		fmt.Println("The decoded value of GRIDU_CONFIGMAP_ENV:", envVar)
	case "second":
		fdat, err := ioutil.ReadFile("/mnt/configmap.txt")
		check(err)
		fmt.Println("Content of the /mnt/configmap.txt:\n", string(fdat))

	}

}
