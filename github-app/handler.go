package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CodeReq struct {
	Code string `json:"code"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	if r.URL.Path == "/callback" {
		code := r.URL.Query().Get("code")
		log.Println(code)

		codeBytes, _ := json.Marshal(&CodeReq{Code: code})
		log.Println(string(codeBytes))
		// reader := bytes.NewReader(codeBytes)
		// reader := bytes.NewReader([]byte(code))

		req, _ := http.NewRequest(http.MethodPost,
			fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code), nil)

		req.Header.Add("Accept", "application/vnd.github.fury-preview+json")
		res, err := http.DefaultClient.Do(req)

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
