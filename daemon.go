package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"sync"
)

func RunDaemon(port int) {
	mdfd := MuttDisplayFilterDaemon{
		Links: map[string]string{},
	}

	router := http.NewServeMux()
	router.HandleFunc("POST /new", mdfd.New)
	router.HandleFunc("GET /{id}", mdfd.RedirectPage)

	log.Printf("Listening on port %d...", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Printf("Failed to listen: %v", err)
	}
}

type MuttDisplayFilterDaemon struct {
	Links map[string]string
	mu    sync.RWMutex
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
	url, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
	} else {
		id := mdfd.randomString(6)
		log.Printf("Creating mapping of %s to %s", id, url)
		mdfd.mu.Lock()
		mdfd.Links[id] = string(url)
		mdfd.mu.Unlock()
		w.WriteHeader(200)
		io.WriteString(w, id)
	}
}

func (mdfd *MuttDisplayFilterDaemon) RedirectPage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Printf("Serving /%s", id)

	mdfd.mu.RLock()
	url, ok := mdfd.Links[id]
	mdfd.mu.RUnlock()

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
		<input id="go" type="button" value="Go (5)" disabled />
		<script>
			document.getElementById('url_edit').focus();
			document.onkeypress = e => e.keyCode === 13 && false;
			let i = 5;
			const goButton = document.getElementById('go');
			const interval = setInterval(() => {
				i--;
				if (i > 0) {
					goButton.value = 'Go ('+i+')';
				} else {
					const redirect = () => window.location.replace(document.getElementById('url_edit').value);

					goButton.value = 'Go';
					goButton.removeAttribute('disabled');
					goButton.addEventListener('click', redirect);
					document.onkeypress = e => e.keyCode === 13 && redirect() && false;
					clearInterval(interval);
				}
			}, 1000);
		</script>
	</body>
	</html>`, url)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	io.WriteString(w, respHtml)
}
