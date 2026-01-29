package main


import (
    "context"
    "log"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    pb "github.com/abelmalu/golang-posts/post/proto/pb"
)

func main() {
    r := gin.Default()
	

    // Connect to gRPC server once
    conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("failed to connect to gRPC server: %v", err)
    }
    defer conn.Close()

    grpcClient := pb.NewPostServiceClient(conn)

    // HTTP POST /posts -> gRPC call
    r.POST("/posts", func(c *gin.Context) {
        // Extract userID from JWT or context (example: from header)
     // Replace with JWT extraction logic

        var reqBody struct {
            Title   string `json:"title" binding:"required"`
            Content string `json:"content" binding:"required"`
        }

        if err := c.ShouldBindJSON(&reqBody); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "status":  "error",
                "message": "Invalid request body",
            })
            return
        }

        // Build gRPC request
        grpcReq := &pb.CreatePostRequest{
            Title:   reqBody.Title,
            Content: reqBody.Content,
        }

        // Add timeout
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        // Call gRPC server
        grpcResp, err := grpcClient.CreatePost(ctx, grpcReq)
        if err != nil {
            log.Printf("gRPC call failed: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "status":  "error",
                "message": "Internal server error",
            })
            return
        }

        // Send response back to frontend
        c.JSON(http.StatusOK, gin.H{
            "status":  grpcResp.Status,
            "message": grpcResp.Message,
        })
    })

    r.Run(":8080")
}
