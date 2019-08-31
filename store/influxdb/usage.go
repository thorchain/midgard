package influxdb

import (
	"math"
)

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
	// Find the usage stats, for all time
	query := "SELECT * FROM swaps_usage GROUP BY target"
	resp, err := in.Query(query)
	if err != nil {
		return
	}

	if len(resp) > 0 && len(resp[0].Series) > 1 {
		// ensure rune if first in the series
		if resp[0].Series[1].Tags["target"] == "rune" {
			//resp[0].Series[0], resp[0].Series[1] = resp[0].Series[1], resp[0].Series[0]
		}
		for _, series := range resp[0].Series {
			cols := series.Columns
			vals := series.Values[0]

			if series.Tags["target"] == "rune" {
				totalRuneTx, _ := getIntValue(cols, vals, "total_rune_tx")
				usage.TotalTx += totalRuneTx
				volRune, _ := getFloatValue(cols, vals, "rune_sum")
				usage.TotalVolAT = math.Abs(volRune)
				feeRune, _ := getFloatValue(cols, vals, "rune_fee_sum")
				usage.TotalEarned += math.Abs(feeRune)
			} else {
				totalTokenTx, _ := getIntValue(cols, vals, "total_token_tx")
				usage.TotalTx += totalTokenTx
				// we get the amount of rune in this case too because we're
				// getting total rune volume.
				volRune, _ := getFloatValue(cols, vals, "rune_sum")
				usage.TotalVolAT += math.Abs(volRune)
				volToken, _ := getFloatValue(cols, vals, "token_sum")
				feeToken, _ := getFloatValue(cols, vals, "token_fee_sum")
				usage.TotalEarned += (feeToken / volToken) * volRune
				// Round to nearest 8 decimal points
				usage.TotalEarned = math.Floor(usage.TotalEarned*100000000) / 100000000
			}
		}
		// round to 8 decimal places
		usage.TotalVolAT = math.Floor(usage.TotalVolAT*100000000) / 100000000
	}

	// Find the usage stats, for 30d
	query = "SELECT total_rune_tx, total_token_tx, token_sum, rune_sum FROM swaps_usage WHERE time > now() -30d"
	resp, err = in.Query(query)
	if err != nil {
		return
	}

	if len(resp) > 0 && len(resp[0].Series) > 0 {
		cols := resp[0].Series[0].Columns
		vals := resp[0].Series[0].Values[0]
		total, _ := getIntValue(cols, vals, "total_token_tx")
		usage.MonthlyTx += total
		total, _ = getIntValue(cols, vals, "total_rune_tx")
		usage.MonthlyTx += total
	}

	// Find the usage stats, for 1 day
	query = "SELECT * FROM swaps_usage WHERE time > now() - 1d GROUP BY target"
	resp, err = in.Query(query)
	if err != nil {
		return
	}

	if len(resp) > 0 && len(resp[0].Series) > 0 {
		for _, series := range resp[0].Series {
			cols := series.Columns
			vals := series.Values[0]
			if series.Tags["target"] == "rune" {
				totalRuneTx, _ := getIntValue(cols, vals, "total_rune_tx")
				usage.DailyTx += totalRuneTx
				volRune, _ := getFloatValue(cols, vals, "rune_sum")
				usage.TotalVol24 += math.Abs(volRune)
			} else {
				totalTokenTx, _ := getIntValue(cols, vals, "total_token_tx")
				usage.DailyTx += totalTokenTx
				// we get the amount of rune in this case too because we're
				// getting total rune volume.
				volRune, _ := getFloatValue(cols, vals, "rune_sum")
				usage.TotalVol24 += math.Abs(volRune)
			}
		}
		// round to 8 decimal places
		usage.TotalVolAT = math.Floor(usage.TotalVolAT*100000000) / 100000000
	}

	// Find total active users
	query = "SELECT token_fee_sum FROM swaps_usage GROUP BY from_address"
	resp, err = in.Query(query)
	if err != nil {
		return
	}
	if len(resp) > 0 {
		usage.TotalUsers = int64(len(resp[0].Series))
	}

	// Find monthly active users
	query = "SELECT token_fee_sum FROM swaps_usage WHERE time > now() -30d GROUP BY from_address"
	resp, err = in.Query(query)
	if err != nil {
		return
	}
	if len(resp) > 0 {
		usage.MonthlyActiveUsers = int64(len(resp[0].Series))
	}

	// Find daily active users
	query = "SELECT token_fee_sum FROM swaps_usage WHERE time > now() -1d GROUP BY from_address"
	resp, err = in.Query(query)
	if err != nil {
		return
	}
	if len(resp) > 0 {
		usage.DailyActiveUsers = int64(len(resp[0].Series))
	}

	return
}
