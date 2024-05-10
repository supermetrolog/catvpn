### Local testing

    docker-compose up -d

- docker-compose exec server bash
    - make server

- docker-compose exec client bash
    - make client
- docker-compose exec client bash
    - make resource

- docker-compose exec server bash
    
- docker-compose exec client bash
    - ping -M do -I tun0 -s 1300 1.1.1.1
    - curl --interface tun0 --connect-timeout 3 http://172.29.0.2:8080/hi?sdawdawd
    - traceroute -i tun0 172.29.0.2
    - traceroute -i tun0 --icmp vpn.diani.ru // для сайтов вне докера



- docker-compose exec resource bash
    - tcpdump -p icmp