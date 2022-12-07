package btc_zombie

import (
	"io/ioutil"
)

func ReadAbiJson() (string, error) {
	b, err := ioutil.ReadFile("./contracts/btc_zombie/abi.json")

	if err != nil {
		return "", err
	}

	return string(b), nil
}
