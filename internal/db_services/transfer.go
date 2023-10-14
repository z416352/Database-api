package db_services

var timeframeToCollectionMap = map[string]string{
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

var collectionToTimeframeMap = map[string]string{
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

func TimeframesToCollections(timeframes interface{}) interface{} {
	switch timeframes.(type) {
	case string:
		return timeframeToCollectionMap[timeframes.(string)]

	case []string:
		var res []string

		for _, timeframe := range timeframes.([]string) {
			res = append(res, timeframeToCollectionMap[timeframe])
		}
		return res

	default:
		return nil
	}

}

func CollectionsToTimeframes(collections interface{}) interface{} {
	switch collections.(type) {
	case string:
		return collectionToTimeframeMap[collections.(string)]

	case []string:
		var res []string

		for _, coll := range collections.([]string) {
			res = append(res, collectionToTimeframeMap[coll])
		}
		return res

	default:
		return nil
	}
}
