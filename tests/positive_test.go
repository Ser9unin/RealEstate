package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var testUsers []User
var testHouses []House
var testFlats []Flat

func init() {
	testUsers = fakeUsersRegister()
	testHouses = fakeHouses()
}

func (ts *TestSuite) TestPositiveSet() {
	var res *http.Response

	ts.Run("test register", func() {
		for i, user := range testUsers {
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
			// заношу полученный UUID в тестувую модель пользователей для следующих тестов
			testUsers[i].UserID = checkResponseData(ts, res, user.UserID)
		}
	})

	ts.Run("test login by uuid", func() {
		for i, user := range testUsers {
			jsonUser, err := json.Marshal(user)
			ts.Require().NoError(err)
			res, err = ts.sendRequest("login", "", jsonUser)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Println(err)
					return
				}
			}()
			checkErrAndCode(ts, err, res)
			// заношу полученный JWT в тестувую модель пользователей для следующих тестов
			testUsers[i].JWT = checkResponseData(ts, res, user.JWT)
		}
	})

	ts.Run("test login by email", func() {
		for _, user := range testUsers {
			// убираю UUID из модели пользователя что бы убедиться что аутентификация производится по Email
			user.UserID = ""
			jsonUser, err := json.Marshal(user)
			ts.Require().NoError(err)
			res, err = ts.sendRequest("login", "", jsonUser)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Println(err)
					return
				}
			}()
			checkErrAndCode(ts, err, res)
			checkResponseData(ts, res, user.JWT)
		}
	})

	ts.Run("test house create", func() {
		for i, house := range testHouses {
			user := testUsers[0]
			jsonHouse, err := json.Marshal(house)
			ts.Require().NoError(err)
			res, err = ts.sendRequest("house/create", user.JWT, jsonHouse)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Println(err)
					return
				}
			}()
			checkErrAndCode(ts, err, res)

			err = json.NewDecoder(res.Body).Decode(&house)
			ts.Require().NoError(err)
			ts.Require().NotNil(house.ID)
			// заношу данные по дому возвращенному из БД в тестовый набор для следующих тестов
			testHouses[i] = house
		}
	})

	// создаём тестовый набор квартир на основе домов, которые заполнили в тесте "test houses create"
	testFlats = fakeFlats(testHouses)
	ts.Run("test flat create", func() {
		for i, flat := range testFlats {
			user := testUsers[0]
			jsonFlat, err := json.Marshal(flat)
			ts.Require().NoError(err)
			res, err = ts.sendRequest("flat/create", user.JWT, jsonFlat)
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

			testFlats[i] = flat
		}
	})

	ts.Run("test flat update", func() {
		for i, flat := range testFlats {
			user := testUsers[0]
			flat.Status = "approved"
			jsonFlat, err := json.Marshal(flat)
			ts.Require().NoError(err)
			res, err = ts.sendRequest("flat/update", user.JWT, jsonFlat)
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

			testFlats[i] = flat
		}
	})

	ts.Run("test get flats by house ID", func() {
		for _, house := range testHouses {
			user := testUsers[0]
			flats := Flats{}
			jsonHouseID, err := json.Marshal(house.ID)
			ts.Require().NoError(err)

			requestStr := fmt.Sprintf("house/%d", house.ID)
			res, err = ts.sendRequest(requestStr, user.JWT, jsonHouseID)
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Println(err)
					return
				}
			}()
			checkErrAndCode(ts, err, res)

			err = json.NewDecoder(res.Body).Decode(&flats)
			ts.Require().NoError(err)
			ts.Require().NotNil(flats)
		}
	})
}
