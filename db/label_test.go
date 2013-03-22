package db

import (
	"fmt"
	"testing"
	"time"
)

func TestSaveLabel(t *testing.T) {
	saveLabel(t, setupSQLitePersist())
	//saveLabel(t, setupPGPersist())
}

func saveLabel(t *testing.T, pers *Persister) {
	defer pers.DeletePersistance()

	var user = pers.NewUser("Antoine", time.Now().UTC(), -5, "antoine@grondin.com")
	var author = pers.NewAuthor("aybabtme", user)
	_ = author.Save()
	var post = pers.NewPost(author.Id(),
		"My first post",
		"Hello World",
		"fake.url/to/image.jpg",
		time.Now().UTC())
	_ = post.Save()

	expected, err := post.AddLabel("food")
	if err != nil {
		t.Error("Couldn't create label", err)
		return
	}

	if expected.Id() == -1 {
		t.Error("Received a label that wasn't saved")
		return
	}

	expected.SetName("delicious food")

	if err := expected.Save(); err != nil {
		t.Error("Save failed", err)
		return
	}

	if expected.Id() != 1 {
		t.Errorf("Id should be 1 at this point, but was %d", expected.Id())
		return
	}

	expected.Save()
	if expected.Id() != 1 {
		t.Errorf("Id should be 1 at this point, but was %d", expected.Id())
		return
	}

	labels, err := post.Labels()
	if err != nil {
		t.Error("Couldn't get labels back from the post", err)
		return
	}

	if len(labels) != 1 {
		t.Errorf("Created only 1 label but received %d", len(labels))
		return
	}

	var actual = labels[0]

	if actual.Id() != expected.Id() {
		t.Errorf("Id expected <%d> but was <%d>", expected.Id(), actual.Id())
		return
	}

	if actual.Name() != expected.Name() {
		t.Errorf("After saving, expected name=<%s> but was <%s>",
			expected.Name(), actual.Name())
		return
	}

}

func TestFindByIdLabel(t *testing.T) {
	findByIdLabel(t, setupSQLitePersist())
	// TODO fix this, it crashes for some reasons
	//findByIdLabel(t, setupPGPersist())
}

func findByIdLabel(t *testing.T, pers *Persister) {
	defer pers.DeletePersistance()

	for i := int64(1); i < 10; i++ {
		var user = pers.NewUser("Antoine", time.Now().UTC(), -5, "antoine@grondin.com")
		var author = pers.NewAuthor("aybabtme", user)
		_ = author.Save()
		var post = pers.NewPost(author.Id(),
			fmt.Sprintf("Title #%d", i),
			fmt.Sprintf("Content #%d", i),
			fmt.Sprintf("ImageUrl #%d", i),
			time.Now().UTC())
		_ = post.Save()
		expected, err := post.AddLabel(fmt.Sprintf("cool topic #%d", i))
		if err != nil {
			t.Error("Couldn't create a label to begin with.")
		}
		actual, err := pers.FindLabelById(expected.Id())

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
	defer pers.DeletePersistance()
	var labelCount = int64(9)

	for i := int64(1); i <= labelCount; i++ {
		var user = pers.NewUser("Antoine", time.Now().UTC(), -5, "antoine@grondin.com")
		var author = pers.NewAuthor("aybabtme", user)
		_ = author.Save()
		var post = pers.NewPost(author.Id(),
			fmt.Sprintf("Title #%d", i),
			fmt.Sprintf("Content #%d", i),
			fmt.Sprintf("ImageUrl #%d", i),
			time.Now().UTC())
		_ = post.Save()
		_, err := post.AddLabel(fmt.Sprintf("cool topic #%d", i))
		if err != nil {
			t.Error("Couldn't create a label to begin with.")
		}
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
}

func TestLabelIdIncrements(t *testing.T) {
	labelIdIncrements(t, setupSQLitePersist())
	// TODO PG doesnt work
	// idIncrements(t, setupPGPersist())
}

func labelIdIncrements(t *testing.T, pers *Persister) {
	defer pers.DeletePersistance()

	for i := int64(1); i < 10; i++ {
		var user = pers.NewUser("Antoine", time.Now().UTC(), -5, "antoine@grondin.com")
		var author = pers.NewAuthor("aybabtme", user)
		_ = author.Save()
		var post = pers.NewPost(author.Id(),
			fmt.Sprintf("Title #%d", i),
			fmt.Sprintf("Content #%d", i),
			fmt.Sprintf("ImageUrl #%d", i),
			time.Now().UTC())
		_ = post.Save()
		label, err := post.AddLabel(fmt.Sprintf("cool topic #%d", i))
		if err != nil {
			t.Error("Couldn't create a label to begin with.")
		}

		if label.Id() != i {
			t.Errorf("Id expected %d but was %d", i, label.Id())
		}
	}
}
