package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	userProfile "github.com/AnNosov/simple_user_api/internal/entity"
	"github.com/AnNosov/simple_user_api/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type UserRoutes struct {
	uus usecase.UserActionUseCase
}

type usersResponse struct {
	Users []userProfile.User `json:"users"`
}

func NewUserActionRoutes(router *chi.Mux, t usecase.UserActionUseCase) {
	r := &UserRoutes{t}
	router.Get("/getAllUsers", r.GetUsers) // метод для теста
	router.Post("/createUser", r.CreateUser)
	router.Post("/makeFriends", r.MakeFriends)
	router.Delete("/user", r.DeleteUser)
	router.Get("/friends/{user_id}", r.GetFriends)
	router.Put("/{user_id}", r.UpdateUserAge)
}

func getRequestContent(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	content, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return content, nil
}

func setErrorResponse(w http.ResponseWriter, status int, response string) {
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func (ur *UserRoutes) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := ur.uus.Users()
	if err != nil {
		setErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonAnswer, err := json.Marshal(usersResponse{users})
	if err != nil {
		setErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonAnswer)

}

func (ur *UserRoutes) CreateUser(w http.ResponseWriter, r *http.Request) {

	content, err := getRequestContent(w, r)
	if err != nil {
		setErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var u *userProfile.User
	var userId struct {
		UserId int `json:"CreatedUserId"`
	}

	if err := json.Unmarshal(content, &u); err != nil {
		setErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if u.Name == "" || u.Age < 1 {
		setErrorResponse(w, http.StatusBadRequest, "error: incorrect profile")
		return
	}

	err = ur.uus.CreateUser(u)
	if err != nil {
		setErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userId.UserId = u.Id
	jsonAnswer, err := json.Marshal(userId)
	if err != nil {
		setErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonAnswer)
}

func (ur *UserRoutes) MakeFriends(w http.ResponseWriter, r *http.Request) {

	var friendsIds struct {
		SourceId int `json:"source_id"`
		TargetId int `json:"target_id"`
	}

	content, err := getRequestContent(w, r)
	if err != nil {
		setErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.Unmarshal(content, &friendsIds); err != nil {
		setErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = ur.uus.MakeFriends(friendsIds.SourceId, friendsIds.TargetId); err != nil {
		setErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sourceName, err := ur.uus.GetUserName(friendsIds.SourceId)
	if err != nil {
		setErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	targetName, err := ur.uus.GetUserName(friendsIds.TargetId)
	if err != nil {
		setErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	answer := []byte(sourceName + " and " + targetName + " are friends")
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}

func (ur *UserRoutes) DeleteUser(w http.ResponseWriter, r *http.Request) {

	var userIdObject struct {
		TargetId int `json:"target_id"`
	}

	content, err := getRequestContent(w, r)
	if err != nil {
		setErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.Unmarshal(content, &userIdObject); err != nil {
		setErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	deletedNameUser, err := ur.uus.GetUserName(userIdObject.TargetId)
	if err != nil {
		setErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err = ur.uus.DeleteUser(userIdObject.TargetId); err != nil {
		setErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	answer := deletedNameUser + " has been deleted"

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func (ur *UserRoutes) GetFriends(w http.ResponseWriter, r *http.Request) {

	var userFriendsObject struct {
		UserFriendNames []string `json:"UserFriendNames"`
	}

	userFriendsObject.UserFriendNames = make([]string, 0)

	userIdStr := chi.URLParam(r, "user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		setErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userFriendsObject.UserFriendNames, err = ur.uus.GetFriends(userId)
	if err != nil {
		setErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonAnswer, err := json.Marshal(userFriendsObject.UserFriendNames)
	if err != nil {
		setErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonAnswer)
}

func (ur *UserRoutes) UpdateUserAge(w http.ResponseWriter, r *http.Request) {

	var userAgeObject struct {
		NewAge int `json:"new_age"`
	}

	userIdStr := chi.URLParam(r, "user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		setErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	content, err := getRequestContent(w, r)
	if err != nil {
		setErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.Unmarshal(content, &userAgeObject); err != nil {
		setErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = ur.uus.UpdateUserAge(userId, userAgeObject.NewAge); err != nil {
		setErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User age updated successfully."))
}
