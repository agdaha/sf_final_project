package postcomment

import (
	"log/slog"
	"net/http"
	"reflect"
	"testing"

	censorservice "github.com/agdaha/sf_final_project/api_gateway/internal/clients/censor_service"
	commentservice "github.com/agdaha/sf_final_project/api_gateway/internal/clients/comment_service"
)

func TestNew(t *testing.T) {
	type args struct {
		commentService commentservice.CommentService
		censorService  censorservice.CensorService
		log            *slog.Logger
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.commentService, tt.args.censorService, tt.args.log); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
