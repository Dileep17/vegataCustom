package main

import (
	vegeta "github.com/tsenart/vegeta/lib"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func getTargeter(payloadBytes []byte) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "POST"
		tgt.URL = "https://test.k6.io/login.php"
		tgt.Body = payloadBytes

		header := http.Header{}
		header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		header.Add("Content-Type", "application/x-www-form-urlencoded")
		header.Add("time",string(time.Now().UnixNano()))
		header.Add("cookie","_ga=GA1.2.425423683.1624524066; _lfa=LF1.1.0f034d0348731267.1627992485029; _hjid=2ea888d1-66e6-4e1c-8c76-d34d9d3890b4; messagesUtk=6197f79f5bc441fa896d552c98683ae1; csrf=NTA4ODczNDQ0")
		tgt.Header = header

		return nil
	}
}

func main() {
	requestRate := vegeta.Rate{Freq: 1, Per: time.Second}
	testDuration := 20 * time.Second
	payloadBytes, _ := ioutil.ReadFile("payload.txt")

	targeter := getTargeter(payloadBytes)
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics

	for res := range attacker.Attack(targeter, requestRate, testDuration, "constant rate test") {
		metrics.Add(res)
	}
	metrics.Close()

	textReporter := vegeta.NewTextReporter(&metrics)
	HDRHistogramreporter := vegeta.NewHDRHistogramPlotReporter(&metrics)
	jsonMetrics := vegeta.NewJSONReporter(&metrics)
	textReporter(os.Stdout)
	jsonMetrics(os.Stdout)
	HDRHistogramreporter(os.Stdout)
}