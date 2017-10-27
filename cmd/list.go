// Copyright © 2017 Charles Haynes <ceh@ceh.bz>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"time"

	units "github.com/docker/go-units"
	"github.com/matthazinski/transmission"
	"github.com/spf13/cobra"
)

var torrents string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List torrents",
	Long: `List torrents for a tracker, allows filtering and sorting
For Example:

trr list - list all torrents for default tracker
trr list -sort active - list all torrents sorted by the time they were last active
trr list -filter uploading -sort added,name - list uloading torrents sorted by
    when they were added, and then by name`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("list(server: %s, torrents: %s)\n", server, torrents)
		a := fmt.Sprintf("http://%s/transmission/rpc", server)
		x, err := transmission.New(a, "", "")
		if err != nil {
			fmt.Println(err)
			return
		}
		ts, err := x.GetTorrents()
		if err != nil {
			fmt.Println(err)
			return
		}
		// ID     Done       Have  ETA           Up    Down  Ratio  Status       Name
		//   11    16%   618.8 MB  Unknown      0.0     7.0    0.0  Up & Down    Leo Kottke
		for _, t := range ts {
			fmt.Printf("%5d %3.0f%% %9s %9s %7s %7s %4.1f %-13s %s\n",
				t.ID,
				t.PercentDone*100.0,
				units.HumanSize(float64(t.Have())),
				ETA(t),
				units.HumanSize(float64(t.RateUpload)),
				units.HumanSize(float64(t.RateDownload)),
				t.UploadRatio,
				Status(t),
				t.Name)
			if err != nil {
				fmt.Println(err)
			}
		}
	},
}

// ETA prints a human readable eta for the torrent
func ETA(t *transmission.Torrent) string {
	if t.LeftUntilDone == 0 {
		return ""
	}
	if t.Eta < 0 && t.RateDownload > 0 {
		t.Eta = time.Duration(t.LeftUntilDone / t.RateDownload)
	}
	if t.Eta < 0 {
		return "∞"
	}
	return units.HumanDuration(t.Eta * time.Second)
}

// Status prints a human readable status for the torrent
func Status(t *transmission.Torrent) string {
	if t.ErrorString != "" {
		return t.ErrorString
	}
	if t.Error != 0 {
		return "Error"
	}
	switch t.Status {
	case transmission.StatusStopped:
		return "Stopped"
	case transmission.StatusCheckPending:
		return "Wait Verify"
	case transmission.StatusChecking:
		return "Verifying"
	case transmission.StatusDownloadPending:
		return "Wait Download"
	case transmission.StatusSeedPending:
		return "Wait Seed"
	case transmission.StatusSeeding, transmission.StatusDownloading:
		if t.RateDownload == 0 && t.RateUpload == 0 {
			return "Idle"
		}
		if t.RateDownload > 0 && t.RateUpload > 0 {
			return "Both"
		}
		if t.RateDownload > 0 {
			return "Downloading"
		}
		if t.RateUpload > 0 {
			return "Uploading"
		}
		return "Wat"
	default:
		return "unknown"
	}
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")
	listCmd.PersistentFlags().StringVar(&torrents, "torrents", "all", "which torrents to operate on")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}