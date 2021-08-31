package api

func Contains(uints []uint, i uint) bool {
	for _, v := range uints {
		if v == i {
			return true
		}
	}

	return false
}
