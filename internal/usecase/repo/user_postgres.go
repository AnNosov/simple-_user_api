package repo

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	userProfile "github.com/AnNosov/simple_user_api/internal/entity" // иначе заменяет на os/user
	"github.com/AnNosov/simple_user_api/pkg/postgres"
	"github.com/lib/pq"
)

const defaultUserCap = 64

type UserRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) Users() ([]userProfile.User, error) {
	rows, err := r.Postgres.DB.Query("select * from  skillbox.users")
	if err != nil {
		return nil, fmt.Errorf("user_postgres - Users: %w", err)
	}
	defer rows.Close()

	users := make([]userProfile.User, 0, defaultUserCap)

	for rows.Next() {
		u := userProfile.User{}

		err := rows.Scan(&u.Id, &u.Name, &u.Age, pq.Array(&u.Friends))
		if err != nil {
			log.Println("user_postgres - Users: ", err)
			continue
		}
		users = append(users, u)
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("users from postgres is empty")
	}
	return users, nil
}

func (r *UserRepo) CheckUser(id int) (bool, error) {
	var check bool
	err := r.Postgres.DB.QueryRow("select exists(select * from  skillbox.users where id = $1)", id).Scan(&check)
	if err != nil {
		return check, fmt.Errorf("user_postgres - CheckUser: %w", err)
	}
	return check, nil
}

func (r *UserRepo) SequenceUserId() (int, error) {
	var id sql.NullInt64
	err := r.Postgres.DB.QueryRow("select max(id) from  skillbox.users").Scan(&id)
	if err == sql.ErrNoRows || !id.Valid {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("user_postgres - SequenceUserId: %w", err)
	}
	return int(id.Int64), nil
}

func (r *UserRepo) CreateUser(u *userProfile.User) error {
	_, err := r.Postgres.DB.Exec("insert into skillbox.users (id, name, age, friends) values ($1, $2, $3, $4)", u.Id, u.Name, u.Age, pq.Array(u.Friends))
	if err != nil {
		return fmt.Errorf("user_postgres - CreateUser: %w", err)
	}
	return nil
}

func (r *UserRepo) GetUserName(id int) (string, error) {
	var name string
	err := r.Postgres.DB.QueryRow("select name from skillbox.users where id = $1", id).Scan(&name)
	if err != nil {
		return "", fmt.Errorf("user_postgres - GetUserName: %w", err)
	}
	return name, nil
}

func (r *UserRepo) GetFriendsIds(id int) ([]int, error) {
	var friends []string
	friendList := make([]int, 0)
	err := r.Postgres.DB.QueryRow("select friends from skillbox.users where id = $1", id).Scan(pq.Array(&friends))
	if err != err {
		return nil, fmt.Errorf("user_postgres - GetFriends: %w", err)
	}

	for _, val := range friends {
		friendId, err := strconv.Atoi(val)
		if err != err {
			return nil, fmt.Errorf("user_postgres - GetFriends: %w", err)
		}
		friendList = append(friendList, friendId)
	}
	return friendList, nil
}

func (r *UserRepo) AddFriend(id, friendId int) error {
	_, err := r.Postgres.DB.Exec("update skillbox.users set friends = array_append(friends, $1) where id = $2", strconv.Itoa(friendId), id)
	if err != nil {
		return fmt.Errorf("user_postgres - AddFriend: %w", err)
	}
	return nil
}

func (r *UserRepo) DeleteFriend(id, friendId int) error {
	_, err := r.Postgres.DB.Exec("update skillbox.users set friends = array_remove(friends, $1) where id = $2", strconv.Itoa(friendId), id)
	if err != nil {
		return fmt.Errorf("user_postgres - DeleteFriend: %w", err)
	}
	return nil
}

func (r *UserRepo) DeleteUser(id int) error {
	_, err := r.Postgres.DB.Exec("delete from skillbox.users where id = $1", id)
	if err != nil {
		return fmt.Errorf("user_postgres - DeleteUser: %w", err)
	}
	return nil
}

func (r *UserRepo) UpdateUserAge(id, age int) error {
	_, err := r.Postgres.DB.Exec("update skillbox.users set age = $1 where id = $2", age, id)
	if err != nil {
		return fmt.Errorf("user_postgres - UpdateUserAge: %w", err)
	}
	return nil
}
