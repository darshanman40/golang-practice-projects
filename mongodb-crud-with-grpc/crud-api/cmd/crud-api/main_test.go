package main

import (
	"context"
	"os"
	"testing"

	"github.com/darshanman40/golang-practice-projects/mongodb-crud-with-grpc/crud-api/api/services"
	blogpb "github.com/darshanman40/golang-practice-projects/mongodb-crud-with-grpc/crud-api/internal/proto"

	// . "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: Re-organize entire test
var errInvalidArgument = status.Errorf(codes.InvalidArgument, "Found Empty Blog ")

var bsServer *services.BlogServiceServer
var mongoCtx context.Context

var testCases = []struct {
	name        string
	req         *blogpb.CreateBlogReq
	expectedErr bool
	err         error
}{
	{
		name: "Create Valid Blog",
		req: &blogpb.CreateBlogReq{
			Blog: &blogpb.Blog{
				AuthorId: "abcddd",
				Title:    "TitleAbcddd",
				Content:  "Aahsfbkbibfiubfiubibc",
			},
		},
		expectedErr: false,
		err:         nil,
	},
	{
		name: "Create Invalid Blog",
		req: &blogpb.CreateBlogReq{
			Blog: nil,
		},
		expectedErr: true,
		err:         errInvalidArgument,
	},
}

func TestMain(m *testing.M) {

	code := m.Run()

	os.Exit(code)

}
