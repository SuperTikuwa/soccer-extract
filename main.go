package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	DIR    = "~/soccer"
	TEAM_L = "FC_AIT"
	TEAM_R = "HELIOS_enemy"
)

type Match struct {
	teamL      string
	teamR      string
	teamLScore int
	teamRScore int
}

func main() {
	out, err := exec.Command("ls", "../../soccer").Output()
	if err != nil {
		log.Fatal(out, err)
	}

	tmp := strings.Split(string(out), "\n")
	logs := make([]string, 0)
	for _, v := range tmp {
		fmt.Println(v)

		if strings.Contains(v, TEAM_L) && strings.Contains(v, TEAM_R) && strings.Contains(v, ".rcg") {
			logs = append(logs, v)
		}
	}

	fmt.Println(logs)

	matches := make([]Match, 0)
	for _, v := range logs {
		matches = append(matches, logNameToMatchData(v))
	}

	fmt.Println(len(matches))

	extract(matches)

}

func logNameToMatchData(logName string) Match {
	tmp := strings.Split(logName, "-vs-")
	reg := regexp.MustCompile(`^[0-9]+-`)
	teamL := reg.ReplaceAllString(tmp[0], "")
	teamR := strings.ReplaceAll(tmp[1], ".rcg", "")

	penaltyReg := regexp.MustCompile(`_[0-9]`)
	if len(penaltyReg.FindAll([]byte(teamL), -1)) >= 2 {
		r := regexp.MustCompile(`_[0-9]$`)
		teamL = r.ReplaceAllString(teamL, "")
		teamR = r.ReplaceAllString(teamR, "")
	}

	l := strings.Split(teamL, "_")
	r := strings.Split(teamR, "_")

	teamL = strings.Join(l[:len(l)-2], "_")
	teamR = strings.Join(r[:len(r)-2], "_")

	scoreL, err := strconv.Atoi(l[len(l)-1])
	if err != nil {
		log.Fatal(err)
	}

	scoreR, err := strconv.Atoi(r[len(r)-1])
	if err != nil {
		log.Fatal(err)
	}

	return Match{teamL, teamR, scoreL, scoreR}
}

func extract(matches []Match) {
	sumL := 0
	winL := 0
	sumR := 0
	winR := 0
	draw := 0

	for _, v := range matches {
		sumL += v.teamLScore
		sumR += v.teamRScore
		if v.teamLScore > v.teamRScore {
			winL++
		} else if v.teamLScore < v.teamRScore {
			winR++
		} else {
			draw++
		}
	}

	fmt.Printf("Team: %s\n", TEAM_L)
	fmt.Printf("Win: %d\n", winL)
	fmt.Printf("Lose: %d\n", winR)
	fmt.Printf("Draw: %d\n", draw)
	fmt.Printf("Score: %d\n", sumL)
	fmt.Printf("Average: %f\n", float64(sumL)/float64(len(matches)))
	fmt.Println()
	fmt.Printf("Team: %s\n", TEAM_R)
	fmt.Printf("Win: %d\n", winR)
	fmt.Printf("Lose: %d\n", winL)
	fmt.Printf("Draw: %d\n", draw)
	fmt.Printf("Score: %d\n", sumR)
	fmt.Printf("Average: %f\n", float64(sumR)/float64(len(matches)))

}
