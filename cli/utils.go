package cli

import (
	"errors"
	"fmt"
	Openapi "jcli/client"
	"os"
	"strings"
	"time"
)

func NewHTTPClient() *Openapi.ClientWithResponses {
	client, err := Openapi.NewClientWithResponses(jocker_engine_url)
	if err != nil {
		fmt.Println("Internal error: ", err)
		os.Exit(1)
	}
	return client
}

type ResponseWithCode interface {
	StatusCode() int
}

func verify_response(response ResponseWithCode, expected_status int, err error) error {
	if err != nil {
		fmt.Println("Could not connect to jocker engine daemon: ", err)
		return err
	}

	if response.StatusCode() != expected_status {
		fmt.Println("Jocker engine returned unsuccesful statuscode: ", response.StatusCode())
		return errors.New("unsuccesful statuscode")
	}
	return nil
}

// HumanDuration returns a human-readable approximation of a duration
// (eg. "About a minute", "4 hours ago", etc.).
func HumanDuration(d time.Duration) string {
	if seconds := int(d.Seconds()); seconds < 1 {
		return "Less than a second"
	} else if seconds == 1 {
		return "1 second"
	} else if seconds < 60 {
		return fmt.Sprintf("%d seconds", seconds)
	} else if minutes := int(d.Minutes()); minutes == 1 {
		return "About a minute"
	} else if minutes < 60 {
		return fmt.Sprintf("%d minutes", minutes)
	} else if hours := int(d.Hours() + 0.5); hours == 1 {
		return "About an hour"
	} else if hours < 48 {
		return fmt.Sprintf("%d hours", hours)
	} else if hours < 24*7*2 {
		return fmt.Sprintf("%d days", hours/24)
	} else if hours < 24*30*2 {
		return fmt.Sprintf("%d weeks", hours/24/7)
	} else if hours < 24*365*2 {
		return fmt.Sprintf("%d months", hours/24/30)
	}
	return fmt.Sprintf("%d years", int(d.Hours())/24/365)
}

func Cell(word string, max_len int) string {
	word_length := len(word)

	if word_length <= max_len {
		return word + Sp(max_len-word_length) + Sp(2)
	} else {
		return word[:max_len] + ".."
	}
}

func Sp(n int) string {
	return strings.Repeat(" ", n)
}
