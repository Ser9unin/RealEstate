package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Ser9unin/RealEstate/internal/register"
	"github.com/Ser9unin/RealEstate/internal/render"
	repository "github.com/Ser9unin/RealEstate/internal/storage/repo"
)

type api struct {
	storage *repository.Queries
	logger  Logger
}

func newAPI(storage *repository.Queries, logger Logger) api {
	return api{
		storage: storage,
		logger:  logger,
	}
}

func (a *api) greetings(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>This is my real estate service!</h1>"))
}

// type dummyuser struct {
// 	Role string `json:"user_type"`
// }

// func (a *api) dummyLogin(w http.ResponseWriter, r *http.Request) {
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Error reading request body", http.StatusBadRequest)
// 		return
// 	}
// 	defer r.Body.Close()

// 	user := dummyuser{}
// 	err = json.Unmarshal(body, &user)
// 	if err != nil {
// 		render.ErrorJSON(w, r, http.StatusBadRequest, err, "Invalid request payload")
// 		return
// 	}

// }

func (a *api) register(w http.ResponseWriter, r *http.Request) {
	newUser := register.NewUserService(a.storage, a.logger)
	newUser.Register(w, r)
}

func (a *api) login(w http.ResponseWriter, r *http.Request) {
	newUser := register.NewUserService(a.storage, a.logger)
	newUser.Login(w, r)
}

func (a *api) houseCreate(w http.ResponseWriter, r *http.Request) {
	var newHouse repository.House
	err := json.NewDecoder(r.Body).Decode(&newHouse)
	if err != nil {
		a.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusBadRequest, err, "")
		return
	}
	newHouse, err = a.storage.NewHouse(r.Context(), newHouse)
	if err != nil {
		a.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("can't create new house"), "")
		return
	}
	render.ResponseJSON(w, r, http.StatusOK, newHouse)
}

type Flats struct {
	FlatsList []repository.Flat `json:"flats"`
}

func (a *api) houseFlats(w http.ResponseWriter, r *http.Request) {
	houseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		a.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusBadRequest, err, "")
		return
	}

	role := w.Header().Get("role")
	flatsList, err := a.storage.FlatsList(r.Context(), houseID, role)
	if err != nil {
		a.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("can't get flats list"), "")
		return
	}

	FlatsResponse := Flats{
		FlatsList: flatsList,
	}
	render.ResponseJSON(w, r, http.StatusOK, FlatsResponse)
}

func (a *api) houseSubscribe(w http.ResponseWriter, r *http.Request) {

}

func (a *api) flatCreate(w http.ResponseWriter, r *http.Request) {
	var newFlat repository.Flat
	err := json.NewDecoder(r.Body).Decode(&newFlat)
	if err != nil {
		a.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusBadRequest, err, "")
		return
	}
	newFlat, err = a.storage.NewFlat(r.Context(), newFlat)
	if err != nil {
		a.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("can't create flat"), "")
		return
	}
	render.ResponseJSON(w, r, http.StatusOK, newFlat)
}

func (a *api) flatUpdate(w http.ResponseWriter, r *http.Request) {
	var updFlat repository.Flat
	err := json.NewDecoder(r.Body).Decode(&updFlat)
	if err != nil {
		a.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusBadRequest, err, "")
		return
	}
	updFlat, err = a.storage.UpdateFlatStatus(r.Context(), "moderator", updFlat.Status, updFlat.HouseId, updFlat.ID)
	if err != nil {
		a.logger.Error(err.Error())
		render.ErrorJSON(w, r, http.StatusNotFound, fmt.Errorf("no such flat"), "")
		return
	}
	render.ResponseJSON(w, r, http.StatusOK, updFlat)
}
