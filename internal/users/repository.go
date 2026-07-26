package users

import "sync"

type userRepository struct {
	users  map[int]*User
	nextID int
	mu     sync.Mutex
}

var repository userRepository = userRepository{users: make(map[int]*User), nextID: 1}

func (u *userRepository) createUser(dto createUserDTO) *User {
	u.mu.Lock()
	user := &User{u.nextID, dto.Name}
	u.users[u.nextID] = user
	u.nextID++
	u.mu.Unlock()

	return user
}

func (u *userRepository) getUserById(id int) *User {
	u.mu.Lock()
	user := u.users[id]
	u.mu.Unlock()

	return user
}
