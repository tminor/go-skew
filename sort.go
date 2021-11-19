package skew

func countingSort(arr []*[]rune, significancePlace int) (sorted []*[]rune) {
	sorted = make([]*[]rune, len(arr))

	max := rune(0)
	for _, v := range arr {
		if (*v)[significancePlace] > max {
			max = (*v)[significancePlace]
		}
	}

	count := make([]int, max + 1)

	for _, v := range arr {
		count[(*v)[significancePlace]]++
	}

	for i := rune(1); i <= max; i++ {
		count[i] += count[i - 1]
	}

	for i := len(arr) - 1; i >= 0; i-- {
		j := (*arr[i])[significancePlace]
		count[j]--
		sorted[count[j]] = arr[i]
	}

	return sorted
}

func radixSort(arr []*[]rune) (sorted []*[]rune) {
	sorted = countingSort(arr, 0)

	sorted = doRadixSort(sorted, 1)

	return sorted
}

func doRadixSort(arr []*[]rune, pos int) []*[]rune {
	bucketSizes := make(map[rune]int)

	for _, v := range arr {
		bucketSizes[(*v)[pos - 1]]++
	}

	bucket := make([]*[]rune, 0)
	for i, j := 0, 0; i < len(arr); i++ {
		if j == len(bucket) && j != 0 {
			bucket = countingSort(bucket, pos)
			if len(*bucket[0]) - 1 < pos {
				bucket = doRadixSort(bucket, pos + 1)
			}
			for i , bi := i - j, 0; i < j; i++ {
				arr[i] = bucket[bi]
				bi++
			}
			bucket = make([]*[]rune, 0)
			j = 0
		} else if j == 0 {
			bucket = make([]*[]rune, bucketSizes[(*arr[i])[pos - 1]])
			bucket[j] = arr[i]
			j++
		} else {
			bucket[j] = arr[i]
			j++
		}
	}

	return arr
}
