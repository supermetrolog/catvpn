package tuntap

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"github.com/supermetrolog/myvpn/internal/common"
	"github.com/supermetrolog/myvpn/internal/helpers/command"
	"net"
)

type TunFactory struct {
}

func New() *TunFactory {
	return &TunFactory{}
}

func (t *TunFactory) Create(subnet net.IPNet, mtu int) (common.Tun, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}

	iface, err := water.New(config)

	if err != nil {
		return nil, fmt.Errorf("create tun iface error: %v", err)
	}

	logrus.Debugf("Created interface with name: %s", iface.Name())

	logrus.Debugf("Назначаем размер MTU: %s, для созданного интерфейса: %s", subnet.IP, iface.Name())

	cmd := fmt.Sprintf("ip link set dev %s mtu %d", iface.Name(), mtu)
	out, err := command.RunCommand(cmd)
	if err != nil {
		return nil, fmt.Errorf("set iface mtu error: out: %s, error: %w", out, err)
	}

	logrus.Debugf("Назначаем IP адресс: %s, для созданного интерфейса: %s", subnet.IP, iface.Name())

	cmd = fmt.Sprintf("ip addr add %s dev %s", subnet.String(), iface.Name())
	out, err = command.RunCommand(cmd)
	if err != nil {
		return nil, fmt.Errorf("set ip addr error: out: %s, error: %w", out, err)
	}

	logrus.Debugln("Включаем созданный интерфейс")

	cmd = fmt.Sprintf("ip link set dev %s up", iface.Name())
	out, err = command.RunCommand(cmd)
	if err != nil {
		return nil, fmt.Errorf("enable created iface error: out: %s, error: %w", out, err)
	}

	//log.Printf("Маршрутизируем пир: %s, для созданного интерфейса: %s\n", ip, iface.Name())
	//
	//cmd = fmt.Sprintf("ip addr add dev %s local %s peer %s", iface.Name(), ip, "10.1.1.2")
	//out, err = command.RunCommand(cmd)
	//if err != nil {
	//	log.Println(out)
	//	return nil, err
	//}
	//
	//log.Printf("Маршрутизируем подсеть через пир\n")
	//
	//cmd = fmt.Sprintf("ip route change %s via %s dev %s", "10.1.1.0/24", "10.1.1.2", iface.Name())
	//out, err = command.RunCommand(cmd)
	//if err != nil {
	//	log.Println(out)
	//	return nil, err
	//}

	return iface, nil
}
