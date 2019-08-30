package influxdb

import "fmt"

type UsageData struct {
	DailyActiveUsers   int64   `json:"dau"`
	MonthlyActiveUsers int64   `json:"mau"`
	TotalUsers         int64   `json:"totalUsers"`
	DailyTx            int64   `json:"dailyTx"`
	MonthlyTx          int64   `json:"monthlyTx"`
	TotalTx            int64   `json:"totalTx"`
	TotalVolAT         float64 `json:"totalVolAT"`
	TotalVol24         float64 `json:"totalVol24"`
	TotalStaked        float64 `json:"totalStaked"`
	TotalEarned        float64 `json:"totalEarned"`
}

func (in Client) GetUsageData() (usage UsageData, err error) {
	// Find the usage stats
	query := "SELECT * FROM swaps_usage GROUP BY target"
	resp, err := in.Query(query)
	if err != nil {
		return
	}

	if len(resp) > 0 && len(resp[0].Series) > 1 && len(resp[0].Series[0].Values) > 0 {
		var runeCols, tokenCols []string
		var runeVals, tokenVals []interface{}
		if resp[0].Series[0].Tags["target"] == "rune" {
			runeCols = resp[0].Series[0].Columns
			runeVals = resp[0].Series[0].Values[0]
			tokenCols = resp[0].Series[1].Columns
			tokenVals = resp[0].Series[1].Values[0]
		} else {
			runeCols = resp[0].Series[1].Columns
			runeVals = resp[0].Series[1].Values[0]
			tokenCols = resp[0].Series[0].Columns
			tokenVals = resp[0].Series[0].Values[0]
		}
		totalTokenTx, _ := getIntValue(tokenCols, tokenVals, "total_token_tx")
		totalRuneTx, _ := getIntValue(runeCols, runeVals, "total_rune_tx")
		usage.TotalTx = totalTokenTx + totalRuneTx
		volToken, _ := getFloatValue(tokenCols, tokenVals, "token_sum")
		volRune, _ := getFloatValue(runeCols, runeVals, "rune_sum")
		usage.TotalVolAT = volToken + volRune
	}

	fmt.Printf("Results: %+v\n", resp)

	return
}
