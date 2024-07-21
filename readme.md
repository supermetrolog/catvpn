### Local testing

    docker-compose up -d

- docker-compose exec server bash
    - make server
- docker-compose exec client bash
    - make client
- docker-compose exec resource bash
    - make resource


- docker-compose exec client bash
    - ping -I tun0 172.29.0.2
    - ping -M do -I tun0 -s 1300 1.1.1.1
    - curl --interface tun0 --connect-timeout 3 http://172.29.0.2:8080/hi?testmessage
    - traceroute -i tun0 172.29.0.2
    - traceroute -i tun0 --icmp google.com // для сайтов вне докера

- docker-compose exec resource bash
    - tcpdump -p icmp