// +build libasound2

package main

import (
	"github.com/jphastings/jan-poka/pkg/tracker"
	"github.com/jphastings/jan-poka/pkg/tts"
	"github.com/jphastings/jan-poka/pkg/tts/googletts"
)

func init() {
	configurables = append(configurables, configurable{
		"Audio",
		func() bool { return environment.UseAudio },
		configureAudio,
	})
}

func configureAudio() (tracker.OnTracked, error) {
	ttsEngine, err := googletts.New()
	if err != nil {
		return nil, err
	}

	return tts.TrackedCallback(ttsEngine), nil
}
