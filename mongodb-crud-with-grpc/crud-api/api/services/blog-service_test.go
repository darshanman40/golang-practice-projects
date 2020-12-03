package services

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	blogpb "github.com/darshanman40/golang-practice-projects/mongodb-crud-with-grpc/crud-api/internal/proto"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errInvalidArgument = status.Errorf(codes.InvalidArgument, "Found Empty Blog ")
var errInvalidID = status.Errorf(codes.InvalidArgument, "the provided hex string is not a valid ObjectID")

var bsServer *BlogServiceServer
var mongoCtx context.Context
var tempID string

const threeSeconds = time.Duration(3) * time.Second

var createTestCases = []struct {
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

var readTestCases = []struct {
	name        string
	req         *blogpb.ReadBlogReq
	expectedErr bool
	err         error
}{
	{
		name: "Read Non-existant Blog",
		req: &blogpb.ReadBlogReq{
			Id: "123123",
		},
		expectedErr: true,
		err:         errInvalidID,
	},
}

func setup() error {
	if mongoCtx == nil {
		log.Println("Connecting to MongoDB...")
		mongoCtx = context.Background()

	}
	if bsServer == nil {
		db, err := mongo.Connect(mongoCtx,
			options.Client().ApplyURI("mongodb://localhost:27017"),
			options.Client().SetConnectTimeout(threeSeconds), // TODO: not working currently
		)
		if err != nil {
			log.Println("Failed to connect Mongodb ")
			log.Fatal(err)
			return err
		}
		log.Println("Connected to Mongodb ")

		bsServer = &BlogServiceServer{}
		bsServer.Init(db, "blog_test")
	}
	if err := bsServer.db.Ping(mongoCtx, nil); err != nil {
		bsServer.db, err = mongo.Connect(mongoCtx,
			options.Client().ApplyURI("mongodb://localhost:27017"),
			options.Client().SetConnectTimeout(threeSeconds), // TODO: redundunt, Have to remove later
		)
		if err != nil {
			fmt.Println("Failed to connect Mongodb ")
			log.Fatal(err)
			return err
		}
		fmt.Println("Connected to Mongodb ")

	}
	return nil
}

func endSetup() {
	fmt.Println("Disconnecting Mongodb...")
	bsServer.db.Disconnect(mongoCtx)
}

func getBlogID() {
	setup()
}

func TestCreateBlog(t *testing.T) {
	if err := setup(); err != nil {
		t.FailNow()
	}

	for _, tc := range createTestCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := NewGomegaWithT(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ctx := context.Background()
			response, err := bsServer.CreateBlog(ctx, tc.req)

			// assert results expectations
			if tc.expectedErr {
				g.Expect(response).To(BeNil(), "Result should be nil")
				g.Expect(err).To(Equal(tc.err))
			} else {
				g.Expect(response.Blog.AuthorId).To(Equal(tc.req.Blog.AuthorId))
				g.Expect(response.Blog.Title).To(Equal(tc.req.Blog.Title))
				g.Expect(response.Blog.Content).To(Equal(tc.req.Blog.Content))
			}
		})
	}
	endSetup()
}

func TestReadBlog(t *testing.T) {
	setup()

	for _, tc := range readTestCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := NewGomegaWithT(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ctx := context.Background()
			response, err := bsServer.ReadBlog(ctx, tc.req)
			if tc.expectedErr {
				g.Expect(response).To(BeNil(), "Result should be nil")
				g.Expect(err).To(Equal(tc.err))
			} else {
				g.Expect(response.Blog.Id).To(Equal(tc.req.Id))

			}
		})
	}
	endSetup()
}
