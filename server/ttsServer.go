package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"raiden_fumo/lang_notebook_core/ai"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/gordonklaus/portaudio"
	"github.com/tosone/minimp3"
)

type TtsServer struct {
	context *oto.Context
	openaiApiKey string
}

func (server *TtsServer) Initialize() {
	portaudio.Initialize()
}

func (server *TtsServer) Close() {
	portaudio.Terminate()
}

func (server *TtsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	jsonBody := map[string]string{}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		http.Error(w, "Couldn't parse request data", http.StatusBadRequest)
		return
	}

	text := jsonBody["text"]
	sound, err := ai.ConvertTextToSpeech(server.openaiApiKey, text, "alloy")
	chk(err)

	var dec *minimp3.Decoder
	var data []byte
	if dec, data, err = minimp3.DecodeFull(sound); err != nil {
		log.Fatal(err)
	}

	rate, channels := dec.SampleRate, dec.Channels

	op := &oto.NewContextOptions{
		SampleRate: rate,
		ChannelCount: channels,
		Format: oto.FormatSignedInt16LE,
	}
	if server.context == nil {
		if server.context, _, err = oto.NewContext(op); err != nil {
			log.Fatal(err)
		}
	}

	var player = server.context.NewPlayer(bytes.NewReader(data))
	player.Play()

	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	err = player.Close()
	if err != nil {
			panic("player.Close failed: " + err.Error())
	}

	w.Write([]byte("Yo-hoo"))
}

func makeTtsServer(openaiApiKey string) *TtsServer {
	ttsServer := &TtsServer{
		openaiApiKey: openaiApiKey,
	}
	ttsServer.Initialize()
	// defer ttsServer.Close()
	return ttsServer
}
