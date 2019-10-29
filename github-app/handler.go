package function

import (
	"io/ioutil"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	if r.Body != nil {
		defer r.Body.Close()
	}

	body, _ := ioutil.ReadFile("function/index.html")

	w.WriteHeader(http.StatusOK)
	w.Write(body)

}
