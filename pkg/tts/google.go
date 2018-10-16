package tts

import (
	"bytes"
	"context"
	"io/ioutil"

	"cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

type googleTTS struct {
	client *texttospeech.Client
	ctx    context.Context
}

func NewGoogle() (*googleTTS, error) {
	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &googleTTS{
		ctx:    ctx,
		client: client,
	}, nil
}

func (g *googleTTS) Speak(text string) error {
	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: text,
			},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-GB",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_MALE,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_LINEAR16,
		},
	}

	resp, err := g.client.SynthesizeSpeech(g.ctx, &req)
	if err != nil {
		return err
	}

	err = play(ioutil.NopCloser(bytes.NewReader(resp.AudioContent)))
	if err != nil {
		return err
	}

	return nil
}
