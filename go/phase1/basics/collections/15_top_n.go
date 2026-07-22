package collections

import "sort"

type Item struct {
	Name      string
	Frequency int
}

func displayTopN(input []Item, topN int) []Item {

	// input and topN check
	if topN <= 0 {
		return []Item{}
	}

	// 修正：ガードの後に容量を確保
	validItems := make([]Item, 0, len(input))
	// frequency negative check
	for _, e := range input {
		if e.Frequency >= 0 {
			validItems = append(validItems, e)
		}
	}

	// validItems check
	if len(validItems) == 0 {
		return []Item{}
	}

	sort.Slice(validItems, func(i, j int) bool {
		if validItems[i].Frequency == validItems[j].Frequency {
			return validItems[i].Name < validItems[j].Name // もし同じ場合は名前順の昇順
		}
		return validItems[i].Frequency > validItems[j].Frequency
	})

	return validItems[:min(topN, len(validItems))] // 修正：topNもしくはlenの少ない方を採用する
}
