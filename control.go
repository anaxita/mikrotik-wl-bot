package main

import (
	"errors"
	"fmt"
	"github.com/go-routeros/routeros"
	"net"
)

type RouterController struct {
}

func NewRouterController() *RouterController {
	return &RouterController{}
}

func (rc *RouterController) AddIP(ip net.IP) error {
	conn, err := routeros.Dial(routerAddr, routerLogin, routerPassword)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Run("/ip/firewall/address-list/add", "=list=WL", "=address="+ip.String())
	if err != nil {
		return err
	}

	return nil
}

func (rc *RouterController) RemoveIP(ip net.IP) error {
	conn, err := routeros.Dial(routerAddr, routerLogin, routerPassword)
	if err != nil {
		return err
	}
	defer conn.Close()

	findIP, err := conn.Run("/ip/firewall/address-list/print", fmt.Sprintf("address=%s", ip.String()), ".proplist=.id")
	if err != nil {
		return err
	}

	if len(findIP.Re) <= 0 {
		return errors.New("ip is not found")
	}

	ipID, ok := findIP.Re[0].Map[".id"]
	if !ok {
		return errors.New("ip is not found")
	}

	_, err = conn.Run("/ip/firewall/address-list/remove", "=.id="+ipID)
	if err != nil {
		return err
	}

	return nil

}
