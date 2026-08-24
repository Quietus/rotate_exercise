package rotate

func LeftRotate(a []byte) []byte {
	var finalBit byte = 0
	n := len(a)
	for i := 0; i < n; i++ {
		temp := a[i]
		if temp >= 128 {
			temp = temp - 128
			if i == 0 {
				finalBit = 1
			} else {
				a[i-1] = a[i-1] + 1
			}
		}
		a[i] = temp << 1
	}
	if finalBit == 1 {
		a[n-1] = a[n-1] + 1
	}
	return a
}

func RightRotate(a []byte) []byte {
	var finalBit byte = 0
	n := len(a)
	for i := n - 1; i >= 0; i-- {
		temp := a[i]
		if temp%2 == 1 {
			temp = temp - 1
			if i == n-1 {
				finalBit = 128
			} else {
				a[i+1] = a[i+1] + 128
			}
		}
		a[i] = temp >> 1
	}
	if finalBit == 128 {
		a[0] = a[0] + 128
	}
	return a
}
