package middlewares

import (
	"testing"
	"time"
)

func TestRateRPS(t *testing.T) {
	st := time.Now()
	rps := NewRPSLimiter(3)
	t.Run("rps first chunk in parallel", func(t *testing.T) {
		data := []struct {
			t    time.Time
			want bool
		}{
			{
				t:    st.Add(200 * time.Millisecond),
				want: true,
			},
			{
				t:    st.Add(300 * time.Millisecond),
				want: true,
			},
			{
				t:    st.Add(600 * time.Millisecond),
				want: true,
			},
			{
				t:    st.Add(700 * time.Millisecond),
				want: false,
			},
		}

		for _, d := range data {
			got := rps.Allow()
			if got != d.want {
				t.Fatalf("expected the result to be %v but got %v", d.want, got)
			}
		}
	})
	t.Run("run mocked requests to test rps 2", func(t *testing.T) {
		time.Sleep(1200 * time.Millisecond)

		st = time.Now()
		data := []struct {
			t    time.Time
			want bool
		}{
			{
				t:    st.Add(200 * time.Millisecond),
				want: true,
			},
			{
				t:    st.Add(300 * time.Millisecond),
				want: true,
			},
			{
				t:    st.Add(600 * time.Millisecond),
				want: true,
			},
			{
				t:    st.Add(700 * time.Millisecond),
				want: false,
			},
		}
		for _, d := range data {
			got := rps.Allow()
			if got != d.want {
				t.Fatalf("expected the result to be %v but got %v", d.want, got)
			}
		}
		return
	})

}
