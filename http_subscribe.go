package main

import (
	"fmt"
    "io"
    "log"
    "net/http"
    "net/url"
)


func subscribe(channelID string) {
    hubURL := "https://pubsubhubbub.appspot.com/subscribe"
    callback := "https://yt_analiser.pedrofonsecacruz.com.br/youtube/callback"
    topic := "https://www.youtube.com/xml/feeds/videos.xml?channel_id=" + channelID

    data := url.Values{}
    data.Set("hub.mode", "subscribe")
    data.Set("hub.topic", topic)
    data.Set("hub.callback", callback)
    data.Set("hub.verify", "async")

    resp, err := http.PostForm(hubURL, data)
    if err != nil {
        log.Fatal("Subscribe error:", err)
    }
    defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)
    fmt.Println("Subscription response:", string(body))
}

func youtubeCallback(w http.ResponseWriter, r *http.Request) {
    // Verification
    challenge := r.URL.Query().Get("hub.challenge")
    if challenge != "" {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(challenge))
        return
    }

    // Notification
    defer r.Body.Close()
    body, _ := io.ReadAll(r.Body)
    fmt.Println("Got YouTube notification:", string(body))
}