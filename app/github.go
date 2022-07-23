package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"
)

var wgGithub sync.WaitGroup

type latestReleaseResponse struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
}

func CheckNewRealeases() {
	for _, k := range Chains.Chain {
		if k.Info.Github != "" {
			ri := Parse(k.Info.Github)
			if ri.Domain == "github.com" {
				lRR := new(latestReleaseResponse)
				wgGithub.Add(1)
				go lRR.Monitor(ri)
			} else {
				log.Printf("Repo %s Not Found or can't be parsed", k.Info.Github)
			}
		}
	}
	wgGithub.Wait()
}

type repoInfo struct {
	Domain           string
	Owner            string
	RepoName         string
	latestReleaseTag string
}

// Parse github repo from config
func Parse(u string) repoInfo {
	//remove / at the end
	if u[len(u)-1:] == "/" {
		u = u[:len(u)-1]
	}

	var re = regexp.MustCompile(`^(https|git)(:\/\/|@)(?P<Domain>[^\/:]+)[\/:](?P<Owner>[^\/:]+)\/(?P<RepoName>.+)$`)
	match := re.FindStringSubmatch(u)

	results := map[string]string{}
	for i, name := range match {
		results[re.SubexpNames()[i]] = name
	}

	return repoInfo{Domain: results["Domain"], Owner: results["Owner"], RepoName: results["RepoName"]}
}

// Monitor check github releases
func (target *latestReleaseResponse) Monitor(ri repoInfo) {
	defer wgGithub.Done()
	for {
		eP := fmt.Sprintf("https://api.github.com/repos/" + ri.Owner + "/" + ri.RepoName + "/releases/latest")
		resp, err := http.Get(eP)

		if err != nil {
			log.Printf("Something went wrong when tried to get repo: %s", err)
		}

		err = json.NewDecoder(resp.Body).Decode(target)
		if err != nil {
			log.Printf("Something went wrong when tried to parse response: %s", err)
		}

		switch ri.latestReleaseTag {
		case "":
			ri.latestReleaseTag = target.TagName
		case target.TagName:
		default:
			releaseDesc := ""
			if target.Body != "" {
				releaseDesc = fmt.Sprintf("Release desc: %s", target.Body)
			}
			e := Event{
				TagName:     target.TagName,
				RepoName:    ri.RepoName,
				ReleaseDesc: releaseDesc,
			}
			e.NewAlertGithubRelease().Send()
			ri.latestReleaseTag = target.TagName
		}
		// Rate limit 60 per hour , we will do only 12 requests per hour
		time.Sleep(300 * time.Second)
	}
}
