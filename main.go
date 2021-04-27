package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func getTimeStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func minInt(x, y uint64) uint64 {
	if x > y {
		return y
	} else {
		return x
	}
}

func maxInt(x, y uint64) uint64 {
	if x > y {
		return x
	} else {
		return y
	}
}

func minByte(x, y byte) byte {
	if x > y {
		return y
	} else {
		return x
	}
}

func maxByte(x, y byte) byte {
	if x > y {
		return x
	} else {
		return y
	}
}

func sumBytes(bytes []byte) byte {
	result := byte(0)
	for _, v := range bytes {
		result += v
	}
	return result
}

func readStatusCode(body string) uint16 {
	leftOfStatus := strings.Index(body, " ")
	statusCodeAsString := body[leftOfStatus+1 : (leftOfStatus + 1 + 4)]
	res, _ := strconv.ParseUint(strings.TrimSpace(statusCodeAsString), 10, 64)
	statusCodeResult := uint16(res)
	return statusCodeResult
}

const requestDelay = 2 //seconds
const maximumByte byte = 255

type Response struct {
	TimeInMs     uint64 `json:"timeInMs"`
	Body         string `json:"body"`
	StatusCode   uint16 `json:"statusCode"`
	ResponseSize byte   `json:"responseSize(bytes)"`
}

type Request struct {
	Address         string `json:"address"`
	AddressStripped string `json:"addressStripped"`
	Route           string `json:"requestedRoute"`
}

type Profiler struct {
	RequestCount           *uint64  `json:"totalRequestCount"`
	FailedRequestCount     *uint64  `json:"failedRequestCount"`
	SuccessfulRequestCount *uint64  `json:"successfulRequestCount"`
	ErrorCodes             []uint16 `json:"errorCodes,omitempty"`

	Request Request `json:"originalRequest"`

	MinResponseSize byte `json:"minResponseSize"`
	MaxResponseSize byte `json:"maxResponseSize"`

	MeanResponseTime    uint64 `json:"meanResponseTime"`
	MedianResponseTime  uint64 `json:"medianResponseTime"`
	FastestResponseTime uint64 `json:"fastestResponseTime"`
	SlowestResponseTime uint64 `json:"slowestResponseTime"`
}

func (profiler *Profiler) String() string {
	return fmt.Sprintf("Percentage of successful requests (%d / %d) = %%%d", *profiler.SuccessfulRequestCount, *profiler.RequestCount, (*profiler.SuccessfulRequestCount / *profiler.RequestCount)*100)
}

func (res *Response) Print() {
	json.NewEncoder(os.Stdout).Encode(res)
}

func (profiler *Profiler) Print() {
	json.NewEncoder(os.Stdout).Encode(profiler)
}

func (req *Request) GetResponse() (Response, error) {
	start := getTimeStamp()
	cf := &tls.Config{
		Rand:               rand.Reader,
		InsecureSkipVerify: true,
	}
	cf.Time = func() time.Time { return time.Now() }
	ssl, err := tls.Dial("tcp", net.JoinHostPort(req.AddressStripped, "443"), cf)
	if err != nil {
		return Response{}, err
	}
	needsWorldWideWeb := strings.Contains(req.Address, "www.")
	tcpRequestString := ""
	if needsWorldWideWeb {
		tcpRequestString = "GET " + req.Route + " HTTP/1.1\r\nHost: www." + req.AddressStripped + "\r\nConnection: close\r\nContent-Type: */* charset=utf-8;\r\nUpgrade-Insecure-Requests: 1\r\nAccept: */*\r\nX-Originating-URL: " + req.AddressStripped + "\r\n\r\n"
	} else {
		tcpRequestString = "GET " + req.Route + " HTTP/1.1\r\nHost: " + req.AddressStripped + "\r\nConnection: close\r\nContent-Type: */* charset=utf-8;\r\nUpgrade-Insecure-Requests: 1\r\nAccept: */*\r\nX-Originating-URL: " + req.AddressStripped + "\r\n\r\n"
	}
	_, err = ssl.Write([]byte(tcpRequestString))
	if err != nil {
		return Response{}, err
	}
	result, err := ioutil.ReadAll(ssl)
	if err != nil {
		return Response{}, err
	}
	responseResult := Response{}
	elapsed := getTimeStamp() - start
	responseResult.TimeInMs = uint64(elapsed)
	responseResult.Body = string(result)
	responseResult.ResponseSize = sumBytes(result)
	responseResult.StatusCode = readStatusCode(responseResult.Body)
	defer ssl.Close()
	return responseResult, nil
}

func (profiler *Profiler) StartProfiling() {
	responseTimes := make([]uint64, 0)
	var responseTimeRunningSum uint64 = 0
	reqCount := *profiler.RequestCount
	for reqCount > 0 {
		currentResponse, err := profiler.Request.GetResponse()
		if err != nil {
			//request failed
			*profiler.FailedRequestCount++
		} else if currentResponse.StatusCode != 200 {
			profiler.ErrorCodes = append(profiler.ErrorCodes, currentResponse.StatusCode)
			*profiler.FailedRequestCount++
		} else {
			profiler.MinResponseSize = minByte(profiler.MinResponseSize, currentResponse.ResponseSize)
			profiler.MaxResponseSize = maxByte(profiler.MaxResponseSize, currentResponse.ResponseSize)

			profiler.SlowestResponseTime = maxInt(profiler.SlowestResponseTime, currentResponse.TimeInMs)
			profiler.FastestResponseTime = minInt(profiler.FastestResponseTime, currentResponse.TimeInMs)

			responseTimes = append(responseTimes, currentResponse.TimeInMs)

			responseTimeRunningSum += currentResponse.TimeInMs

			*profiler.SuccessfulRequestCount++
			time.Sleep(requestDelay * time.Second)
		}
		reqCount--
	}

	sort.Slice(responseTimes, func(i, j int) bool {
		return i < j
	})

	if len(responseTimes) > 0 {
		isEven := len(responseTimes)%2 == 0
		if isEven {
			if len(responseTimes) == 2 {
				profiler.MedianResponseTime = (responseTimes[0] + responseTimes[1]) / 2
			} else {
				leftIdx := (len(responseTimes) / 2) - 1
				rightIdx := len(responseTimes) / 2
				profiler.MedianResponseTime = (responseTimes[leftIdx] + responseTimes[rightIdx]) / 2
			}
		} else {
			medianIdx := len(responseTimes) / 2
			profiler.MedianResponseTime = responseTimes[medianIdx]
		}
	} else {
		profiler.MedianResponseTime = math.MaxUint64
	}
	if *profiler.SuccessfulRequestCount == 0 {
		profiler.MeanResponseTime = math.MaxUint64
	} else {
		profiler.MeanResponseTime = responseTimeRunningSum / uint64(*profiler.SuccessfulRequestCount)
	}
}

func main() {
	url := flag.String("url", "", "The url to request")
	profile := flag.Uint64("profile", 0, "The amount of requests to profile")
	flag.Parse()
	if *profile < 0 {
		log.Println("Your profile count parameter must be positive!")
		os.Exit(1)
	} else if len(*url) == 0 {
		log.Fatal("You must provide a -url argument. Example -url=\"cloudflare.com\"")
	} else {
		address, route, addressStripped := parseAddressAndRoute(*url)
		request := &Request{ Address: address, AddressStripped: addressStripped, Route: route }
		if *profile == 0 {
			response, err := request.GetResponse()
			if err != nil {
				log.Println("There was an error handling your request")
				log.Fatal(err.Error())
				os.Exit(1)
			}
			response.Print()
		} else {
			profiler := &Profiler{ RequestCount: profile, FailedRequestCount: new(uint64), SuccessfulRequestCount: new(uint64),
				MinResponseSize: maximumByte, MaxResponseSize: 0, FastestResponseTime: math.MaxUint64, SlowestResponseTime: 0, Request: *request, ErrorCodes: make([]uint16, 0) }
			profiler.StartProfiling()
			profiler.Print()
			//percentage stat
			log.Println(profiler)
		}
	}
	os.Exit(0)
}

func parseAddressAndRoute(url string) (string, string, string) {
	address := url[:strings.LastIndex(url, ".")]
	route := ""
	for idx := strings.LastIndex(url, "."); idx < len(url); idx++ {
		if url[idx:idx+1] != "/" {
			address += url[idx : idx+1]
		} else {
			route = url[idx:]
			break
		}
	}
	addressStripped := address
	if (strings.HasPrefix(address, "http") && strings.Contains(address, "www.")) || strings.HasPrefix(address, "www.") {
		addressStripped = address[strings.Index(address, ".")+1:]
	} else if strings.HasPrefix(address, "http") {
		addressStripped = address[strings.Index(address, ":")+3:]
	}
	routeIdx := strings.Index(addressStripped, "/")
	if routeIdx != -1 {
		route = addressStripped[routeIdx:]
	}
	if route == "" {
		route = "/"
	}
	return address, route, addressStripped
}
