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
				user.Role = "unexpected"
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
}

func (ts *TestSuite) TestNegativeHousesAndFlats() {
	var res *http.Response
	// создаю в базе клиента и модератора
	client := User{
		Email:    "client@client.com",
		HashPass: "password",
		Role:     "client",
	}
	client.JWT = createUser(ts, client)

	admin := User{
		Email:    "moderator@moderator.com",
		HashPass: "password",
		Role:     "moderator",
	}
	admin.JWT = createUser(ts, admin)

	admin2 := User{
		Email:    "moderator2@moderator.com",
		HashPass: "password",
		Role:     "moderator",
	}
	admin2.JWT = createUser(ts, admin2)

	// проверяем что пользователь со статусом client не может создать дома
	ts.Run("users registered and login as client fail house create", func() {
		for _, house := range badHouses {
			jsonHouse, err := json.Marshal(house)
			ts.Require().NoError(err)
			res, err = ts.sendRequest("house/create", client.JWT, jsonHouse)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Println(err)
					return
				}
			}()
			checkUnautorized(ts, err, res)
		}
	})

	// создаю дома под модератором
	for _, house := range badHouses {
		jsonHouse, err := json.Marshal(house)
		ts.Require().NoError(err)
		res, err = ts.sendRequest("house/create", admin.JWT, jsonHouse)
		defer func() {
			if err := res.Body.Close(); err != nil {
				log.Println(err)
				return
			}
		}()
		checkErrAndCode(ts, err, res)
	}

	// создаю квартиры, т.к. эта функция доступна всем то проблемы не возникнет
	badFlats := fakeFlats(badHouses)
	ts.Run("flat create", func() {
		for _, flat := range badFlats {
			jsonFlat, err := json.Marshal(flat)
			ts.Require().NoError(err)
			res, err = ts.sendRequest("flats/create", client.JWT, jsonFlat)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Println(err)
					return
				}
			}()
			checkErrAndCode(ts, err, res)
		}
	})

	ts.Run("test flat update", func() {
		for _, flat := range testFlats {
			flat.Status = "on moderate"
			jsonFlat, err := json.Marshal(flat)
			ts.Require().NoError(err)
			res, err = ts.sendRequest("flat/update", admin.JWT, jsonFlat)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Println(err)
					return
				}
			}()
			checkErrAndCode(ts, err, res)

			err = json.NewDecoder(res.Body).Decode(&flat)
			ts.Require().NoError(err)
			ts.Require().NotNil(flat.ID)

			// проверяем что другой модератор не может поменять статус квартиры
			flat.Status = "approved"
			jsonFlat, err = json.Marshal(flat)
			ts.Require().NoError(err)
			res, err = ts.sendRequest("flat/update", admin2.JWT, jsonFlat)
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
