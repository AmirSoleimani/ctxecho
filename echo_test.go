package ctxecho_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/AmirSoleimani/ctxecho"
)

type (
	ctxStringKey string
	ctxIntKey    int
	ctxAnyKey    any
)

type testStruct struct {
	Username string
}

func TestInspect(t *testing.T) {

	testCases := []struct {
		Input func() context.Context
		Want  map[any]any
	}{
		{
			Input: func() context.Context {
				return context.Background()
			},
			Want: map[any]any{},
		},
		{
			Input: func() context.Context {
				return context.WithValue(context.Background(), ctxStringKey("mykey"), "myvalue")
			},
			Want: map[any]any{
				ctxStringKey("mykey"): "myvalue",
			},
		},
		{
			Input: func() context.Context {
				ctx := context.WithValue(context.Background(), ctxStringKey("mykey"), "myvalue")
				ctx = context.WithValue(ctx, ctxIntKey(126), 1000)
				return ctx
			},
			Want: map[any]any{
				ctxStringKey("mykey"): "myvalue",
				ctxIntKey(126):        1000,
			},
		},
		{
			Input: func() context.Context {
				ctx := context.WithValue(context.Background(), ctxStringKey("mykey"), "myvalue")
				ctx = context.WithValue(ctx, ctxIntKey(126), "1000x")
				ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
				defer cancel()
				return ctx
			},
			Want: map[any]any{
				ctxStringKey("mykey"): "myvalue",
				ctxIntKey(126):        "1000x",
			},
		},
		{
			Input: func() context.Context {
				ctx := context.WithValue(context.Background(), ctxStringKey("mykey"), "myvalue")
				ctx, cancel := context.WithDeadline(ctx, time.Now().Add(1*time.Hour))
				defer cancel()
				return ctx
			},
			Want: map[any]any{
				ctxStringKey("mykey"): "myvalue",
			},
		},
		{
			Input: func() context.Context {
				ctx := context.WithValue(context.Background(), ctxStringKey("mykey"), "myvalue")
				ctx, _ = context.WithDeadline(ctx, time.Now().Add(1*time.Hour))
				ctx, _ = context.WithTimeout(ctx, time.Hour)
				return ctx
			},
			Want: map[any]any{
				ctxStringKey("mykey"): "myvalue",
			},
		},
		{
			Input: func() context.Context {
				ctx := context.WithValue(context.Background(), ctxStringKey("mykey"), "myvalue")
				ctx, _ = context.WithDeadline(ctx, time.Now().Add(1*time.Hour))
				ctx, _ = context.WithTimeout(ctx, time.Hour)
				ctx = context.WithValue(ctx, ctxAnyKey("XXX"), 123)
				return ctx
			},
			Want: map[any]any{
				ctxStringKey("mykey"): "myvalue",
				ctxAnyKey("XXX"):      123,
			},
		},
		{
			Input: func() context.Context {
				shadow := testStruct{
					Username: "amirso",
				}
				return context.WithValue(context.Background(), ctxStringKey("mykey"), shadow)
			},
			Want: map[any]any{
				ctxStringKey("mykey"): testStruct{
					Username: "amirso",
				},
			},
		},
	}

	for _, tc := range testCases {
		ctx := tc.Input()
		got := ctxecho.Inspect(ctx)
		if !reflect.DeepEqual(got, tc.Want) {
			t.Errorf("got %+v, want %+v", got, tc.Want)
		}
	}
}
