package repository

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

const postsTableName = "posts"

type Post struct {
	ID        string `dynamodbav:"id"`
	UserID    string `dynamodbav:"user_id"`
	Content   string `dynamodbav:"content"`
	CreatedAt string `dynamodbav:"created_at"`
}

type DynamoPostRepository struct {
	client *dynamodb.Client
}

func NewDynamoPostRepository(ctx context.Context, endpoint string) (*DynamoPostRepository, error) {
	cfg, err := loadDynamoConfig(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)
	repo := &DynamoPostRepository{client: client}

	if err := repo.ensureTable(ctx); err != nil {
		return nil, err
	}

	return repo, nil
}

func loadDynamoConfig(ctx context.Context, endpoint string) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"), // Região dummy para o local
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint, HostnameImmutable: true}, nil
			}),
		),
	)
}

func (r *DynamoPostRepository) ensureTable(ctx context.Context) error {
	_, err := r.client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: awsString(postsTableName),
	})
	if err == nil {
		return nil
	}

	var nfe *types.ResourceNotFoundException
	if !isResourceNotFound(err, &nfe) {
		return err
	}

	_, err = r.client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: awsString(postsTableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: awsString("id"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: awsString("id"), KeyType: types.KeyTypeHash},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	return err
}

func (r *DynamoPostRepository) CreatePost(ctx context.Context, userID, content string) (string, error) {
	post := Post{
		ID:        uuid.New().String(),
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	item, err := attributevalue.MarshalMap(post)
	if err != nil {
		return "", err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: awsString(postsTableName),
		Item:      item,
	})
	if err != nil {
		return "", err
	}

	return post.ID, nil
}

func (r *DynamoPostRepository) ListPosts(ctx context.Context, userID string, limit int32, cursor string) ([]Post, string, error) {
	input := &dynamodb.ScanInput{
		TableName:        awsString(postsTableName),
		Limit:            awsInt32(limit),
		ExclusiveStartKey: nil,
	}

	if cursor != "" {
		input.ExclusiveStartKey = map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: cursor},
		}
	}

	out, err := r.client.Scan(ctx, input)
	if err != nil {
		return nil, "", err
	}

	var posts []Post
	err = attributevalue.UnmarshalListOfMaps(out.Items, &posts)
	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if out.LastEvaluatedKey != nil {
		if val, ok := out.LastEvaluatedKey["id"].(*types.AttributeValueMemberS); ok {
			nextCursor = val.Value
		}
	}

	return posts, nextCursor, nil
}

func awsString(v string) *string { return &v }
func awsInt32(v int32) *int32 { return &v }
func isResourceNotFound(err error, target **types.ResourceNotFoundException) bool {
	return errors.As(err, target)
}
