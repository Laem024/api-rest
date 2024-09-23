package viewmodel

import "paralelos/model"

type UserViewModel struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
}

// Transformación del Model a ViewModel
func NewUserViewModel(user model.User) UserViewModel {
    return UserViewModel{
        ID:   user.ID,
        Name: user.Name,
		Email: user.Email,
    }
}
