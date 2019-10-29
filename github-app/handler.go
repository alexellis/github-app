package function

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	if r.URL.Path == "/callback" {
		code := r.URL.Query().Get("code")
		reader := bytes.NewBufferString(code)
		res, err := http.Post(fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code), "application/vnd.github.fury-preview+json", reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if res.Body != nil {
			defer res.Body.Close()
			result, _ := ioutil.ReadAll(res.Body)
			w.Write(result)
		}
	}

	res, _ := ioutil.ReadFile("function/index.html")

	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
