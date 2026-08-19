package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/IBM/sarama"
	dto "github.com/abelmalu/golang-posts/like/internal/dtos"
	ierrors "github.com/abelmalu/golang-posts/like/internal/errors"
	"github.com/abelmalu/golang-posts/like/internal/service"
)

type MockLikeRepository struct {
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
	repoErr = ierrors.NewDatabaseError(ierrors.MSGDatabaseError,nil)
	responseLiked = "post liked successfully"
	responseUnliked = "post unliked successfully"
)


func (m *MockLikeRepository) ToggleLike(ctx context.Context, state bool, userID, postID int) (string, error) {
	return m.response, m.err
}

func (m *MockLikeRepository) CreateCacheUser(ctx context.Context, userID int, username, name string) error {
	return nil
}

func (m *MockLikeRepository) CreateCachePost(ctx context.Context, postID int, title string) error {
	return nil
}

func (m *MockLikeRepository) GetPostLikes(ctx context.Context, postID, limit int, cursor string) (*dto.PaginatedPostLikesResponse, error) {
	return nil, nil
}

func (mk *MockKafkaClient) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {

	return mk.partition, int64(mk.offSet), mk.err
}

func TestFollowService_ToggleFollow(t *testing.T) {

	tests := []struct {
		name        string
		state      bool
		userID  int
		postID int
		mockRepo    *MockLikeRepository
		mockKafka   *MockKafkaClient
		checkErr    func (t *testing.T,err error)
		checkResp   func (t *testing.T,resp string)
		wantErr     bool
	}{
		{
			name:        "success_like",
			state:      true,
			userID:  1,
			postID: 2,
			mockRepo: &MockLikeRepository{

				response: responseLiked,
			},
			mockKafka: &MockKafkaClient{

				partition: 1,
				offSet:    2,
			},

			checkResp: func(t *testing.T, resp string) {

				if resp != responseLiked {
					fmt.Println(resp)

					t.Fatalf("unexpected response %v",resp)
				}
			},
		},

		{
			name:        "success_unlike",
			state:      true,
			userID:  1,
			postID: 2,
			mockRepo: &MockLikeRepository{

				response: responseUnliked,
			},
			mockKafka: &MockKafkaClient{

				partition: 1,
				offSet:    2,
			},

			checkResp: func(t *testing.T, resp string) {

				if resp != responseUnliked {

					t.Fatalf("unexpected response %v",resp)
				}
			},
		},
		{
			name:        "kafka error",
			state:      true,
			userID:  1,
			postID: 2,
			wantErr: true,
			mockRepo: &MockLikeRepository{

				response: "liked successfully",
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
			name:        "kafka error",
			state:      true,
			userID:  1,
			postID: 2,
			wantErr: true,
			mockRepo: &MockLikeRepository{

				err: repoErr,
			},
			mockKafka: &MockKafkaClient{

				partition: 1,
				offSet: 0,
			},

			checkErr: func(t *testing.T, err error) {

				if err.Error() != repoErr.Error(){

					t.Fatalf("Unexpected error %v",err)
				}

			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			sv := service.NewLikeService(tt.mockRepo, tt.mockKafka)

			resp, err := sv.ToggleLike(t.Context(), tt.state, tt.userID, tt.postID)

			if tt.wantErr {

				tt.checkErr(t,err)
				return
			}

			tt.checkResp(t,resp.Message)

		})
	}
}
