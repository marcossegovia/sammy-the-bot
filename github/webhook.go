package github

import (
	"io"
	"net/http"
	"fmt"
	"regexp"
	"strconv"
	"log"
	"bytes"
	"time"
	"encoding/json"

	"github.com/MarcosSegovia/sammy-the-bot/sammy"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Hook struct {
	sammy *sammy.Sammy
}

func NewHook(s *sammy.Sammy) *Hook {
	hook := new(Hook)
	hook.sammy = s
	return hook
}

func (h *Hook) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Hook received !")

	r, err := regexp.Compile("/github/hooks/([0-9]+)")
	check(err, "could not set regular expression for github hooks: %v")
	matches := r.FindStringSubmatch(req.URL.Path)
	if matches[1] == "" {
		fmt.Errorf("payload failed to send a valid chatId : %v", matches[1])
	}
	chatId, err := strconv.ParseInt(matches[1], 10, 64)
	user, err := h.sammy.GetUser(chatId)
	if err != nil {
		check(err, "could not get user because: %v")
		return
	}

	switch req.Header.Get("X-GitHub-Event") {
	case "ping":
		h.pingEvent(user, req)
	case "push":
		h.pushEvent(user, req)
	case "pull_request":
		h.pullRequestEvent(user, req)
	}
}

type WebHookPayload struct {
	Payload Payload `json:"payload"`
}

type Payload struct {
	Ref         string `json:"ref"`
	Action      string `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
	Created     bool `json:"created"`
	Deleted     bool `json:"deleted"`
	Forced      bool `json:"forced"`
	CompareUrl  string `json:"compare"`
	Commits     []Commit `json:"commits"`
	HeadCommit  Commit `json:"head_commit"`
	Pusher      Author `json:"pusher"`
}

func (p Payload) BranchName() string {
	r, err := regexp.Compile("refs/heads/(.*)")
	check(err, "could not set regular expression for github hooks: %v")
	matches := r.FindStringSubmatch(p.Ref)
	if matches[1] == "" {
		fmt.Errorf("payload failed to send a valid branch name : %v", matches[1])
	}

	return matches[1]
}

type PullRequest struct {
	Id               int `json:"number"`
	State            string `json:"state"`
	Title            string `json:"title"`
	Author           User `json:"user"`
	Body             string `json:"body"`
	CreatedAt        time.Time `json:"created_at"`
	Url              string `json:"html_url"`
	RequestReviewers []User `json:"requested_reviewers"`
	Merged           bool `json:"merged"`
}

type User struct {
	Id    int `json:"id"`
	Login string `json:"login"`
}

type Commit struct {
	Id        string `json:"id"`
	TreeId    string `json:"tree_id"`
	Message   string `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Author    Author `json:"author"`
	Committer Author `json:"committer"`
	Url       string `json:"url"`
}

type Author struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username,omitempty"`
}

func (h *Hook) pingEvent(user *sammy.User, req *http.Request) {
	var buffer bytes.Buffer
	buffer.WriteString("Your hook has correctly being set ! ")
	buffer.WriteString("\U0001F680")
	msg := tgbotapi.NewMessage(user.ChatId, buffer.String())
	h.sammy.Api.Send(msg)
}

func (h *Hook) pushEvent(user *sammy.User, req *http.Request) {
	var payload Payload
	var buffer bytes.Buffer
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&payload)
	check(err, "could not decode request values because: %v")

	if payload.Deleted {
		buffer.WriteString("\U0000274C")
		buffer.WriteString(payload.Pusher.Name + " has *deleted* branch " + payload.BranchName())
		msg := tgbotapi.NewMessage(user.ChatId, buffer.String())
		msg.ParseMode = "Markdown"
		h.sammy.Api.Send(msg)
		return
	}

	buffer.WriteString("\U00002B06")
	buffer.WriteString(payload.Pusher.Name + " has *pushed* " + strconv.Itoa(len(payload.Commits)) + " commits to " + payload.BranchName() + ": \n")
	for _, commit := range payload.Commits {
		buffer.WriteString("> [" + commit.Id + "](" + commit.Url + ") " + commit.Message + " - " + commit.Committer.Name + "\n")
	}
	if len(payload.Commits) > 1 {
		buffer.WriteString("Go to the last commit >>> [" + payload.HeadCommit.Id + "](" + payload.HeadCommit.Url + ")")
	}

	msg := tgbotapi.NewMessage(user.ChatId, buffer.String())
	msg.ParseMode = "Markdown"
	h.sammy.Api.Send(msg)
}

func (h *Hook) pullRequestEvent(user *sammy.User, req *http.Request) {
	var payload Payload
	var buffer bytes.Buffer
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&payload)
	check(err, "could not decode request values because: %v")

	switch payload.Action {
	case "review_requested":
		buffer.WriteString("\U0001F3A9")
		buffer.WriteString(" " + payload.PullRequest.Author.Login + " has *requested a review* to ")
		for _, reviewer := range payload.PullRequest.RequestReviewers {
			buffer.WriteString("\U0001F46E")
			buffer.WriteString(" " + reviewer.Login + " ")
		}
		buffer.WriteString("\n in pull request [#" + strconv.Itoa(payload.PullRequest.Id) + "]("+ payload.PullRequest.Url + ")")
	case "opened":
		buffer.WriteString("\U0001F3A9")
		buffer.WriteString(payload.PullRequest.Author.Login + " has *opened a pull request* [#" + strconv.Itoa(payload.PullRequest.Id) + "]("+ payload.PullRequest.Url + ") \n")
	case "closed":
		buffer.WriteString("Pull request [#" + strconv.Itoa(payload.PullRequest.Id) + "](" + payload.PullRequest.Url + ") has been closed")
		if payload.PullRequest.Merged {
			buffer.WriteString(" and fully merged ")
			buffer.WriteString("\U00002705")
		}

	}

	msg := tgbotapi.NewMessage(user.ChatId, buffer.String())
	msg.ParseMode = "Markdown"
	h.sammy.Api.Send(msg)
}

func check(err error, msg string) {
	if err != nil {
		log.Printf(msg, err)
	}
}
