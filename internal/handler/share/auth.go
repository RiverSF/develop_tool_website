package share

import (
	"errors"

	"develop_tools/internal/model"
)

var (
	errUUIDRequired = errors.New("uuid required")
	errUserNotFound = errors.New("user not found")
	errForbidden    = errors.New("forbidden")
)

func resolveUserID(uuid string) (int, error) {
	if uuid == "" {
		return 0, errUUIDRequired
	}
	userKey := model.NewUserKeyModel().GetUserIdByKey(uuid)
	if userKey.UserId <= 0 {
		return 0, errUserNotFound
	}
	return userKey.UserId, nil
}

func loadOwnedShare(id, userID int) (*model.Share, error) {
	share := model.NewShareModel().GetShareDataById(id, userID)
	if share.Id == 0 {
		return nil, errForbidden
	}
	return share, nil
}
