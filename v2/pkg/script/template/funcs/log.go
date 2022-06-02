package funcs

import (
	"log"
)

func Log(msg ...string) string { // 除非想特別標明時間，不然產出預設使用: t.Execute(os.Stdout, context)，會直接把結果打印到終端機上
	for _, curMsg := range msg {
		if curMsg != "" {
			log.Println(curMsg)
		}
	}
	return ""
}
