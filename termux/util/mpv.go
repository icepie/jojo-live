package util

import "os/exec"

func Mpv(url string) error {
	mpv := exec.Command("mpv", url)
	err := mpv.Start()
	return err
}
