package bank

import (
	"io"
	"os"
	"strings"
)

func loadToken() (bool, error) {
	file, err := os.Open("./token")
	if err != nil && os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return false, err
	}

	splitted := strings.Split(string(data), "\n")

	accessToken = splitted[0]
	revokeUrl = splitted[1]

	accessToken = strings.TrimSpace(accessToken)
	revokeUrl = strings.TrimSpace(revokeUrl)

	accessToken = strings.Trim(accessToken, "\n")
	revokeUrl = strings.Trim(revokeUrl, "\n")

	accessToken = strings.Trim(accessToken, "\r")
	revokeUrl = strings.Trim(revokeUrl, "\r")

	return true, nil
}

func storeToken() error {
	file, err := os.Create("./token")
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(accessToken)
	if err != nil {
		return err
	}

	_, err = file.WriteString("\n")
	if err != nil {
		return err
	}

	_, err = file.WriteString(revokeUrl)
	if err != nil {
		return err
	}

	return nil
}

func removeToken() error {
	err := os.Remove("./token")
	if err != nil {
		return err
	}

	return nil
}
