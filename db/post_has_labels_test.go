package db

import (
	"fmt"
	"testing"
	"time"
)

func generatePost(p *Persister, i int64) (*Post, error) {
	var user = p.NewUser("Antoine",
		time.Now().UTC(),
		-5,
		"antoine@grondin.com")

	var author = p.NewAuthor("aybabtme", user)
	err := author.Save()
	if err != nil {
		return nil, err
	}
	var post = p.NewPost(author.Id(),
		fmt.Sprintf("Title #%d", i),
		fmt.Sprintf("Content #%d", i),
		fmt.Sprintf("ImageUrl #%d", i),
		time.Now().UTC())

	err = post.Save()
	return post, err
}

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
	defer p.DeletePersistance()

	// Get a post
	post, err := generatePost(p, 0)
	if err != nil {
		t.Error(err)
		return
	}
	// to which we add a label
	label, err := post.AddLabel("potato stories")
	if err != nil {
		t.Error(err)
		return
	}
	// we dont need to assert that its been added, this is done elsewhere
	// so we just remove the label right away
	err = post.RemoveLabel(&label)
	if err != nil {
		t.Error(err)
		return
	}

	// now the post should return an empty list for its labels
	shouldBeEmptyList, err := post.Labels()
	if err != nil {
		t.Error(err)
		return
	}

	if len(shouldBeEmptyList) != 0 {
		t.Errorf("Expected len(list)=<%d> but was <%d>",
			0, len(shouldBeEmptyList))
		return
	}
}

// post.Labels()
func TestAllLabelsOfPost(t *testing.T) {
	allLabelsOfPost(t, setupSQLitePersist())
	//allLabelsOfPost(t, setupPGPersist())
}

func allLabelsOfPost(t *testing.T, p *Persister) {
	defer p.DeletePersistance()

	post, err := generatePost(p, 0)
	if err != nil {
		t.Error(err)
		return
	}

	var expected []Label
	for i := int(1); i < 10; i++ {
		label, _ := post.AddLabel(fmt.Sprintf("cool topic #%d", i))
		expected = append(expected, label)
	}

	actual, err := post.Labels()

	if err != nil {
		t.Error("Couldn't query for Labels", err)
		return
	}

	for i := 0; i < len(expected); i++ {
		if actual[i].Name() != expected[i].Name() {
			t.Errorf("Expected <%s> but was <%s>",
				expected[i].Name(), actual[i].Name())
			return
		}
	}
}

// label.Destroy()
func TestDestroyLabel(t *testing.T) {
	destroyLabel(t, setupSQLitePersist())
	// TODO fix this, it crashes for some reason
	// destroyLabel(t, setupPGPersist())
}

func destroyLabel(t *testing.T, p *Persister) {
	defer p.DeletePersistance()

	for i := int64(1); i < 10; i++ {
		post, err := generatePost(p, 0)
		if err != nil {
			t.Error(err)
			return
		}

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
		actual, err := p.FindLabelById(id)

		if actual != nil {
			t.Errorf("Label should be destroyed but non-nil and id=%d and name=%s",
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
	defer p.DeletePersistance()

	var expected []Post
	var labels []Label
	for i := int(1); i < 10; i++ {

		post, err := generatePost(p, 0)
		if err != nil {
			t.Error(err)
			return
		}

		label, err := post.AddLabel("cat video")

		if err != nil {
			t.Errorf("Couldn't insert label #%d", i)
			t.Error(err)
		} else {
			labels = append(labels, label)
		}
		expected = append(expected, *post)
	}

	for i := 0; i < len(labels)-1; i++ {
		if labels[i].Id() != labels[i+1].Id() {
			t.Errorf("Expected same label id, but were <%d> and <%d>",
				labels[i].Id(), labels[i+1].Id())
			return
		}
	}

	var theLabel = labels[0]

	actual, err := theLabel.Posts()

	if err != nil {
		t.Error("Couldn't query for Posts", err)
		return
	}

	if len(expected) != len(actual) {
		t.Errorf("Received len(actual)=%d but expected len(expected)=%d",
			len(actual), len(expected))
		return
	}

	for i := 0; i < len(expected); i++ {
		if actual[i].Title() != expected[i].Title() {
			t.Errorf("Expected <%s> but was <%s>",
				expected[i].Title(), actual[i].Title())
			return
		}
	}
}
