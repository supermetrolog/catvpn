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
    - ping -I tun0 1.1.1.1
    - curl --interface tun0 http://172.29.0.2:8080/hi?sdawdawd


- docker-compose exec resource bash
    - tcpdump -p icmp