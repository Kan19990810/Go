package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/gordonklaus/portaudio"
)

func main() {
	// 初始化PortAudio
	portaudio.Initialize()
	defer portaudio.Terminate()

	// 定义音频参数
	sampleRate := 44100
	framesPerBuffer := 64
	numChannels := 1
	seconds := 10

	// 创建指定文件夹
	outputDir := "/home/lapwsl/Data"
	os.MkdirAll(outputDir, os.ModePerm)

	// 创建输入流
	stream, err := portaudio.OpenDefaultStream(numChannels, 0, float64(sampleRate), framesPerBuffer, nil)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer stream.Close()

	// 创建WAV文件
	outputFile, err := os.Create(outputDir + "/00001.wav")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// 写入WAV文件头
	writeWavHeader(outputFile, numChannels, sampleRate, 16, uint32(sampleRate*seconds))

	// 开始录音
	err = stream.Start()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	fmt.Println("Recording...")

	// 创建缓冲区存储音频数据
	buffer := make([]int16, framesPerBuffer)

	// 录制指定秒数
	for i := 0; i < sampleRate/framesPerBuffer*seconds; i++ {
		err = stream.Read()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		// 写入缓冲区数据到WAV文件
		binary.Write(outputFile, binary.LittleEndian, buffer)
	}

	// 停止录音
	err = stream.Stop()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	fmt.Println("Recording finished")
}

// writeWavHeader 写入WAV文件头
func writeWavHeader(file *os.File, numChannels, sampleRate, bitsPerSample int, numSamples uint32) {
	var (
		fileSize      uint32 = 44 + numSamples*uint32(numChannels)*uint32(bitsPerSample)/8 - 8
		byteRate      uint32 = uint32(sampleRate) * uint32(numChannels) * uint32(bitsPerSample) / 8
		blockAlign    uint16 = uint16(numChannels) * uint16(bitsPerSample) / 8
		audioFormat   uint16 = 1
		subchunk2Size        = numSamples * uint32(numChannels) * uint32(bitsPerSample) / 8
	)

	binary.Write(file, binary.LittleEndian, []byte("RIFF"))
	binary.Write(file, binary.LittleEndian, fileSize)
	binary.Write(file, binary.LittleEndian, []byte("WAVE"))
	binary.Write(file, binary.LittleEndian, []byte("fmt "))
	binary.Write(file, binary.LittleEndian, uint32(16))
	binary.Write(file, binary.LittleEndian, audioFormat)
	binary.Write(file, binary.LittleEndian, uint16(numChannels))
	binary.Write(file, binary.LittleEndian, uint32(sampleRate))
	binary.Write(file, binary.LittleEndian, byteRate)
	binary.Write(file, binary.LittleEndian, blockAlign)
	binary.Write(file, binary.LittleEndian, uint16(bitsPerSample))
	binary.Write(file, binary.LittleEndian, []byte("data"))
	binary.Write(file, binary.LittleEndian, subchunk2Size)
}
