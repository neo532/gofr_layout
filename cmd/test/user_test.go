package test

import (
	"testing"

	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
)

func TestDemo(t *testing.T) {
	userApiService, _, err := UserApiService()
	if err != nil {
		t.Fatal(err)
	}
	if userApiService == nil {
		t.Fatal("UserApiService is nil")
	}
	user, err := userApiService.GetById(t.Context(), &pb.GetByIdRequest{Id: 1})
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("user is nil")
	}
	t.Logf("user: %v", user)
}
