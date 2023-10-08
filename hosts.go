package hosts

import (
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/no-src/log"
)

func isWindows() bool {
	return runtime.GOOS == "windows"
}

// PrintHosts print the hosts info
func PrintHosts() {
	hostsFile := "/etc/hosts"
	if isWindows() {
		sysRoot := os.Getenv("SYSTEMROOT")
		hostsFile = sysRoot + `\System32\drivers\etc\hosts`
	}

	hostsBytes, err := os.ReadFile(hostsFile)
	if err != nil {
		log.Error(err, "read hosts file error")
		return
	}

	var hosts HostItemList

	hostsStr := string(hostsBytes)
	// index 0 host row
	// index 1 ip address
	// index 2 host list
	// index 3 last host matched
	reg := regexp.MustCompile(`(\S+)(([\t ]+\S+)+)\r?\n?`)
	hostMatches := reg.FindAllStringSubmatch(hostsStr, -1)
	for _, match := range hostMatches {
		if len(match) >= 4 {
			ip := match[1]
			hostnameStr := match[2]
			if strings.Contains(ip, "#") == false {
				hashIndex := strings.Index(hostnameStr, "#")
				if hashIndex > 0 {
					hostnameStr = hostnameStr[:hashIndex]
				}
				// index 0 hostname,contain whitespace
				// index 1 hostname,not contain whitespace
				hostnameReg := regexp.MustCompile(`[\t ]+(\S+)`)
				hostNameMatches := hostnameReg.FindAllStringSubmatch(hostnameStr, -1)
				var hostNameList []string
				for _, hostNameMatch := range hostNameMatches {
					hostNameList = append(hostNameList, hostNameMatch[1])
				}
				hosts = append(hosts, &HostItem{
					IP:           ip,
					HostNameList: hostNameList,
				})

			}
		}
	}

	log.Log("hosts[%s]", hostsFile)
	log.Log("-------------------------------------------------")
	hosts = recombine(hosts)
	for _, item := range hosts {
		log.Log("%s %s", item.IP, strings.Join(item.HostNameList, " "))
	}
}

// recombine try to recombine hosts,group by ip and distinct,order by ip
func recombine(hosts HostItemList) (newHosts HostItemList) {
	maps := make(map[string]*HostItem)
	for _, hostItem := range hosts {
		h := maps[hostItem.IP]
		if h == nil {
			h = &HostItem{
				IP: hostItem.IP,
			}
			newHosts = append(newHosts, h)
			maps[hostItem.IP] = h
		}
		h.HostNameList = append(h.HostNameList, hostItem.HostNameList...)
	}

	for _, hostItem := range newHosts {
		hostItem.HostNameList = distinct(hostItem.HostNameList)
	}

	sort.Sort(sort.Reverse(newHosts))
	return newHosts
}

func distinct(source []string) (target []string) {
	if source == nil {
		return nil
	}
	distinctMap := make(map[string]bool)
	for _, item := range source {
		exist := distinctMap[item]
		if !exist {
			target = append(target, item)
			distinctMap[item] = true
		}
	}
	sort.Strings(target)
	return target
}

type HostItem struct {
	IP           string
	HostNameList []string
}

type HostItemList []*HostItem

func (h HostItemList) Len() int {
	return len(h)
}
func (h HostItemList) Less(i, j int) bool {
	return strings.Compare(h[i].IP, h[j].IP) > 0
}
func (h HostItemList) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
