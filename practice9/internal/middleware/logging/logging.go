package logging

import (
	"io"
	"log"
	"os"
)

func Load() {
	// Открытие файла для логирования
	file, err := os.OpenFile("history.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Не удалось открыть файл лога:", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// Настройка вывода ошибок в файл и стандартный поток вывода
	log.SetOutput(io.MultiWriter(file, os.Stdout))
}
