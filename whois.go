package whois

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

const (
	ZoneWhoisServer   = "whois.iana.org"
	defaulWhoisServer = "whois.arin.net"
	aliasWhoisServer  = ".whois-servers.net"
)

func GetWhois(domain string) (string, error) {
	return GetWhoisTimeout(domain, time.Second*5)
}

func GetWhoisTimeout(domain string, timeout time.Duration) (result string, err error) {
	s, e := GetPossibleWhoisServers(domain, timeout)
	if e != nil {
		err = e
		return
	}
	result = strings.Join(s, ",")
	return
}

func GetPossibleWhoisServers(domain string, timeout time.Duration) (whoisServers []string, err error) {
	parts := strings.Split(domain, ".")
	lenDomain := len(parts)
	if lenDomain < 2 {
		err = fmt.Errorf("Domain(%s) name is wrong!", domain)
		return
	}

	topLevelDomain := parts[lenDomain-1]
	if lenDomain > 2 {
		secondLevelDomain := parts[lenDomain-2] + "." + parts[lenDomain-1]
		whoisServers = append(whoisServers, "whois.nic."+secondLevelDomain)
		whoisServers = append(whoisServers, "whois."+secondLevelDomain)
	}

	whoisServers = append(whoisServers, "whois.nic."+topLevelDomain)
	whoisServers = append(whoisServers, "whois."+topLevelDomain)
	whoisServers = append(whoisServers, defaulWhoisServer)
	whoisServers = append(whoisServers, topLevelDomain+aliasWhoisServer)

	if whoisFromIANA := GetWhoisServerFromIANA(topLevelDomain, timeout); whoisFromIANA != "" {
		whoisServers = append(whoisServers, whoisFromIANA)
	}

	return
}

func GetWhoisServerFromIANA(zone string, timeout time.Duration) string {
	data, err := GetWhoisData(zone, ZoneWhoisServer, timeout)
	if err != nil {
		return ""
	}
	result := ParseWhoisServer(data)
	return result
}

func GetWhoisData(domain, server string, timeout time.Duration) (data string, err error) {
	var (
		connection net.Conn
		buffer     []byte
	)

	connection, err = net.DialTimeout("tcp", net.JoinHostPort(server, "43"), timeout)

	if err != nil {
		return
	}

	defer connection.Close()

	connection.Write([]byte(domain + "\r\n"))

	buffer, err = ioutil.ReadAll(connection)

	if err != nil {
		return
	}

	data = string(buffer[:])
	return
}
