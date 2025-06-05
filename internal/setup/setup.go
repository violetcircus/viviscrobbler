package setup

import (
	"crypto/md5"
	"github.com/violetcircus/viviscrobbler/internal/secret"
	"io"
	"sort"
)

func Setup() {

	//todo
}

// takes a slice of url parameters and creates the signature.
// they are passed as single strings containing both the
// parameter name and its value, to make sorting easier.
func SignSignature(parameters []string) string {
	secrets := secret.GetSecrets()
	h := md5.New()

	sort.Strings(parameters)

	for _, param := range parameters {
		io.WriteString(h, param)
	}
	io.WriteString(h, secrets.Secret)
	result := string(h.Sum(nil))
	return result
}

func authenticate() {

}

func getSession() {

}
