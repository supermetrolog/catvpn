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

	logrus.Debugf("SRC: %s -- DST: %s -- ID: %d -- CHECKSUM: %d -- TOTAL LEN: %d", header.Src, header.Dst, header.ID, header.Checksum, header.TotalLen)
}
