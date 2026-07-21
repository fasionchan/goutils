package main

import (
	"log"
	"net/http"
	"os"

	browserlib "github.com/fasionchan/goutils/libs/browser"
)

func main() {
	mode := "instance"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	apiHandler := http.NotFoundHandler()

	switch mode {
	case "instance":
		browser, err := browserlib.ConnectRodBrowser()
		if err != nil {
			log.Fatal(err)
		}
		defer browser.Close()

		apiHandler = browserlib.NewBrowserApiHandler(browser)
	case "pool":
		pool := browserlib.NewBrowserPoolFromTypedLaunchFunc(browserlib.LaunchRodBrowser)
		defer pool.Close()

		apiHandler = pool
	default:
		log.Fatalf("Invalid mode: %s", mode)
		return
	}

	// id, err := browser.NewTab(browserlib.NewTabWithUrl("https://time.is/zh/"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// controller := browserlib.NewRemoteController(browser, id, &websocket.Upgrader{
	// 	CheckOrigin: func(r *http.Request) bool {
	// 		return true
	// 	},
	// })

	log.Println("Server started on port 8080")
	http.ListenAndServe(":8080", apiHandler)
}
