// Play sound files via button press
// YAML file will consist of what button contains what sound
// sound folder will hold sound files
// sound files will need to be MP3s for now

// we need to install gcc for windows
// need to install 64 bit gcc... why is 32bit default...

// while waiting for gcc to install lets get the sound libraries in

package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"gopkg.in/yaml.v2"
)

// Config
// struct to unmarshal config.yaml to
// must use capital letters to be public
// golang is weird
type Config struct {
	BtnOne struct {
		Title    string `yaml:"title"`
		FileName string `yaml:"filename"`
	} `yaml:"btn_one"`
}

// SndPlayer
// Plays sounds associated with buttons
func SndPlayer(playLocalPath string, playFileName string) error {
	// open mp3 file
	playFilePath := path.Join(
		path.Dir(playLocalPath),
		"sounds\\"+playFileName,
	)

	playFile, err := os.Open(playFilePath)

	if err != nil {
		log.Printf(
			"Error opening file %s: %v",
			playFileName,
			err,
		)
		return err
	}

	// decode file
	decodeFile, err := mp3.NewDecoder(playFile)
	if err != nil {
		log.Printf(
			"Error decoding file %s: %v",
			playFileName,
			err,
		)
		return err
	}

	// use oto universal audio player and play the file
	// will explore more in the options
	// NewPlayer pass
	//  the sample rate
	//  channel number
	//    - mostly 2
	//  bit depth in bytes
	//    - bytes per sample per channel
	//    - mostly 2
	//  bufferSizeInBytes
	//    - size of the context to play the sound
	//    - might need to pass the size of file to this?
	audioPlayer, err := oto.NewPlayer(decodeFile.SampleRate(), 2, 2, int(decodeFile.Length()))
	if err != nil {
		log.Printf(
			"Error playing file %s: %v",
			playFileName,
			err,
		)
		return err
	}

	// copy file data to player buffer...
	log.Printf(
		"Loading %s [%d[bytes]]",
		playFilePath,
		decodeFile.Length(),
	)

	if _, err := io.Copy(audioPlayer, decodeFile); err != nil {
		log.Printf("Failed playing file: %v", err)
		return err
	}

	// close player when done
	// maybe not defer
	audioPlayer.Close()

	return nil
}

// main
// main function of the app
// runs the GUI
func main() {
	// import yaml configuration
	// grabbing this from another go app I wrote

	// reading files sucks
	// lets get the runtime path of the go app
	_, localPath, _, ok := runtime.Caller(0)
	log.Printf("Opening local path for YAML file: %s", localPath)
	if ok == false {
		log.Fatalf("Getting localPath failed")
	}

	// use that to get our "local" path
	localConfigPath := path.Join(path.Dir(localPath), "config.yaml")

	// open the file
	configData, err := ioutil.ReadFile(localConfigPath)
	if err != nil {
		log.Fatalf("error opening yaml file: %v", err)
	}

	// create an empty local config struct
	// pass it by reference along with a buffer to unmarshal
	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		log.Fatalf("error reading yaml: %v", err)
	}

	// initialize app from fyne lib
	app := app.New()

	// set new window with title
	appWindow := app.NewWindow("Soundboard v1")

	// add some content to the app window
	appWindow.SetContent(
		// create a vertical box to hold btn, labels, etc
		widget.NewVBox(
			// label
			widget.NewLabel("List of Buttons!"),

			// first sound button
			// will add in name configuration from YAML
			widget.NewButton(config.BtnOne.Title, func() {
				log.Printf("Playing sound %s", config.BtnOne.FileName)

				// send to sound player
				err := SndPlayer(
					localPath,
					config.BtnOne.FileName,
				)
				if err != nil {
					log.Fatalf("Playing sound failed: %v", err)
				}

			}),

			// quit button
			widget.NewButton("Quit", func() {
				app.Quit()
			}),
		),
	)

	appWindow.ShowAndRun()
}
