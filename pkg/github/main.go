package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Github contains settings to interact with Github
type Github struct {
	Owner        string
	Repository   string
	Username     string
	Token        string
	URL          string
	Version      string
	directory    string
	Branch       string
	remoteBranch string
	User         string
	Email        string
}

// GetDirectory returns the local git repository path
func (g *Github) GetDirectory() (directory string) {
	return g.directory
}

// GetVersion retrieves the version tag from Github Releases
func (g *Github) GetVersion() string {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%s",
		g.Owner,
		g.Repository,
		g.Version)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", "token "+g.Token)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
	}

	v := map[string]string{}
	json.Unmarshal(body, &v)

	if val, ok := v["name"]; ok {
		return val
	}
	fmt.Printf("\u2717 No tag founded from %s\n", url)
	return ""

}

// Init Github struct
func (g *Github) Init(version string) {
	g.Version = version
	g.setDirectory(version)
	g.remoteBranch = fmt.Sprintf("updatecli/%v", version)

}

func (g *Github) setDirectory(version string) {

	directory := fmt.Sprintf("%v/%v/%v/%v", os.TempDir(), g.Owner, g.Repository, g.Version)

	if _, err := os.Stat(directory); os.IsNotExist(err) {

		err := os.MkdirAll(directory, 0755)
		if err != nil {
			fmt.Println(err)
		}
	}

	g.directory = directory

	fmt.Printf("Directory: %v\n", g.directory)
}

// Clean Github working directory
func (g *Github) Clean() {
	os.RemoveAll(g.directory)
}

// Clone run `git clone`
func (g *Github) Clone() string {
	URL := fmt.Sprintf("https://%v:%v@github.com/%v/%v.git",
		g.Username,
		g.Token,
		g.Owner,
		g.Repository)
	_, err := git.PlainClone(g.directory, false, &git.CloneOptions{
		URL:        URL,
		RemoteName: g.Branch,
		Progress:   os.Stdout,
	})

	g.Checkout()

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("\t\t%s downloaded in %s", URL, g.directory)
	return g.directory
}

// Commit run `git commit`
func (g *Github) Commit(file, message string) {
	r, err := git.PlainOpen(g.directory)
	if err != nil {
		fmt.Println(err)
	}

	w, err := r.Worktree()
	if err != nil {
		fmt.Println(err)
	}

	status, err := w.Status()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(status)

	commit, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  g.User,
			Email: g.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	obj, err := r.CommitObject(commit)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(obj)

}

// Checkout create and use a temporary branch
func (g *Github) Checkout() {
	r, err := git.PlainOpen(g.directory)
	if err != nil {
		fmt.Println(err)
	}

	w, err := r.Worktree()
	if err != nil {
		fmt.Println(err)
	}

	branch := "updatecli/" + g.Version
	fmt.Printf("Creating Branch: %v\n", branch)

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Create: true,
		Force:  false,
		Keep:   true,
	})

	if err != nil {
		fmt.Println(err)
	}
}

// Add run `git add`
func (g *Github) Add(file string) {

	fmt.Printf("\t\tAdding file: %s", file)

	r, err := git.PlainOpen(g.directory)
	if err != nil {
		fmt.Println(err)
	}

	w, err := r.Worktree()
	if err != nil {
		fmt.Println(err)
	}

	_, err = w.Add(file)
	if err != nil {
		fmt.Println(err)
	}
}

// Push run `git push`
func (g *Github) Push() {

	r, err := git.PlainOpen(g.directory)
	if err != nil {
		fmt.Println(err)
	}

	URL := fmt.Sprintf("https://%v:%v@github.com/%v/%v.git",
		g.Username,
		g.Token,
		g.Owner,
		g.Repository)

	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: g.remoteBranch,
		URLs: []string{URL},
	})

	err = r.Push(&git.PushOptions{
		RemoteName: g.remoteBranch,
		Progress:   os.Stdout,
	})
	if err != nil {
		fmt.Println(err)
	}

	g.OpenPR()
}

// OpenPR creates a new pull request
func (g *Github) OpenPR() {
	title := fmt.Sprintf("[Updatecli] Bump to version %v", g.Version)

	if g.isPRExist(title) {
		fmt.Println("PR already exist")
		return
	}

	bodyPR := "Please pull these awesome changes in!"

	URL := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls",
		g.Owner,
		g.Repository)

	jsonData := fmt.Sprintf("{ \"title\": \"%v\", \"body\": \"%v\", \"head\": \"%v\", \"base\": \"%v\"}", title, bodyPR, g.remoteBranch, g.Branch)

	var jsonStr = []byte(jsonData)

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonStr))

	req.Header.Add("Authorization", "token "+g.Token)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
	}

	v := map[string]string{}
	err = json.Unmarshal(body, &v)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res.Status)

}

// isPRExist checks if an open pull request already exist based on a title
func (g *Github) isPRExist(title string) bool {

	URL := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls",
		g.Owner,
		g.Repository)

	req, err := http.NewRequest("GET", URL, nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", "token "+g.Token)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
	}

	v := []map[string]string{}

	err = json.Unmarshal(body, &v)

	if err != nil {
		fmt.Println(err)
	}

	for _, v := range v {
		if v["title"] == title {
			return true
		}
	}

	return false
}
