package post

import (
	"testing"
	"time"

	"github.com/lino-network/lino/test"
	post "github.com/lino-network/lino/tx/post"
	"github.com/lino-network/lino/types"

	crypto "github.com/tendermint/go-crypto"
)

// test publish a normal post
func TestNormalPublish(t *testing.T) {
	newAccountTransactionPriv := crypto.GenPrivKeyEd25519()
	newAccountPostPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newUser"
	postID1 := "New Post 1"
	postID2 := "New Post 2"
	// recover some stake
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		crypto.GenPrivKeyEd25519(), newAccountTransactionPriv, newAccountPostPriv, "100")

	test.CreateTestPost(
		t, lb, newAccountName, postID1, 0, newAccountPostPriv, "", "", "", "", "0", baseTime)
	test.CreateTestPost(
		t, lb, newAccountName, postID2, 1, newAccountTransactionPriv, "", "", "", "", "0", baseTime)
}

// test publish a repost
func TestNormalRepost(t *testing.T) {
	newAccountPostPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newUser"
	postID := "New Post"
	repostID := "Repost"
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		crypto.GenPrivKeyEd25519(), crypto.GenPrivKeyEd25519(), newAccountPostPriv, "100")

	test.CreateTestPost(
		t, lb, newAccountName, postID, 0, newAccountPostPriv, "", "", "", "", "0", baseTime)
	test.CreateTestPost(
		t, lb, newAccountName, repostID, 1, newAccountPostPriv,
		newAccountName, postID, "", "", "0", baseTime)

}

// test invalid repost if source post id doesn't exist
func TestInvalidRepost(t *testing.T) {
	newAccountPostPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newUser"
	postID := "New Post"
	repostID := "Repost"
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		crypto.GenPrivKeyEd25519(), crypto.GenPrivKeyEd25519(), newAccountPostPriv, "100")

	postCreateParams := post.PostCreateParams{
		PostID:                  postID,
		Title:                   string(make([]byte, 50)),
		Content:                 string(make([]byte, 1000)),
		Author:                  types.AccountKey(newAccountName),
		RedistributionSplitRate: "0",
	}
	msg := post.NewCreatePostMsg(postCreateParams)
	// reject due to stake
	test.SignCheckDeliver(t, lb, msg, 0, true, newAccountPostPriv, baseTime)
	postCreateParams.SourceAuthor = types.AccountKey(newAccountName)
	postCreateParams.SourcePostID = "invalid"
	postCreateParams.PostID = repostID
	msg = post.NewCreatePostMsg(postCreateParams)
	// invalid source post id
	test.SignCheckDeliver(t, lb, msg, 1, false, newAccountPostPriv, baseTime)
}

// test publish a comment
func TestComment(t *testing.T) {
	newAccountPostPriv := crypto.GenPrivKeyEd25519()
	newAccountName := "newUser"
	postID := "New Post"
	comment := "Comment"
	baseTime := time.Now().Unix() + 3600
	lb := test.NewTestLinoBlockchain(t, test.DefaultNumOfVal)

	test.CreateAccount(t, newAccountName, lb, 0,
		crypto.GenPrivKeyEd25519(), crypto.GenPrivKeyEd25519(), newAccountPostPriv, "100")

	test.CreateTestPost(
		t, lb, newAccountName, postID, 0, newAccountPostPriv, "", "", "", "", "0", baseTime)
	test.CreateTestPost(
		t, lb, newAccountName, comment, 1, newAccountPostPriv,
		"", "", newAccountName, postID, "0", baseTime)
}