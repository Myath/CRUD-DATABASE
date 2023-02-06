package handler

import (
	"errors"
	"log"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jmoiron/sqlx"
)

func rollAlreadyExists(value any) error {
	roll, ok := value.(int)
	if !ok {
		return errors.New("unsupported data given")
	}

	db, err := sqlx.Connect("postgres", "user=postgres password=secret dbname=students_management sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	var student []Student

	if err := db.Select(&student, `SELECT * FROM students`); err != nil {
		log.Fatal(err)
	}

	// var editUser User
	for _, user := range student {
		if user.Roll == roll {
			return errors.New("the roll already exists")
		}
	}
	return nil
}

func emailAlreadyExists(value any) error {
	email, ok := value.(string)
	if !ok {
		return errors.New("unsupported data given")
	}

	db, err := sqlx.Connect("postgres", "user=postgres password=secret dbname=students_management sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	var student []Student

	if err := db.Select(&student, `SELECT * FROM students`); err != nil {
		log.Fatal(err)
	}

	// var editUser User
	for _, user := range student {
		if user.Email == email {
			return errors.New("the email already exists")
		}
	}
	return nil
}

func (s Student) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name,
			validation.Required.Error("The name field is required."),
			validation.Length(3, 32).Error("The name field must be between 3 to 32 characters."),
		),
		validation.Field(&s.Email,
			validation.Required.When(s.ID == 0).Error("The email field is required."),
			is.Email.Error("This email is not valid."),
		),
		validation.Field(&s.Roll,
			validation.Required.When(s.ID == 0).Error("Student roll start from 1"),
			validation.Min(1).Error("Student roll start from 1"),
			validation.Max(200).Error("Only 200 Students are allowed"),
		),
		validation.Field(&s.English,
			validation.Min(0).Error("The lowest mark is 0."),
			validation.Max(100).Error("The highest mark is 100."),
		),
		validation.Field(&s.Bangla,
			validation.Min(0).Error("The lowest mark is 0."),
			validation.Max(100).Error("The highest mark is 100."),
		),
		validation.Field(&s.Mathematics,
			validation.Min(0).Error("The lowest mark is 0."),
			validation.Max(100).Error("The highest mark is 100."),
		),
	)
}

func (s Student) storeAlreadyExistValidate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Email,
			validation.By(emailAlreadyExists),
		),
		validation.Field(&s.Roll,
			validation.By(rollAlreadyExists),
		),
	)
}

func (a LoginAdmin) adminValidation() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Username,
			validation.Required.Error("The username field is required."),
			validation.Length(3, 32).Error("The name field must be between 3 to 32 characters."),
		),
		validation.Field(&a.Password,
			validation.Required.Error("The password field is required."),
			validation.Length(3, 32).Error("The name field must be between 3 to 32 characters."),
		),
	)
}