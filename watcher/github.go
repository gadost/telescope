package watcher

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/gadost/telescope/alert"
)

var wgGithub sync.WaitGroup

type latestReleaseResponse struct {
	URL             string    `json:"url"`
	HTMLURL         string    `json:"html_url"`
	AssetsURL       string    `json:"assets_url"`
	UploadURL       string    `json:"upload_url"`
	TarballURL      string    `json:"tarball_url"`
	ZipballURL      string    `json:"zipball_url"`
	DiscussionURL   string    `json:"discussion_url"`
	ID              int       `json:"id"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Body            string    `json:"body"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Author          struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Assets []struct {
		URL                string    `json:"url"`
		BrowserDownloadURL string    `json:"browser_download_url"`
		ID                 int       `json:"id"`
		NodeID             string    `json:"node_id"`
		Name               string    `json:"name"`
		Label              string    `json:"label"`
		State              string    `json:"state"`
		ContentType        string    `json:"content_type"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		Uploader           struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
	} `json:"assets"`
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
			log.Printf("Something went wrong when tried to parse reponse: %s", err)
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
			alert.NewAlertGithubRelease(target.TagName, ri.RepoName, releaseDesc).Send()
			ri.latestReleaseTag = target.TagName
		}
		// Rate limit 60 per hour , we will do only 12 requests per hour
		time.Sleep(300 * time.Second)
	}
}
