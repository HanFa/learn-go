package main

import (
	"bytes"
	"encoding/base64"
	bili "github.com/JimmyZhangJW/biliStreamClient"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"log"
	"time"
)

const (
	SecretId  = "TODO"
	SecretKey = "TODO"
)

func main() {
	biliClient := bili.New()
	biliClient.Connect(4980320)
	defer biliClient.Disconnect()
	for {
		packBody := <-biliClient.Ch
		switch packBody.Cmd {
		case "DANMU_MSG":
			danmu, err := packBody.ParseDanmu()
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(danmu)
			sanitizedMsg := bili.Sanitize(danmu.Message)
			if bili.IsContainChineseWord(sanitizedMsg) {
				encodedVoice, err := bili.GetVoiceFromTencentCloud(SecretId, SecretKey, bili.DefaultGirlVoice, danmu.Message)
				if err != nil {
					log.Fatalln(err)
				}
				data, err := base64.StdEncoding.DecodeString(encodedVoice)

				streamer, format, err := wav.Decode(bytes.NewReader(data))
				speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
				speaker.Play(streamer)
			}

		case "SEND_GIFT":
			log.Println(packBody.ParseGift())
		case "COMBO_SEND":
			log.Println(packBody.ParseGiftCombo())
		default:
			log.Println(packBody.Cmd)
		}
	}

}
