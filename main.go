package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"ip2region/utils/ip2region"
)

var instance *ip2region.Ip2Region

//go:embed ip2region.db
var dbFile []byte

func init() {
	tmp, err := ip2region.NewWithRaw(dbFile)
	if err != nil {
		panic(err)
	}
	instance = tmp
}

func ReadFromPipeline() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return
	}

	if (stat.Mode() & os.ModeNamedPipe) == 0 {
		//here no pipe
		return
	}

	reader := bufio.NewReader(os.Stdin)
	sc := bufio.NewScanner(reader)
	sc.Split(bufio.ScanLines)

	for {
		if !sc.Scan() {
			break
		}
		line := sc.Text()
		Region(line)
	}
}

func Region(ip string) {
	ip = strings.Trim(ip, " ")
	ipInfo, err := instance.MemorySearch(ip)
	fmt.Printf("%s\t->\t", ip)
	if err == nil {
		fmt.Printf("%s.%s.%s.%s\n", ipInfo.Country, ipInfo.Province, ipInfo.City, ipInfo.ISP)
	} else {
		fmt.Println(err)
	}
}

func main() {
	switch len(os.Args) {
	case 1:
		ReadFromPipeline()
	case 2:
		Region(os.Args[1])
	default:
	}
}
