package cfclient

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListSpaces(t *testing.T) {
	Convey("List Space", t, func() {
		mocks := []MockRoute{
			{"GET", "/v2/spaces", listSpacesPayload, "", 200},
			{"GET", "/v2/spacesPage2", listSpacesPayloadPage2, "", 200},
		}
		setupMultiple(mocks, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		spaces, err := client.ListSpaces()
		So(err, ShouldBeNil)

		So(len(spaces), ShouldEqual, 4)
		So(spaces[0].Guid, ShouldEqual, "8efd7c5c-d83c-4786-b399-b7bd548839e1")
		So(spaces[0].Name, ShouldEqual, "dev")
		So(spaces[0].OrganizationGuid, ShouldEqual, "a537761f-9d93-4b30-af17-3d73dbca181b")
		So(spaces[1].Guid, ShouldEqual, "657b5923-7de0-486a-9928-b4d78ee24931")
		So(spaces[1].Name, ShouldEqual, "demo")
		So(spaces[1].OrganizationGuid, ShouldEqual, "da0dba14-6064-4f7a-b15a-ff9e677e49b2")
		So(spaces[2].Guid, ShouldEqual, "9ffd7c5c-d83c-4786-b399-b7bd54883977")
		So(spaces[2].Name, ShouldEqual, "test")
		So(spaces[2].OrganizationGuid, ShouldEqual, "a537761f-9d93-4b30-af17-3d73dbca181b")
		So(spaces[3].Guid, ShouldEqual, "329b5923-7de0-486a-9928-b4d78ee24982")
		So(spaces[3].Name, ShouldEqual, "prod")
		So(spaces[3].OrganizationGuid, ShouldEqual, "da0dba14-6064-4f7a-b15a-ff9e677e49b2")
	})
}

func TestSpaceOrg(t *testing.T) {
	Convey("Find space org", t, func() {
		setup(MockRoute{"GET", "/v2/org/foobar", orgPayload, "", 200}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		space := &Space{
			Guid:   "123",
			Name:   "test space",
			OrgURL: "/v2/org/foobar",
			c:      client,
		}
		org, err := space.Org()
		So(err, ShouldBeNil)

		So(org.Name, ShouldEqual, "test-org")
		So(org.Guid, ShouldEqual, "da0dba14-6064-4f7a-b15a-ff9e677e49b2")
	})
}

func TestSpaceQuota(t *testing.T) {
	Convey("Get space quota", t, func() {
		setup(MockRoute{"GET", "/v2/space_quota_definitions/9ffd7c5c-d83c-4786-b399-b7bd54883977", spaceQuotaPayload, "", 200}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		space := &Space{
			QuotaDefinitionGuid: "9ffd7c5c-d83c-4786-b399-b7bd54883977",
			c:                   client,
		}

		spaceQuota, err := space.Quota()
		So(err, ShouldBeNil)

		So(spaceQuota.Guid, ShouldEqual, "9ffd7c5c-d83c-4786-b399-b7bd54883977")
		So(spaceQuota.Name, ShouldEqual, "test-2")
		So(spaceQuota.NonBasicServicesAllowed, ShouldEqual, false)
		So(spaceQuota.TotalServices, ShouldEqual, 10)
		So(spaceQuota.TotalRoutes, ShouldEqual, 20)
		So(spaceQuota.MemoryLimit, ShouldEqual, 30)
		So(spaceQuota.InstanceMemoryLimit, ShouldEqual, 40)
		So(spaceQuota.AppInstanceLimit, ShouldEqual, 50)
		So(spaceQuota.AppTaskLimit, ShouldEqual, 60)
		So(spaceQuota.TotalServiceKeys, ShouldEqual, 70)
		So(spaceQuota.TotalReservedRoutePorts, ShouldEqual, 80)
	})
}

func TestSpaceSummary(t *testing.T) {
	Convey("Get space summary", t, func() {
		setup(MockRoute{"GET", "/v2/spaces/494d8b64-8181-4183-a6d3-6279db8fec6e/summary", spaceSummaryPayload, "", 200}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		space := &Space{
			Guid: "494d8b64-8181-4183-a6d3-6279db8fec6e",
			c:    client,
		}

		summary, err := space.Summary()
		So(err, ShouldBeNil)

		So(summary.Guid, ShouldEqual, "494d8b64-8181-4183-a6d3-6279db8fec6e")
		So(summary.Name, ShouldEqual, "test")

		So(len(summary.Apps), ShouldEqual, 1)
		So(summary.Apps[0].Guid, ShouldEqual, "b5f0d1bd-a3a9-40a4-af1a-312ad26e5379")
		So(summary.Apps[0].Name, ShouldEqual, "test-app")
		So(summary.Apps[0].ServiceCount, ShouldEqual, 1)
		So(summary.Apps[0].RunningInstances, ShouldEqual, 1)
		So(summary.Apps[0].Memory, ShouldEqual, 256)
		So(summary.Apps[0].Instances, ShouldEqual, 1)
		So(summary.Apps[0].DiskQuota, ShouldEqual, 512)
		So(summary.Apps[0].State, ShouldEqual, "STARTED")
		So(summary.Apps[0].Diego, ShouldEqual, true)

		So(len(summary.Services), ShouldEqual, 1)
		So(summary.Services[0].Guid, ShouldEqual, "3c5c758c-6b76-46f6-89d5-677909bfc975")
		So(summary.Services[0].Name, ShouldEqual, "test-service")
		So(summary.Services[0].BoundAppCount, ShouldEqual, 1)
	})
}

func TestSpaceRoles(t *testing.T) {
	Convey("Get space roles", t, func() {
		setup(MockRoute{"GET", "/v2/spaces/494d8b64-8181-4183-a6d3-6279db8fec6e/user_roles", spaceRolesPayload, "", 200}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		space := &Space{
			Guid: "494d8b64-8181-4183-a6d3-6279db8fec6e",
			c:    client,
		}

		roles, err := space.Roles()
		So(err, ShouldBeNil)

		So(len(roles), ShouldEqual, 1)
		So(roles[0].Guid, ShouldEqual, "uaa-id-363")
		So(roles[0].Admin, ShouldEqual, false)
		So(roles[0].Active, ShouldEqual, false)
		So(roles[0].DefaultSpaceGuid, ShouldEqual, "")
		So(roles[0].Username, ShouldEqual, "everything@example.com")
		So(roles[0].SpaceRoles, ShouldResemble, []string{"space_developer", "space_manager", "space_auditor"})
		So(roles[0].SpacesUrl, ShouldEqual, "/v2/users/uaa-id-363/spaces")
		So(roles[0].OrganizationsUrl, ShouldEqual, "/v2/users/uaa-id-363/organizations")
		So(roles[0].ManagedOrganizationsUrl, ShouldEqual, "/v2/users/uaa-id-363/managed_organizations")
		So(roles[0].BillingManagedOrganizationsUrl, ShouldEqual, "/v2/users/uaa-id-363/billing_managed_organizations")
		So(roles[0].AuditedOrganizationsUrl, ShouldEqual, "/v2/users/uaa-id-363/audited_organizations")
		So(roles[0].ManagedSpacesUrl, ShouldEqual, "/v2/users/uaa-id-363/managed_spaces")
		So(roles[0].AuditedSpacesUrl, ShouldEqual, "/v2/users/uaa-id-363/audited_spaces")
	})
}

func TestAssociateSpaceAuditorByUsername(t *testing.T) {
	Convey("Associate auditor by username", t, func() {
		setup(MockRoute{"PUT", "/v2/spaces/bc7b4caf-f4b8-4d85-b126-0729b9351e56/auditors", associateSpaceAuditorPayload, "", 201}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		space := &Space{
			Guid: "bc7b4caf-f4b8-4d85-b126-0729b9351e56",
			c:    client,
		}

		newSpace, err := space.AssociateAuditorByUsername("user-name")
		So(err, ShouldBeNil)
		So(newSpace.Guid, ShouldEqual, "bc7b4caf-f4b8-4d85-b126-0729b9351e56")
	})
}

func TestRemoveSpaceAuditorByUsername(t *testing.T) {
	Convey("Remove auditor by username", t, func() {
		setup(MockRoute{"DELETE", "/v2/spaces/bc7b4caf-f4b8-4d85-b126-0729b9351e56/auditors", "", "", 204}, t)
		defer teardown()
		c := &Config{
			ApiAddress: server.URL,
			Token:      "foobar",
		}
		client, err := NewClient(c)
		So(err, ShouldBeNil)

		space := &Space{
			Guid: "bc7b4caf-f4b8-4d85-b126-0729b9351e56",
			c:    client,
		}

		err = space.RemoveAuditorByUsername("user-name")
		So(err, ShouldBeNil)
	})
}
