package main

/*
https://developer.fyne.io/started/

After go get:
	go mod tidy

*/
import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Note struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Text    string `json:"text"`
}

func getNotesCount(baseUrl string) string {
	resp, err := http.Get(baseUrl + "getNoteCount")
	if err == nil {
		body, bodyErr := io.ReadAll(resp.Body)
		if bodyErr == nil {
			return string(body)
		}
	}
	return "Неизвестно"
}

func getNoteById(baseUrl, id string) string {
	resp, err := http.Get(baseUrl + "readNote?id=" + id)
	if err == nil && resp.StatusCode == 200 {
		decoder := json.NewDecoder(resp.Body)
		var note Note
		err := decoder.Decode(&note)
		if err == nil {
			ans := "Имя: " + note.Name + "\nФамилия: " + note.Surname + "\nЗаметка: " + note.Text
			return ans
		}
	}
	return "Ошибка!\nИмя: ?\nФамилия: ?\nЗаметка: ?"
}

func sendNote(baseUrl, name, surname, text string) bool {
	values := map[string]string{"name": name, "surname": surname, "text": text}
	jsonValue, _ := json.Marshal(values)
	resp, err := http.Post(baseUrl+"createNote", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err, resp)
		return false

	} else {
		return true

	}
}

func clearShell() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	myBaseUrl := "http://127.0.0.1:4862/"
	var name, surname, text, isContinue, noteID string

	scanner := bufio.NewScanner(os.Stdin)
	var sl []Note
	flag := false
	for true {
		if flag {
			break
		}

		for true {
			fmt.Print("Ввести новую заметку?\ny-да\tn-nope (exit)\ta-вывести заметку\n>>")
			scanner.Scan()
			isContinue = scanner.Text()

			if isContinue == "n" {
				flag = true
				break
			} else if isContinue == "a" {
				clearShell()
				fmt.Println("Всего заметок: ", getNotesCount(myBaseUrl))
				fmt.Println("Введите ID")
				scanner.Scan()
				noteID = scanner.Text()
				clearShell()
				fmt.Println(getNoteById(myBaseUrl, noteID))
			} else if isContinue == "y" {
				clearShell()
				fmt.Println("Введите имя:")
				scanner.Scan()
				name = scanner.Text()
				fmt.Println("Введите фамилию:")
				scanner.Scan()
				surname = scanner.Text()
				fmt.Println("Введите текст:")
				scanner.Scan()
				text = scanner.Text()
				clearShell()
				fmt.Print("Surname: "+surname+"\nName: "+name+"\ntext: ", text, "\nСтатус отправки:")
				if sendNote(myBaseUrl, name, surname, text) {
					fmt.Println("Успешно.")
				} else {
					fmt.Println("Ошибка!")
				}
				fmt.Println("Всего заметок: ", getNotesCount(myBaseUrl))
				row := Note{name, surname, text}
				sl = append(sl, row)
			}

		}

	}
	clearShell()
	fmt.Println("Bye!")
}
