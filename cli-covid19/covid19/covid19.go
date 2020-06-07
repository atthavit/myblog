package covid19

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/guptarohit/asciigraph"
)

type Type string

var (
	Confirmed    Type = "confirmed"
	Deaths       Type = "deaths"
	Hospitalized Type = "hospitalized"
	Recovered    Type = "recovered"
	Types             = []Type{
		Confirmed,
		Deaths,
		Hospitalized,
		Recovered,
	}
	URL = "https://covid19.th-stat.com/api/open/"
)

type Data struct {
	Confirmed       int `json:"Confirmed"`
	Deaths          int `json:"Deaths"`
	Hospitalized    int `json:"Hospitalized"`
	Recovered       int `json:"Recovered"`
	NewConfirmed    int `json:"NewConfirmed"`
	NewRecovered    int `json:"NewRecovered"`
	NewHospitalized int `json:"NewHospitalized"`
	NewDeaths       int `json:"NewDeaths"`
}

func PrintToday() error {
	resp, err := http.Get(URL + "today")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var data Data
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	fmt.Printf(
		"%s - New Confirmed: %d, New Deaths: %d, New Recovered: %d, Total Confirmed: %d\n",
		time.Now().Format("Jan 2 2006"),
		data.NewConfirmed,
		data.NewDeaths,
		data.NewRecovered,
		data.Confirmed,
	)
	return nil
}

func Plot(t Type) error {
	resp, err := http.Get(URL + "timeline")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var data struct {
		List []Data `json:"Data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	days := 30
	points := make([]float64, days)
	length := len(data.List)
	for i := 0; i < days; i++ {
		var p int
		d := data.List[length-days+i]
		switch t {
		case Confirmed:
			p = d.Confirmed
		case Deaths:
			p = d.Deaths
		case Hospitalized:
			p = d.Hospitalized
		case Recovered:
			p = d.Recovered
		}
		points[i] = float64(p)
	}
	graph := asciigraph.Plot(points, asciigraph.Height(10))
	fmt.Println(strings.ToUpper(string(t)))
	fmt.Println(graph)
	return nil
}
