package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/disk"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	InfoStore "alphamon_agent/src/lib"
	
)

func main() {
	//var wg sync.WaitGroup˜
	//MakeRequest()
  
	epochRequest()
	go MakeServer()
	ticker := time.NewTicker(10 * time.Second)


	// for every `tick` that our `ticker`
	// emits, we print `tock`
	for _ = range ticker.C {
		MakeRequest()
	}
	// wg.Wait()

}
func MakeServer() {


	fmt.Println("Start Thread")
	http.HandleFunc("/connect", homeLink)
	log.Fatal(http.ListenAndServe(":8081", nil))
	fmt.Println("Exit Thread")

}
func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Welcome home!")
	MakeRequest()
}
func epochRequest() {

	fmt.Printf("Sending epoch register request!")
type Config struct {
		Server struct {
			Port     string `yaml:"port"`
			Host     string `yaml:"host"`
			Key		 string `yaml:"key"`
		} `yaml:"server"`
	}
	f, err := os.Open("/etc/alphamon/server_config.yml")
	if err != nil {
		log.Fatalln(err)
	}

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalln(err)
	}
	//Send request to register with server

	devName := getHostName()
	fmt.Printf("DevName "+devName)
	
	registerMsg := map[string]interface{}{

		"deviceId": "",
		"name":   devName,
	}

   bytesRepresentation, err := json.Marshal(registerMsg)
	if err != nil {
		log.Fatalln(err)
	}
	InfoStore.Write(cfg.Server.Host)
	//http.Post(cfg.Server.Host+"/register", "application/json", bytes.NewBuffer(bytesRepresentation))
	fmt.Println("Key is :",cfg.Server.Key)
	client := &http.Client{}
	req, err := http.NewRequest("POST", cfg.Server.Host+"/register", bytes.NewBuffer(bytesRepresentation))
	req.Header.Add("Key", cfg.Server.Key)
	req.Header.Add("Content-Type","Application/json")
	resp, err := client.Do(req)
	if err ==nil {
	
	defer resp.Body.Close()
	}

}

func MakeRequest() {
	fmt.Printf("Sending request")
	type Config struct {
		Server struct {
			Port     string `yaml:"port"`
			Host     string `yaml:"host"`
			DeviceId string `yaml:"serial"`
		} `yaml:"server"`
	}
	f, err := os.Open("/etc/alphamon/server_config.yml")
	if err != nil {
		log.Fatalln(err)
	}

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalln(err)
	}

	v, _ := mem.VirtualMemory()
	cpu_use, _ := cpu.Percent(0, false)
	platform, family, version, _ := host.PlatformInformation()
	is := fmt.Sprintf("%.2f", v.UsedPercent)
	cpu_percent := fmt.Sprintf("%.2f", cpu_use[0])
        disk_use := getDiskUsage()
	message := map[string]interface{}{

		"deviceId": cfg.Server.DeviceId,
		"osName":   platform + family + version,
		"cpuUsage": cpu_percent,
		"memUsage": is,
		"diskUsage":disk_use,
		"timestamp":time.Now(),
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(cfg.Server.Host)

	resp, err := http.Post(cfg.Server.Host+"/test", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}
	respdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	responseString := string(respdata)
	fmt.Println(responseString)
}

func getDiskUsage() (diskUsed float64) {
	parts, _ := disk.Partitions(false)

	var usage []*disk.UsageStat
        var diskU float64

	for _, part := range parts {
		u, _ := disk.Usage(part.Mountpoint)
		usage = append(usage, u)
		fmt.Printf( fmt.Sprintf("%f",u.UsedPercent))
                fmt.Printf("\n")
                fmt.Printf(u.Path)
                fmt.Printf("\n")

			diskU= u.UsedPercent;	
	
			break
	}
        return diskU

}


func getHostName()(hostnmae string){



	devInfo,_ := host.Info()

	fmt.Printf(devInfo.Hostname)
	return devInfo.Hostname
}
