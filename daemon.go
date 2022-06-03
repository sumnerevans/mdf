package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"

	"github.com/gorilla/mux"
)

func RunDaemon(port int) {
	mdfd := MuttDisplayFilterDaemon{
		Links: map[string]string{},
	}

	router := mux.NewRouter()
	router.HandleFunc("/new", mdfd.New)
	router.HandleFunc("/{id}", mdfd.RedirectPage)
	http.Handle("/", router)

	log.Printf("Listening on port %d...", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Printf("Failed to listen: %v", err)
	}
}

type MuttDisplayFilterDaemon struct {
	Links map[string]string
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var numLetters = big.NewInt(int64(len(letters)))

func (mdfd *MuttDisplayFilterDaemon) randomString(length int) string {
	s := make([]rune, length)
	for i := range s {
		n, _ := rand.Int(rand.Reader, numLetters)
		s[i] = letters[n.Int64()]
	}
	return string(s)
}

func (mdfd *MuttDisplayFilterDaemon) New(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	url, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
	} else {
		id := mdfd.randomString(6)
		log.Printf("Creating mapping of %s to %s", id, url)
		mdfd.Links[id] = string(url)
		w.WriteHeader(200)
		io.WriteString(w, id)
	}
}

func (mdfd *MuttDisplayFilterDaemon) RedirectPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Printf("Serving /%s", id)

	url, ok := mdfd.Links[id]

	if !ok {
		log.Printf("%s not found", id)
		w.WriteHeader(404)
		return
	}
	log.Printf("Serving redirect page for %s -> %s", id, url)

	respHtml := fmt.Sprintf(`<!doctype html>
	<html>
	<body>
		<textarea id="url_edit" rows="10" style="width: 100%%;">%s</textarea>
		<br />
		<input id="go" type="button" value="Go" />
		<script>
			const redirect = () => window.location.replace(document.getElementById('url_edit').value);
			document.getElementById('go').addEventListener('click', redirect);
			document.onkeypress = e => e.keyCode === 13 && redirect() && false;
			document.getElementById('url_edit').focus();
		</script>
	</body>
	</html>`, url)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	io.WriteString(w, respHtml)
}
