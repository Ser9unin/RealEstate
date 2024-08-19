package tests

import (
	"encoding/json"
	"log"
	"net/http"
)

func (ts *TestSuite) TestNegativeSet() {
	var res *http.Response

	ts.Run("test register fail", func() {
		for i, user := range badUsers {
			switch {
			case i/2 == 0:
				user.Email = ""
			case i/3 == 0:
				user.HashPass = ""
			case i/5 == 0:
				user.Role = ""
			default:
				user.Email = ""
				user.HashPass = ""
				user.Role = ""
			}
			jsonUser, err := json.Marshal(user)
			ts.Require().NoError(err)
			res, err = ts.sendRequest("register", "", jsonUser)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Println(err)
					return
				}
			}()
			checkBadRequest(ts, err, res)
		}
	})

	ts.Run("test login", func() {
		for _, user := range badUsers {
			// т.к. мы не создали пользователей в базе в процессе test register
			// все запросы будут валиться с ошибкой 500, хотя возможно правильнее сделать с 400 или 404
			jsonUser, err := json.Marshal(user)
			ts.Require().NoError(err)
			res, err = ts.sendRequest("login", "", jsonUser)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Println(err)
					return
				}
			}()
			checkSeverError(ts, err, res)
		}
	})

	ts.Run("test fail house create", func() {
		// т.к. нет зарегистрированных пользователей, то невозможно будет создать дома
		for _, house := range badHouses {
			user := badUsers[0]
			jsonHouse, err := json.Marshal(house)
			ts.Require().NoError(err)
			res, err = ts.sendRequest("house/create", user.JWT, jsonHouse)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Println(err)
					return
				}
			}()
			checkUnautorized(ts, err, res)
		}
	})
	// создаю пользователей в базе как клиентов
	for i, user := range badUsers {
		user.Role = "client"
		jsonUser, err := json.Marshal(user)
		ts.Require().NoError(err)
		res, err = ts.sendRequest("register", "", jsonUser)
		defer func() {
			if err := res.Body.Close(); err != nil {
				log.Println(err)
				return
			}
		}()
		checkErrAndCode(ts, err, res)
		badUsers[i].UserID = checkResponseData(ts, res, user.UserID)
		res, err = ts.sendRequest("login", "", jsonUser)
		defer func() {
			if err := res.Body.Close(); err != nil {
				log.Println(err)
				return
			}
		}()
		checkErrAndCode(ts, err, res)
		badUsers[i].JWT = checkResponseData(ts, res, user.JWT)
	}

	ts.Run("users registered and login as client fail house create", func() {
		for _, house := range badHouses {
			user := badUsers[0]
			jsonHouse, err := json.Marshal(house)
			ts.Require().NoError(err)
			res, err = ts.sendRequest("house/create", user.JWT, jsonHouse)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Println(err)
					return
				}
			}()
			checkUnautorized(ts, err, res)
		}
	})
}
