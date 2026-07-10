package service

func isPhoneValid(phone string) bool {
	hasDigit := false
	prevHyphen := false
	for _, r := range phone {
		switch {
		case r >= '0' && r <= '9':
			hasDigit = true
			prevHyphen = false
		case r == ' ':
			prevHyphen = false
		case r == '-':
			if prevHyphen {
				return false
			}
			prevHyphen = true
		default:
			return false
		}
	}
	return hasDigit
}
