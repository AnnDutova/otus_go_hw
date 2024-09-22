package programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/goccy/go-json"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	file := bufio.NewScanner(r)
	file.Split(bufio.ScanLines)

	var user User
	dotDomain := "." + domain
	for file.Scan() {
		if err := json.Unmarshal(file.Bytes(), &user); err != nil {
			return result, err
		}

		if strings.HasSuffix(user.Email, dotDomain) {
			atIndex := strings.Index(user.Email, "@")
			if atIndex != -1 {
				result[strings.ToLower(user.Email[(atIndex+1):])]++
			}
		}
	}

	return result, nil
}
