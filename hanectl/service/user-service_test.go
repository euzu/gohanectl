package service

import (
	"github.com/stretchr/testify/assert"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/test/mock_test"
	"testing"
)

func TestFindByUsername(t *testing.T) {

	repoMock := new(mock_test.UserRepoMock)
	srv := UserService {
		userRepo: repoMock,
	}
	repoMock.On("FindByUsername", "sarcon").Return(&model.User{Username: "sarcon"}, nil)

	user, err := srv.FindByUsername("sarcon")
	assert.Nil(t, err)
	assert.NotNil(t, user)

	repoMock.AssertExpectations(t)
}

