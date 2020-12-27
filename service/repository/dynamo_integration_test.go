// +build integration

package repository_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TaylorOno/bookmarker/cmd/config"
	"github.com/TaylorOno/bookmarker/service/repository"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	// The canonical name of the image on Docker Hub
	_dockerImage = "amazon/dynamodb-local"
)

var _ = Describe("Dynamo Integration", func() {
	var (
		dynamoRepo *repository.Dynamo
		closer     func()
		err        error
	)

	BeforeSuite(func() {
		port := randomPort()
		closer, err = startDynamoContainer(port)
		if err != nil {
			Fail(err.Error())
		}

		session, err := config.NewAWSSessions("id", "secret", "us-west-2", fmt.Sprintf("http://localhost:%v", port))
		if err != nil {
			Fail(err.Error())
		}

		client := config.NewDynamoClient(session)
		dynamoRepo = repository.NewDynamoRepository(client, "testBookmark")
		loadTestData(dynamoRepo)
	})

	AfterSuite(func() {
		closer()
	})

	Context("CreateBookmark", func() {
		It("Calls dynamoRepo save with user bookmarks", func() {
			UserBookmark := testBookmark("Malcolm", "War and Peace", "IN_PROGRESS", 7)
			result, err := dynamoRepo.CreateBookmark(context.Background(), UserBookmark)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.UserId).To(Equal("Malcolm"))
		})
	})

	Context("DeleteBookmark", func() {
		It("Calls dynamoRepo delete with user and book", func() {
			UserBookmark := testBookmark("Malcolm", "Crime and Punishment", "IN_PROGRESS", 7)
			dynamoRepo.CreateBookmark(context.Background(), UserBookmark)
			dynamoRepo.DeleteBookmark(context.Background(), "Malcolm", "Crime and Punishment")
			_, err := dynamoRepo.GetBookmark(context.Background(), "Malcolm", "Crime and Punishment")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("not found"))
		})
	})

	Context("GetBookmark", func() {
		It("Calls dynamoRepo get with user and book", func() {
			UserBookmark := testBookmark("Malcolm", "Pride and Prejudice", "IN_PROGRESS", 38)
			dynamoRepo.CreateBookmark(context.Background(), UserBookmark)
			result, err := dynamoRepo.GetBookmark(context.Background(), "Malcolm", "Pride and Prejudice")
			Expect(err).ToNot(HaveOccurred())
			Expect(result.Page).To(Equal(38))
		})
	})

	Context("GetBookmarks", func() {
		It("Calls dynamoRepo query with user and filter", func() {
			result, err := dynamoRepo.GetBookmarks(context.Background(), "Malcolm", "IN_PROGRESS", 30)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(result)).To(Equal(30))
		})
	})
})

func loadTestData(dynamo *repository.Dynamo) {
	bookmarks := bookmarksFromFile()
	for _, bookmark := range bookmarks {
		_, err := dynamo.CreateBookmark(context.Background(), bookmark)
		if err != nil {
			Fail("failed to load test data")
		}
	}
}

func randomPort() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(8999) + 1000
}

func testBookmark(name, book, status string, page int) repository.UserBookmark {
	return repository.UserBookmark{
		UserId:      name,
		LastUpdated: time.Now().UTC().Format("2006-01-02T15:04:05Z07:00.000"),
		Book:        book,
		Status:      status,
		Page:        page,
	}
}

func bookmarksFromFile() []repository.UserBookmark {
	var bookmarks []repository.UserBookmark
	file, err := os.Open("test_data/data.csv")
	if err != nil {
		Fail(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bookmark := strings.Split(scanner.Text(), ",")
		page, _ := strconv.Atoi(bookmark[2])
		bookmarks = append(bookmarks, repository.UserBookmark{
			UserId:      bookmark[0],
			LastUpdated: time.Now().UTC().Format("2006-01-02T15:04:05Z07:00.000"),
			Book:        bookmark[1],
			Page:        page,
			Status:      bookmark[3],
		})
	}

	if err = scanner.Err(); err != nil {
		Fail(err.Error())
	}

	return bookmarks
}

// Create a DynamoDB Docker container mapped to the specified TCP port.
func startDynamoContainer(port int) (closer func(), err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()

	dockerClient := getDynamoClient()

	pullLatestImage(ctx, dockerClient)

	// Create Container
	containerCfg := container.Config{Image: _dockerImage}
	hostCfg := container.HostConfig{
		AutoRemove: true,
		PortBindings: nat.PortMap{
			"8000/tcp": []nat.PortBinding{
				{HostPort: fmt.Sprintf("%d/tcp", port)},
			},
		},
	}

	resp := createContainer(dockerClient, containerCfg, hostCfg)

	closeContainer := createCloser(dockerClient, resp)
	err = startContainer(dockerClient, resp)

	return closeContainer, nil
}

func createCloser(dockerClient *client.Client, resp container.ContainerCreateCreatedBody) func() {
	return func() {
		timeout := 5 * time.Second
		err := dockerClient.ContainerStop(context.Background(), resp.ID, &timeout)
		if err != nil {
			fmt.Printf("failed to stop container: %s", err.Error())
		}

		options := types.ContainerRemoveOptions{RemoveVolumes: true, RemoveLinks: true, Force: true}
		if err := dockerClient.ContainerRemove(context.Background(), resp.ID, options); err != nil {
			fmt.Printf("failed to remove container: %s", err.Error())
		}
	}
}

func createContainer(dockerClient *client.Client, containerCfg container.Config, hostCfg container.HostConfig) container.ContainerCreateCreatedBody {
	containerName := fmt.Sprintf("test_dynamodb_%d", time.Now().Unix())
	resp, err := dockerClient.ContainerCreate(context.Background(), &containerCfg, &hostCfg, nil, nil, containerName)
	if err != nil {
		Fail(err.Error())
	}

	return resp
}

func startContainer(dockerClient *client.Client, resp container.ContainerCreateCreatedBody) error {
	// Spin up the container using the defined configurations.
	err := dockerClient.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		Fail(err.Error())
	}

	time.Sleep(1 * time.Second)

	return err
}

func pullLatestImage(ctx context.Context, dockerClient *client.Client) {
	reader, err := dockerClient.ImagePull(ctx, _dockerImage, types.ImagePullOptions{})
	if err != nil {
		Fail(err.Error())
	}
	io.Copy(os.Stdout, reader)
}

func getDynamoClient() *client.Client {
	// Get docker Client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		Fail(err.Error())
	}
	return dockerClient
}
