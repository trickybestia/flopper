// Taken from https://github.com/bwmarrin/dgvoice
/*
Copyright (c) 2015, Bruce
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of dgvoice nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

/*******************************************************************************
 * This is very experimental code and probably a long way from perfect or
 * ideal.  Please provide feed back on areas that would improve performance
 *
 */

// Package dgvoice provides opus encoding and audio file playback for the
// Discordgo package.
package dgvoice

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

// NOTE: This API is not final and these are likely to change.

// Technically the below settings can be adjusted however that poses
// a lot of other problems that are not handled well at this time.
// These below values seem to provide the best overall performance
const (
	channels            int = 2                   // 1 for mono, 2 for stereo
	frameRate           int = 48000               // audio sampling rate
	frameSize           int = 960                 // count of uint16 in one channel in audio frame
	maxBytes            int = (frameSize * 2) * 2 // max size of opus data
	nanosecondsPerFrame int = 1_000_000_000 / (frameRate / frameSize)

	messageSkip   int = 0
	messagePause  int = 1
	messageResume int = 2
)

type PlaybackController struct {
	messageChannel chan int
	elapsedFrames  uint32
	skipped        bool
	paused         bool
}

func NewPlaybackController() *PlaybackController {
	return &PlaybackController{
		messageChannel: make(chan int, 10),
	}
}

func (playbackController *PlaybackController) Reset() {
	playbackController.paused = false
	playbackController.skipped = false
	playbackController.elapsedFrames = 0
}

func (playbackController *PlaybackController) Skip() error {
	if playbackController.skipped {
		return errors.New("already skipped")
	}

	playbackController.skipped = true
	playbackController.messageChannel <- messageSkip

	return nil
}

func (playbackController *PlaybackController) Pause() error {
	if playbackController.paused {
		return errors.New("already paused")
	}

	playbackController.paused = true
	playbackController.messageChannel <- messagePause

	return nil
}

func (playbackController *PlaybackController) Resume() error {
	if !playbackController.paused {
		return errors.New("already playing")
	}

	playbackController.paused = false
	playbackController.messageChannel <- messageResume

	return nil
}

func (playbackController *PlaybackController) Paused() bool {
	return playbackController.paused
}

func (playbackController *PlaybackController) ElapsedTime() time.Duration {
	return time.Duration(int64(playbackController.elapsedFrames) * int64(nanosecondsPerFrame))
}

// OnError gets called by dgvoice when an error is encountered.
// By default logs to STDERR
var OnError = func(str string, err error) {
	prefix := "dgVoice: " + str

	if err != nil {
		os.Stderr.WriteString(prefix + ": " + err.Error())
	} else {
		os.Stderr.WriteString(prefix)
	}
}

func createEncoder() (*gopus.Encoder, error) {
	encoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)

	if err != nil {
		return nil, err
	}

	encoder.SetBitrate(96000)

	return encoder, nil
}

func sendPCMFrame(encoder *gopus.Encoder, v *discordgo.VoiceConnection, pcmFrame []int16) {
	opus, err := encoder.Encode(pcmFrame, frameSize, maxBytes)

	if err != nil {
		OnError("Encoding Error", err)

		return
	}

	v.OpusSend <- opus
}

func PlayAudio(v *discordgo.VoiceConnection, ffmpegInputArg string, playbackController *PlaybackController) {
	ffmpegCommand := exec.Command("ffmpeg", "-i", ffmpegInputArg, "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")

	ffmpegout, err := ffmpegCommand.StdoutPipe()

	if err != nil {
		OnError("StdoutPipe Error", err)

		return
	}

	if err = ffmpegCommand.Start(); err != nil {
		OnError("RunStart Error", err)

		return
	}

	// prevent memory leak from residual ffmpeg streams
	defer func() {
		ffmpegCommand.Process.Kill()
		ffmpegCommand.Process.Wait()
	}()

	setSpeaking := func(speaking bool) {
		err = v.Speaking(speaking)

		if err != nil {
			OnError("Couldn't set speaking", err)
		}
	}

	setSpeaking(true)

	defer setSpeaking(false)

	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 524288)
	encoder, err := createEncoder()

	if err != nil {
		OnError("NewEncoder Error", err)
		return
	}

	for {
		select {
		case message := <-playbackController.messageChannel:
			switch message {
			case messageSkip:
				return
			case messagePause:
				setSpeaking(false)

				newMessage := <-playbackController.messageChannel

				if newMessage == messageSkip {
					return
				}

				setSpeaking(true)
			}
		default:
			pcmFrame := make([]int16, frameSize*channels)
			err = binary.Read(ffmpegbuf, binary.LittleEndian, &pcmFrame)

			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return
			}

			if err != nil {
				OnError("error reading from ffmpeg stdout", err)
				return
			}

			sendPCMFrame(encoder, v, pcmFrame)

			playbackController.elapsedFrames++
		}
	}
}
