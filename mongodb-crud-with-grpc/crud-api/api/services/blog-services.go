package services

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	blogpb "github.com/darshanman40/golang-practice-projects/mongodb-crud-with-grpc/crud-api/internal/proto"
)

// BlogServiceServer ...
type BlogServiceServer struct {
	db     *mongo.Client
	blogdb *mongo.Collection
}

var errorDetailBuilder strings.Builder

//BlogItem ...
type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

//Init ...
func (s *BlogServiceServer) Init(db *mongo.Client, collectionName string) {
	s.db = db
	s.blogdb = db.Database("mydb").Collection(collectionName)
}

//ReadBlog read blog article
func (s *BlogServiceServer) ReadBlog(ctx context.Context, req *blogpb.ReadBlogReq) (*blogpb.ReadBlogRes, error) {

	errorDetailBuilder.Reset()
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	result := s.blogdb.FindOne(ctx, bson.M{"_id": oid})
	data := BlogItem{}

	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	response := &blogpb.ReadBlogRes{
		Blog: &blogpb.Blog{
			Id: oid.Hex(),
		},
	}
	return response, nil
}

//CreateBlog create new article
func (s *BlogServiceServer) CreateBlog(ctx context.Context, req *blogpb.CreateBlogReq) (*blogpb.CreateBlogRes, error) {
	blog := req.GetBlog()
	errorDetailBuilder.Reset()

	if blog == nil {
		errorDetailBuilder.WriteString("Found Empty Blog ")
		return nil, status.Errorf(codes.InvalidArgument, errorDetailBuilder.String())
	}

	data := BlogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	result, err := s.blogdb.InsertOne(ctx, data)
	if err != nil {
		errorDetailBuilder.WriteString("Internal Error: ")
		errorDetailBuilder.WriteString(err.Error())
		return nil, status.Errorf(codes.Internal, errorDetailBuilder.String())
	}
	oid := result.InsertedID.(primitive.ObjectID)
	blog.Id = oid.Hex()

	return &blogpb.CreateBlogRes{Blog: blog}, nil
}

//UpdateBlog update existing article
func (s *BlogServiceServer) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogReq) (*blogpb.UpdateBlogRes, error) {
	blog := req.GetBlog()
	errorDetailBuilder.Reset()

	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		errorDetailBuilder.WriteString("Could not convert the supplied blog id to a MongoDB ObjectId: ")
		errorDetailBuilder.WriteString(err.Error())
		return nil, status.Errorf(codes.InvalidArgument, errorDetailBuilder.String())
	}

	update := bson.M{
		"author_id": blog.GetAuthorId(),
		"title":     blog.GetTitle(),
		"content":   blog.GetContent(),
	}

	filter := bson.M{"_id": oid}
	result := s.blogdb.FindOneAndUpdate(ctx, filter, bson.M{"$set": update},
		options.FindOneAndUpdate().SetReturnDocument(1))

	decoded := BlogItem{}
	if err := result.Decode(&decoded); err != nil {
		errorDetailBuilder.WriteString("Could not find blog with supplied ID: ")
		errorDetailBuilder.WriteString(err.Error())
		return nil, status.Error(codes.NotFound, errorDetailBuilder.String())
	}
	return &blogpb.UpdateBlogRes{
		Blog: &blogpb.Blog{
			Id:       decoded.ID.Hex(),
			AuthorId: decoded.AuthorID,
			Title:    decoded.Title,
			Content:  decoded.Content,
		},
	}, nil
}

//DeleteBlog Delete the blog by ID
func (s *BlogServiceServer) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogReq) (*blogpb.DeleteBlogRes, error) {

	errorDetailBuilder.Reset()

	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		errorDetailBuilder.WriteString("Could not convert to ObjectId: ")
		errorDetailBuilder.WriteString(err.Error())
		return nil, status.Error(codes.InvalidArgument, errorDetailBuilder.String())
	}

	if _, err := s.blogdb.DeleteOne(ctx, bson.M{"_id": oid}); err != nil {
		errorDetailBuilder.WriteString("Could not find/delete blog with id  ")
		errorDetailBuilder.WriteString(req.GetId())
		errorDetailBuilder.WriteString(": ")
		errorDetailBuilder.WriteString(err.Error())
		return nil, status.Error(codes.NotFound, errorDetailBuilder.String())
	}
	return &blogpb.DeleteBlogRes{
		Success: true,
	}, nil
}

//ListBlogs retrieve/stream list of blogs
func (s *BlogServiceServer) ListBlogs(req *blogpb.ListBlogsReq, stream blogpb.BlogService_ListBlogsServer) error {

	data := &BlogItem{}
	errorDetailBuilder.Reset()
	cursor, err := s.blogdb.Find(context.Background(), bson.M{})
	if err != nil {
		errorDetailBuilder.WriteString("Unknown Internal error: ")
		errorDetailBuilder.WriteString(err.Error())
		return status.Error(codes.Internal, errorDetailBuilder.String())
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		if err := cursor.Decode(data); err != nil {
			errorDetailBuilder.WriteString("Could not decode data: ")
			errorDetailBuilder.WriteString(err.Error())
			return status.Error(codes.Unavailable, errorDetailBuilder.String())
		}

		stream.Send(&blogpb.ListBlogsRes{
			Blog: &blogpb.Blog{
				Id:       data.ID.Hex(),
				AuthorId: data.AuthorID,
				Content:  data.Content,
				Title:    data.Title,
			},
		})
	}

	if err := cursor.Err(); err != nil {
		errorDetailBuilder.WriteString("Unknown cursor error: ")
		errorDetailBuilder.WriteString(err.Error())
		return status.Error(codes.Internal, errorDetailBuilder.String())
	}

	return nil
}
