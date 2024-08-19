package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/stretchr/testify/require"
)

func setDataBase(ctx context.Context) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode= %s",
		"0.0.0.0",
		"5433",
		"dev",
		"somepass",
		"apartmentsdb",
		"disable")

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("unable to start db: ", err)
	}
	db.SetMaxOpenConns(1000) // Максимальное количество открытых подключений
	db.SetMaxIdleConns(50)   // Максимальное количество простаивающих подключений
	db.SetConnMaxLifetime(5 * time.Second)
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal("context cancelled", err)
	}
	return db, nil
}

func TestHouse(t *testing.T) {
	ctx := context.Background()
	db, err := setDataBase(ctx)
	require.NoError(t, err)
	defer db.Close()

	clearHousesTable(db)

	storage := New(db)
	testHouses := fakeHouses()
	fmt.Println(len(testHouses))
	t.Run("positive test house add to db", func(t *testing.T) {
		for _, house := range testHouses {
			// проверяем запрос на создание дома, должен вернуть данные о доме с датой создания
			resultCreate, err := storage.NewHouse(ctx, house)
			require.NoError(t, err)
			require.NotNil(t, resultCreate.ID)
			require.Equal(t, house.Address, resultCreate.Address)
			require.Equal(t, house.Year, resultCreate.Year)
			require.Equal(t, house.Developer, resultCreate.Developer)
			require.NotNil(t, resultCreate.CreateAt)
			require.NotNil(t, resultCreate.UpdateAt)

			house.ID = resultCreate.ID
		}
	})

	// делаем запрос данных о доме ожидаем успешный ответ т.к. запрашиваем существующий ID
	t.Run("positive test get houses list from DB", func(t *testing.T) {
		// запрашиваем список всех домов
		resultHouses, errDB := storage.HousesList(ctx)
		require.NoError(t, errDB)
		for i, house := range testHouses {
			require.NotNil(t, resultHouses[i].ID)
			require.Equal(t, house.Address, resultHouses[i].Address)
			require.Equal(t, house.Year, resultHouses[i].Year)
			require.Equal(t, house.Developer, resultHouses[i].Developer)
			require.NotNil(t, resultHouses[i].CreateAt)
		}
	})

	t.Run("positive test get by single house from DB", func(t *testing.T) {
		resultHouses, errDB := storage.HousesList(ctx)
		require.NoError(t, errDB)
		for _, house := range resultHouses {
			// запрашиваем по одному дому
			resultHouse, errDB := storage.House(ctx, house.ID)
			require.NoError(t, errDB)
			require.NotNil(t, resultHouse.ID)
			require.Equal(t, house.Address, resultHouse.Address)
			require.Equal(t, house.Year, resultHouse.Year)
			require.Equal(t, house.Developer, resultHouse.Developer)
			require.NotNil(t, resultHouse.CreateAt)
		}
	})
	clearHousesTable(db)
}

func clearHousesTable(db *sql.DB) {
	_, err := db.Exec("DELETE from houses;")
	if err != nil {
		log.Println("Ошибка удаления данных:", err)
		return
	}
	log.Println("Таблица houses очищена")
}

func TestFlat(t *testing.T) {
	ctx := context.Background()
	db, err := setDataBase(ctx)
	require.NoError(t, err)
	defer db.Close()

	clearHousesTable(db)
	clearFaltsTable(db)

	storage := New(db)
	testHouses := fakeHouses()
	for _, house := range testHouses {
		_, err := storage.NewHouse(ctx, house)
		require.NoError(t, err)
	}

	testHousesFromDB, err := storage.HousesList(ctx)
	require.NoError(t, err)
	// создаём базу квартир на основе базы домов
	testFlats := fakeFlats(testHousesFromDB)
	t.Run("positive test flat add to db", func(t *testing.T) {
		for _, flat := range testFlats {
			// проверяем запрос на создание квартиры, должен вернуть данные о квартире с датой создания
			resultCreate, err := storage.NewFlat(ctx, flat)
			require.NoError(t, err)
			require.NotNil(t, resultCreate.ID)
			require.Equal(t, flat.HouseID, resultCreate.HouseID)
			require.Equal(t, flat.Price, resultCreate.Price)
			require.Equal(t, flat.Rooms, resultCreate.Rooms)
			require.Equal(t, resultCreate.Status, "created")
		}

		// после создания квартир в домах БД должно обновиться поле updated_at проверяем,
		// для этого делаем запрос списка домов снова
		updTestHouses, errDB := storage.HousesList(ctx)
		require.NoError(t, errDB)
		for i, house := range updTestHouses {
			require.NotEqual(t, house.UpdateAt, testHouses[i].UpdateAt)
		}
	})

	// делаем запрос данных о доме ожидаем успешный ответ т.к. запрашиваем существующий ID
	t.Run("positive test moderator request flat from DB", func(t *testing.T) {
		for _, house := range testHousesFromDB {
			testFlatsFromDB, err := storage.FlatsList(ctx, house.ID, moderator)
			require.NoError(t, err)
			for _, flat := range testFlatsFromDB {
				// проверяем запрос квартиры, должен вернуть данные о квартире с датой создания
				resultRequest, err := storage.Flat(ctx, moderator, flat.HouseID, flat.ID)
				require.NoError(t, err)
				require.Equal(t, flat.ID, resultRequest.ID)
				require.Equal(t, flat.HouseID, resultRequest.HouseID)
				require.Equal(t, flat.Price, resultRequest.Price)
				require.Equal(t, flat.Rooms, resultRequest.Rooms)
				require.Equal(t, resultRequest.Status, "created")
			}
		}
	})

	t.Run("negative test user request flat from DB", func(t *testing.T) {
		for _, house := range testHousesFromDB {
			testFlatsFromDB, err := storage.FlatsList(ctx, house.ID, moderator)
			require.NoError(t, err)
			for _, flat := range testFlatsFromDB {
				// проверяем запрос на создание квартиры, должен вернуть данные о квартире с датой создания
				_, err := storage.Flat(ctx, "user", flat.HouseID, flat.ID)
				require.Error(t, err)
			}
		}
	})

	t.Run("positive test moderator update status", func(t *testing.T) {
		user := User{
			UserID:   "cae36e0f-69e5-4fa8-a179-a52d083c5549",
			Email:    "moderator@moderator.com",
			HashPass: "password",
			Role:     "moderator",
		}

		for _, house := range testHousesFromDB {
			testFlatsFromDB, err := storage.FlatsList(ctx, house.ID, moderator)
			require.NoError(t, err)
			for _, flat := range testFlatsFromDB {
				// проверяем запрос квартиры, должен вернуть данные о квартире с датой создания
				resultRequest, err := storage.UpdateFlatStatus(ctx, user, "on moderate", flat)
				require.NoError(t, err)
				require.Equal(t, flat.ID, resultRequest.ID)
				require.Equal(t, flat.HouseID, resultRequest.HouseID)
				require.Equal(t, flat.Price, resultRequest.Price)
				require.Equal(t, flat.Rooms, resultRequest.Rooms)
				require.Equal(t, resultRequest.Status, "on moderate")
			}
		}
	})

	t.Run("negative test moderator update status", func(t *testing.T) {
		user := User{
			UserID:   "cae36e0f-69e5-4fa8-a179-a52d083c566",
			Email:    "moderator2@moderator.com",
			HashPass: "password",
			Role:     "moderator",
		}

		for _, house := range testHousesFromDB {
			testFlatsFromDB, err := storage.FlatsList(ctx, house.ID, moderator)
			require.NoError(t, err)
			for _, flat := range testFlatsFromDB {
				// проверяем запрос квартиры, должен вернуть данные о квартире с датой создания
				_, err := storage.UpdateFlatStatus(ctx, user, "approved", flat)
				require.Error(t, err)
			}
		}
	})

	clearHousesTable(db)
	clearFaltsTable(db)
}

func clearFaltsTable(db *sql.DB) {
	_, err := db.Exec("Delete from flats;")
	if err != nil {
		log.Println("Ошибка удаления данных:", err)
		return
	}
	log.Println("Таблица flats очищена")
}

func TestUsers(t *testing.T) {
	ctx := context.Background()
	db, err := setDataBase(ctx)
	require.NoError(t, err)
	defer db.Close()
	storage := New(db)
	testUsers := fakeUsersRegister()
	fmt.Println(len(testUsers))

	t.Run("positive test new user to db", func(t *testing.T) {
		for _, user := range testUsers {
			// проверяем запрос на создание дома, должен вернуть данные о доме с датой создания
			resultCreate, err := storage.NewUser(ctx, user)
			require.NoError(t, err)
			require.Equal(t, user.UserID, resultCreate.UserID)
			require.Equal(t, user.Role, resultCreate.Role)
		}
	})

	t.Run("positive test get user by ID and pass from DB", func(t *testing.T) {
		for _, user := range testUsers {
			// запрашиваем по одному дому
			resultLogin, err := storage.UserByID(ctx, user.UserID)
			require.NoError(t, err)
			require.Equal(t, user.UserID, resultLogin.UserID)
			require.Equal(t, user.Role, resultLogin.Role)
		}
	})

	t.Run("positive test get user by Email and pass from DB", func(t *testing.T) {
		for _, user := range testUsers {
			// запрашиваем по одному дому
			resultLogin, err := storage.UserByEmail(ctx, user.Email)
			require.NoError(t, err)
			require.Equal(t, user.UserID, resultLogin.UserID)
			require.Equal(t, user.Role, resultLogin.Role)
		}
	})

	t.Run("positive test check user by ID and pass from DB", func(t *testing.T) {
		for _, user := range testUsers {
			// запрашиваем по одному дому
			resultCheck, err := storage.UserByIDAndRole(ctx, user.UserID, user.Role)
			require.NoError(t, err)
			require.True(t, resultCheck)
		}
	})

	_, err = db.Exec("TRUNCATE users;")
	if err != nil {
		log.Println("Ошибка удаления данных:", err)
		return
	}
	log.Println("Таблица users очищена")
}
