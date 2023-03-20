package service

import (
	"log"
	"worker/model"
	"worker/operator"
)

type UserService struct {
	mongoOperator *operator.MongoOperator
}

func NewUserService(mongoOperator *operator.MongoOperator) *UserService {
	return &UserService{mongoOperator: mongoOperator}
}

func (us *UserService) ShowUser(userName string) (*model.User, error) {
	user, err := us.mongoOperator.GetUser(userName)
	if err != nil {
		log.Println("ShowUser Failed: ", err)
		return nil, err
	}

	return user, nil
}
