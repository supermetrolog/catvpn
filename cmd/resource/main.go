package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	_ "github.com/supermetrolog/myvpn/internal/logger"
	"net/http"
)

func main() {
	runTestServer("")
}

func runTestServer(ip string) {
	http.HandleFunc("/hi", func(writer http.ResponseWriter, request *http.Request) {
		logrus.Infof(request.RemoteAddr)

		//req, _ := http.NewRequest("GET", "https://google.com", nil)
		//res, _ := http.DefaultClient.Do(req)
		//defer res.Body.Close()
		//io.Copy(writer, res.Body)

		_, err := writer.Write([]byte(fmt.Sprintf("hi %s Mesage: %s \n ", request.RemoteAddr, request.URL.String())))

		if err != nil {
			logrus.Error(err)
		}

		return
	})

	err := http.ListenAndServe(fmt.Sprintf("%s:8080", ip), nil)

	if err != nil {
		logrus.Error(err)
	}
}
