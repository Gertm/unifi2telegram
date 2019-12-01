package main

import (
	"fmt"
	"regexp"
	"testing"
)

func TestRegexp(t *testing.T) {
	cameraRegexp, _ = regexp.Compile("Camera\\[[a-zA-Z0-9]*\\|(\\w*)\\]")
	m := cameraRegexp.FindStringSubmatch("1574309874.201 2019-11-21 05:17:54.201/CET: INFO   [uv.recording.info] Camera[F09FC214766D|Voordeur] MOTION ENDED motion:1168, rec_id:5dd60fe8c5e2254c674486eb in AnalyticsEvtBus-4")
	fmt.Printf("%q\n", m)
	if m[1] != "Voordeur" {
		t.Errorf("Did not find voordeur.")
	}
}
