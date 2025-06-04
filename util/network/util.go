package network

var current Network

func SetCurrent(network Network) {
	current = network
}

func Current() Network {
	return current
}

func CurrentIs(networks ...Network) bool {
	return Is(Current(), networks...)
}

func Is(current Network, networks ...Network) bool {
	for _, network := range networks {
		if network == current {
			return true
		}
	}
	return false
}
