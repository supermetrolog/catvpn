package routeconfigurator

import (
	"fmt"
	"github.com/supermetrolog/myvpn/internal/helpers/command"
	"net"
)

type ServerTrafficRouteConfigurator struct {
}

func NewServerTrafficRouteConfigurator() *ServerTrafficRouteConfigurator {
	return &ServerTrafficRouteConfigurator{}
}

func (t *ServerTrafficRouteConfigurator) RouteToSubnet(subnet net.IPNet) error {
	CIDR := subnet.String()

	// Разрешаем транзин пакетов между сетевыми интерфейсами
	cmd := fmt.Sprintf("sysctl -w net.ipv4.ip_forward=1")
	out, err := command.RunCommand(cmd)
	if err != nil {
		return fmt.Errorf("enable forward error: out: %s, error: %w", out, err)
	}

	// Меняем (маскируем) srcIp на Ip внешнего интерфейса если srcIp == SUBNET и dstIp != SUBNET
	cmd = fmt.Sprintf("iptables -t nat -A POSTROUTING -s %s ! -d %s -j MASQUERADE", CIDR, CIDR)
	out, err = command.RunCommand(cmd)
	if err != nil {
		return fmt.Errorf("iptables setup POSTROUTING error: out: %s, error: %w", out, err)
	}

	// Разрешаем транзин пакетов из подсети tun интерфейса
	cmd = fmt.Sprintf("iptables -A FORWARD -s %s -m state --state RELATED,ESTABLISHED -j ACCEPT", CIDR)
	//cmd = fmt.Sprintf("iptables -A FORWARD -s %s -j ACCEPT", CIDR)
	out, err = command.RunCommand(cmd)
	if err != nil {
		return fmt.Errorf("iptables setup FORWARD with state error: out: %s, error: %w", out, err)
	}

	// Разрешаем транзин пакетов в подсеть tun интерфейса
	cmd = fmt.Sprintf("iptables -A FORWARD -d %s -j ACCEPT", CIDR)
	out, err = command.RunCommand(cmd)
	if err != nil {
		return fmt.Errorf("iptables setup FORWARD error: out: %s, error: %w", out, err)
	}

	return nil
}
