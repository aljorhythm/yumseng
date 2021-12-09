package main

import (
	"fmt"
	"github.com/aljorhythm/yumseng/utils"
	"net/http"
	"net/url"
)

func getAllowedOrigins() []string {
	reactTsPort := 3000
	localUiDev := fmt.Sprintf("http://localhost:%d", reactTsPort)
	return []string{localUiDev}
}

func checkSameOrigin(r *http.Request) bool {
	origin := r.Header["Origin"]
	if len(origin) == 0 {
		return true
	}
	u, err := url.Parse(origin[0])
	if err != nil {
		return false
	}
	return utils.EqualASCIIFold(u.Host, r.Host)
}

func requestOriginIsInList(list []string, request *http.Request) bool {
	origin := request.Header["Origin"]
	if len(origin) == 0 {
		return true
	}

	u, err := url.Parse(origin[0])
	if err != nil {
		panic(err)
	}

	key := u.Host
	for _, s := range list {
		if key == s {
			return true
		}
	}
	return false
}
