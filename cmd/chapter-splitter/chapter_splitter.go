package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// global variables
var chapters_path string
var audio_path string

type Timestamp struct {
	Time time.Duration
	Name string
}

type SongSlice struct {
	Name string
	Start int
	Stop int
}

func parse_timestamps(path string) ([]Timestamp, error) {
	timestamps := make([]Timestamp, 0)
	var timestamp Timestamp

	file, err := os.Open(path)
	if err != nil {
		return timestamps, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		// split line into [<time>, <song name>]
		line := strings.SplitN(scanner.Text(), " ", 2)
		if len(line) < 2 {
			continue
		}

		// parse timestamp hours:minutes:seconds
		//t, err := time.ParseDuration(line[0])
		//if err != nil {
		//	fmt.Println("parse duration error:", err)
		//	return timestamps, err
		//}

		colons := strings.Count(line[0], ":")
		if colons == 1 {
			var m,s int
			fmt.Sscanf(line[0], "%d:%d", &m, &s)
			timestamp.Time = time.Duration(m) * time.Minute +
				time.Duration(s) * time.Second
		} else if colons == 2 {
			var h,m,s int
			fmt.Sscanf(line[0], "%d:%d:%d", &h, &m, &s)
			timestamp.Time = time.Duration(h) * time.Hour +
				time.Duration(m) * time.Minute +
				time.Duration(s) * time.Second
		} else {
			fmt.Println("skipping", line, ", cannot parse time")
			continue
		}

		timestamp.Name = strings.Replace(strings.TrimSpace(line[1]), "/", "-", -1)
		timestamps = append(timestamps, timestamp)
	}

	return timestamps, nil
}

func copy_audio_slice(src string, dst string, start float64, duration float64) error {
	cmd := exec.Command("ffmpeg", "-y", "-i", src, "-ss", fmt.Sprint(start), "-t", fmt.Sprint(duration), "-c", "copy", dst)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
	}

	return err
}

func main() {
	flag.StringVar(&chapters_path, "chapters", "", "Path to the chapters text file.")
	flag.StringVar(&audio_path, "audio", "", "Path to the audio file.")
	flag.Parse()

	fmt.Println("Chapter splitter: splitting file", audio_path, "with chapters", chapters_path);

	timestamps, err := parse_timestamps(chapters_path)
	if err != nil {
		fmt.Println("Error parsing timestamps:", err)
		return
	}

	N := len(timestamps)
	for i, timestamp := range timestamps {
		start := timestamp.Time.Seconds() + 1
		// TODO: this will work unless the last song on the file is actually longer
		// than 1 hour
		duration := 3600.0
		if i + 1 < N {
			duration = timestamps[i + 1].Time.Seconds() - start
		}
		name := timestamp.Name
		fmt.Println("Song:", name, "start:", start, "duration:", duration)

		// ffmpeg -ss start -i audio_path name.out -c copy -t duration name-out.opus
		err = copy_audio_slice(audio_path, fmt.Sprintf("out/%s.opus", name), start, duration)
		if err != nil {
			fmt.Println("skipping due to error:", err)
		}
	}
}
