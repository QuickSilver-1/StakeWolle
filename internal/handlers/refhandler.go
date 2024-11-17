package handlers

import (
	"context"
	"net/http"
	dh "referal/internal/dbhandlers"
	"referal/pkg/db"
	"referal/pkg/server"
	"strconv"
	"strings"
	"sync"
	"time"
)

func GenRef(w http.ResponseWriter, r *http.Request) {
	tokenS := r.Header.Get("Authorization")

	tokenS = strings.Split(tokenS, " ")[1]
	email, err := server.DecodeJWT(tokenS)

	if err != nil {
		server.AnswerHandler(w, 500, err.Error())
		return
	}

	code := genHash(email + time.Now().String())
	out := make(chan string)

	go db.DB.Query("get", out, dh.DBData{
		UserEmail: email,
	})

	dayS := r.FormValue("day")
	day, err := strconv.Atoi(dayS)

	if err != nil {
		server.AnswerHandler(w, 400, "Неверный формат данных")
		return
	}

	exist := <-out

	if exist != "" {
		server.AnswerHandler(w, 400, "Перед генерацией нового кода необходимо удалить старый")
		return
	}

	wg := &sync.WaitGroup{}

	var rout chan string
	wg.Add(1)
	go func() {
		defer wg.Done()
		res := <-out

		if res != "success" {
			server.AnswerHandler(w, 500, "Ошибка записи в базу данных")
			return
		}

		ctx := context.Background()
		rout = make(chan string)
		go db.NewKey(ctx, email, code, rout)
	}()

	out = make(chan string)
	go db.DB.Query("generate", out, dh.DBData{
		RefString: code,
		DayExpires: time.Now().AddDate(0, 0, day),
		UserEmail: email})

	wg.Wait()
	res := <-rout

	if res != "success" {
		server.AnswerHandler(w, 500, res)
	}

	server.AnswerHandler(w, 200, code)
}

func DelRef(w http.ResponseWriter, r *http.Request) {
	tokenS := r.Header.Get("Authorization")

	tokenS = strings.Split(tokenS, " ")[1]
	email, err := server.DecodeJWT(tokenS)

	ctx := context.Background()
	rout := make(chan string)
	go db.DelKey(ctx, email, rout)
	
	if err != nil {
		server.AnswerHandler(w, 500, err.Error())
		return
	}

	out := make(chan string)
	go db.DB.Query("delete", out, dh.DBData{
		UserEmail: email,
	})

	res := <-rout
	if res == "exist" {
		rout = make(chan string)
		go db.DelKey(ctx, email, rout)
	}

	ans := <-out
	sliseAns := strings.Split(ans, " ")
	code, _ := strconv.Atoi(sliseAns[0])
	value := strings.Join(sliseAns[1:], " ")

	res = <-rout
	if res != "success" {
		server.AnswerHandler(w, 500, res)
		return
	}
	
	server.AnswerHandler(w, code, value)
}

func GetCode(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	ctx := context.Background()
	rout := make(chan string)
	go db.KeyExist(ctx, email, rout)

	out := make(chan string)
	check := make(chan string)

	exist := <-rout
	if exist == "exist" {
		rout = make(chan string)
		db.KeyExist(ctx, email, rout)
		
		code := <-rout
		server.AnswerHandler(w, 200, code)
		return
	}

	go db.DB.Query("get", out, dh.DBData{
		UserEmail: email,
	})
	go db.DB.Query("check", check, dh.DBData{
		UserEmail: email,
	})
	
	var code string
	var ref bool = true
	select {
	case code = <-out:
		if code == "" {
			ref = false
			exist := <-check

			if exist == "" {
				server.AnswerHandler(w, 400, "Пользователя с такой почтой не существует")
				return
		}
	}
	
	case exist := <-check:
		if exist == "" {
			server.AnswerHandler(w, 400, "Пользователя с такой почтой не существует")
			return
		}

		code = <-out
		if code == "" {
			ref = false
		}
	}

	if !ref {
		server.AnswerHandler(w, 400, "У этого пользователя нет реферального кода")
		return
	}

	server.AnswerHandler(w, 200, code)
}

func GetRefs(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	ctx := context.Background()
	rout := make(chan string)
	go db.KeyExist(ctx, id, rout)

	check := make(chan string)

	exist := <-rout
	if exist == "exist" {
		rout = make(chan string)
		db.KeyExist(ctx, id, rout)
		
		code := <-rout
		server.AnswerHandler(w, 200, code)
		return
	}

	go db.DB.Query("checkid", check, dh.DBData{
		UserId: id,
	})

	out := make(chan string)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	pass := <-check
	if pass == "" {
		server.AnswerHandler(w, 400, "Пользователя с таким id не существует")
		return
	}
 
	go func() {
		defer wg.Done()
		users := []string{}

		for ref := range out {
			users = append(users, ref)
		}

		if len(users) == 0 {
			server.AnswerHandler(w, 400, "У данного пльзователя нет рефералов")
			return
		}

		server.AnswerHandler(w, 200, users)

		rout = make(chan string)
		go db.NewKey(ctx, id, users, rout)
	}()

	go db.DB.Query("referals", out, dh.DBData{
		UserId: id,
	})

	wg.Wait()

	res := <-rout
	if res != "success" {
		server.AnswerHandler(w, 500, res)
	}
}