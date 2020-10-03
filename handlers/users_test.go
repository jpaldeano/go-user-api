package handlers_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jpaldi/go-user-api/handlers"
	"github.com/jpaldi/go-user-api/mongo"
)

type mockDatabase struct {
	createUser func(ctx context.Context, nickname string, firstname string, lastname string, password string, email string, country string) (*mongo.User, error)
	updateUser func(ctx context.Context, guid string, nickname string, firstname string, lastname string, password string, email string, country string) (*mongo.User, error)
	removeUser func(ctx context.Context, guid string) error
	getUsers   func(ctx context.Context, params url.Values) ([]*mongo.User, error)
}

func (m mockDatabase) CreateUser(ctx context.Context, nickname string, firstname string, lastname string, password string, email string, country string) (*mongo.User, error) {
	return m.createUser(ctx, nickname, firstname, lastname, password, email, country)
}

func (m mockDatabase) GetUsers(ctx context.Context, params url.Values) ([]*mongo.User, error) {
	return m.getUsers(ctx, params)
}

func (m mockDatabase) RemoveUser(ctx context.Context, guid string) error {
	return m.removeUser(ctx, guid)
}

func (m mockDatabase) UpdateUser(ctx context.Context, guid string, nickname string, firstname string, lastname string, password string, email string, country string) (*mongo.User, error) {
	return m.updateUser(ctx, guid, nickname, firstname, lastname, password, email, country)
}

func createPOSTRequest(method string, path string, body string) *http.Request {
	r, _ := http.NewRequest(method, path, strings.NewReader(body))

	r.Header.Set("Content-Type", "application/json")
	return r
}

func mockInsertUserInDatabaseOK() mockDatabase {
	return mockDatabase{
		createUser: func(ctx context.Context, nickname string, firstname string, lastname string, password string, email string, country string) (*mongo.User, error) {
			// Only testing Nickname but these fields should be all tested
			return &mongo.User{Nickname: nickname}, nil
		},
	}
}
func TestCreateUser(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name               string
		request            *http.Request
		expectedStatusCode int
		expectedResponse   string
		database           mockDatabase
	}{
		{
			name: "should return a 200 if the user is created",
			request: createPOSTRequest(http.MethodPost, "/users",
				`{
					"nickname": "test",
					"email": "test@email.uk",
					"first_name": "test",
					"last_name": "test",
					"password": "test",
					"country": "UK"}`),
			database:           mockInsertUserInDatabaseOK(),
			expectedResponse:   "{\"id\":\"\",\"nickname\":\"test\",\"first_name\":\"\",\"last_name\":\"\",\"password\":\"\",\"email\":\"\",\"country\":\"\"}\n",
			expectedStatusCode: 200,
		},

		{
			name: "should report fields if these are missing in the POST body and return a 400",
			request: createPOSTRequest(http.MethodPost, "/users",
				`{
					"nickname": "test",
					"email": "test@email.uk",
					"first_name": "test",
					"last_name": "test",
					"country": "UK"}`),
			database:           mockInsertUserInDatabaseOK(),
			expectedResponse:   "{\"validationError\":{\"password\":[\"The password field is required!\"]}}\n",
			expectedStatusCode: 400,
		},
	} {

		t.Run(tt.name, func(t *testing.T) {
			handler := handlers.Handler{
				Database: tt.database,
			}

			w := httptest.NewRecorder()

			handler.CreateUser(w, tt.request)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				t.Fatalf("couldn't read response body: got %s , err %s", body, err.Error())
			}

			if string(body) != tt.expectedResponse {
				t.Fatalf("wrong error: got %s want %s", body, tt.expectedResponse)
			}

			if resp.StatusCode != tt.expectedStatusCode {
				t.Fatalf("wrong error: got %d want %d", resp.StatusCode, tt.expectedStatusCode)
			}

		})
	}
}

func TestUpdateUser(t *testing.T) {
	// TODO
}

func TestRemoveUser(t *testing.T) {
	// TODO
}

func TestGetUsers(t *testing.T) {
	// TODO
}
