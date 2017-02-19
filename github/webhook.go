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
	}
}

type WebHookPayload struct {
	Payload Payload `json:"payload"`
}

type Payload struct {
	Commits    []Commit `json:"commits"`
	CompareUrl string `json:"compare"`
	HeadCommit Commit `json:"head_commit"`
	Pusher     Author `json:"pusher"`
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
	buffer.Write([]byte{240, 159, 154, 128})
	msg := tgbotapi.NewMessage(user.ChatId, buffer.String())
	h.sammy.Api.Send(msg)
}

func (h *Hook) pushEvent(user *sammy.User, req *http.Request) {
	var payload Payload
	var buffer bytes.Buffer
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&payload)
	check(err, "could not decode request values because: %v")
	buffer.Write([]byte{226, 172, 134})
	buffer.WriteString(payload.Pusher.Name + " has **pushed** " + strconv.Itoa(len(payload.Commits)) + " commits, some are: \n")
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

func check(err error, msg string) {
	if err != nil {
		log.Printf(msg, err)
	}
}
