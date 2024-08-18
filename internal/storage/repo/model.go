package repository

import "time"

type House struct {
	ID        int       `json:"house_id,omitempty" faker:"-"`
	Address   string    `json:"address,omitempty" faker:"-"`
	Year      int       `json:"year,omitempty" faker:"boundary_start=1900, boundary_end=2024"`
	Developer string    `json:"developer,omitempty" faker:"word"`
	CreateAt  time.Time `json:"created_at,omitempty" faker:"-"`
	UpdateAt  time.Time `json:"updated_at,omitempty" faker:"-"`
}

type Flat struct {
	ID      int    `json:"id" faker:"-"`
	HouseId int    `json:"house_id" faker:"-"`
	Price   int    `json:"price,omitempty" faker:"boundary_start=1000000, boundary_end=200000000"`
	Rooms   int    `json:"rooms,omitempty" faker:"boundary_start=1, boundary_end=6"`
	Status  string `json:"status,omitempty" faker:"-"`
}

type User struct {
	UserID   string `json:"id,omitempty" faker:"uuid_hyphenated"`
	Email    string `json:"email,omitempty" faker:"email"`
	HashPass string `json:"password,omitempty" faker:"password"`
	Role     string `json:"user_type,omitempty" faker:"-"`
}
