package countryAlpha2

import (
	"io/ioutil"
	"strings"
)

func CA2() map[string]string {
	aC2Entity := make(map[string]string)
	data, err := ioutil.ReadFile("countryAlpha2/countries.txt")
	if err != nil {
		panic(err)
	}
	slice := strings.Split(string(data), "\n")
	for _, v := range slice {
		v = strings.Replace(v, "\r", "", -1) //удаление каретки \r
		r := strings.Split(v, " ")
		aC2Entity[r[0]] = v[3:]
	}
	return aC2Entity
}
