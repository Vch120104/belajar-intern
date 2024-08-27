package test

import (
	"fmt"
	"testing"
	"time"
)

func Test_CarWashScreen(t *testing.T) {
	// config.InitEnvConfigs(true, "")
	// tx := config.InitDB()

	// timeDiff := 1

	// hours := math.Floor(float64(timeDiff))
	// minutes := (float64(timeDiff) - hours) * 100

	// duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
	currentTime := time.Now()

	// newTime := currentTime.Add(duration)
	newTime := currentTime

	fmt.Print(float32(newTime.Hour()) + float32(newTime.Minute())/60)
}
