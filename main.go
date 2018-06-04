package main

import "fmt"
import "net/http"
import "io/ioutil"
import "time"
import s "strings"
import "os"
import "bytes"
import (
	"encoding/json"
	"io"
)

type slackRequest struct {
	Text string `json:"text"`
}

type Game struct {
	ID                       string      `json:"ID"`
	ScheduleStatus           string      `json:"scheduleStatus"`
	OriginalDate             interface{} `json:"originalDate"`
	OriginalTime             interface{} `json:"originalTime"`
	DelayedOrPostponedReason interface{} `json:"delayedOrPostponedReason"`
	Date                     string      `json:"date"`
	Time                     string      `json:"time"`
	AwayTeam                 struct {
		ID           string `json:"ID"`
		City         string `json:"City"`
		Name         string `json:"Name"`
		Abbreviation string `json:"Abbreviation"`
	} `json:"awayTeam"`
	HomeTeam struct {
		ID           string `json:"ID"`
		City         string `json:"City"`
		Name         string `json:"Name"`
		Abbreviation string `json:"Abbreviation"`
	} `json:"homeTeam"`
	Location string `json:"location"`
}

type Scores struct {
	Scoreboard struct {
		LastUpdatedOn string `json:"lastUpdatedOn"`
		GameScore     []struct {
			Game struct {
				ID                       string      `json:"ID"`
				ScheduleStatus           string      `json:"scheduleStatus"`
				OriginalDate             interface{} `json:"originalDate"`
				OriginalTime             interface{} `json:"originalTime"`
				DelayedOrPostponedReason interface{} `json:"delayedOrPostponedReason"`
				Date                     string      `json:"date"`
				Time                     string      `json:"time"`
				AwayTeam                 struct {
					ID           string `json:"ID"`
					City         string `json:"City"`
					Name         string `json:"Name"`
					Abbreviation string `json:"Abbreviation"`
				} `json:"awayTeam"`
				HomeTeam struct {
					ID           string `json:"ID"`
					City         string `json:"City"`
					Name         string `json:"Name"`
					Abbreviation string `json:"Abbreviation"`
				} `json:"homeTeam"`
				Location string `json:"location"`
			} `json:"game"`
			IsUnplayed    string      `json:"isUnplayed"`
			IsInProgress  string      `json:"isInProgress"`
			IsCompleted   string      `json:"isCompleted"`
			PlayStatus    interface{} `json:"playStatus"`
			AwayScore     string      `json:"awayScore,omitempty"`
			HomeScore     string      `json:"homeScore,omitempty"`
			InningSummary struct {
				Inning []struct {
					Number    string `json:"@number"`
					AwayScore string `json:"awayScore"`
					HomeScore string `json:"homeScore"`
				} `json:"inning"`
			} `json:"inningSummary"`
			CurrentIntermission string `json:"currentIntermission,omitempty"`
			CurrentInning       string `json:"currentInning,omitempty"`
			CurrentInningHalf   string `json:"currentInningHalf,omitempty"`
		} `json:"gameScore"`
	}
}

func main() {

	dat, err := ioutil.ReadFile("./config")

	configs := s.Split(string(dat), "\n")
	configs = append(configs[:2], configs[2+1:]...)
	for _, element := range configs {
		envName := s.Split(element, "=")
		os.Setenv(envName[0], envName[1])
	}

	username := os.Getenv("username")
	password := os.Getenv("password")
	slackUrl := os.Getenv("slackUrl")

	currentTime := time.Now().Local().Format("20060102")
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.mysportsfeeds.com/v1.2/pull/mlb/2018-regular/scoreboard.json?fordate="+currentTime, nil)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("weep womp")
	}

	var scoreBoard Scores
	bodyText, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(bodyText, &scoreBoard); err != nil {
		panic(err)
	}

	var gameSummary = ""
	for _, element := range scoreBoard.Scoreboard.GameScore {
		gameSummary += element.Game.HomeTeam.Name + ": " + element.HomeScore + " vs " + element.Game.AwayTeam.Name + ": " + element.AwayScore + "\n"
	}


	scoresPost := slackRequest{Text: gameSummary}
	fmt.Println(slackRequest(scoresPost))
	b := new(bytes.Buffer)

	json.NewEncoder(b).Encode(scoresPost)

	resp1, _ := http.Post(slackUrl, "application/json; charset=utf-8", b)
	io.Copy(os.Stdout, resp1.Body)

}
