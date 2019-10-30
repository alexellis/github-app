package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type CodeReq struct {
	Code string `json:"code"`
}

type AppTemplate struct {
	AppID         string
	AppURL        string
	AppName       string
	PEM           string
	WebhookSecret string
	Response      string
}

type AppResult struct {
	ID            int    `json:"id"`
	PEM           string `json:"pem"`
	URL           string `json:"html_url"`
	Name          string `json:"name"`
	WebhookSecret string `json:"webhook_secret"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	if r.URL.Path == "/callback" {
		code := r.URL.Query().Get("code")

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

			appRes := AppResult{}

			err := json.Unmarshal(result, &appRes)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var outBuffer bytes.Buffer

			tmpl, err := template.ParseFiles("result.html")
			err = tmpl.Execute(&outBuffer, AppTemplate{
				AppID:         fmt.Sprintf("%d", appRes.ID),
				AppName:       appRes.Name,
				AppURL:        appRes.URL,
				PEM:           appRes.PEM,
				WebhookSecret: appRes.WebhookSecret,
				Response:      string(result),
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write(outBuffer.Bytes())

			return
		}
		return
	}

	res, err := ioutil.ReadFile("./index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
