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

// takes a map of url parameters and creates the signature.
func SignSignature(parameters map[string]string) string {
	secrets := secret.GetSecrets()

	// get parameter keys and sort them
	keys := make([]string, len(parameters))
	i := 0
	for k := range parameters {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	h := md5.New()
	for _, key := range keys {
		io.WriteString(h, key)
		io.WriteString(h, parameters[key])
	}
	io.WriteString(h, secrets.Secret)
	result := string(h.Sum(nil))
	return result
}

func GetToken() {

}
