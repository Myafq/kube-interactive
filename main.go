package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//go build -ldflags "-X main.hash=$(date -j -f "%a %b %d %T %Z %Y" "`date`" "+%s")"
var hash string = "123456789"

func main() {
	expectation := "possible commands: config, workload, ingress"
	// check for hackers
	_, kubeenv := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	if !kubeenv {
		fmt.Println("Trying to hack a test, huh? You shouldn't probably do that.")
		os.Exit(1)
	}
	// subcommands and args
	config := flag.NewFlagSet("config", flag.ExitOnError)
	configTask := config.String("t", "Zero", "Number of ConfigMap task")
	wl := flag.NewFlagSet("workloads", flag.ExitOnError)
	wlTask := wl.String("t", "Zero", "Number of Workloads task")

	if len(os.Args) < 2 {
		fmt.Println(expectation)
		os.Exit(1)
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/hostname", getHostname)
	go http.ListenAndServe(":8084", nil)

	switch os.Args[1] {
	case "config":
		config.Parse(os.Args[2:])
		ConfigCheck(*configTask)
	case "workloads":
		wl.Parse(os.Args[2:])
		WorkLoads(*wlTask)
	default:
		fmt.Println(expectation)
		os.Exit(1)
	}
	time.Sleep(120 * time.Minute)
}
func WorkLoads(t string) {
	switch t {
	case "first":
		svc, svcExists := os.LookupEnv("serviceName")
		if !svcExists {
			fmt.Println("serviceName environment variable doesn't exist! Fix your specification.")
			os.Exit(1)
		}
		cluster := make(map[string]bool)
		for len(cluster) < 3 {
			fmt.Println("Looking for cluster members on", "http://"+svc+":8084")
			time.Sleep(5 * time.Second)
			ips, _ := net.LookupHost(svc)

			for ip := range ips {
				time.Sleep(1 * time.Second)
				clusterMember, err := http.Get("http://" + string(ip) + ":8084/hostname")
				if err != nil {
					fmt.Println("Error occured while discovering cluster members:", err)
					continue
				}
				defer clusterMember.Body.Close()
				body, _ := ioutil.ReadAll(clusterMember.Body)
				cluster[string(body)] = true
				currState := ""
				for k := range cluster {
					currState += k
				}
				fmt.Println("Current cluster members:", currState)
			}
		}
		fmt.Println("We've got 3 instances of application online!")
		hm, _ := os.Hostname()
		if hm[len(hm)-2:] == "-2" {
			hasher := sha1.New()
			hasher.Write([]byte("sts" + t + hash))
			answer := hex.EncodeToString(hasher.Sum(nil))
			fmt.Println("This is third instance of statefullset! So here is ne answer:", answer[:8])
		}
	}
}
func ConfigCheck(t string) {
	fmt.Println("This is the", t, "task in configmap topic.")
	switch t {
	case "first":
		envVar := os.Getenv("GRIDU_CONFIGMAP_ENV")
		fmt.Println("The value of GRIDU_CONFIGMAP_ENV is", envVar)
		if envVar == "KUBERNETES_IS_VERY_FUN" {
			hasher := sha1.New()
			hasher.Write([]byte(t + hash))
			answer := hex.EncodeToString(hasher.Sum(nil))
			fmt.Println("Everything is looks like expected. Here is the correct answer:", answer[:8])
		} else {
			fmt.Println("Env variable value is not equal to KUBERNETES_IS_VERY_FUN \nTry to change your configmap or pod specification.")
		}

	case "second":
		fdat, err := ioutil.ReadFile("/mnt/GRIDU_CONFIGMAP_ENV")
		check(err)
		fmt.Println("Content of the /mnt/GRIDU_CONFIGMAP_ENV:\n", string(fdat))
		if string(fdat) == "KUBERNETES_IS_VERY_FUN" {
			hasher := sha1.New()
			hasher.Write([]byte("second" + hash))
			answer := hex.EncodeToString(hasher.Sum(nil))
			fmt.Println("Everything is looks like expected. Here is the correct answer:", answer[:8])
		} else {
			fmt.Println("File content is not equal to KUBERNETES_IS_VERY_FUN \nTry to change your configmap or pod specification.")
		}
	case "third":
		envVar := os.Getenv("GRIDU_SECRET_ENV")
		fmt.Println("The value of GRIDU_SECRET_ENV is", envVar)
		if envVar == "KUBERNETES_IS_VERY_SECURE" {
			hasher := sha1.New()
			hasher.Write([]byte(t + hash))
			answer := hex.EncodeToString(hasher.Sum(nil))
			fmt.Println("Everything is looks like expected. Here is the correct answer:", answer[:8])
		} else {
			fmt.Println("Env variable value is not equal to KUBERNETES_IS_VERY_SECURE \nTry to change your secret or pod specification.")
		}
	case "fourth":
		fdat, err := ioutil.ReadFile("/mnt/GRIDU_SECRET_ENV")
		check(err)
		fmt.Println("Content of the /mnt/GRIDU_SECRET_ENV:\n", string(fdat))
		if string(fdat) == "KUBERNETES_IS_VERY_SECURE" {
			hasher := sha1.New()
			hasher.Write([]byte(t + hash))
			answer := hex.EncodeToString(hasher.Sum(nil))
			fmt.Println("Everything is looks like expected. Here is the correct answer:", answer[:8])
		} else {
			fmt.Println("File content is not equal to KUBERNETES_IS_VERY_SECURE \nTry to change your secret or pod specification.")
		}
	}

}
func getHostname(w http.ResponseWriter, r *http.Request) {
	hm, _ := os.Hostname()
	fmt.Fprintf(w, hm)
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
