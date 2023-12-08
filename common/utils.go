package common

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
)

func GetDataPath() string {
	var cachePath, _ = os.UserCacheDir()
	return filepath.Join(cachePath, "audius_cli_player_test")
}

func GetDurationText(len int) string {
	hourNum := len / 3600
	minNum := int(math.Mod(float64(len), 3600) / 60)
	secNum := math.Mod(float64(len), 60)
	secMod := ""

	if secNum < 10 {
		secMod = "0"
	}

	if hourNum == 0 {
		return fmt.Sprintf("%v:%v%v", minNum, secMod, secNum)
	} else {
		minMod := ""

		if minNum < 10 {
			minMod = "0"
		}
		return fmt.Sprintf("%v:%v%v:%v%vs", hourNum, minMod, minNum, secMod, secNum)
	}
}

func GetLengthText(len int) string {
	hourNum := len / 3600
	minNum := int(math.Mod(float64(len), 3600) / 60)
	secNum := math.Mod(float64(len), 60)
	secMod := ""

	if secNum < 10 {
		secMod = "0"
	}

	if hourNum == 0 {
		return fmt.Sprintf("%vm %v%vs", minNum, secMod, secNum)
	} else {
		minMod := ""

		if minNum < 10 {
			minMod = "0"
		}
		return fmt.Sprintf("%vh %v%vm %v%vs", hourNum, minMod, minNum, secMod, secNum)
	}
}
