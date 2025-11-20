package types

type Student struct {
	Id    int64
	Name  string `validate:"required"`
	Age   int    `validate:"required,min=0"`
	Email string `validate:"required,email"`
}
