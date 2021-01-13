package rtc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/livekit/livekit-server/pkg/rtc/types"
	"github.com/livekit/livekit-server/pkg/rtc/types/typesfakes"
	"github.com/livekit/livekit-server/proto/livekit"
)

func TestIsReady(t *testing.T) {
	tests := []struct {
		state livekit.ParticipantInfo_State
		ready bool
	}{
		{
			state: livekit.ParticipantInfo_JOINING,
			ready: false,
		},
		{
			state: livekit.ParticipantInfo_JOINED,
			ready: true,
		},
		{
			state: livekit.ParticipantInfo_ACTIVE,
			ready: true,
		},
		{
			state: livekit.ParticipantInfo_DISCONNECTED,
			ready: false,
		},
	}

	for _, test := range tests {
		t.Run(test.state.String(), func(t *testing.T) {
			p := &ParticipantImpl{
				state: test.state,
			}
			assert.Equal(t, test.ready, p.IsReady())
		})
	}
}

func TestTrackPublishEvents(t *testing.T) {
	p := newParticipantForTest("test")
	track := &typesfakes.FakePublishedTrack{}
	track.IDReturns("id")
	published := false
	updated := false
	p.OnTrackUpdated(func(p types.Participant, track types.PublishedTrack) {
		updated = true
	})
	p.OnTrackPublished(func(p types.Participant, track types.PublishedTrack) {
		published = true
	})
	p.handleTrackPublished(track)

	assert.True(t, published)
	assert.False(t, updated)
	assert.Len(t, p.publishedTracks, 1)

	track.OnCloseArgsForCall(0)()
	assert.Len(t, p.publishedTracks, 0)
	assert.True(t, updated)
}

func newParticipantForTest(name string) *ParticipantImpl {
	p, _ := NewParticipant(&typesfakes.FakePeerConnection{}, &typesfakes.FakeSignalConnection{}, name, ReceiverConfig{})
	return p
}
