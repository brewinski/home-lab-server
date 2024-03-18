package vq_test

import (
	"testing"

	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/vq"
)

func TestDetectLadderChanges(t *testing.T) {
	type args struct {
		old vq.GetLadderResponseBody
		new vq.GetLadderResponseBody
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "detect changes when the ladder changes positions change",
			args: args{
				old: vq.GetLadderResponseBody{
					Records: []vq.LadderRecord{
						{
							ID: "1",
							Fields: vq.LadderFields{
								Rank:           "1",
								TeamNameLookup: "Aces",
							},
						},
						{
							ID: "2",
							Fields: vq.LadderFields{
								Rank:           "2",
								TeamNameLookup: "APG",
							},
						},
					},
					Offset: "",
				},
				new: vq.GetLadderResponseBody{
					Records: []vq.LadderRecord{
						{
							ID: "1",
							Fields: vq.LadderFields{
								Rank:           "1",
								TeamNameLookup: "APG",
							},
						},
						{
							ID: "2",
							Fields: vq.LadderFields{
								Rank:           "2",
								TeamNameLookup: "Aces",
							},
						},
					},
					Offset: "",
				},
			},
			want: true,
		},
		{
			name: "detest changes when team points change",
			args: args{
				old: vq.GetLadderResponseBody{
					Records: []vq.LadderRecord{
						{
							ID: "1",
							Fields: vq.LadderFields{
								Rank:              "1",
								TeamNameLookup:    "Aces",
								CompetitionPoints: "10",
							},
						},
						{
							ID: "2",
							Fields: vq.LadderFields{
								Rank:              "2",
								TeamNameLookup:    "APG",
								CompetitionPoints: "9",
							},
						},
					},
					Offset: "",
				},
				new: vq.GetLadderResponseBody{
					Records: []vq.LadderRecord{
						{
							ID: "1",
							Fields: vq.LadderFields{
								Rank:              "1",
								TeamNameLookup:    "Aces",
								CompetitionPoints: "10",
							},
						},
						{
							ID: "2",
							Fields: vq.LadderFields{
								Rank:              "2",
								TeamNameLookup:    "APG",
								CompetitionPoints: "13",
							},
						},
					},
					Offset: "",
				},
			},
			want: true,
		},
		{
			name: "ignore changes to anything that doesn't impact the ladder ordering.",
			args: args{
				old: vq.GetLadderResponseBody{
					Records: []vq.LadderRecord{
						{
							ID: "1",
							Fields: vq.LadderFields{
								Rank:           "1",
								TeamNameLookup: "Aces",
							},
						},
						{
							ID: "2",
							Fields: vq.LadderFields{
								Rank:           "2",
								TeamNameLookup: "APG",
								TotalSetsA:     "something",
							},
						},
					},
					Offset: "",
				},
				new: vq.GetLadderResponseBody{
					Records: []vq.LadderRecord{
						{
							ID: "1",
							Fields: vq.LadderFields{
								Rank:           "1",
								TeamNameLookup: "Aces",
							},
						},
						{
							ID: "2",
							Fields: vq.LadderFields{
								Rank:           "2",
								TeamNameLookup: "APG",
								TotalSetsA:     "something else",
							},
						},
					},
					Offset: "",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vq.DetectLadderChanges(tt.args.old, tt.args.new); got != tt.want {
				t.Errorf("DetectLadderChanges() = %v, want %v", got, tt.want)
			}
		})
	}
}
