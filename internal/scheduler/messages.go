package scheduler

import (
	"fmt"
	"math/rand/v2"
)

var messageTemplates = []string {
	"sekarang jadwalnya %s nih! semangat ya",
	"ser, sekarang waktunya %s",
	"sekarang waktunya %s loh, gak lupa kan?",
	"jangan lupa sama jadwal %s nya ya",
	"lagi dmn? sekarang jadwalnya %s",
}

func pickMessage(activityName string) string {
	template := messageTemplates[rand.IntN(len(messageTemplates))]
	return fmt.Sprintf(template, activityName)
}