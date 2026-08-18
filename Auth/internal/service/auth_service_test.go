package service_test

import (
	"context"
	"database/sql"
	"net/url"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/abelmalu/golang-posts/Auth/internal/dto"
	model "github.com/abelmalu/golang-posts/Auth/internal/models"
	"github.com/abelmalu/golang-posts/Auth/internal/service"
	"github.com/abelmalu/golang-posts/platform"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc/metadata"
)

type MockAuthRepository struct {
	User *model.User
	err  error
}

type MockKafkaClient struct{}

type MockMinioClient struct{}

func (m *MockKafkaClient) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return 0, 0, nil
}

func (m *MockKafkaClient) SendMessages(msgs []*sarama.ProducerMessage) error {
	return nil
}

func (m *MockKafkaClient) Close() error {
	return nil
}

func (m *MockKafkaClient) TxnStatus() sarama.ProducerTxnStatusFlag {
	return sarama.ProducerTxnFlagReady
}

func (m *MockKafkaClient) IsTransactional() bool {
	return false
}

func (m *MockKafkaClient) BeginTxn() error {
	return nil
}

func (m *MockKafkaClient) CommitTxn() error {
	return nil
}

func (m *MockKafkaClient) AbortTxn() error {
	return nil
}

func (m *MockKafkaClient) AddOffsetsToTxn(offsets map[string][]*sarama.PartitionOffsetMetadata, groupId string) error {
	return nil
}

func (m *MockKafkaClient) AddMessageToTxn(msg *sarama.ConsumerMessage, groupId string, metadata *string) error {
	return nil
}

func (m *MockMinioClient) PresignedGetObject(ctx context.Context, bucketName, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error) {
	return &url.URL{Scheme: "https", Host: "abebe.com", Path: "/mock"}, nil
}

func (m *MockMinioClient) PresignedPostPolicy(ctx context.Context, p *minio.PostPolicy) (u *url.URL, formData map[string]string, err error) {
	return &url.URL{Scheme: "https", Host: "abebe.com", Path: "/mock"}, map[string]string{}, nil
}

var logger = platform.InitZapLogger()

func (m *MockAuthRepository) Register(ctx context.Context, user *model.User) (*model.User, error) {
	return m.User, nil
}

func (m *MockAuthRepository) Login(ctx context.Context, userName, password string) (*model.User, error) {
	return nil, nil
}

func (m *MockAuthRepository) Logout(ctx context.Context, tokenID string) error {
	return nil
}

func (m *MockAuthRepository) StoreRefreshTokens(userID int, refreshToken string, expiresAt time.Time, clientType string) (sql.Result, error) {
	return nil, nil
}

func (m *MockAuthRepository) RevokeRefreshToken(refreshToken string) error {
	return nil
}

func (m *MockAuthRepository) GetRefreshToken(refreshToken string) (*model.RefreshToken, error) {
	return nil, nil
}

func (m *MockAuthRepository) GetUserByID(ID int) (*model.User, error) {
	return nil, nil
}

func (m *MockAuthRepository) SearchUser(ctx context.Context, username, cursor string, limit int) (*dto.PaginatedResponse, error) {
	return nil, nil
}

func (m *MockAuthRepository) UpdateFailedLoginAttempts(ctx context.Context, user *model.User) (*model.User, error) {
	return nil, nil
}

func (m *MockAuthRepository) TemporaryLockUntil(ctx context.Context, user *model.User) (*model.User, error) {
	return nil, nil
}

func (m *MockAuthRepository) IncreaseFollowCount(ctx context.Context, followerID, followingID int) error {
	return nil
}

func (m *MockAuthRepository) DecreaseFollowCount(ctx context.Context, followerID, followingID int) error {
	return nil
}

func (m *MockAuthRepository) GetUserProfile(ctx context.Context, userID int64) (*model.User, error) {
	return nil, nil
}

func (m *MockAuthRepository) UpdateUserAvatar(ctx context.Context, avatar string, userID int64) (*model.User, error) {
	return nil, nil
}

func TestAuthService_Register_Success(t *testing.T) {
	mr := &MockAuthRepository{
		User: &model.User{
			Name:     "hellow",
			Username: "helloww",
			Email:    "email",
			Password: "pass@344#P",
		},

	}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-client-type", "web"))

	sv := service.NewAuthService(mr, nil, &MockKafkaClient{}, logger, &MockMinioClient{})

	user := &model.User{
		Name:     "hellow",
		Username: "helloww",
		Email:    "email",
		Password: "pass@344#P",
	}

	cu, _, err := sv.Register(ctx, user)

	if err != nil {

		t.Fatalf("unexpected error: %v", err)
	}

	if cu != mr.User {

		t.Fatalf("the returned user doesn't match")
	}

}
