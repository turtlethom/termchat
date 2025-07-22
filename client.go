package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	fmt.Print("Enter your name: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	fmt.Fprintf(conn, "%s\n", name)

	app := tview.NewApplication()

	chatView := tview.NewTextView()
	// chatView is *tview.TextView here; it has Write() method

	inputField := tview.NewInputField()
	// inputField is *tview.InputField here; it has SetDoneFunc, GetText, SetText

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := inputField.GetText()
			if strings.TrimSpace(text) != "" {
				fmt.Fprintf(conn, "%s\n", text)
			}
			inputField.SetText("")
		}
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatView, 0, 1, false).
		AddItem(inputField, 3, 1, true)

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			line := scanner.Text()
			app.QueueUpdateDraw(func() {
				chatView.Write([]byte(line + "\n"))
			})
		}
	}()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
