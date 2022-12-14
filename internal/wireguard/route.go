package wireguard

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/gluetun/internal/netlink"
)

func (w *Wireguard) addRoutes(link netlink.Link, destinations []*net.IPNet,
	firewallMark int) (err error) {
	for _, dst := range destinations {
		err = w.addRoute(link, dst, firewallMark)
		if err == nil {
			continue
		}

		ipv6Dst := dst.IP.To4() == nil
		if ipv6Dst && strings.Contains(err.Error(), "permission denied") {
			w.logger.Errorf("cannot add route for IPv6 due to a permission denial. "+
				"Ignoring and continuing execution; "+
				"Please report to https://github.com/qdm12/gluetun/issues/998 if you find a fix. "+
				"Full error string: %s", err)
			continue
		}
		return fmt.Errorf("adding route for destination %s: %w", dst, err)
	}
	return nil
}

func (w *Wireguard) addRoute(link netlink.Link, dst *net.IPNet,
	firewallMark int) (err error) {
	route := &netlink.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       dst,
		Table:     firewallMark,
	}

	err = w.netlink.RouteAdd(route)
	if err != nil {
		return fmt.Errorf(
			"cannot add route for link %s, destination %s and table %d: %w",
			link.Attrs().Name, dst, firewallMark, err)
	}

	return err
}
