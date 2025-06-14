package request

type UpdateUserRequest struct {
	Email       string `json:"email" validate:"omitempty,email"`
	Name        string `json:"name" validate:"omitempty,min=2,max=100"`
	Phone       string `json:"phone" validate:"omitempty,min=10,max=20"`
	DateOfBirth string `json:"date_of_birth" validate:"omitempty,date"`
	Gender      string `json:"gender" validate:"omitempty,oneof='male' 'female' ''"`
	Address     string `json:"address" validate:"omitempty,max=255"`
	City        string `json:"city" validate:"omitempty,max=100"`
	Country     string `json:"country" validate:"omitempty,max=100"`
	Avatar      string `json:"avatar" validate:"omitempty,url,max=255"`
	Bio         string `json:"bio" validate:"omitempty"`
	ParentName  string `json:"parent_name" validate:"omitempty,max=100"`
	ParentPhone string `json:"parent_phone" validate:"omitempty,min=10,max=20"`
	ParentEmail string `json:"parent_email" validate:"omitempty,email,max=100"`
}
