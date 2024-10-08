package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

var (
	testUsers  []User
	testHouses []House
	testFlats  []Flat
	badUsers   []User
	badHouses  []House
)

func init() {
	testUsers = fakeUsersRegister()
	testHouses = fakeHouses()
	badUsers = fakeUsersRegister()
	badHouses = fakeHouses()
}

func (ts *TestSuite) sendRequest(targetURL, token string, payload []byte) (*http.Response, error) {
	url := fmt.Sprintf("http://apartments/%s", targetURL)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if token != "" {
		req.Header.Add("Authorization", "bearer "+token)
	}

	return http.DefaultClient.Do(req)
}

func checkErrAndCode(ts *TestSuite, err error, res *http.Response) {
	ts.Require().NoError(err)
	ts.Require().Equal(res.StatusCode, http.StatusOK)
}

func checkBadRequest(ts *TestSuite, err error, res *http.Response) {
	ts.Require().NoError(err)
	ts.Require().Equal(res.StatusCode, http.StatusBadRequest)
}

func checkSeverError(ts *TestSuite, err error, res *http.Response) {
	ts.Require().NoError(err)
	ts.Require().Equal(res.StatusCode, http.StatusInternalServerError)
}

func checkUnautorized(ts *TestSuite, err error, res *http.Response) {
	ts.Require().NoError(err)
	ts.Require().Equal(res.StatusCode, http.StatusUnauthorized)
}

func checkResponseData(ts *TestSuite, res *http.Response, checkUnit string) string {
	err := json.NewDecoder(res.Body).Decode(&checkUnit)
	ts.Require().NoError(err)
	ts.Require().NotNil(checkUnit)

	return checkUnit
}

func createUser(ts *TestSuite, user User) string {
	fmt.Println(user)
	jsonUser, err := json.Marshal(user)
	ts.Require().NoError(err)
	res, err := ts.sendRequest("register", "", jsonUser)
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Println(err)
			return
		}
	}()
	checkErrAndCode(ts, err, res)
	checkResponseData(ts, res, user.UserID)
	res, err = ts.sendRequest("login", "", jsonUser)
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Println(err)
			return
		}
	}()
	checkErrAndCode(ts, err, res)
	return checkResponseData(ts, res, user.JWT)
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
