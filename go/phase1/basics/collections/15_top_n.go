package collections

import "sort"

type Item struct {
	Name      string
	Frequency int
}

func displayTopN(input []Item, topN int) []Item {

	validItems := make([]Item, 0, len(input))

	// input and topN check
	if topN <= 0 || len(input) == 0 {
		return []Item{}
	}

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

	items := make([]Item, 0, topN)

	for i := range topN {
		if len(validItems) < topN { // validItemsがtopNよりも少ない場合は、validItemsを全件表示
			items = append(items, validItems...)
			return items
		}

		items = append(items, validItems[i])
	}
	return items
}
