package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Video Info Struct...
type VideoInfo struct {
	CodecName          string `json:"codec_name"`
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	DisplayAspectRatio string `json:"display_aspect_ratio"`
	FrameRate          string `json:"r_frame_rate"`
	BitRate            string `json:"bit_rate"`
}

// Audio Info Struct...
type AudioInfo struct {
	CodecName  string `json:"codec_name"`
	SampleRate string `json:"sample_rate"`
	Channels   int    `json:"channels"`
	BitRate    string `json:"bit_rate"`
}

// Main Func...
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./video_info <file_path>")
		return
	}

	filePath := os.Args[1]

	ffprobePath, err := findFFprobe()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	videoInfo, displayAspectRatio, videoErr := getVideoInfo(filePath, ffprobePath)
	audioInfo, audioErr := getAudioInfo(filePath, ffprobePath)

	if videoErr != nil {
		fmt.Printf("Ошибка при получении информации о видео: %v\n", videoErr)
		return
	}

	data := buildVideoInfoString(videoInfo, displayAspectRatio)
	data += buildAudioInfoString(audioInfo)
	detailedVideoInfo := buildFullVideoInfoString(videoInfo, displayAspectRatio)
	detailedAudioInfo := buildFullAudioInfoString(audioInfo)

	if audioErr == nil {
		data = "Общая информация:\n---------------\n" + data
	}

	if videoErr == nil {
		data += detailedVideoInfo
	}

	if audioErr == nil {
		data += detailedAudioInfo
	}

	if err := openInNotepad(data); err != nil {
		fmt.Printf("Ошибка при открытии в блокноте: %v\n", err)
	} else {
		os.Exit(0)
	}
}

// Find ffprobe utility func...
func findFFprobe() (string, error) {
	ffprobeName := "ffprobe"
	paths := []string{ffprobeName}

	if path, err := exec.LookPath(ffprobeName); err == nil {
		return path, nil
	}

	if ext := filepath.Ext(ffprobeName); ext == "" {
		paths = append(paths, ffprobeName+getExecutableExtension())
	}

	// Поиск в стандартных путях
	standardPaths := []string{"/usr/bin", "/usr/local/bin", "/usr/local/ffmpegs", "C:/Program Files/ffmpegs"}
	for _, dir := range standardPaths {
		for _, p := range paths {
			if _, err := os.Stat(filepath.Join(dir, p)); err == nil {
				return filepath.Join(dir, p), nil
			}
		}
	}

	return "", fmt.Errorf("ffprobe не найден. Убедитесь, что он установлен и добавлен в PATH.")
}

func getExecutableExtension() string {
	if os.PathSeparator == '\\' {
		return ".exe" // Для Windows
	}
	return "" // Для Linux и macOS
}

func buildVideoInfoString(info *VideoInfo, displayAspectRatio string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s, %dx%d, %s, %s кадр/с, %s кбит/с\n",
		getCodecName(info.CodecName), info.Width, info.Height, getDisplayAspectRatio(displayAspectRatio), info.FrameRate, info.BitRate))
	return builder.String()
}

func buildFullVideoInfoString(info *VideoInfo, displayAspectRatio string) string {
	var builder strings.Builder
	builder.WriteString("\nВидео дорожка:\n---------------\n")
	builder.WriteString(fmt.Sprintf("Codec ID: %s\n", info.CodecName))
	builder.WriteString(fmt.Sprintf("Width: %dpx\n", info.Width))
	builder.WriteString(fmt.Sprintf("Height: %dpx\n", info.Height))
	builder.WriteString(fmt.Sprintf("Display Aspect Ratio: %s\n", displayAspectRatio))
	builder.WriteString(fmt.Sprintf("Frame Rate: %s\n", info.FrameRate))
	builder.WriteString(fmt.Sprintf("Bit Rate: %s\n", info.BitRate))
	return builder.String()
}

func getCodecName(codecName string) string {
	switch codecName {
	case "h264":
		return "MPEG4/ISO/AVC"
	case "h265":
		return "HEVC/H.265"
	default:
		return codecName
	}
}

func getDisplayAspectRatio(aspectRatio string) string {
	if aspectRatio == "" {
		return "N/A"
	}
	return aspectRatio
}

func buildAudioInfoString(info *AudioInfo) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s, %s кГц, %d ch, %s Кбит/с\n",
		info.CodecName, info.SampleRate, info.Channels, info.BitRate))
	return builder.String()
}

func buildFullAudioInfoString(info *AudioInfo) string {
	var builder strings.Builder
	builder.WriteString("\nАудио дорожка:\n---------------\n")
	builder.WriteString(fmt.Sprintf("Codec ID: %s\n", info.CodecName))
	builder.WriteString(fmt.Sprintf("Sampling Rate: %s\n", info.SampleRate))
	builder.WriteString(fmt.Sprintf("Channels: %d\n", info.Channels))
	builder.WriteString(fmt.Sprintf("Bit Rate: %s\n", info.BitRate))
	return builder.String()
}

func openInNotepad(data string) error {
	tmpfile, err := ioutil.TempFile("", "audio_info_*.txt")
	if err != nil {
		return err
	}
	defer tmpfile.Close()

	if _, err := tmpfile.WriteString(data); err != nil {
		return err
	}

	// Команда скрытия окна командной строки
	hideCmd := exec.Command("cmd", "/C", "start", "notepad.exe", tmpfile.Name())

	// Скрытие окна командной строки
	if err := hideCmd.Start(); err != nil {
		return err
	}

	// Ждём завершения скрытия окна командной строки
	if err := hideCmd.Wait(); err != nil {
		return err
	}

	return nil
}

func getVideoInfo(filePath, ffprobePath string) (*VideoInfo, string, error) {
	cmd := exec.Command(ffprobePath,
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=codec_name,width,height,display_aspect_ratio,r_frame_rate,bit_rate",
		"-of", "json",
		filePath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, "", err
	}

	var info struct {
		Streams []VideoInfo `json:"streams"`
	}

	if err := json.Unmarshal(output, &info); err != nil {
		return nil, "", err
	}

	if len(info.Streams) > 0 {
		return &info.Streams[0], info.Streams[0].DisplayAspectRatio, nil
	}

	return nil, "", fmt.Errorf("видео дорожка не найдена")
}

func getAudioInfo(filePath, ffprobePath string) (*AudioInfo, error) {
	cmd := exec.Command(ffprobePath,
		"-v", "error",
		"-select_streams", "a:0",
		"-show_entries", "stream=codec_name,sample_rate,channels,bit_rate",
		"-of", "json",
		filePath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var info struct {
		Streams []AudioInfo `json:"streams"`
	}

	if err := json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	if len(info.Streams) > 0 {
		return &info.Streams[0], nil
	}

	return nil, fmt.Errorf("аудио дорожка не найдена")
}
