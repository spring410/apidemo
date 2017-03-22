package mysqldb

import (
	"strings"
)

const (
	Error_1049 = `Error 1049` //For exmaple:Error Code: 1049. Unknown database 'xxxx'
	Error_1149 = `Error 1146` //For exmaple:Error Code: 1146. Table 'accounts.userss' doesn't exist
)

func DatabaseUnkown(m string) bool {
	return strings.Contains(m, Error_1049)
}

func TableNotExist(m string) bool {
	return strings.Contains(m, Error_1149)
}
