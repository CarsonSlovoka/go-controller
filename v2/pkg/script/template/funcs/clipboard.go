package funcs

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"log"
)

// PasteStr paste a string, support UTF-8,
// write the string to clipboard and tap `cmd + v`
func PasteStr(str string) string {
	robotgo.PasteStr(str)
	return fmt.Sprintf("write the %q to clipboard and paste it.", str)
}

// WriteAll write string to clipboardæ¨å¥½ ä¸ç Hello World!!!ðæ¨å¥½ ä¸ç Hello World!!!ðæ¨å¥½ ä¸ç Hello World!!!ðæ¨å¥½ ä¸ç Hello World!!!ðæ¨å¥½ ä¸ç Hello World!!!ðæ¨å¥½ ä¸ç Hello World!!!ðæ¨å¥½ ä¸ç Hello World!!!ð
func WriteAll(text string) string {
	err := robotgo.WriteAll(text)
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("write string to clipboard: %q", text)
}
