package main

import (
	"context"
	"log"
	"net/http"
	"os"

	browserlib "github.com/fasionchan/goutils/libs/browser"
	browsermcp "github.com/fasionchan/goutils/libs/browser/mcp"
)

func main() {
	mode := "instance"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	apiHandler := http.NotFoundHandler()
	opts := browserlib.NewBrowserLaunchOptionsFromEnv(os.Getenv)

	apiPrefix := os.Getenv("API_PREFIX")
	mcpPath := os.Getenv("MCP_PATH")
	if mcpPath == "" {
		mcpPath = browsermcp.DefaultMCPPath
	}

	switch mode {
	case "instance":
		browser, err := browserlib.LaunchRodBrowserForManager(context.Background(), opts)
		if err != nil {
			log.Fatal(err)
		}
		defer browser.Close()

		apiHandler = browserlib.NewBrowserApiHandler(browser).NewChiOpenApiRouter(apiPrefix)

		mcpServer, err := browsermcp.NewBrowserMcpServer(browser, browsermcp.WithPath(mcpPath))
		if err != nil {
			log.Fatal(err)
		}
		apiHandler = mcpServer.MountOnto(apiHandler)
		log.Printf("MCP mounted at %s (sse: %s/sse, streamable: %s)", mcpPath, mcpPath, mcpPath)
	case "pool":
		pool := browserlib.NewBrowserPoolFromTypedLaunchFunc(opts, browserlib.LaunchRodBrowserForManager)
		defer pool.Close()

		apiHandler = pool.NewChiOpenApiRouter(apiPrefix)
		log.Printf("MCP skipped in pool mode (instance only)")
	default:
		log.Fatalf("Invalid mode: %s", mode)
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret != "" {
		apiHandler = browserlib.JwtPathValidator(secret)(apiHandler)
	}

	handler := withOptionalDemoStatic(apiHandler)

	log.Println("Server started on port 8080")
	http.ListenAndServe(":8080", handler)
}

// withOptionalDemoStatic serves built demo assets under /demo/ when DEMO_STATIC_DIR is set.
func withOptionalDemoStatic(api http.Handler) http.Handler {
	demoDir := os.Getenv("DEMO_STATIC_DIR")
	if demoDir == "" {
		return api
	}

	mux := http.NewServeMux()
	mux.Handle("/demo/", http.StripPrefix("/demo/", http.FileServer(http.Dir(demoDir))))
	mux.Handle("/", api)
	log.Printf("Serving demo static files from %s at /demo/", demoDir)
	return mux
}
