package encode

type frequencyTable map[byte]int

// Record how often each byte pattern occurs in the input
func newFrequencyTable(data []byte) frequencyTable {
	counts := frequencyTable{}
	for _, b := range data {
		if _, ok := counts[b]; !ok {
			counts[b] = 0
		}
		counts[b]++
	}

	return counts
}
