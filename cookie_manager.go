package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func SaveSessionId(w http.ResponseWriter, id int) {
	http.SetCookie(w, &http.Cookie{Name: "id", Value: fmt.Sprintf("%d", id), Expires: time.Now().Add(time.Hour)})
}

func GetSessionId(r *http.Request, id *int) {
	cookie, err := r.Cookie("id")
	if err != nil {
		log.Println("Error retrieving id cookie ", err)
		return
	}

	saved, err := parseInt(cookie.Value)
	if err != nil {
		log.Println("Error parsing id cookie into int ", err)
		return
	}

	*id = saved
}

/*func SaveResponse(res Response, w http.ResponseWriter) {
	expiration := time.Now().Add(10 * time.Minute)

	http.SetCookie(w, &http.Cookie{Name: "wave", Value: fmt.Sprintf("%d", res.wave), Expires: expiration})
	http.SetCookie(w, &http.Cookie{Name: "id", Value: fmt.Sprintf("%d", res.id), Expires: expiration})
	http.SetCookie(w, &http.Cookie{Name: "condition", Value: fmt.Sprintf("%d", res.condition), Expires: expiration})

	for i := range res.targets {
		q := res.targets[i]
		http.SetCookie(w, &http.Cookie{Name: fmt.Sprintf("q%d_text", i), Value: q.text, Expires: expiration})
		http.SetCookie(w, &http.Cookie{Name: fmt.Sprintf("q%d_s1", i), Value: q.s1, Expires: expiration})
		http.SetCookie(w, &http.Cookie{Name: fmt.Sprintf("q%d_s2", i), Value: q.s2, Expires: expiration})
	}
}

func GetSeedResponse(r *http.Request, res *Response) {
	waveCookie, err := r.Cookie("wave")
	if err != nil {
		return
	}

	wave, err := parseInt(waveCookie.Value)
	if err != nil {
		return
	}
}*/
