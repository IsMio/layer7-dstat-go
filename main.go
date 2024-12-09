package main

import (
	"ddosbackend/module"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jessevdk/go-flags"
	"github.com/shirou/gopsutil/net"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

var (
	m          module.Flag
	conn       = int32(0)
	netinfoold []net.IOCountersStat
	upgrader   = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	b runtime.MemStats
)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		netinfo, _ := net.IOCounters(true)

		if len(netinfoold) == 0 {
			netinfoold = netinfo
		}
		bytesend := netinfo[m.NIC].BytesSent
		byterecv := netinfo[m.NIC].BytesRecv
		atomic.AddUint64(&bytesend, -netinfoold[m.NIC].BytesSent)
		atomic.AddUint64(&byterecv, -netinfoold[m.NIC].BytesRecv)
		netinfoold = netinfo
		runtime.ReadMemStats(&b)
		msg := module.Module{
			Inspeed:     byterecv,
			Outspeed:    bytesend,
			Connection:  conn,
			Memoryusage: b.Sys / 1024 / 1024,
		}
		jsonstu, err := json.Marshal(msg)
		if err != nil {
			fmt.Println(err)
			return
		}
		ws.WriteMessage(websocket.TextMessage, jsonstu)
		time.Sleep(1 * time.Second)

	}

}
func work(w http.ResponseWriter, r *http.Request) {
	conn++
	fmt.Fprintf(w, "success")

}
func timer() {
	for {
		time.Sleep(1 * time.Second)
		conn = 0
	}
}

func main() {
	var m module.Flag
	_, err := flags.Parse(&m)
	if err != nil {
		return
	}
	netinfo, _ := net.IOCounters(true)
	if m.List {

		for i, v := range netinfo {
			fmt.Println(i, v.Name)

		}
	} else if m.Addr == "" || m.NIC == 0 {
		fmt.Println("Need -a and -n use -h for help")
		return
	}
	fmt.Printf("Server start at %s select listen NIC %s\n", m.Addr, netinfo[m.NIC].Name)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", WebSocketHandler)
	mux.HandleFunc("/attack", work)
	mux.Handle("/", http.FileServer(http.Dir("./html")))
	srv := &http.Server{
		Addr:    m.Addr,
		Handler: mux,
	}
	go timer()
	go srv.ListenAndServe()
	select {}

}
