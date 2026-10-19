// Команда admin создаёт администратора с паролем, скрытым при вводе.
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"studyblog/internal/app"
	"studyblog/internal/config"
	"studyblog/internal/service"
)

// main выводит понятную ошибку создания администратора.
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
}

// run запрашивает данные, регистрирует аккаунт и выдаёт административную роль.
func run() error {
	c, err := config.Load(true)
	if err != nil {
		return err
	}

	a, err := app.New(c)
	if err != nil {
		return err
	}

	defer a.DB.Close()
	reader := bufio.NewReader(os.Stdin)
	read := func(prompt string) (string, error) {
		fmt.Print(prompt)
		line, err := reader.ReadString('\n')
		return strings.TrimSpace(line), err
	}

	var in service.Registration
	in.Login, err = read("Логин (3–24 символа: a-z, 0-9, _): ")
	if err != nil {
		return err
	}

	in.Name, err = read("Имя (2–60 символов): ")
	if err != nil {
		return err
	}

	in.Email, err = read("Электронная почта: ")
	if err != nil {
		return err
	}

	fmt.Print("Пароль (не менее 10 символов, ввод скрыт): ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return err
	}

	fmt.Print("Повторите пароль: ")
	repeat, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return err
	}

	if string(password) != string(repeat) {
		return fmt.Errorf("пароли не совпадают")
	}

	in.Password = string(password)
	_, err = a.Service.RegisterAdmin(context.Background(), in)
	if err != nil {
		return err
	}

	fmt.Println("Администратор создан.")
	return nil
}
