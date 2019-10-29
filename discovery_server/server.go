package main

import (
	"net/http"
	"strings"
	"fmt"
)

type HostList struct {
	hosts []string
}

// Server list constructor
func NewHostList() *HostList{
	// Construct with empty list
	return &HostList{
		hosts: []string{},
	}
}

// Server list register host method
func (hostList *HostList) registerHost(w http.ResponseWriter, req *http.Request){
	host := req.URL.Query().Get("host")
	hostList.hosts = append(hostList.hosts, host)
}

// Server list get hosts method
func (hostList *HostList) getHosts(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(w, "%s", strings.Join(hostList.hosts[:], "\n"))
}

func main() {
	hostList := NewHostList()
	http.HandleFunc("/register_host", hostList.registerHost)
	http.HandleFunc("/get_hosts", hostList.getHosts)

	// TODO: Need to set up a method to continually ping all of the hosts in the hostlist to see if they are still active
	http.ListenAndServe(fmt.Sprintf(":%d", 3020), nil)
}
