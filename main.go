package main

import "fmt"
import "net/http"
import "io/ioutil"
import "time"
import "os"
import "bytes"
import (
	"encoding/json"
)

type slackRequest struct {
	Text string `json:"text"`
}

type Configuration struct {
	Username    string
	Password 	string
	SlackUrl 	string
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
	doPost()
}

func doPost() {

	configuration, err := readConfig()


	if(err != nil){
		fmt.Println("weep womp")
	}


	username := configuration.Username
	password := configuration.Password
	slackUrl := configuration.SlackUrl

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

	fmt.Println(slackUrl)
	resp1, err := http.Post(slackUrl, "application/json; charset=utf-8", b)
	fmt.Println(err);
	fmt.Println(resp1);
}


func readConfig() (Configuration, error) {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(configuration)
	return configuration, err
}
