package repository

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDynamoPostRepository_CreatesTable(t *testing.T) {
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}

	ctx := context.Background()
	repo, err := NewDynamoPostRepository(ctx, endpoint)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	n := postsTableName
	// Verifica se a tabela foi criada
	out, err := repo.client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: &n,
	})
	assert.NoError(t, err)
	assert.Equal(t, postsTableName, *out.Table.TableName)
}

func TestCreatePost(t *testing.T) {
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}

	ctx := context.Background()
	repo, err := NewDynamoPostRepository(ctx, endpoint)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	require.NoError(t, cleanupTable(t.Context(), repo))

	userID := "test-user"
	content := "Hello, DynamoDB!"

	postID, err := repo.CreatePost(ctx, userID, content)
	assert.NoError(t, err)
	assert.NotEmpty(t, postID)
}

func TestListPosts(t *testing.T) {
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}

	ctx := context.Background()
	repo, err := NewDynamoPostRepository(ctx, endpoint)
	assert.NoError(t, err)
	require.NoError(t, cleanupTable(t.Context(), repo))

	// Cria alguns posts
	userID := "list-user"
	for i := 0; i < 3; i++ {
		content := "Post " + strconv.Itoa(i)
		_, err := repo.CreatePost(ctx, userID, content)
		assert.NoError(t, err)
	}

	// Lista os posts
	posts, nextCursor, err := repo.ListPosts(ctx, userID, 10, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, posts)
	assert.Len(t, posts, 3)
	assert.Empty(t, nextCursor) // Não deve ter mais páginas nesse caso

	// Verifica conteúdo
	for _, post := range posts {
		assert.Equal(t, userID, post.UserID)
		assert.NotEmpty(t, post.Content)
	}
}

func cleanupTable(ctx context.Context, repo *DynamoPostRepository) error {
	// Faz um scan pra pegar todos os IDs
	out, err := repo.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: awsString(postsTableName),
		ProjectionExpression: awsString("id"),
	})
	if err != nil {
		return err
	}

	// Deleta cada item pelo ID
	for _, item := range out.Items {
		var p Post
		err := attributevalue.UnmarshalMap(item, &p)
		if err != nil {
			return err
		}
		_, err = repo.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: awsString(postsTableName),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: p.ID},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
