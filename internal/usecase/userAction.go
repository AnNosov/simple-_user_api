package usecase

import (
	"fmt"
	"log"

	userProfile "github.com/AnNosov/simple_user_api/internal/entity"
	"github.com/AnNosov/simple_user_api/internal/usecase/repo"
)

type UserActionUseCase struct {
	repo repo.UserRepo
}

func New(r repo.UserRepo) *UserActionUseCase {
	return &UserActionUseCase{
		repo: r,
	}
}

func checkFriend(friendsList []int, friendId int) bool {
	for _, val := range friendsList {
		if val == friendId {
			return true
		}
	}
	return false
}

func (uc *UserActionUseCase) Users() ([]userProfile.User, error) {
	users, err := uc.repo.Users()
	if err != nil {
		return nil, fmt.Errorf("UserAction - Users: %w", err)
	}
	return users, nil
}

func (uc *UserActionUseCase) CreateUser(u *userProfile.User) error {
	id, err := uc.repo.SequenceUserId()
	if err != nil {
		return fmt.Errorf("UserAction - CreateUser: %w", err)
	}
	u.Id = id + 1
	err = uc.repo.CreateUser(u)
	if err != nil {
		return fmt.Errorf("UserAction - CreateUser: %w", err)
	}
	return nil
}

func (uc *UserActionUseCase) MakeFriends(sourceId, targetId int) error {
	checkSourceId, err := uc.repo.CheckUser(sourceId)
	if err != nil {
		return fmt.Errorf("UserAction - MakeFriends: %w", err)
	}

	checkTargetId, err := uc.repo.CheckUser(targetId)
	if err != nil {
		return fmt.Errorf("UserAction - MakeFriends: %w", err)
	}

	if !checkSourceId || !checkTargetId {
		return fmt.Errorf("sourceId or targetId was not found")
	}

	sFds, err := uc.repo.GetFriendsIds(sourceId)
	if err != nil {
		return fmt.Errorf("UserAction - MakeFriends: %w", err)
	}

	tFds, err := uc.repo.GetFriendsIds(targetId)
	if err != nil {
		return fmt.Errorf("UserAction - MakeFriends: %w", err)
	}

	if !checkFriend(sFds, targetId) {
		if err = uc.repo.AddFriend(sourceId, targetId); err != nil {
			return fmt.Errorf("UserAction - MakeFriends: %w", err)
		}
	}

	if !checkFriend(tFds, sourceId) {
		if err = uc.repo.AddFriend(targetId, sourceId); err != nil {
			return fmt.Errorf("UserAction - MakeFriends: %w", err)
		}
	}

	return nil
}

func (uc *UserActionUseCase) DeleteUser(id int) error {
	checkId, err := uc.repo.CheckUser(id)
	if err != nil {
		return fmt.Errorf("UserAction - DeleteUser: %w", err)
	}
	if !checkId {
		return fmt.Errorf("ID was not found")
	}

	friendList, err := uc.repo.GetFriendsIds(id)
	if err != nil {
		return fmt.Errorf("UserAction - DeleteUser: %w", err)
	}

	for _, val := range friendList {
		if err = uc.repo.DeleteFriend(val, id); err != nil {
			log.Println("UserAction - DeleteUser: ", err.Error())
			continue
		}
	}

	if err = uc.repo.DeleteUser(id); err != nil {
		return fmt.Errorf("UserAction - DeleteUser: %w", err)
	}
	return nil
}

func (uc *UserActionUseCase) GetFriends(id int) ([]string, error) {

	friendIds, err := uc.repo.GetFriendsIds(id)
	if err != nil {
		return nil, fmt.Errorf("UserAction - GetFriends: %w", err)
	}

	if len(friendIds) == 0 {
		return nil, fmt.Errorf("UserAction - GetFriends: friend list is empty")
	}

	friendList := make([]string, 0)

	for _, val := range friendIds {
		valName, err := uc.repo.GetUserName(val)
		if err != nil {
			log.Println("UserAction - GetFriends: ", err.Error())
		}

		if err = uc.repo.DeleteFriend(val, id); err != nil {
			log.Println("UserAction - GetFriends: ", err.Error())
			continue
		}
		friendList = append(friendList, valName)
	}
	return friendList, nil
}

func (uc *UserActionUseCase) UpdateUserAge(id, age int) error {
	checkId, err := uc.repo.CheckUser(id)
	if err != nil {
		return fmt.Errorf("UserAction - UpdateUserAge: %w", err)
	}
	if !checkId {
		return fmt.Errorf("ID was not found")
	}

	if err := uc.repo.UpdateUserAge(id, age); err != nil {
		return fmt.Errorf("UserAction - UpdateUserAge: %w", err)
	}
	return nil
}

func (uc *UserActionUseCase) GetUserName(id int) (string, error) {
	userName, err := uc.repo.GetUserName(id)
	if err != nil {
		return "", fmt.Errorf("UserAction - GetUserName: %w", err)
	}
	return userName, nil
}
