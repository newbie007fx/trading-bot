package crypto

var notifyCounter uint = 0

func SetNotifyCounter(counter uint) {
	if counter > notifyCounter {
		notifyCounter = counter
	}
}

func IsShouldNotify() bool {
	if notifyCounter == 0 {
		return false
	}

	notifyCounter--
	return true
}
