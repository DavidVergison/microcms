package savearticle

import (
	"errors"
	"testing"

	"github.com/DavidVergison/microcms/articles/entities"
	customError "github.com/DavidVergison/microcms/articles/errors"
	"github.com/DavidVergison/microcms/tests"
	"github.com/go-test/deep"
)

func TestSaveHandler(t *testing.T) {

	/*
		GIVEN a user "anon" who can write to a site "gogole" and "fabecook"
		AND a article to save on "gogole"
		AND this is an update
		WHEN saving the article
		THEN there is no error
		AND the function response is false
	*/
	t.Run("A user can update an article", func(t *testing.T) {
		// GIVEN
		article := tests.TestGetArticle()
		site1 := tests.TestGetSite()
		site2 := tests.TestGetSite()
		article.User = site1.Username
		site2.Username = site1.Username
		article.Location.Site = site1.Site
		article.ReadingKey = site1.ReadingKey

		articleDbAdapter := fakeArticleDynamoDBAdapter{Answer: true}
		userDbAdapter := fakeUserDynamoDBAdapter{
			Answer: []entities.Site{
				site1,
				site2,
			},
		}

		// WHEN
		err := SaveArticle(&articleDbAdapter, &userDbAdapter, article)

		//THEN
		if err != nil {
			t.Fatalf("Unexpected error : %s", err.Error())
		}
		if diff := deep.Equal(article, articleDbAdapter.SavedArticle); diff != nil {
			t.Error(diff)
		}
	})

	/*
		GIVEN a user "anon" who can write to a site "gogole" and "fabecook"
		AND a article to save on "azamon"
		WHEN saving the article
		THEN there is an error
		AND the error is CustomErrForbidden
	*/
	t.Run("A user can not create an article on any site", func(t *testing.T) {
		// GIVEN
		article := tests.TestGetArticle()
		site1 := tests.TestGetSite()
		site2 := tests.TestGetSite()

		articleDbAdapter := fakeArticleDynamoDBAdapter{Answer: true}
		userDbAdapter := fakeUserDynamoDBAdapter{
			Answer: []entities.Site{
				site1,
				site2,
			},
		}

		// WHEN
		err := SaveArticle(&articleDbAdapter, &userDbAdapter, article)

		//THEN
		if err == nil {
			t.Fatal("Should answer with an error")
		}
		_, ok := err.(customError.CustomErrorForbidden)
		if !ok {
			t.Fatalf("Should answer with a CustomErrForbidden error, got %s", err.Error())
		}
	})

	/*
		GIVEN a user "anon" who can write to a site "gogole"
		AND a article to save on "gogole"
		AND there is a pb with the database when writing the article
		WHEN saving the article
		THEN there is an error
		AND the error is CustomDbError
	*/
	t.Run("When there is a DB error, the response is a CustomDbError", func(t *testing.T) {
		// GIVEN
		article := tests.TestGetArticle()
		site1 := tests.TestGetSite()
		site2 := tests.TestGetSite()
		article.User = site1.Username
		site2.Username = site1.Username
		article.Location.Site = site1.Site

		articleDbAdapter := fakeArticleDynamoDBAdapter{Error: errors.New("Sample Error")}
		userDbAdapter := fakeUserDynamoDBAdapter{
			Answer: []entities.Site{
				site1,
				site2,
			},
		}

		// WHEN
		err := SaveArticle(&articleDbAdapter, &userDbAdapter, article)

		//THEN
		if err == nil {
			t.Fatal("Should answer with an error")
		}
		_, ok := err.(customError.CustomErrorWithTheDatabase)
		if !ok {
			t.Fatalf("Should answer with a CustomErrorWithTheDatabase error, got %s", err.Error())
		}
	})
}

type fakeArticleDynamoDBAdapter struct {
	SavedArticle entities.Article
	Answer       bool
	Error        error
}

func (f *fakeArticleDynamoDBAdapter) SaveArticle(article entities.Article) (bool, error) {
	f.SavedArticle = article
	if f.Error != nil {
		return false, f.Error
	}
	return f.Answer, nil
}

type fakeUserDynamoDBAdapter struct {
	SavedUsername string
	Answer        []entities.Site
	Error         error
}

func (f *fakeUserDynamoDBAdapter) GetSitesByUser(username string) ([]entities.Site, error) {
	f.SavedUsername = username
	if f.Error != nil {
		return []entities.Site{}, f.Error
	}
	return f.Answer, f.Error
}
