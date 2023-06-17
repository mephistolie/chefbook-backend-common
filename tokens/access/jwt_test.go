package access

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

const (
	privateKeyPath = "../../tests/certs/test_rsa"

	testEmail            = "test@test.com"
	testRole             = "moderator"
	testSubscriptionPlan = "test_plan"
	testTtl              = 30 * time.Minute
)

var (
	testUserId   = uuid.New()
	testNickname = "test"
)

func TestJwtGeneration(t *testing.T) {
	input := Payload{
		UserId:           testUserId,
		Email:            testEmail,
		Nickname:         &testNickname,
		Role:             testRole,
		SubscriptionPlan: testSubscriptionPlan,
		Deleted:          true,
	}

	producer, publicKey, err := NewProducer(privateKeyPath)
	if err != nil {
		t.Fatal(err)
	}

	parser := NewParserByKey(publicKey)

	token, err := producer.Produce(input, testTtl)
	if err != nil {
		t.Fatal(err)
	}

	output, err := parser.Parse(token)
	if err != nil {
		t.Fatal(err)
	}

	if input.UserId != output.UserId ||
		input.Email != output.Email ||
		*input.Nickname != *output.Nickname ||
		input.Role != output.Role ||
		input.SubscriptionPlan != output.SubscriptionPlan {
		t.Fail()
	}
}
