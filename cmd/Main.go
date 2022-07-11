package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"password-storage/internal/app"
	Consts "password-storage/internal/pkg/consts"
	spec "password-storage/internal/pkg/models"
)

type server struct {
}

func (s server) validateSession(ctx echo.Context) error {
	sessionId, err := ctx.Cookie("sessionId")
	if err != nil {
		return err
	}

	sessionValid := app.IsSessionValid(sessionId.Value)

	if sessionValid != true {
		err := json.NewEncoder(ctx.Response()).Encode(fmt.Errorf("session expired"))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s server) PostApiV1CreateNewPasswordData(ctx echo.Context) error {
	err := s.validateSession(ctx)
	if err != nil {
		return err
	}
	newPassDataRequest := spec.NewPasswordDataRequest{}
	err = json.NewDecoder(ctx.Request().Body).Decode(&newPassDataRequest)
	if err != nil {
		return err
	}
	ch := make(chan spec.NewPasswordDataResponse)
	go func(newPassDataRequest spec.NewPasswordDataRequest) {
		ch <- app.CreateNewPasswordData(newPassDataRequest)
	}(newPassDataRequest)
	err = json.NewEncoder(ctx.Response()).Encode(<-ch)
	if err != nil {
		return err
	}
	return nil
}

func (s server) GetApiV1GetAllPasswordData(ctx echo.Context) error {
	err := s.validateSession(ctx)
	if err != nil {
		return err
	}
	ch := make(chan []spec.AllPasswordData)
	go func() {
		ch <- app.GetAllPasswordData()
	}()
	err = json.NewEncoder(ctx.Response()).Encode(<-ch)
	if err != nil {
		return err
	}
	return nil
}

func (s server) GetApiV1GetPasswordDataByKeyKey(ctx echo.Context, key string) error {
	err := s.validateSession(ctx)
	if err != nil {
		return err
	}
	ch := make(chan spec.PasswordDataByKey)
	go func(key interface{}) {
		ch <- app.GetPasswordDataByKey(key)
	}(key)
	err = json.NewEncoder(ctx.Response()).Encode(<-ch)
	if err != nil {
		return err
	}
	return nil
}

func (s server) PostApiV1SignIn(ctx echo.Context) error {
	login := spec.LoginDataRequest{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&login)
	if err != nil {
		return err
	}
	if *login.User != Consts.DefaultUser || *login.Password != Consts.DefaultPassword {
		return fmt.Errorf("incorrect login or password")
	}
	ch := make(chan string)
	go func() {
		ch <- app.CreateNewSession()
	}()

	sessionId := <-ch
	response := spec.LoginDataResponse{
		SessionId: &sessionId,
	}
	err = json.NewEncoder(ctx.Response()).Encode(response)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	//Вызываем бота
	go func() {
		app.TelegramBot()
	}()

	e := echo.New()
	spec.RegisterHandlers(e, server{})

	defer func(Con *pg.DB) {
		err := Con.Close()
		if err != nil {
			fmt.Println("Не удалось закрыть соединение")
		}
	}(app.Con)

	e.Logger.Fatal(e.Start(":8080"))
}
