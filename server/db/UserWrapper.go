package db

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"
)

// Get FrameIDs available for download for a given User
func GetAvailableFramesForDownload(userID int) ([]int32, error) {
	return global.GetFrameIDsForDownload(context.Background(), int32(userID))
}

func GetPrevileges(userName string, password string) ([]string, error) {
	user, err := global.GetUserDetails(context.Background(), userName)
	if err != nil {
		return nil, err
	}

	temp := md5.Sum([]byte(password))
	var hash strings.Builder
	for i := 0; i < len(temp); i++ {
		hash.Write([]byte(fmt.Sprintf("%02X", temp[i])))
	}
	if strings.Compare(hash.String(), user.Password) == 0 {
		permissions := strings.Split(user.Permissions, ",")
		return permissions, nil
	}
	return nil, fmt.Errorf("password wrong for user %s", userName)
}
