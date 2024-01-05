/*
 * Author: fasion
 * Created time: 2024-01-05 09:24:00
 * Last Modified by: fasion
 * Last Modified time: 2024-01-05 11:19:44
 */

package netutils

import (
	"encoding/json"
	"net/url"
)

func JoinRawUrl(base string, target string) (string, error) {
	url, err := JoinUrlPro(nil, nil, base, target)
	if err != nil {
		return "", nil
	}
	return url.String(), nil
}

func JoinUrlWithRaw(base *url.URL, target string) (*url.URL, error) {
	return JoinUrlPro(base, nil, "", target)
}

func JoinUrlPro(base, target *url.URL, rawBase, rawTarget string) (joined *url.URL, err error) {
	if target == nil {
		target, err = url.Parse(rawTarget)
		if err != nil {
			return nil, err
		}
	}

	if target.Scheme != "" {
		return target, nil
	}

	if base == nil {
		base, err = url.Parse(rawBase)
		if err != nil {
			return nil, err
		}
	}

	return JoinUrl(base, target), nil
}

func JoinUrl(base, target *url.URL) *url.URL {
	if target.Scheme != "" {
		return target
	}

	result := *target

	result.Scheme = base.Scheme
	result.User = base.User
	result.Host = base.Host
	result.Path = base.Path

	return result.JoinPath(target.Path)
}

func UrlInfo(_url *url.URL) map[string]any {
	password, hasPassword := _url.User.Password()
	return map[string]any{
		"Scheme":      _url.Scheme,
		"Username":    _url.User.Username(),
		"HasPassword": hasPassword,
		"Password":    password,
		"Host":        _url.Host,
		"Hostname":    _url.Hostname(),
		"Port":        _url.Port(),
		"Path":        _url.Path,
		"RawPath":     _url.RawPath,
		"Query":       _url.Query(),
		"RawQuery":    _url.RawQuery,
		"Frame":       _url.Fragment,
		"RawFrame":    _url.RawFragment,

		"Opaque": _url.Opaque,
	}
}

func JsonMarshalUrlInfo(_url *url.URL, prefix, indent string) (result []byte) {
	result, _ = json.MarshalIndent(UrlInfo(_url), prefix, indent)
	return
}

func JsonifyUrlInfo(_url *url.URL, prefix, indent string) string {
	return string(JsonMarshalUrlInfo(_url, prefix, indent))
}

func ParseAndJsonifyUrl(_url string) (string, error) {
	parsed, err := url.Parse(_url)
	if err != nil {
		return "", err
	}
	return JsonifyUrlInfo(parsed, "", "  "), nil
}
