package db

import (
	"fmt"
	"testing"
	"time"
)

// post.AddLabel(string)
func TestAddLabelToPost(t *testing.T) {
	addLabelToPost(t, setupSQLitePersist())
	//addLabelToPost(t, setupPGPersist())
}

func addLabelToPost(t *testing.T, p *Persister) {

}

// post.RemoveLabel(Label)
func TestRemoveLabelFromPost(t *testing.T) {
	removeLabelFromPost(t, setupSQLitePersist())
	//removeLabelFromPost(t, setupPGPersist())
}

func removeLabelFromPost(t *testing.T, p *Persister) {

}

// post.Labels()
func TestAllLabelsOfPost(t *testing.T) {
	allLabelsOfPost(t, setupSQLitePersist())
	//allLabelsOfPost(t, setupPGPersist())
}

func allLabelsOfPost(t *testing.T, p *Persister) {

}

// label.Destroy()
func TestDestroyLabel(t *testing.T) {
	destroyLabel(t, setupSQLitePersist())
	// TODO fix this, it crashes for some reason
	// destroyLabel(t, setupPGPersist())
}

func destroyLabel(t *testing.T, pers *Persister) {
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
		expected, err := post.AddLabel(fmt.Sprintf("cool topic #%d", i))
		if err != nil {
			t.Error("Couldn't create a label to begin with.")
		}
		var id = expected.Id()
		expected.Save()

		err = expected.Destroy()
		if err != nil {
			t.Error("Couldn't delete the label", err)
		}
		actual, err := pers.FindLabelById(id)

		if actual.Id() != -1 {
			t.Errorf("Label should be destroyed but id=%d and name=%s",
				actual.Id(), actual.Name())
		}

		if err == nil {
			t.Error("An error should have been raised")
		}

	}
}

// label.Posts()
func TestAllPostsOfLabel(t *testing.T) {
	allPostsOfLabel(t, setupSQLitePersist())
	//allPostsOfLabel(t, setupPGPersist())
}

func allPostsOfLabel(t *testing.T, p *Persister) {

}
