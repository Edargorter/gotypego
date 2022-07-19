//GoTypeGo

package main

import (
	"fmt"
	"os"
	"os/exec"
//	"strings"
	"io/ioutil"
//	"time"
	"log"
	"runtime"
	"golang.org/x/exp/constraints"
	"golang.org/x/term"
)

var esc = map[string]string{"reset" : "\u001b[0m",
							"bg_yellow" : "\u001b[43m",
							"bg_blue" : "\u001b[44m",
							"bg_white" : "\u001b[47;1m",
							"bg_green" : "\u001b[42m",
							"bg_red" : "\u001b[41m",
							"green" : "\u001b[32m",
							"black" : "\u001b[30m",
							"red" : "\u001b[31m",
							"backspace" : "\b\033[K",
							"cursorleft" : "\x1b[1D"}

type TGame struct {
	tg_name string
	tg_filename string
}

//Game vars 
var score int
var miss int
var word string
var text string
var tg TGame
var cmd_buf = make([]byte, 1)
var cmd_mode = false

//OS vars 
var old_state *term.State
var os_cmds = make(map[string] string)

//Helpful funcs 
func min[T constraints.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}

func max[T constraints.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

//Clear screen
func cls() {
	cmd := exec.Command(os_cmds["clear"])
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func get_n_string(s string, n int) string {
	nstr := ""
	for i := 0; i < n; i++ {
		nstr += s
	}
	return nstr
}

func display() {
	cls()
	fmt.Printf("%s%s%s", esc["bg_green"], esc["black"], text[:score]) //Green
	fmt.Printf("%s%s", esc["bg_red"], text[score:score+miss]) //Red 
	fmt.Printf("%s%s\r\n\r\n", esc["reset"], text[score+miss:]) //Reset 
	fmt.Printf(" > %s", word)
}

func quit() {
	term.Restore(int(os.Stdin.Fd()), old_state)
	os.Exit(0)
}

func play(tg TGame, player string) {
	file_bytes, err := ioutil.ReadFile(tg.tg_filename)
	if err != nil {
		panic(err)
	}

	text = string(file_bytes)[:len(file_bytes)-1]

	display()

	score = 0
	miss = 0
	word = ""

	for ;; {
		//Read one byte 
		_, err := os.Stdin.Read(cmd_buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		c := cmd_buf[0]
		if cmd_mode {
			fmt.Print("nothing\r\n")
		} else { //Test character 
			switch c {
				//^C SIGINT -> quit
				case 0x3:
					quit()

				//Backspace character
				case 0x7f, 0x8:
					if len(word) > 0 {
						word = word[:len(word) - 1]
						fmt.Print(esc["backspace"])
						miss = max(miss - 1, 0)
					}

				//Increase green (score) or red 
				default:
					if c >= 0x20 && c <= 0x7e {
						//Correct next char
						if miss == 0 && c == text[score] {
							score++
							if score == len(text) {
								cls()
								fmt.Print("Victory!\r\n")
								quit()
							}
							word += string(c)
							if c == ' ' {
								fmt.Print(get_n_string(esc["backspace"], len(word)))
								word = ""
							}
						} else if miss < len(text) - score {
							miss++
							word += string(c)
						}
					}
			}
			display()
		}
	}
}

func main() {
	//Detect OS and set commands 
	host_os := runtime.GOOS

	if host_os == "Windows" {
		os_cmds["clear"] = "cls"
		os_cmds["remove"] = "del"
	} else {
		os_cmds["clear"] = "clear"
		os_cmds["remove"] = "rm"
	}

	//Terminal Raw Mode 
	prev_state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf(err.Error())
	}
	old_state = prev_state

	test_file := "texts/test_text.txt"
	player := "player1"
	test_game := TGame{tg_name: "test", tg_filename: test_file}
	play(test_game, player)
}
