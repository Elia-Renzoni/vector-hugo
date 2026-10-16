package functions_test

import (
	"testing"

	"github.com/vector-hugo/internal/functions"
	"github.com/vector-hugo/internal/memory"
)

func TestList(t *testing.T) {
	t.Run("TestListLPUSH", func(t *testing.T) {
		arena, err := memory.NewArena()
		if err != nil {
			t.Fail()
		}

		targets := []string{"abba", "bbaac"}

		_, err = functions.Lpush(arena, targets)
		if err != nil {
			t.Fail()
		}

		for _, target := range targets {
			ok, _, err := arena.Scan([]byte(target))
			if err != nil {
				t.Fail()
			}

			if !ok {
				t.Fatal("expected true got false")
			}
		}
	})

	t.Run("TestListRPUSH", func(t *testing.T) {
		arena, err := memory.NewArena()
		if err != nil {
			t.Fail()
		}

		targets := []string{"ffff", "cccd"}

		_, err = functions.Rpush(arena, targets)
		if err != nil {
			t.Fail()
		}

		for _, target := range targets {
			ok, _, err := arena.Scan([]byte(target))
			if err != nil {
				t.Fail()
			}

			if !ok {
				t.Fatal("expected true got false")
			}
		}

		_, n, _ := arena.Scan(nil)
		if n != 2 {
			t.Fatalf("expected 2 got %d", n)
		}
	})

	t.Run("TestListFree", func(t *testing.T) {
		arena, err := memory.NewArena()
		if err != nil {
			t.Fail()
		}

		targets := []string{"1230", "sjsjl"}

		_, err = functions.Lpush(arena, targets)
		if err != nil {
			t.Fail()
		}

		for range targets {
			arena.Free(nil, false)
		}

		_, n, _ := arena.Scan(nil)
		if n != 0 {
			t.Fatalf("expected 0 got %d", n)
		}
	})

	t.Run("TestListFree2", func(t *testing.T) {
		arena, err := memory.NewArena()
		if err != nil {
			t.Fail()
		}

		targets := []string{"30322", "foobarmock"}

		_, err = functions.Lpush(arena, targets)
		if err != nil {
			t.Fail()
		}

		for range targets {
			arena.Free(nil, true)
		}

		ok, _, _ := arena.Scan([]byte("30322"))
		ok, _, _ = arena.Scan([]byte("foobarmock"))
		if ok != false {
			t.Fatalf("expected 0 got %t", ok)
		}
	})

	t.Run("TestListLLEN", func(t *testing.T) {
		arena, err := memory.NewArena()
		if err != nil {
			t.Fail()
		}

		targets := []string{"1230", "sjsjl", "mockbarmock", "eeeee"}

		_, err = functions.Lpush(arena, targets)
		if err != nil {
			t.Fail()
		}

		n, err := functions.Llen(arena)
		if err != nil {
			t.Fail()
		}

		if n != len(targets) {
			t.Fatalf("expected %d got %d", len(targets), n)
		}

	})

	t.Run("TestListRPOP", func(t *testing.T) {
		arena, err := memory.NewArena()
		if err != nil {
			t.Fail()
		}

		targets := []string{"aaa", "bbbb", "cccc", "ddd"}
		n, err := functions.Lpush(arena, targets)
		if err != nil {
			t.Fail()
		}

		if n != len(targets) {
			t.Fail()
		}

		data, err := functions.Rpop(arena)
		if err != nil {
			t.Fail()
		}

		t.Logf("result: %s", data)

		if data != targets[len(targets)-1] {
			t.Fatalf("expected %s but got %s", targets[len(targets)-1], data)
		}
	})

	t.Run("TestListLPOP", func(t *testing.T) {
		arena, err := memory.NewArena()
		if err != nil {
			t.Fail()
		}
		targets := []string{"aaa", "bbbb", "cccc", "ddd"}
		n, err := functions.Lpush(arena, targets)
		if err != nil {
			t.Fail()
		}
		if n != len(targets) {
			t.Fail()
		}

		data, err := functions.Lpop(arena)
		if err != nil {
			t.Fail()
		}

		t.Logf("result: %s", data)

		if data != targets[0] {
			t.Fatalf("expected %s but got %s", targets[0], data)
		}
	})
}
