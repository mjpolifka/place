package main

import "testing"

func TestParseArgsAndRun(t *testing.T) {
	// test not enough args
	err := parseArgsAndRun([]string{"fake program name"})
	if err != nil {
		if err.Error() != "Not enough arguments, show help" {
			t.Error("want: Not enough arguments, show help | got:", err)
		}
		return
	}
	t.Error("want: FAIL | got: PASS")
}

func TestMove(t *testing.T) {
	// test not enough args
	t.Run("not-enough-args", func(t *testing.T) {
		args := []string{"fake program name", "notepad", "0", "0", "800"} // missing height
		err := move(args)
		if err != nil {
			if err.Error() != "Not enough args for 'move', show help" {
				t.Error("want: Not enough args for 'move', show help | got:", err)
			}
			return
		}
		t.Error("want: FAIL, got: PASS")
	})
	// normalizeProcessName covered in validation_test
	// test x not int
	t.Run("x-not-int", func(t *testing.T) {
		args := []string{"fake program name", "notepad", "t", "0", "800", "800"} // missing height
		err := move(args)
		if err != nil {
			if err.Error() != "strconv.Atoi: parsing \"t\": invalid syntax" {
				t.Error("want: strconv.Atoi: parsing \"t\": invalid syntax | got:", err)
			}
			return
		}
		t.Error("want: FAIL, got: PASS")
	})
	// test y not int
	t.Run("y-not-int", func(t *testing.T) {
		args := []string{"fake program name", "notepad", "0", "t", "800", "800"} // missing height
		err := move(args)
		if err != nil {
			if err.Error() != "strconv.Atoi: parsing \"t\": invalid syntax" {
				t.Error("want: strconv.Atoi: parsing \"t\": invalid syntax | got:", err)
			}
			return
		}
		t.Error("want: FAIL, got: PASS")
	})
	// test width not int
	t.Run("width-not-int", func(t *testing.T) {
		args := []string{"fake program name", "notepad", "0", "0", "t", "800"} // missing height
		err := move(args)
		if err != nil {
			if err.Error() != "strconv.Atoi: parsing \"t\": invalid syntax" {
				t.Error("want: strconv.Atoi: parsing \"t\": invalid syntax | got:", err)
			}
			return
		}
		t.Error("want: FAIL, got: PASS")
	})
	// test height not int
	t.Run("height-not-int", func(t *testing.T) {
		args := []string{"fake program name", "notepad", "0", "0", "800", "t"} // missing height
		err := move(args)
		if err != nil {
			if err.Error() != "strconv.Atoi: parsing \"t\": invalid syntax" {
				t.Error("want: strconv.Atoi: parsing \"t\": invalid syntax | got:", err)
			}
			return
		}
		t.Error("want: FAIL, got: PASS")
	})
	// validateDimensions covered in validation_test
	// moveWindow covered in windows_calls_test
}

func TestCreate(t *testing.T) {
	// test too many args
	t.Run("too-many-args", func(t *testing.T) {
		t.Chdir(t.TempDir())
		args := []string{"fake program name", "create", "desktop", "too_many"}
		t.Log("len(args):", len(args))
		if err := create(args); err != nil {
			if err.Error() != "Too many args for 'create'" {
				t.Error("want: Too many args for 'create' | got:", err)
			}
			return
		}
		t.Error("want: FAIL | got: PASS")
	})
	// test not enough args
	t.Run("not-enough-args", func(t *testing.T) {
		t.Chdir(t.TempDir())
		args := []string{"fake program name", "create"}
		t.Log("len(args):", len(args))
		if err := create(args); err != nil {
			if err.Error() != "Not enough args for 'create'" {
				t.Error("want: Not enough args for 'create' | got:", err)
			}
			return
		}
		t.Error("want: FAIL | got: PASS")
	})
}
