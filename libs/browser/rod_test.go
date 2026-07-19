package browser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
	"github.com/stretchr/testify/require"
)

func TestRodBrowser(t *testing.T) {
	browser, err := ConnectRodBrowser()
	if err != nil {
		t.Fatalf("failed to connect browser: %v", err)
	}
	defer browser.Close()

	tab1, err := browser.NewTab(NewNewTabOptions(NewTabWithUrl("https://www.baidu.com")))
	if err != nil {
		t.Fatalf("failed to new tab: %v", err)
	}

	tabs, err := browser.ListTabs()
	if err != nil {
		t.Fatalf("failed to list tabs: %v", err)
	}

	ids := tabs.Ids()

	require.True(t, ids.Contain(tab1.Id))

	tab2, err := browser.NewTab(NewNewTabOptions(NewTabWithUrl("https://www.qq.com")))
	if err != nil {
		t.Fatalf("failed to new tab: %v", err)
	}

	time.Sleep(5 * time.Second)

	tabs, err = browser.ListTabs()
	if err != nil {
		t.Fatalf("failed to list tabs: %v", err)
	}

	ids = tabs.Ids()

	require.True(t, ids.Contain(tab2.Id))

	require.Equal(t, 2, len(tabs))

	json.NewEncoder(os.Stdout).Encode(tabs)

	if err := browser.SwitchToTab(tab1.Id); err != nil {
		t.Fatalf("failed to switch to tab: %v", err)
		return
	}

	texts, err := browser.GetTexts(tab1.Id, "body", SelectorTypeCss)
	if err != nil {
		t.Fatalf("failed to get texts: %v", err)
	}

	fmt.Println(texts)

	htmls, err := browser.GetHtmls(tab1.Id, "body", SelectorTypeCss)
	if err != nil {
		t.Fatalf("failed to get htmls: %v", err)
	}

	fmt.Println(htmls)

	err = browser.SetCookies(tab1.Id, []*http.Cookie{
		{
			Name:     "test",
			Value:    "test",
			Domain:   "www.baidu.com",
			Path:     "/",
			Secure:   true,
			HttpOnly: false,
			SameSite: http.SameSiteNoneMode,
			Expires:  time.Now().Add(1 * time.Hour),
		},
	})
	if err != nil {
		t.Fatalf("failed to set cookies: %v", err)
		return
	}

	time.Sleep(1 * time.Second)

	cookies, err := browser.GetCookies(tab1.Id)
	if err != nil {
		t.Fatalf("failed to get cookies: %v", err)
	}

	for _, cookie := range cookies {
		fmt.Println(cookie.Name)
		json.NewEncoder(os.Stdout).Encode(cookie)
	}

	screenshot, err := browser.Screenshot(tab1.Id)
	if err != nil {
		t.Fatalf("failed to screenshot: %v", err)
	}

	os.WriteFile("screenshot.png", screenshot, 0644)
}

func TestRodBrowserCookies(t *testing.T) {
	browser, err := ConnectRodBrowser()
	if err != nil {
		t.Fatalf("failed to connect browser: %v", err)
		return
	}
	defer browser.Close()

	tab, err := browser.NewTab(NewNewTabOptions(NewTabWithUrl("https://www.baidu.com")))
	if err != nil {
		t.Fatalf("failed to new tab: %v", err)
		return
	}

	err = browser.SetCookies(tab.Id, []*http.Cookie{
		{
			Name:   "test",
			Value:  "test",
			Domain: "www.baidu.com",
			Path:   "/",
			Secure: true,
		},
	})
	if err != nil {
		t.Fatalf("failed to set cookies: %v", err)
		return
	}

	cookies, err := browser.GetCookies(tab.Id)
	if err != nil {
		t.Fatalf("failed to get cookies: %v", err)
		return
	}

	for _, cookie := range cookies {
		fmt.Println(cookie.Name)
		json.NewEncoder(os.Stdout).Encode(cookie)
	}

	names := types.Strings(stl.Map(cookies, func(cookie *http.Cookie) string {
		return cookie.Name
	}))
	fmt.Println(names.Join(", "))

	require.True(t, names.Contain("test"))
}

func TestRodBrowserPrintToPdf(t *testing.T) {
	browser, err := ConnectRodBrowser()
	if err != nil {
		t.Fatalf("failed to connect browser: %v", err)
		return
	}
	defer browser.Close()

	tab, err := browser.NewTab(NewNewTabOptions(NewTabWithUrl("https://www.baidu.com")))
	if err != nil {
		t.Fatalf("failed to new tab: %v", err)
		return
	}

	time.Sleep(1 * time.Second)

	reader, err := browser.PrintToPdf(tab.Id)
	if err != nil {
		t.Fatalf("failed to print to pdf: %v", err)
		return
	}
	defer reader.Close()

	file, err := os.Create("pdf.pdf")
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
		return
	}
	defer file.Close()

	io.Copy(file, reader)
}

func TestRodBrowserStartScreencast(t *testing.T) {
	defer time.Sleep(time.Second)

	browser, err := ConnectRodBrowser()
	if err != nil {
		t.Fatalf("failed to connect browser: %v", err)
		return
	}
	defer browser.Close()

	tab, err := browser.NewTab(NewNewTabOptions(NewTabWithUrl("https://time.is/zh/")))
	if err != nil {
		t.Fatalf("failed to new tab: %v", err)
		return
	}

	stream, err := browser.StartScreencast(tab.Id, nil)
	if err != nil {
		t.Fatalf("failed to start screencast: %v", err)
		return
	}
	defer stream.Close()

	startTime := time.Now()

	var i int
	for frame := range stream.BytesChan {
		elapsed := time.Since(startTime)
		if elapsed > 10*time.Second {
			break
		}

		os.WriteFile(fmt.Sprintf("frames/%05d.jpg", i), frame, 0644)
		i++
	}
}
