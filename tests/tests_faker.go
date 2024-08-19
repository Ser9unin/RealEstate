package tests

import (
	"fmt"
	"log"

	"github.com/go-faker/faker/v4"
)

func fakeNewFlat(house House) Flat {
	fakeFlat := new(Flat)
	err := faker.FakeData(fakeFlat)
	if err != nil {
		log.Fatalf(fmt.Sprintf("can't create fake data: %s", err.Error()))
	}
	fakeFlat.HouseID = house.ID

	return *fakeFlat
}

func fakeFlats(houses []House) []Flat {
	fakeFlats := make([]Flat, 0, 110)
	for _, house := range houses {
		for i := 0; i < 10; i++ {
			fakeFlat := fakeNewFlat(house)
			fakeFlats = append(fakeFlats, fakeFlat)
		}
	}

	return fakeFlats
}

func fakeNewHouse() House {
	fakeHouse := new(House)
	err := faker.FakeData(fakeHouse)
	if err != nil {
		log.Fatalf(fmt.Sprintf("can't create fake data: %s", err.Error()))
	}
	fakerAddress := faker.GetRealAddress()

	fakeHouse.Address = fakerAddress.Address

	return *fakeHouse
}

func fakeHouses() []House {
	fakeHouses := make([]House, 0, 10)
	capHouses := cap(fakeHouses)
	for i := 0; i < capHouses; i++ {
		fakeHouse := fakeNewHouse()
		fakeHouses = append(fakeHouses, fakeHouse)
	}
	return fakeHouses
}

func fakeUsersRegister() []User {
	fakeUsers := make([]User, 0, 10)
	capHouses := cap(fakeUsers)
	for i := 0; i < capHouses; i++ {
		fakeUser := fakeNewUser()
		fakeUsers = append(fakeUsers, fakeUser)
	}
	return fakeUsers
}

func fakeNewUser() User {
	fakeUser := new(User)
	err := faker.FakeData(fakeUser)
	if err != nil {
		log.Fatalf(fmt.Sprintf("can't create fake data: %s", err.Error()))
	}

	fakeUser.Role = "moderator"

	return *fakeUser
}
