package fixtures

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/planetary-social/scuttlego-pub/service/domain"
	"github.com/planetary-social/scuttlego/service/domain/feeds/message"
	"github.com/planetary-social/scuttlego/service/domain/identity"
	"github.com/planetary-social/scuttlego/service/domain/refs"
)

func SomePublicIdentity() identity.Public {
	return SomePrivateIdentity().Public()
}

func SomePrivateIdentity() identity.Private {
	v, err := identity.NewPrivate()
	if err != nil {
		panic(err)
	}
	return v
}

func SomePositiveInt() int {
	return 1 + rand.Intn(math.MaxInt-1)
}

func SomeTime() time.Time {
	return time.Unix(rand.Int63(), 0)
}

func SomeSecretKeySeed() domain.SecretKeySeed {
	v, err := domain.NewSecretKeySeed()
	if err != nil {
		panic(err)
	}
	return v
}

func SomeRefFeed() refs.Feed {
	return refs.MustNewFeed(fmt.Sprintf("@%s.ed25519", randomBase64(32)))
}

func SomeRefMessage() refs.Message {
	return refs.MustNewMessage(fmt.Sprintf("%%%s.sha256", randomBase64(32)))
}

func SomeRefIdentity() refs.Identity {
	return refs.MustNewIdentity(fmt.Sprintf("@%s.ed25519", randomBase64(32)))
}

func SomeMessage() message.Message {
	return SomeMessageWithFeedSequence(SomeRefFeed(), SomeSequence())
}

func SomeDuration() time.Duration {
	return time.Duration(somePositiveInt()) * time.Second
}

func SomeMessageWithFeedSequence(feed refs.Feed, sequence message.Sequence) message.Message {
	var previous *refs.Message
	if !sequence.IsFirst() {
		tmp := SomeRefMessage()
		previous = &tmp
	}

	return message.MustNewMessage(
		SomeRefMessage(),
		previous,
		sequence,
		SomeRefIdentity(),
		feed,
		SomeTime(),
		SomeContent(),
		SomeRawMessage(),
	)
}

func SomeSequence() message.Sequence {
	return message.MustNewSequence(somePositiveInt())
}

func SomeContent() message.Content {
	return message.MustNewContent(SomeRawContent(), nil, nil)
}

func SomeRawMessage() message.RawMessage {
	return message.MustNewRawMessage(someBytes())
}

func SomeRawContent() message.RawContent {
	return message.MustNewRawContent(someBytes())
}

func Directory(t testing.TB) string {
	name, err := os.MkdirTemp("", "scuttlego-test")
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		err := os.RemoveAll(name)
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(cleanup)

	return name
}

func somePositiveInt() int {
	return 1 + rand.Intn(math.MaxInt-1)
}

func randomBase64(bytes int) string {
	return base64.StdEncoding.EncodeToString(SomeBytesOfLength(bytes))
}

func someBytes() []byte {
	return SomeBytesOfLength(10 + rand.Intn(100))
}

func SomeBytesOfLength(n int) []byte {
	r := make([]byte, n)
	if _, err := cryptorand.Read(r); err != nil {
		panic(err)
	}
	return r
}
