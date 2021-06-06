package helper

import "strings"

var msgs map[string]string = map[string]string{
	"required": "{name} tidak boleh kosong",
	"email":    "format {name} tidak valid",
}

func GetValidationMsg(name, tag string) string {
	msg, ok := msgs[tag]
	if !ok {
		msg = "{name} tidak valid"
	}

	return strings.ReplaceAll(msg, "{name}", name)
}
