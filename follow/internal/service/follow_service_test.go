package service_test
import (
	"context"
	"testing"

	"github.com/IBM/sarama"
	dto "github.com/abelmalu/golang-posts/follow/internal/dtos"
	ierrors "github.com/abelmalu/golang-posts/follow/internal/errors"
	"github.com/abelmalu/golang-posts/follow/internal/service"
)

type MockFollowRepository struct {
	response string
	err      error
}

type MockKafkaClient struct {
	partition int32
	offSet    int32
	err       error
}

var (
	kafkaErr =  ierrors.NewInternalError(ierrors.ErrorMessage("Kafka Sending Error"), nil)
	responseFollowed = "unfollowed successfully"
	responseUnfollowed = "followed successfully"
	repoError = ierrors.NewDatabaseError(ierrors.MSGDatabaseError, nil)
)

func (m *MockFollowRepository) ToggleFollow(ctx context.Context, state bool, followerID, followingID int) (string, error) {
	return m.response, m.err
}

func (m *MockFollowRepository) GetUserFollowers(ctx context.Context, followingID, limit int, cursor string) (*dto.PaginatedFollowersResponse, error) {
	return nil, nil
}

func (m *MockFollowRepository) GetUserUserFollowings(ctx context.Context, followerId, limit int, cursor string) (*dto.PaginatedFollowingsResponse, error) {
	return nil, nil
}

func (m *MockFollowRepository) CreateCacheUser(ctx context.Context, userID int, username, name string) error {
	return nil
}

func (mk *MockKafkaClient) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {

	return mk.partition, int64(mk.offSet), mk.err
}

func TestFollowService_ToggleFollow(t *testing.T) {

	tests := []struct {
		name        string
		follow      bool
		followerID  int
		followingID int
		mockRepo    *MockFollowRepository
		mockKafka   *MockKafkaClient
		checkErr    func (t *testing.T,err error)
		checkResp   func (t *testing.T,resp string)
		wantErr     bool
	}{
		{
			name:        "success_follow",
			follow:      true,
			followerID:  1,
			followingID: 2,
			mockRepo: &MockFollowRepository{

				response: responseFollowed,
			},
			mockKafka: &MockKafkaClient{

				partition: 1,
				offSet:    2,
			},

			checkResp: func(t *testing.T, resp string) {

				if resp != responseFollowed {

					t.Fatalf("unexpected response %v",resp)
				}
			},
		},

		{
			name:        "success_unfollow",
			follow:      true,
			followerID:  1,
			followingID: 2,
			mockRepo: &MockFollowRepository{

				response: responseUnfollowed,
			},
			mockKafka: &MockKafkaClient{

				partition: 1,
				offSet:    2,
			},

			checkResp: func(t *testing.T, resp string) {

				if resp != responseUnfollowed {

					t.Fatalf("unexpected response %v",resp)
				}
			},
		},
		{
			name:        "kafka error",
			follow:      true,
			followerID:  1,
			followingID: 2,
			wantErr: true,
			mockRepo: &MockFollowRepository{

				response: "followed successfully",
			},
			mockKafka: &MockKafkaClient{

				err:  kafkaErr,
			},

			checkErr: func(t *testing.T, err error) {

				if err.Error() != kafkaErr.Error(){

					t.Fatalf("Unexpected error %v",err)
				}

			},
		},

		{
			name:        "repo error",
			follow:      true,
			followerID:  1,
			followingID: 2,
			wantErr: true,
			mockRepo: &MockFollowRepository{

				err: repoError,
			},
			mockKafka: &MockKafkaClient{

				partition: 1,
				offSet:    2,
			},

			checkErr: func(t *testing.T, err error) {

				if err.Error() != repoError.Error(){

					t.Fatalf("Unexpected error %v",err)
				}

			},
		},
	}

	for _, tt := range tests {

	
		t.Run(tt.name, func(t *testing.T) {

			sv := service.NewFollowService(tt.mockRepo, tt.mockKafka)

			resp, err := sv.ToggleFollow(t.Context(), tt.follow, tt.followerID, tt.followingID)

			if tt.wantErr {

				tt.checkErr(t,err)
				return
			}

			tt.checkResp(t,resp.Message)

		})
	}
}
