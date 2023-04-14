package access

import (
	"testing"
	"time"
)

const (
	privateKeyPath = "../../tests/certs/test_rsa"
	publicKeyPath  = "../../tests/certs/test_rsa.pub"

	testUserId  = "1"
	testEmail   = "test@test.com"
	testRole    = "moderator"
	testPremium = true
	testTtl     = 30 * time.Minute
)

var testNickname = "test"

func TestJwtGeneration(t *testing.T) {
	input := Payload{
		UserId:   testUserId,
		Email:    testEmail,
		Nickname: &testNickname,
		Role:     testRole,
		Premium:  testPremium,
	}

	producer, err := NewProducer(privateKeyPath)
	if err != nil {
		t.Fatal(err)
	}

	parser, err := NewParser(publicKeyPath)
	if err != nil {
		t.Fatal(err)
	}

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
		input.Premium != output.Premium {
		t.Fail()
	}
}
