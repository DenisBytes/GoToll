// obu = on-board unit

package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/DenisBytes/GoToll/types"
	"github.com/gorilla/websocket"
)

const (
	wsEndpoint = "ws://127.0.0.1:30000/ws"
)

var (
	sendInterval = time.Second * 5
)

func genOBUIDs(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(999999)
	}
	return ids
}

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func main() {
	obuIDs := genOBUIDs(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i < len(obuIDs); i++ {
			lat, long := genLatLong()
			data := types.OBUData{
				OBUID: obuIDs[i],
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(sendInterval)
	}
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
