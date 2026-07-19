package projectdata

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddVideo(t *testing.T) {
	// GIVEN a new ProjectData and ProtoVideo
	// WHEN passing the ProtoVideo to ProjectData.AddVideo
	// THEN ProjectData.videos will contain the video

	pd := ProjectData{}
	pv := ProtoVideo{
		Name:        "Cheesy pick-up lines",
		Description: "A collection of cheesy pick-up lines from random places",
		Path:        "/home/bob/Videos/cheesy-pickups.mp4",
	}

	pd.AddVideo(pv)

	require.Equal(t, 1, len(pd.videos), "Wrong projectdata.videos length")
	require.Equal(t, pv.Name, pd.videos[0].name, "projectdata.videos contains wrong video (name)")
	require.Equal(t, pv.Description, pd.videos[0].desc, "projectdata.videos contains wrong video (description)")
	require.Equal(t, pv.Path, pd.videos[0].path, "projectdata.videos contains wrong video (paths)")
}

func TestDeleteVideo(t *testing.T) {
	// GIVEN a ProjectData with a video
	// WHEN calling DeleteVideo with the ID of an existing video
	// THEN deletes the video

	v := video{
		id:   "abc123",
		name: "Cheesy pick-up lines",
		desc: "A collection of cheesy pick-up lines from random places",
		path: "/home/bob/Videos/cheesy-pickups.mp4",
	}

	pd := ProjectData{
		videos: []video{v},
	}

	require.Equal(t, 1, len(pd.videos), "Test setup failed")
	err := pd.DeleteVideo(v.id)
	require.NoError(t, err, "Error trying to delete video")
	require.Equal(t, 0, len(pd.videos), "Video not deleted")
}
