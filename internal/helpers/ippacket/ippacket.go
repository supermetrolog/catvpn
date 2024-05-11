package ippacket

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/ipv4"
)

func LogHeader(frame []byte) {
	header, err := ipv4.ParseHeader(frame)

	if err != nil {
		logrus.Warnf("Parse ip header error: %v", err)
		return
	}

	logrus.Infof("SRC: %s -> DST: %s", header.Src, header.Dst)
	logrus.Debugf("PROTOCOL: %d; ID: %d; CHECKSUM: %d; TOTAL LEN: %d;", header.Protocol, header.ID, header.Checksum, header.TotalLen)
}
