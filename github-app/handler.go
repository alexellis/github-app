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
		req, _:=http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions",code), reader)
		req.Header.Add("Accept","application/vnd.github.fury-preview+json")
		res, err := http.Post(fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code), , reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if res.Body != nil {
			defer res.Body.Close()
			result, _ := ioutil.ReadAll(res.Body)
			w.Write(result)
			return
		}
		return
	}

	res, err := ioutil.ReadFile("./index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
