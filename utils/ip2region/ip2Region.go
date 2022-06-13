package ip2region

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	INDEX_BLOCK_LENGTH  = 12
	TOTAL_HEADER_LENGTH = 8192
)

type Ip2Region struct {
	// super block index info
	firstIndexPtr int64
	lastIndexPtr  int64
	totalBlocks   int64

	dbBinStr []byte
}

type IpInfo struct {
	CityId   int64
	Country  string
	Region   string
	Province string
	City     string
	ISP      string
}

func (ip IpInfo) String() string {
	return strconv.FormatInt(ip.CityId, 10) + "|" + ip.Country + "|" + ip.Region + "|" + ip.Province + "|" + ip.City + "|" + ip.ISP
}

func getIpInfo(cityId int64, line []byte) *IpInfo {

	lineSlice := strings.Split(string(line), "|")
	ipInfo := &IpInfo{}
	length := len(lineSlice)
	ipInfo.CityId = cityId
	if length < 5 {
		for i := 0; i <= 5-length; i++ {
			lineSlice = append(lineSlice, "")
		}
	}

	ipInfo.Country = lineSlice[0]
	ipInfo.Region = lineSlice[1]
	ipInfo.Province = lineSlice[2]
	ipInfo.City = lineSlice[3]
	ipInfo.ISP = lineSlice[4]
	return ipInfo
}

func New(path string) (*Ip2Region, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ret := &Ip2Region{
		dbBinStr: raw,
	}
	ret.MemoryInit()
	return ret, nil
}
func NewWithRaw(raw []byte) (*Ip2Region, error) {
	ret := &Ip2Region{
		dbBinStr: raw,
	}
	ret.MemoryInit()
	return ret, nil
}

func (ip2r *Ip2Region) MemoryInit() error {
	ip2r.firstIndexPtr = getLong(ip2r.dbBinStr, 0)
	ip2r.lastIndexPtr = getLong(ip2r.dbBinStr, 4)
	ip2r.totalBlocks = (ip2r.lastIndexPtr-ip2r.firstIndexPtr)/INDEX_BLOCK_LENGTH + 1
	return nil
}

func (ip2r *Ip2Region) MemorySearch(ipStr string) (*IpInfo, error) {
	if ip2r.totalBlocks == 0 {
		return nil, fmt.Errorf("ip2region file is not initialized")
	}

	ip, err := StrIP2Int(ipStr)
	if err != nil {
		return nil, err
	}

	h := ip2r.totalBlocks
	var dataPtr, l int64
	for l <= h {

		m := (l + h) >> 1
		p := ip2r.firstIndexPtr + m*INDEX_BLOCK_LENGTH
		sip := getLong(ip2r.dbBinStr, p)
		if ip < sip {
			h = m - 1
		} else {
			eip := getLong(ip2r.dbBinStr, p+4)
			if ip > eip {
				l = m + 1
			} else {
				dataPtr = getLong(ip2r.dbBinStr, p+8)
				break
			}
		}
	}
	if dataPtr == 0 {
		return nil, errors.New("not found")
	}

	dataLen := ((dataPtr >> 24) & 0xFF)
	dataPtr = (dataPtr & 0x00FFFFFF)
	ipInfo := getIpInfo(getLong(ip2r.dbBinStr, dataPtr), ip2r.dbBinStr[(dataPtr)+4:dataPtr+dataLen])
	return ipInfo, nil
}

func getLong(b []byte, offset int64) int64 {
	val := (int64(b[offset]) |
		int64(b[offset+1])<<8 |
		int64(b[offset+2])<<16 |
		int64(b[offset+3])<<24)
	return val
}

func StrIP2Int(IpStr string) (int64, error) {
	bits := strings.Split(IpStr, ".")
	if len(bits) != 4 {
		return 0, errors.New("ip format error")
	}
	var sum int64
	for i, n := range bits {
		bit, _ := strconv.ParseInt(n, 10, 64)
		sum += bit << uint(24-8*i)
	}
	return sum, nil
}
