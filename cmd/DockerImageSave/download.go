package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/dustin/go-humanize"
)

func downloadFile(afile string, animateProgress bool) bool {
	client := grab.NewClient()
	pwd, err := os.Getwd()
	if err != nil {
		return false
	}
	req, _ := grab.NewRequest(pwd, afile)
	req.SkipExisting = false
	req.NoResume = false

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	if resp.HTTPResponse != nil {
		fmt.Printf("  %v\n", resp.HTTPResponse.Status)
	} else {
		return false
	}

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			if animateProgress {
				fmt.Printf("  transferred %v / %v (%.2f%%)\t\t\r",
					humanize.Bytes(uint64(resp.BytesComplete())),
					humanize.Bytes(uint64(resp.Size)),
					100*resp.Progress())
			}

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		return false
	}
	//newName := pwd + strings.ReplaceAll(strings.Replace(resp.Filename, pwd, "", 1), ":", "_") // it might break on windows if we don't clean the path first
	//err = os.Rename(resp.Filename, newName)
	//if err != nil {
	//	fmt.Printf("Problems renaming file %v \n", resp.Filename)
	//}
	fmt.Printf("Download saved to %v \n", resp.Filename)
	return true
}
