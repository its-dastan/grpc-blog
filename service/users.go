package service

import (
	"errors"
	"fmt"
	"github.com/its-dastan/grpc-blog/db"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const (
	usersCollection = "users"
)

type User struct {
	ID           bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string        `json:"firstName,omitempty" bson:"firstName,omitempty"`
	Email        string        `json:"email,omitempty" bson:"email,omitempty"`
	MobileNumber int64         `json:"mobile_number,omitempty" bson:"mobile_number,omitempty"`
	Password     string        `json:"password,omitempty" bson:"password,omitempty"`
}

func (user *User) Register(server *AuthServer) (string, error) {
	s, c := db.Connect(usersCollection)
	defer s.Close()

	count, err := c.Find(bson.M{"email": user.Email}).Count()
	if err != nil {
		return "", fmt.Errorf("cannot get Users from the database %w", err)
	}
	if count > 1 {
		return "", errors.New("the email already exists")
	}

	userCopy := &User{
		Email:    user.Email,
		Password: user.Password,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("cannot hash password: %w", err)
	}
	user.Password = string(hashedPassword)
	err = c.Insert(user)
	if err != nil {
		return "", errors.New("internal error! please try again later")
	}
	return userCopy.Login(server)
}

func (user *User) Login(server *AuthServer) (string, error) {
	s, c := db.Connect(usersCollection)
	defer s.Close()

	var result *User
	err := c.Find(bson.M{"email": user.Email}).One(&result)
	if err != nil {
		return "", fmt.Errorf("cannot find the email: %w", err)
	}
	if result == nil {
		return "", fmt.Errorf("invalid email id")
	}

	err = comparePassword(result.Password, user.Password)
	if err != nil {
		return "", fmt.Errorf("wrong password")
	}

	err = c.Find(bson.M{"email": user.Email}).One(result)
	if err != nil {
		return "", fmt.Errorf("cannot get the user from the db : %w", err)
	}
	token, err := server.JWTManager.Generate(result)
	log.Println(token)
	return token, nil
}

func comparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
