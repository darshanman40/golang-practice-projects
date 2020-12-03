package mocks

import (
	"context"

	blogpb "github.com/darshanman40/golang-practice-projects/mongodb-crud-with-grpc/crud-api/internal/proto"
	mock "github.com/stretchr/testify/mock"
)

type BlogServiceHelper struct {
	mock.Mock
}

func (_m *BlogServiceHelper) ReadBlog(ctx context.Context, req *blogpb.ReadBlogReq) (*blogpb.ReadBlogRes, error) {
	ret := _m.Called(ctx, req)
	var r0 *blogpb.ReadBlogRes
	if rf, ok := ret.Get(0).(func(context.Context, *blogpb.ReadBlogReq) (*blogpb.ReadBlogRes, error)); ok {
		r0, _ = rf(ctx, req)
	}
	// else {
	// 	r0 = ret.Get(0).(blogpb.ReadBlogReq)
	// }

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
