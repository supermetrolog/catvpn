package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	runTestServer("")
}

func runTestServer(ip string) {
	http.HandleFunc("/hi", func(writer http.ResponseWriter, request *http.Request) {
		log.Println(request.RemoteAddr)

		//req, _ := http.NewRequest("GET", "https://google.com", nil)
		//
		//res, _ := http.DefaultClient.Do(req)
		//
		//defer res.Body.Close()
		//
		//io.Copy(writer, res.Body)

		writer.Write([]byte(fmt.Sprintf("hi %s Mesage: %s \n ", request.RemoteAddr, request.URL.String())))
		//writer.Write([]byte("### Local testing\n\n    docker-compose up -d\n\n- docker-compose exec server bash\n    - make server\n\n- docker-compose exec client bash\n    - make client\n- docker-compose exec client bash\n    - make resource\n\n- docker-compose exec server bash\n    \n- docker-compose exec client bash\n    - ping -M do -I tun0 -s 1300 1.1.1.1\n    - curl --interface tun0 --connect-timeout 3 http://172.29.0.2:8080/hi?sdawdawd\n\n\n- docker-compose exec resource bash\n    - tcpdump -p icmp"))
		return
	})
	err := http.ListenAndServe(fmt.Sprintf("%s:8080", ip), nil)
	if err != nil {
		log.Println(err)
	}
}
