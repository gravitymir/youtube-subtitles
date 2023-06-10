

![License](https://img.shields.io/badge/license-MIT-green?logo=github)
![Lines of code](https://img.shields.io/tokei/lines/github/gravitymir/ytsubtitles?logo=github)
![Github Repository Size](https://img.shields.io/github/repo-size/gravitymir/ytsubtitles?logo=github)
[![Github forks](https://img.shields.io/github/forks/gravitymir/ytsubtitles?logo=github)](https://github.com/gravitymir/ytsubtitles/network/members)
![Repository Top Language](https://img.shields.io/github/languages/top/gravitymir/ytsubtitles?logo=github)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gravitymir/ytsubtitles/master?logo=go&logoColor=4DD6FD)
![Go Report](https://goreportcard.com/badge/github.com/gravitymir/ytsubtitles)
![GitHub contributors](https://img.shields.io/github/contributors/gravitymir/ytsubtitles?logo=github)

![GitHub release (latest by date)](https://img.shields.io/github/v/release/gravitymir/ytsubtitles?logo=go&logoColor=4DD6FD)
![GitHub last commit](https://img.shields.io/github/last-commit/gravitymir/ytsubtitles?logo=github)
[![GitHub open issues](https://img.shields.io/github/issues/gravitymir/ytsubtitles?logo=github)](https://github.com/gravitymir/ytsubtitles/issues)
[![GitHub closed issues](https://img.shields.io/github/issues-closed/gravitymir/ytsubtitles?logo=github)](https://github.com/gravitymir/ytsubtitles/issues)

[![GitHub Repo stars](https://img.shields.io/github/stars/gravitymir/ytsubtitles?label=ytsubtitles&logo=github&color=505050&logoColor=fff)](https://github.com/gravitymir/ytsubtitles)
[![GitHub User's stars](https://img.shields.io/github/stars/gravitymir?label=gravitymir&logo=github&color=505050&logoColor=fff)](https://github.com/gravitymir)





# ytsubtitles
#### Module which help you scrap subtitles (captions) from YouTube
## Getting Started

``` shell
#terminal | console
go get -u github.com/gravitymir/ytsubtitles
```

## Usage examples

``` go
package main

import (
	"fmt"
	yts "github.com/gravitymir/ytsubtitles"
)

func main() {
	//you can use types of videoID
	// 1 https://www.youtube.com/watch?v=videoID
	// 2 www.youtube.com/watch?v=videoID
	// 3 youtube.com/watch?v=videoID
	// 4 https://www.youtube.com/watch?v=videoID&t=215s
	// 5 https://youtu.be/videoID?t=215
	// 6 videoID
	// any string with videoID
	if mySbt, err := yts.Get("https://youtu.be/CLkkj3aka4g?t=215"); err != nil {
		fmt.Println(err)
	} else {
		if sbtText, err := mySbt.PlainText("Russian"); err == nil {
			fmt.Printf("%+v", mySbt.Subtitles)
			fmt.Println(string(sbtText))
		}
		if sbtJSON, err := mySbt.Json("Russian"); err == nil {
			fmt.Println(string(sbtJSON))
		}
		if sbtJSONPretty, err := mySbt.JsonPretty("Russian"); err == nil {
			fmt.Println(string(sbtJSONPretty))
		}
		fmt.Printf("%+v", mySbt.Tracks)
	}
}