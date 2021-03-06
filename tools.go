package whois

import (
	"regexp"
	"strings"
)

func parser(re *regexp.Regexp, group int, data string) (result []string) {
	found := re.FindAllStringSubmatch(data, -1)
	if len(found) > 0 {
		for _, one := range found {
			if len(one) >= 2 && len(one[group]) > 0 {
				result = appendIfMissing(result, one[group])
			}
		}
	}
	return
}

func ParseWhoisServer(whois string) string {
	data := parser(regexp.MustCompile(`(?m:^)(\s+)?(?i)(Whois server|whois):\s+(.*?)(\s|$)`), 3, whois)
	res := ""
	if len(data) > 0 {
		res = data[0]
	}
	return res
}

func ParseReferServer(whois string) string {
	return parser(regexp.MustCompile(`(?i)(refer):\s+(.*?)(\s|$)`), 2, whois)[0]
}

//Parse uniq name servers from whois
func ParseNameServers(whois string) []string {
	return parser(regexp.MustCompile(`(?i)(Name Server|nserver):\s+(.*?)(\s|$)`), 2, whois)
}

//Parse uniq domain status(codes) from whois
func ParseDomainStatus(whois string) []string {
	return parser(regexp.MustCompile(`(?i)(Domain )?(Status|state):\s+(.*?)(\s|$)`), 3, whois)
}

func IsWhoisDataCorrect(whois string) bool {
	if len(whois) == 0 {
		return false
	}
	hasDomain := parser(regexp.MustCompile(`(?i)(Domain Name|domain):\s+(.*?)(\s|$)`), 2, whois)
	return len(hasDomain) > 0
}

func ParseNofound(whois string) bool {
	data := parser(regexp.MustCompile(`(?i)(Not found|No match for|No entries found)(\s|$)`), 1, whois)
	if len(data) > 0 {
		return true
	}
	return false
}

func appendIfMissing(slice []string, i string) []string {
	i = strings.ToLower(i)

	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func whoisWeight(whois string) (res int) {
	res += len(parser(regexp.MustCompile(`(?i)(address)(\s|$)`), 1, whois))
	res += len(parser(regexp.MustCompile(`(?i)(phone)(\s|$)`), 1, whois))
	res += len(parser(regexp.MustCompile(`(?i)(e-mail)(\s|$)`), 1, whois))
	res += len(parser(regexp.MustCompile(`(?i)(Admin)(\s|$)`), 1, whois))
	res += len(parser(regexp.MustCompile(`(?i)(Tech)(\s|$)`), 1, whois))
	res += len(parser(regexp.MustCompile(`(?i)(Registrant)(\s|$)`), 1, whois))
	res += len(parser(regexp.MustCompile(`(?i)(person)(\s|$)`), 1, whois))
	res += len(parser(regexp.MustCompile(`(?i)(PostalCode)(\s|$)`), 1, whois))
	res += len(parser(regexp.MustCompile(`(?i)(CountryCode)(\s|$)`), 1, whois))
	res += len(parser(regexp.MustCompile(`(?i)(Fax)(\s|$)`), 1, whois))
	return
}
