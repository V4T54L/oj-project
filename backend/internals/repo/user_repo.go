package repo

import (
	"algo-arena-be/internals/models"
	"context"
	"errors"
)

type userRepo struct {
	mockUser *models.UserInfo
}

func (u *userRepo) GetUserByUserID(
	ctx context.Context, userId int,
) (*models.UserInfo, error) {
	if userId != 1 {
		return nil, errors.New("user not found")
	}

	return u.mockUser, nil
}

func (u *userRepo) CreateUser(ctx context.Context, username, email, avatarUrl, hashedPassword string) (int, error) {
	return 1, nil
}

func (u *userRepo) GetUserByCreds(ctx context.Context, username, hashedPassword string) (*models.UserInfo, error) {
	return u.mockUser, nil
}
