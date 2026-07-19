package main

import (
	"log"
	"net/http"

	browserlib "github.com/fasionchan/goutils/libs/browser"
)

func main() {
	browser, err := browserlib.ConnectRodBrowser()
	if err != nil {
		log.Fatal(err)
	}
	defer browser.Close()

	apiHandler := browserlib.NewBrowserApiHandler(browser)

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
