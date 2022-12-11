package main

/*
https://developer.fyne.io/started/

After go get:
	go mod tidy

*/
import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io"
	"log"
	"net/http"
)

type Note struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Text    string `json:"text"`
}

func main() {
	myBaseUrl := "http://127.0.0.1:4862/"
	myApp := app.New()
	myWindow := myApp.NewWindow("Note4U")

	entryUrl := widget.NewEntry()
	entryUrl.SetText(myBaseUrl)
	entryName, entrySurname, entryText, entryID := widget.NewEntry(), widget.NewEntry(), widget.NewEntry(), widget.NewEntry()
	labelSendStatus := widget.NewLabel("Статус отправки:")
	labelNotesCount := widget.NewLabel("Заметок на сервере: ??")

	btnUpdateNotesCount := widget.NewButtonWithIcon("Принудительно обновить кол-во заметок", theme.ViewRefreshIcon(), func() {
		isOk := false
		resp, err := http.Get(entryUrl.Text + "getNoteCount")
		if err == nil {
			body, bodyErr := io.ReadAll(resp.Body)
			if bodyErr == nil {
				isOk = true
				labelNotesCount.SetText("Заметок на сервере " + string(body))
			}
		}
		if !isOk {
			dialog.ShowInformation("Внимание", "Ошибка обновления кол-ва заметок.\nПроверьте адрес сервера!", myWindow)
			labelNotesCount.SetText("Заметок на сервере: ??")
		}
	})
	btnSend := widget.NewButtonWithIcon("Отправка заметки на сервер", theme.UploadIcon(), func() {
		values := map[string]string{"name": entryName.Text, "surname": entrySurname.Text, "text": entryText.Text}
		jsonValue, _ := json.Marshal(values)
		resp, err := http.Post(entryUrl.Text+"createNote", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			labelSendStatus.SetText("Статус отправки: ОШИБКА!")
			dialog.ShowInformation("Внимание", "Ошибка отправки заметки.\nПроверьте адрес сервера!", myWindow)
			log.Fatal(err, resp)
		} else {
			labelSendStatus.SetText("Статус отправки: Успешно!\nИмя: " + entryName.Text + "\nФамилия: " + entrySurname.Text + "\nТекст: " + entryText.Text)
			fmt.Println("Send")

		}

		resp, err = http.Get(entryUrl.Text + "getNoteCount")
		if err == nil {
			body, bodyErr := io.ReadAll(resp.Body)
			if bodyErr == nil {
				labelNotesCount.SetText("Заметок на сервере (ID отправленной заметки): " + string(body))
			}

		}

	})

	labelGetNote := widget.NewLabel("Имя: ??\nФамилия: ??\nЗаметка: ??")
	btnGet := widget.NewButtonWithIcon("Получение заметки с сервера", theme.DownloadIcon(), func() {
		resp, err := http.Get(entryUrl.Text + "readNote?id=" + entryID.Text)
		isOk := false
		if err == nil && resp.StatusCode == 200 {
			decoder := json.NewDecoder(resp.Body)
			var note Note
			err := decoder.Decode(&note)
			if err == nil {
				labelGetNote.SetText("Имя: " + note.Name + "\nФамилия: " + note.Surname + "\nЗаметка: " + note.Text)
				isOk = true
			}
		} else {
			fmt.Print("ERR: ")
			fmt.Println(err, resp)
		}
		if !isOk {
			dialog.ShowInformation("Внимание", "Ошибка получения заметки.\nПроверьте ID и адрес сервера!", myWindow)
			labelGetNote.SetText("Имя: ??\nФамилия: ??\nЗаметка: ??")
		}

	})
	myWindow.SetContent(container.NewVBox(
		widget.NewLabel("Сервер ( \"http://ip:port/\" ):"),
		entryUrl,
		widget.NewSeparator(),
		widget.NewLabel("Введите имя:"),
		entryName,
		widget.NewLabel("Введите фамилию:"),
		entrySurname,
		widget.NewLabel("Введите текст заметки:"),
		entryText,
		btnSend,
		labelSendStatus,
		widget.NewSeparator(),
		btnUpdateNotesCount,
		labelNotesCount,
		widget.NewSeparator(),
		widget.NewLabel("Введите ID заметки:"),
		entryID,
		btnGet,
		labelGetNote,
	))
	myApp.Settings().SetTheme(theme.DarkTheme())
	myWindow.ShowAndRun()
}
