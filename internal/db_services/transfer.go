package db_services

import "unicode"

// collectionMap[<timeframe>] = collection
var collectionMap = map[string]string{
	"1m":  "m1",
	"3m":  "m3",
	"5m":  "m5",
	"15m": "m15",
	"30m": "m30",
	"1h":  "h1",
	"2h":  "h2",
	"4h":  "h4",
	"6h":  "h6",
	"8h":  "h8",
	"12h": "h12",
	"1d":  "d1",
	"3d":  "d3",
	"1w":  "w1",
	"1M":  "M1", // 30 days
}

// timeframeMap[<collection>] = timeframe
var timeframeMap = map[string]string{
	"m1":  "1m",
	"m3":  "3m",
	"m5":  "5m",
	"m15": "15m",
	"m30": "30m",
	"h1":  "1h",
	"h2":  "2h",
	"h4":  "4h",
	"h6":  "6h",
	"h8":  "8h",
	"h12": "12h",
	"d1":  "1d",
	"d3":  "3d",
	"w1":  "1w",
	"M1":  "1M", // 30 days
}

func Transfer_TF_Coll(i interface{}) interface{} {
	if s, ok := i.(string); ok {
		if unicode.IsDigit(rune(s[0])) {
			return collectionMap[s]
		} else {
			return timeframeMap[s]
		}
	}

	if s, ok := i.([]string); ok {
		var res []string

		if unicode.IsDigit(rune(s[0][0])) {
			for _, timeframe := range s {
				res = append(res, collectionMap[timeframe])
			}
		} else {
			for _, coll := range s {
				res = append(res, timeframeMap[coll])
			}
		}

		return res
	}

	return nil
}
