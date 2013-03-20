package db

import (
	"fmt"
	"testing"
)

func TestNewLabel(t *testing.T) {
	newLabel(t, setupSQLitePersist())
	//newLabel(t, setupPGPersist())
}

func newLabel(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()

	var label = persist.NewLabel("justin bieber topic")

	if label == nil {
		t.Error("Receive a nil label")
	}
}

func TestSaveLabel(t *testing.T) {
	saveLabel(t, setupSQLitePersist())
	//saveLabel(t, setupPGPersist())
}

func saveLabel(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()

	var label = persist.NewLabel("food label")

	if label == nil {
		t.Error("Receive a nil label")
	}

	if label.Id() != -1 {
		t.Error("Id should be of -1 at this point")
	}

	if err := label.Save(); err != nil {
		t.Error("Save failed", err)
	}

	if label.Id() != 1 {
		t.Error("Id should be 1 at this point")
	}
}

func TestDestroyLabel(t *testing.T) {
	destroyLabel(t, setupSQLitePersist())
	// TODO fix this, it crashes for some reason
	// destroyLabel(t, setupPGPersist())
}

func destroyLabel(t *testing.T, pers *Persister) {
	defer pers.DeletePersistance()

	for i := int64(1); i < 100; i++ {

		var expected = pers.NewLabel(fmt.Sprintf("cool topic #%d", i))

		var id = expected.Id()
		expected.Save()

		expected.Destroy()
		actual, err := pers.FindLabelById(id)

		if actual != nil {
			t.Error("Label shouldnt exist in DB after destroy")
		}

		if err == nil {
			t.Error("An error should have been raised")
		}

	}
}

func TestFindByIdLabel(t *testing.T) {
	findByIdLabel(t, setupSQLitePersist())
	// TODO fix this, it crashes for some reasons
	//findByIdLabel(t, setupPGPersist())
}

func findByIdLabel(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()
	for i := int64(1); i < 100; i++ {
		var expected = persist.NewLabel(fmt.Sprintf("cool topic #%d", i))
		expected.Save()

		actual, err := persist.FindLabelById(expected.Id())

		if err != nil {
			t.Errorf("Error while querying label %d: %v", i, err)
			return
		}

		if actual.Name() != expected.Name() {
			t.Errorf("Expected <%s> but was <%s>\n",
				expected.Name(), actual.Name())
			return
		}
	}
}

func TestFindAllLabel(t *testing.T) {
	findAllLabel(t, setupSQLitePersist())
	//findAllLabel(t, setupPGPersist())
}

func findAllLabel(t *testing.T, pers *Persister) {

	var labelCount = int64(10)

	for i := int64(1); i <= labelCount; i++ {
		var label = pers.NewLabel(fmt.Sprintf("cool topic #%d", i))
		label.Save()
	}

	labels, err := pers.FindAllLabels()
	if err != nil {
		t.Errorf("Couldn't query labels although just saved %d",
			labelCount)
	}

	if labels == nil {
		t.Errorf("Saved %d labels but query returns none",
			labelCount)
	}

	if int64(len(labels)) != labelCount {
		t.Errorf("Saved and expected <%d> posts, was <%d>",
			labelCount, len(labels))
	}

	for idx, label := range labels {
		if label.Id() != int64(idx)+int64(1) {
			t.Errorf("Expected <%d> but was <%d>", idx+1, label.Id())
		}
	}

	defer pers.DeletePersistance()
}

func TestLabelIdIncrements(t *testing.T) {
	labelIdIncrements(t, setupSQLitePersist())
	// TODO PG doesnt work
	// idIncrements(t, setupPGPersist())
}

func labelIdIncrements(t *testing.T, persist *Persister) {
	defer persist.DeletePersistance()

	for i := int64(1); i < 100; i++ {
		var label = persist.NewLabel(fmt.Sprintf("cool topic #%d", i))

		if label.Id() != -1 {
			t.Error("Id should be of -1 at this point")
		}

		label.Save()

		if label.Id() != i {
			t.Errorf("Id expected %d but was %d", i, label.Id())
		}
	}
}
