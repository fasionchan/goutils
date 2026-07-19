package browser

import (
	"fmt"
	"os"
	"testing"

	"github.com/gorilla/websocket"
)

func TestRemoteControllerWebSocket(t *testing.T) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	var i int
	for {
		typ, message, err := conn.ReadMessage()
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(typ, len(message))

		os.WriteFile(fmt.Sprintf("frames/%05d.jpg", i), message, 0644)
		i++
	}
}
