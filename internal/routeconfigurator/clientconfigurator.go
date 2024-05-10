package routeconfigurator

type ClientTrafficRouteConfigurator struct {
}

func NewClientTrafficRouteConfigurator() *ClientTrafficRouteConfigurator {
	return &ClientTrafficRouteConfigurator{}
}

func (t *ClientTrafficRouteConfigurator) RouteToIface(ifaceName string) error {
	//log.Printf("Назначаем форвардинг для созданного интерфейса: %s\n", ifaceName)
	//
	//cmd := fmt.Sprintf("sysctl -w net.ipv4.ip_forward=1")
	//out, err := command.RunCommand(cmd)
	//if err != nil {
	//	return fmt.Errorf("sysctl forward error: out: %s, error: %w", out, err)
	//}
	//
	//cmd = fmt.Sprintf("iptables -t nat -A POSTROUTING -o tun0 -j MASQUERADE")
	//out, err = command.RunCommand(cmd)
	//if err != nil {
	//	return fmt.Errorf("iptables nat postrouting error: out: %s, error: %w", out, err)
	//}
	//
	//cmd = fmt.Sprintf("iptables -I FORWARD 1 -i tun0 -m state --state RELATED,ESTABLISHED -j ACCEPT")
	//out, err = command.RunCommand(cmd)
	//if err != nil {
	//	return fmt.Errorf("iptables forward -i tun error: out: %s, error: %w", out, err)
	//}
	//
	//cmd = fmt.Sprintf("iptables -I FORWARD 1 -o tun0 -j ACCEPT")
	//out, err = command.RunCommand(cmd)
	//if err != nil {
	//	return fmt.Errorf("iptables forward -o tun error: out: %s, error: %w", out, err)
	//}
	return nil
}
