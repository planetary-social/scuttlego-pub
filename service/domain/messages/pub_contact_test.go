package known

import (
	"testing"

	"github.com/planetary-social/scuttlego/fixtures"
	"github.com/planetary-social/scuttlego/service/domain/feeds"
	"github.com/stretchr/testify/require"
)

func TestPubFollow_ImplementsContactMessage(t *testing.T) {
	require.Implements(t, new(feeds.ContactMessage), MustNewPubFollow(fixtures.SomeRefIdentity()))
}
