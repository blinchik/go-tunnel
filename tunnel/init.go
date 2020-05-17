package tunnel

import (
	"fmt"
	"os"
	"time"
)

var home, err = os.UserHomeDir()

//Recoverer autostart a function after crushing
func Recoverer(maxPanics int, id string, f func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover ", id)
			fmt.Println(err)

			if maxPanics == 0 {
				panic("TOO MANY PANICS")
			} else {
				time.Sleep(2 * time.Second)
				go Recoverer(maxPanics-1, id, f)
			}
		}
	}()
	fmt.Println("start ", id)

	f()
}

func init() {

}
