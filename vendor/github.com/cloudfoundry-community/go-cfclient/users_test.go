package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListUsers(t *testing.T) {
	Convey("List Users", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/users", listUsersPayload, "", 200},
			{"GET", "/v2/usersPage2", listUsersPayloadPage2, "", 200},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		users, err := client.ListUsers()
		So(err, ShouldBeNil)

		So(len(users), ShouldEqual, 4)
		So(users[0].Guid, ShouldEqual, "ccec6d06-5f71-48a0-a4c5-c91a1d9f2fac")
		So(users[0].Username, ShouldEqual, "testUser1")
		So(users[1].Guid, ShouldEqual, "f97f5699-c920-4633-aa23-bd70f3db0808")
		So(users[1].Username, ShouldEqual, "testUser2")
		So(users[2].Guid, ShouldEqual, "cadd6389-fcf6-4928-84f0-6153556bf693")
		So(users[2].Username, ShouldEqual, "testUser3")
		So(users[3].Guid, ShouldEqual, "79c854b0-c12a-41b7-8d3c-fdd6e116e385")
		So(users[3].Username, ShouldEqual, "testUser4")
	})
}

func TestGetUserByUsername(t *testing.T) {
	Convey("Get User by Username", t, func() {
		user1 := User{Guid: "ccec6d06-5f71-48a0-a4c5-c91a1d9f2fac", Username: "testUser1"}
		user2 := User{Guid: "f97f5699-c920-4633-aa23-bd70f3db0808", Username: "testUser2"}
		user3 := User{Guid: "cadd6389-fcf6-4928-84f0-6153556bf693", Username: "testUser3"}
		user4 := User{Guid: "79c854b0-c12a-41b7-8d3c-fdd6e116e385", Username: "testUser4"}
		users := Users{user1, user2, user3, user4}

		So(users.GetUserByUsername("testUser1"), ShouldResemble, user1)
		So(users.GetUserByUsername("testUser2"), ShouldResemble, user2)
		So(users.GetUserByUsername("testUser3"), ShouldResemble, user3)
		So(users.GetUserByUsername("testUser4"), ShouldResemble, user4)
	})
}
