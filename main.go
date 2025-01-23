package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
)

type Task struct {
	Type        string   `json:"type"`
	CountryCode string   `json:"country_code,omitempty"`     // Optional for tasks "add" "remove"
	Domains     []string `json:"domains,omitempty"`          // Optional for tasks "add" "remove"
	Entries     []string `json:"entries,omitempty"`          // Optional for "delete" tasks
	SrcCountry  string   `json:"src_country_code,omitempty"` // Optional for "copy" tasks
	DstCountry  string   `json:"dst_country_code,omitempty"` // Optional for "copy" tasks
}

type TaskList struct {
	Src   string `json:"src"`
	Dst   string `json:"dst"`
	Tasks []Task `json:"tasks"`
}

func main() {
	taskFile := flag.String("t", "./tasks.json", "Specify the location of the tasks.json file")
	help := flag.Bool("h", false, "Show usage information")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	file, err := os.ReadFile(*taskFile)
	if err != nil {
		fmt.Println("Error reading task file:", err)
		return
	}

	var taskList TaskList
	if err := json.Unmarshal(file, &taskList); err != nil {
		fmt.Println("Error unmarshalling task file:", err)
		return
	}

	data, err := os.ReadFile(taskList.Src)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	geositeList := &routercommon.GeoSiteList{}
	if err := proto.Unmarshal(data, geositeList); err != nil {
		fmt.Println("Error unmarshalling data:", err)
		return
	}

	for _, task := range taskList.Tasks {
		switch task.Type {
		case "add", "remove":
			for _, entry := range geositeList.Entry {
				if task.CountryCode == "*" || strings.EqualFold(entry.CountryCode, task.CountryCode) {
					if task.Type == "remove" {
						for _, domainToRemove := range task.Domains {
							var filteredDomains []*routercommon.Domain
							for _, domain := range entry.Domain {
								if domain.GetValue() != domainToRemove {
									filteredDomains = append(filteredDomains, domain)
								} else {
									fmt.Printf("Removed %s from %s\n", domainToRemove, entry.CountryCode)
								}
							}
							entry.Domain = filteredDomains
						}
					} else if task.Type == "add" {
						for _, domainToAdd := range task.Domains {
							exists := false
							for _, domain := range entry.Domain {
								if domain.GetValue() == domainToAdd {
									exists = true
									fmt.Printf("%s already exists in %s\n", domainToAdd, entry.CountryCode)
									break
								}
							}
							if !exists {
								newDomain := &routercommon.Domain{
									Type:  routercommon.Domain_RootDomain,
									Value: domainToAdd,
								}
								entry.Domain = append(entry.Domain, newDomain)
								fmt.Printf("Added %s to %s\n", domainToAdd, entry.CountryCode)
							}
						}
					}
				}
			}
		case "copy":
			var srcEntry *routercommon.GeoSite
			var dstEntry *routercommon.GeoSite

			for _, entry := range geositeList.Entry {
				if strings.EqualFold(entry.CountryCode, task.SrcCountry) {
					srcEntry = entry
				}
				if strings.EqualFold(entry.CountryCode, task.DstCountry) {
					dstEntry = entry
				}
			}

			if dstEntry == nil {
				dstEntry = &routercommon.GeoSite{CountryCode: task.DstCountry}
				geositeList.Entry = append(geositeList.Entry, dstEntry)
			}

			if srcEntry != nil {
				existingDomains := make(map[string]struct{})
				for _, dstDomain := range dstEntry.Domain {
					existingDomains[dstDomain.GetValue()] = struct{}{}
				}

				for _, domain := range srcEntry.Domain {
					if _, exists := existingDomains[domain.GetValue()]; !exists {
						dstEntry.Domain = append(dstEntry.Domain, domain)
						fmt.Printf("Copied %s from %s to %s\n", domain.GetValue(), task.SrcCountry, task.DstCountry)
					}
				}
			}
		case "delete":
			for i := 0; i < len(geositeList.Entry); i++ {
				entry := geositeList.Entry[i]
				for _, countryCodeToDelete := range task.Entries {
					if strings.EqualFold(entry.CountryCode, countryCodeToDelete) {
						fmt.Printf("Deleting entry for %s\n", entry.CountryCode)
						geositeList.Entry = append(geositeList.Entry[:i], geositeList.Entry[i+1:]...)
						i--
						break
					}
				}
			}
		}
	}

	/*
		for _, entry := range geositeList.Entry {
			sort.Slice(entry.Domain, func(i, j int) bool {
				return entry.Domain[i].GetValue() < entry.Domain[j].GetValue()
			})
		}
	*/

	modifiedData, err := proto.Marshal(geositeList)
	if err != nil {
		fmt.Println("Error marshalling modified data:", err)
		return
	}

	if err := os.WriteFile(taskList.Dst, modifiedData, 0644); err != nil {
		fmt.Println("Error writing modified file:", err)
		return
	}

	fmt.Printf("Modified geosite.dat has been saved to '%s'.\n", taskList.Dst)
}
