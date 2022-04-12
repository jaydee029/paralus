package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RafayLabs/rcloud-base/pkg/common"
	"github.com/RafayLabs/rcloud-base/pkg/query"
	commonv3 "github.com/RafayLabs/rcloud-base/proto/types/commonpb/v3"
	v3 "github.com/RafayLabs/rcloud-base/proto/types/commonpb/v3"
	userv3 "github.com/RafayLabs/rcloud-base/proto/types/userpb/v3"
	"github.com/google/uuid"
)

func performUserBasicChecks(t *testing.T, user *userv3.User, uuuid string) {
	_, err := uuid.Parse(user.GetMetadata().GetOrganization())
	if err == nil {
		t.Error("org in metadata should be name not id")
	}
	_, err = uuid.Parse(user.GetMetadata().GetPartner())
	if err == nil {
		t.Error("partner in metadata should be name not id")
	}
}

func performUserBasicAuthzChecks(t *testing.T, mazc mockAuthzClient, uuuid string, roles []*userv3.ProjectNamespaceRole) {
	if len(mazc.cp) > 0 {
		for i, u := range mazc.cp[len(mazc.cp)-1].Policies {
			if u.Sub != "u:user-"+uuuid {
				t.Errorf("invalid sub in policy sent to authz; expected '%v', got '%v'", "u:user-"+uuuid, u.Sub)
			}
			if u.Obj != roles[i].Role {
				t.Errorf("invalid obj in policy sent to authz; expected '%v', got '%v'", roles[i].Role, u.Obj)
			}
			if roles[i].Namespace != nil {
				if u.Ns != fmt.Sprint(*roles[i].Namespace) {
					t.Errorf("invalid ns in policy sent to authz; expected '%v', got '%v'", fmt.Sprint(roles[i].Namespace), u.Ns)
				}
			} else {
				if u.Ns != "*" {
					t.Errorf("invalid ns in policy sent to authz; expected '%v', got '%v'", "*", u.Ns)
				}
			}
			if roles[i].Project != nil {
				if u.Proj != *roles[i].Project {
					t.Errorf("invalid proj in policy sent to authz; expected '%v', got '%v'", roles[i].Project, u.Proj)
				}
			} else {
				if u.Proj != "*" {
					t.Errorf("invalid proj in policy sent to authz; expected '%v', got '%v'", "*", u.Proj)
				}
			}
		}
	}
}

func TestCreateUser(t *testing.T) {
	db, mock := getDB(t)
	defer db.Close()

	ap := &mockAuthProvider{}
	mazc := mockAuthzClient{}
	us := NewUserService(ap, db, &mazc, nil, common.CliConfigDownloadData{}, getLogger(), true)

	uuuid := uuid.New().String()
	puuid, ouuid := addParterOrgFetchExpectation(mock)
	mock.ExpectBegin()
	mock.ExpectCommit()

	user := &userv3.User{
		Metadata: &v3.Metadata{Partner: "partner-" + puuid, Organization: "org-" + ouuid, Name: "user-" + uuuid},
		Spec:     &userv3.UserSpec{},
	}
	user, err := us.Create(context.Background(), user)
	if err != nil {
		t.Fatal("could not create user:", err)
	}
	performUserBasicChecks(t, user, uuuid)
	if user.GetMetadata().GetName() != "user-"+uuuid {
		t.Errorf("expected name 'user-%v'; got '%v'", uuuid, user.GetMetadata().GetName())
	}
	performBasicAuthProviderChecks(t, *ap, 1, 0, 1, 0)
}

func TestCreateUserWithRole(t *testing.T) {
	tt := []struct {
		name       string
		role       bool
		project    bool
		namespace  bool
		dbname     string
		scope      string
		shouldfail bool
	}{
		{"just role", true, false, false, "authsrv_accountresourcerole", "system", false},
		{"just role org scope", true, false, false, "authsrv_accountresourcerole", "organization", false},
		{"just project", false, true, false, "authsrv_accountrole", "system", true},         // no role creation without role
		{"just namespace", false, false, true, "authsrv_accountrole", "system", true},       // no role creation without role,
		{"project and namespace", false, true, true, "authsrv_accountrole", "system", true}, // no role creation without role,
		{"project and role", true, true, false, "authsrv_projectaccountresourcerole", "project", false},
		{"project role namespace", true, true, true, "authsrv_projectaccountresourcerole", "project", false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := getDB(t)
			defer db.Close()

			ap := &mockAuthProvider{}
			mazc := mockAuthzClient{}
			us := NewUserService(ap, db, &mazc, nil, common.CliConfigDownloadData{}, getLogger(), true)

			uuuid := uuid.New().String()

			puuid, ouuid := addParterOrgFetchExpectation(mock)

			mock.ExpectBegin()
			ruuid := addResourceRoleFetchExpectation(mock, tc.scope)
			role := &userv3.ProjectNamespaceRole{}
			if tc.role {
				role.Role = idname(ruuid, "role")
			}
			if tc.project {
				pruuid := addFetchIdExpectation(mock, "project")
				role.Project = &pruuid
			}
			if tc.namespace {
				var ns int64 = 7
				role.Namespace = &ns
			}
			mock.ExpectQuery(fmt.Sprintf(`INSERT INTO "%v"`, tc.dbname)).
				WithArgs().WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New().String()))
			mock.ExpectCommit()

			user := &userv3.User{
				Metadata: &v3.Metadata{Partner: "partner-" + puuid, Organization: "org-" + ouuid, Name: "user-" + uuuid},
				Spec:     &userv3.UserSpec{ProjectNamespaceRoles: []*userv3.ProjectNamespaceRole{role}},
			}

			user, err := us.Create(context.Background(), user)
			if tc.shouldfail {
				if err == nil {
					// TODO: check for proper error messages
					t.Fatal("expected user not to be created, but was created")
				} else {
					return
				}
			}
			if err != nil {
				t.Fatal("could not create user:", err)
			}
			performUserBasicChecks(t, user, uuuid)
			if user.GetMetadata().GetName() != "user-"+uuuid {
				t.Errorf("expected name 'user-%v'; got '%v'", uuuid, user.GetMetadata().GetName())
			}

			performBasicAuthProviderChecks(t, *ap, 1, 0, 1, 0)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock := getDB(t)
	defer db.Close()

	ap := &mockAuthProvider{}
	mazc := mockAuthzClient{}
	us := NewUserService(ap, db, &mazc, nil, common.CliConfigDownloadData{}, getLogger(), true)

	// performing update
	uuuid := addUserIdFetchExpectation(mock)
	puuid, ouuid := addParterOrgFetchExpectation(mock)
	mock.ExpectBegin()
	_ = addUserRoleMappingsUpdateExpectation(mock, uuuid)
	addUserGroupMappingsUpdateExpectation(mock, uuuid)
	ruuid := addResourceRoleFetchExpectation(mock, "project")
	pruuid := addFetchExpectation(mock, "project")
	mock.ExpectQuery(`INSERT INTO "authsrv_projectaccountresourcerole"`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New().String()))
	mock.ExpectCommit()

	var ns int64 = 7
	user := &userv3.User{
		Metadata: &v3.Metadata{Partner: "partner-" + puuid, Organization: "org-" + ouuid, Name: "user-" + uuuid},
		Spec:     &userv3.UserSpec{ProjectNamespaceRoles: []*userv3.ProjectNamespaceRole{{Project: idnamea(pruuid, "project"), Namespace: &ns, Role: idname(ruuid, "role")}}},
	}
	user, err := us.Update(context.Background(), user)
	if err != nil {
		t.Fatal("could not create user:", err)
	}
	performUserBasicChecks(t, user, uuuid)
	if user.GetMetadata().GetName() != "user-"+uuuid {
		t.Errorf("expected name 'user-%v'; got '%v'", uuuid, user.GetMetadata().GetName())
	}
	performBasicAuthProviderChecks(t, *ap, 0, 1, 0, 0)
}

func TestUserGetByName(t *testing.T) {
	db, mock := getDB(t)
	defer db.Close()

	ap := &mockAuthProvider{}
	mazc := mockAuthzClient{}
	us := NewUserService(ap, db, &mazc, nil, common.CliConfigDownloadData{}, getLogger(), true)

	puuid := uuid.New().String()
	ouuid := uuid.New().String()
	guuid := uuid.New().String()
	ruuid := uuid.New().String()
	pruuid := uuid.New().String()

	uuuid := addUserFetchExpectation(mock)
	mock.ExpectQuery(`SELECT "group"."id".* FROM "authsrv_group" AS "group" JOIN authsrv_groupaccount ON authsrv_groupaccount.group_id="group".id WHERE .authsrv_groupaccount.account_id = '` + uuuid + `'`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"name"}).
		AddRow("group-" + guuid))
	mock.ExpectQuery(`SELECT authsrv_resourcerole.name as role, authsrv_group.name as group FROM "authsrv_grouprole" JOIN authsrv_resourcerole ON authsrv_resourcerole.id=authsrv_grouprole.role_id JOIN authsrv_group ON authsrv_group.id=authsrv_grouprole.group_id WHERE`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"role", "group"}).AddRow("role-"+ruuid, "group-"+guuid))

	mock.ExpectQuery(`SELECT authsrv_resourcerole.name as role, authsrv_project.name as project, authsrv_group.name as group FROM "authsrv_projectgrouprole" JOIN authsrv_resourcerole ON authsrv_resourcerole.id=authsrv_projectgrouprole.role_id JOIN authsrv_project ON authsrv_project.id=authsrv_projectgrouprole.project_id JOIN authsrv_group ON authsrv_group.id=authsrv_projectgrouprole.group_id WHERE`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"role", "project"}).AddRow("role-"+ruuid, "project-"+puuid))
	mock.ExpectQuery(`SELECT authsrv_resourcerole.name as role, authsrv_project.name as project, namespace_id as namespace, authsrv_group.name as group FROM "authsrv_projectgroupnamespacerole"`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"role", "project"}).AddRow("role-"+ruuid, "project-"+puuid))

	mock.ExpectQuery(`SELECT authsrv_resourcerole.name as role FROM "authsrv_accountresourcerole" JOIN authsrv_resourcerole ON authsrv_resourcerole.id=authsrv_accountresourcerole.role_id WHERE .authsrv_accountresourcerole.account_id = '` + uuuid + `'`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("role-" + ruuid))

	mock.ExpectQuery(`SELECT distinct authsrv_resourcerole.name as role, authsrv_project.name as project FROM "authsrv_projectaccountresourcerole" JOIN authsrv_resourcerole ON authsrv_resourcerole.id=authsrv_projectaccountresourcerole.role_id JOIN authsrv_project ON authsrv_project.id=authsrv_projectaccountresourcerole.project_id WHERE .authsrv_projectaccountresourcerole.account_id = '` + uuuid + `'`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"role", "project"}).AddRow("role-"+ruuid, "project-"+pruuid))
	mock.ExpectQuery(`SELECT authsrv_resourcerole.name as role, authsrv_project.name as project, namespace_id as namespace FROM "authsrv_projectaccountnamespacerole" JOIN authsrv_resourcerole ON authsrv_resourcerole.id=authsrv_projectaccountnamespacerole.role_id JOIN authsrv_project ON authsrv_project.id=authsrv_projectaccountnamespacerole.project_id WHERE .authsrv_projectaccountnamespacerole.account_id = '` + uuuid + `'`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"role", "project", "namespace"}).AddRow("role-"+ruuid, "project-"+pruuid, 9))

	user := &userv3.User{
		Metadata: &v3.Metadata{Partner: "partner-" + puuid, Organization: "org-" + ouuid, Name: "user-" + uuuid},
	}
	user, err := us.GetByName(context.Background(), user)
	if err != nil {
		t.Fatal("could not get user:", err)
	}
	performUserBasicChecks(t, user, uuuid)
	if user.GetMetadata().GetName() != "user-"+uuuid {
		t.Errorf("invalid email for user, expected '%v'; got '%v'", "user-"+uuuid, user.GetMetadata().GetName())
	}
	if len(user.GetSpec().GetGroups()) != 1 {
		t.Errorf("invalid number of groups returned for user, expected 2; got '%v'", len(user.GetSpec().GetGroups()))
	}
	if len(user.GetSpec().GetProjectNamespaceRoles()) != 6 {
		t.Errorf("invalid number of roles returned for user, expected 3; got '%v'", len(user.GetSpec().GetProjectNamespaceRoles()))
	}
	if user.GetSpec().GetProjectNamespaceRoles()[2].GetNamespace() != 9 {
		t.Errorf("invalid namespace in role returned for user, expected 9; got '%v'", user.GetSpec().GetProjectNamespaceRoles()[2].Namespace)
	}
	performBasicAuthProviderChecks(t, *ap, 0, 0, 0, 0)
}

func TestUserGetInfo(t *testing.T) {
	db, mock := getDB(t)
	defer db.Close()

	ap := &mockAuthProvider{}
	mazc := mockAuthzClient{}
	us := NewUserService(ap, db, &mazc, nil, common.CliConfigDownloadData{}, getLogger(), false)

	uuuid := uuid.New().String()
	fakeuuuid := uuid.New().String()
	puuid := uuid.New().String()
	ouuid := uuid.New().String()
	guuid := uuid.New().String()
	ruuid := uuid.New().String()
	pruuid := uuid.New().String()

	mock.ExpectQuery(`SELECT "identities"."id", "identities"."schema_id", .*WHERE .traits ->> 'email' = 'user-` + uuuid + `'.`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"id", "traits"}).AddRow(uuuid, []byte(`{"email":"johndoe@provider.com", "first_name": "John", "last_name": "Doe", "organization_id": "`+ouuid+`", "partner_id": "`+puuid+`", "description": "My awesome user"}`)))
	mock.ExpectQuery(`SELECT "group"."id".* FROM "authsrv_group" AS "group" JOIN authsrv_groupaccount ON authsrv_groupaccount.group_id="group".id WHERE .authsrv_groupaccount.account_id = '` + uuuid + `'`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"name"}).
		AddRow("group-" + guuid))
	mock.ExpectQuery(`SELECT authsrv_resourcerole.name as role, authsrv_group.name as group FROM "authsrv_grouprole" JOIN authsrv_resourcerole ON authsrv_resourcerole.id=authsrv_grouprole.role_id JOIN authsrv_group ON authsrv_group.id=authsrv_grouprole.group_id WHERE`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"role", "group"}).AddRow("role-"+ruuid, "group-"+guuid))
	mock.ExpectQuery(`SELECT authsrv_resourcerole.name as role, authsrv_project.name as project, authsrv_group.name as group FROM "authsrv_projectgrouprole" JOIN authsrv_resourcerole ON authsrv_resourcerole.id=authsrv_projectgrouprole.role_id JOIN authsrv_project ON authsrv_project.id=authsrv_projectgrouprole.project_id JOIN authsrv_group ON authsrv_group.id=authsrv_projectgrouprole.group_id WHERE`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"role", "project"}).AddRow("role-"+ruuid, "project-"+pruuid))
	mock.ExpectQuery(`SELECT authsrv_resourcerole.name as role, authsrv_project.name as project, namespace_id as namespace, authsrv_group.name as group FROM "authsrv_projectgroupnamespacerole"`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"role", "project"}).AddRow("role-"+ruuid, "project-"+pruuid))
	mock.ExpectQuery(`SELECT authsrv_resourcerole.name as role FROM "authsrv_accountresourcerole" JOIN authsrv_resourcerole ON authsrv_resourcerole.id=authsrv_accountresourcerole.role_id WHERE .authsrv_accountresourcerole.account_id = '` + uuuid + `'`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("role-" + ruuid))
	mock.ExpectQuery(`SELECT distinct authsrv_resourcerole.name as role, authsrv_project.name as project FROM "authsrv_projectaccountresourcerole" JOIN authsrv_resourcerole ON authsrv_resourcerole.id=authsrv_projectaccountresourcerole.role_id JOIN authsrv_project ON authsrv_project.id=authsrv_projectaccountresourcerole.project_id WHERE .authsrv_projectaccountresourcerole.account_id = '` + uuuid + `'`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"role", "project"}).AddRow("role-"+ruuid, "project-"+pruuid))
	mock.ExpectQuery(`SELECT authsrv_resourcerole.name as role, authsrv_project.name as project, namespace_id as namespace FROM "authsrv_projectaccountnamespacerole" JOIN authsrv_resourcerole ON authsrv_resourcerole.id=authsrv_projectaccountnamespacerole.role_id JOIN authsrv_project ON authsrv_project.id=authsrv_projectaccountnamespacerole.project_id WHERE .authsrv_projectaccountnamespacerole.account_id = '` + uuuid + `'`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"role", "project", "namespace"}).AddRow("role-"+ruuid, "project-"+pruuid, 9))
	mock.ExpectQuery(`SELECT "resourcerole"."id" FROM "authsrv_resourcerole" AS "resourcerole" WHERE .name = 'role-` + ruuid + `'. AND .trash = FALSE.`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(ruuid, "role-"+ruuid))
	mock.ExpectQuery(`SELECT authsrv_resourcepermission.name as name FROM "authsrv_resourcepermission" JOIN authsrv_resourcerolepermission ON authsrv_resourcerolepermission.resource_permission_id=authsrv_resourcepermission.id WHERE .authsrv_resourcerolepermission.resource_role_id = '` + ruuid + `'. AND .authsrv_resourcepermission.trash = FALSE. AND .authsrv_resourcerolepermission.trash = FALSE.`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("account.read").AddRow("account.write"))

	user := &userv3.User{
		Metadata: &v3.Metadata{Partner: "partner-" + puuid, Organization: "org-" + ouuid, Name: "user-" + fakeuuuid},
	}
	ctx := context.WithValue(context.Background(), common.SessionDataKey, &commonv3.SessionData{Username: "user-" + uuuid})
	userinfo, err := us.GetUserInfo(ctx, user)

	if err != nil {
		t.Fatal("could not get user:", err)
	}

	if userinfo.Metadata.Name != "johndoe@provider.com" {
		t.Errorf("incorrect username; expected '%v', got '%v'", "johndoe@provider.com", userinfo.Metadata.Name)
	}
	if userinfo.Spec.FirstName != "John" {
		t.Errorf("incorrect first name; expected '%v', got '%v'", "John", userinfo.Spec.FirstName)
	}
	if userinfo.Spec.LastName != "Doe" {
		t.Errorf("incorrect last name; expected '%v', got '%v'", "Doe", userinfo.Spec.LastName)
	}
	if len(userinfo.Spec.Groups) != 1 {
		t.Errorf("incorrect number of groups; expected '%v', got '%v'", 1, len(userinfo.Spec.Groups))
	}
	if userinfo.Spec.Groups[0] != "group-"+guuid {
		t.Errorf("incorrect group name; expected '%v', got '%v'", "group-"+guuid, userinfo.Spec.Groups[0])
	}
	if len(userinfo.Spec.Permissions) != 6 {
		t.Errorf("incorrect number of permissions; expected '%v', got '%v'", 6, len(userinfo.Spec.Permissions))
	}
	if len(userinfo.Spec.Permissions[0].Permissions) != 2 {
		t.Errorf("incorrect number of permissions; expected '%v', got '%v'", 2, len(userinfo.Spec.Permissions[0].Permissions))
	}

}

func TestUserGetById(t *testing.T) {
	db, mock := getDB(t)
	defer db.Close()

	ap := &mockAuthProvider{}
	mazc := mockAuthzClient{}
	us := NewUserService(ap, db, &mazc, nil, common.CliConfigDownloadData{}, getLogger(), true)

	uuuid := uuid.New().String()
	puuid := uuid.New().String()
	ouuid := uuid.New().String()
	pruuid := uuid.New().String()

	// lookup by id
	mock.ExpectQuery(`SELECT "identities"."id",.* FROM "identities" WHERE .*id = '` + uuuid + `'`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"id", "traits"}).AddRow(uuuid, []byte(`{"email":"johndoe@provider.com", "first_name": "John", "last_name": "Doe", "organization_id": "`+ouuid+`", "partner_id": "`+puuid+`", "description": "My awesome user"}`)))

	guuid := addUsersGroupFetchExpectation(mock, uuuid)
	addGroupRoleMappingsFetchExpectation(mock, guuid, pruuid)
	addUserRoleMappingsFetchExpectation(mock, uuuid, pruuid)

	user := &userv3.User{
		Metadata: &v3.Metadata{Partner: "partner-" + puuid, Organization: "org-" + ouuid, Id: uuuid},
	}
	user, err := us.GetByID(context.Background(), user)
	if err != nil {
		t.Fatal("could not get user:", err)
	}
	performUserBasicChecks(t, user, uuuid)
	if len(user.GetSpec().GetGroups()) != 1 {
		t.Errorf("invalid number of groups returned for user, expected 1; got '%v'", len(user.GetSpec().GetGroups()))
	}
	if len(user.GetSpec().GetProjectNamespaceRoles()) != 6 {
		t.Errorf("invalid number of roles returned for user, expected 6; got '%v'", len(user.GetSpec().GetProjectNamespaceRoles()))
	}
	if user.GetSpec().GetProjectNamespaceRoles()[2].GetNamespace() != 7 {
		t.Errorf("invalid namespace in role returned for user, expected 7; got '%v'", user.GetSpec().GetProjectNamespaceRoles()[2].Namespace)
	}

	performBasicAuthProviderChecks(t, *ap, 0, 0, 0, 0)
}

func TestUserList(t *testing.T) {
	tests := []struct {
		name     string
		q        string
		limit    int64
		offset   int64
		orderBy  string
		order    string
		role     string
		group    string
		projects []string
		utype    string
	}{
		{"simple list", "", 50, 20, "", "", "", "", []string{}, ""},
		{"simple list with type", "", 50, 20, "", "", "", "", []string{}, "password"},
		{"sorted list", "", 50, 20, "email", "asc", "", "", []string{}, ""},
		{"sorted list without dir", "", 50, 20, "email", "", "", "", []string{}, ""},
		{"sorted list with q", "filter-query", 50, 20, "email", "asc", "", "", []string{}, ""},
		{"sorted list with role", "", 50, 20, "email", "asc", "role-name", "", []string{}, ""},
		{"sorted list with role and group", "", 50, 20, "email", "asc", "role-name", "group-name", []string{}, ""},
		{"sorted list with q and role", "filter-query", 50, 20, "email", "asc", "role-name", "", []string{}, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock := getDB(t)
			defer db.Close()

			ap := &mockAuthProvider{}
			mazc := mockAuthzClient{}
			us := NewUserService(ap, db, &mazc, nil, common.CliConfigDownloadData{}, getLogger(), true)

			uuuid1 := uuid.New().String()
			uuuid2 := uuid.New().String()
			pruuid := uuid.New().String()

			puuid, ouuid := addParterOrgFetchExpectation(mock)
			q := ""
			if tc.q != "" {
				q = ` AND .traits ->> 'email' ILIKE '%` + tc.q + `%'. OR .traits ->> 'first_name' ILIKE '%` + tc.q + `%'. OR .traits ->> 'last_name' ILIKE '%` + tc.q + `%'. `
			}
			order := ""
			if tc.orderBy != "" {
				order = `ORDER BY "traits ->> '` + tc.orderBy + `' `
			}
			if tc.order != "" {
				order = order + tc.order + `" `
			}
			if tc.role != "" {
				addFetchExpectation(mock, "resourcerole")
			}
			if tc.group != "" {
				addFetchExpectation(mock, "group")
			}
			if tc.role != "" || tc.group != "" || len(tc.projects) != 0 {
				addSentryLookupExpectation(mock, []string{uuuid1, uuuid2}, puuid, ouuid)
				mock.ExpectQuery(`SELECT "identities"."id", .*WHERE .id IN .'` + uuuid1 + `', '` + uuuid2 + `'..  ` + q + order + `LIMIT ` + fmt.Sprint(tc.limit) + ` OFFSET ` + fmt.Sprint(tc.offset)).
					WithArgs().WillReturnRows(sqlmock.NewRows([]string{"id", "traits"}).
					AddRow(uuuid1, []byte(`{"email":"johndoe@provider.com", "first_name": "John", "last_name": "Doe", "organization_id": "`+ouuid+`", "partner_id": "`+puuid+`", "description": "My awesome user"}`)).
					AddRow(uuuid2, []byte(`{"email":"johndoe@provider.com", "first_name": "John", "last_name": "Doe", "organization_id": "`+ouuid+`", "partner_id": "`+puuid+`", "description": "My awesome user"}`)))
			} else {
				if tc.utype != "" {
					mock.ExpectQuery(`SELECT "identities"."id", .*, "identity_credential"."id" AS "identity_credential__id", .*FROM "identities" LEFT JOIN "identity_credentials" AS "identity_credential" ON ."identity_credential"."identity_id" = "identities"."id". LEFT JOIN "identity_credential_types" AS "identity_credential__identity_credential_type" ON ."identity_credential__identity_credential_type"."id" = "identity_credential"."identity_credential_type_id". WHERE .name = '` + tc.utype + `'. LIMIT 50 OFFSET 20`).
						WithArgs().WillReturnRows(sqlmock.NewRows([]string{"id", "traits"}).
						AddRow(uuuid1, []byte(`{"email":"johndoe@provider.com", "first_name": "John", "last_name": "Doe", "organization_id": "`+ouuid+`", "partner_id": "`+puuid+`", "description": "My awesome user"}`)).
						AddRow(uuuid2, []byte(`{"email":"johndoe@provider.com", "first_name": "John", "last_name": "Doe", "organization_id": "`+ouuid+`", "partner_id": "`+puuid+`", "description": "My awesome user"}`)))
				} else {
					mock.ExpectQuery(`SELECT "identities"."id".* LIMIT 50 OFFSET 20$`).
						WithArgs().WillReturnRows(sqlmock.NewRows([]string{"id", "traits"}).
						AddRow(uuuid1, []byte(`{"email":"johndoe@provider.com", "first_name": "John", "last_name": "Doe", "organization_id": "`+ouuid+`", "partner_id": "`+puuid+`", "description": "My awesome user"}`)).
						AddRow(uuuid2, []byte(`{"email":"johndoe@provider.com", "first_name": "John", "last_name": "Doe", "organization_id": "`+ouuid+`", "partner_id": "`+puuid+`", "description": "My awesome user"}`)))
				}
			}

			guuid := addUsersGroupFetchExpectation(mock, uuuid1)
			addGroupRoleMappingsFetchExpectation(mock, guuid, pruuid)
			addUserRoleMappingsFetchExpectation(mock, uuuid1, pruuid)

			guuid = addUsersGroupFetchExpectation(mock, uuuid2)
			addGroupRoleMappingsFetchExpectation(mock, guuid, pruuid)
			addUserRoleMappingsFetchExpectation(mock, uuuid2, pruuid)

			qo := &commonv3.QueryOptions{
				Q:            tc.q,
				Limit:        tc.limit,
				Offset:       tc.offset,
				OrderBy:      tc.orderBy,
				Order:        tc.order,
				Organization: ouuid,
				Partner:      puuid,
				Role:         tc.role,
				Group:        tc.group,
				Type:         tc.utype,
			}

			userlist, err := us.List(context.Background(), query.WithOptions(qo))
			if err != nil {
				t.Fatal("could not list users:", err)
			}
			if userlist.Metadata.Count != 2 {
				t.Fatalf("incorrect number of users returned, expected 2; got %v", userlist.Metadata.Count)
			}
			if userlist.Items[0].Metadata.Name != "johndoe@provider.com" || userlist.Items[1].Metadata.Name != "johndoe@provider.com" {
				t.Errorf("incorrect user names returned when listing; expected '%v' and '%v'; got '%v' and '%v'", "johndoe@provider.com", "johndoe@provider.com", userlist.Items[0].Metadata.Name, userlist.Items[1].Metadata.Name)
			}
			if len(userlist.Items[0].GetSpec().GetGroups()) != 1 {
				t.Errorf("invalid number of groups returned for user, expected 1; got '%v'", len(userlist.Items[0].GetSpec().GetGroups()))
			}

			if len(userlist.Items[0].GetSpec().GetProjectNamespaceRoles()) != 6 {
				t.Errorf("invalid number of roles returned for user, expected 6; got '%v'", len(userlist.Items[0].GetSpec().GetProjectNamespaceRoles()))
			}
			if userlist.Items[0].GetSpec().GetProjectNamespaceRoles()[2].GetNamespace() != 7 {
				t.Errorf("invalid namespace in role returned for user, expected 7; got '%v'", userlist.Items[0].GetSpec().GetProjectNamespaceRoles()[2].Namespace)
			}

			performBasicAuthProviderChecks(t, *ap, 0, 0, 0, 0)

		})
	}
}

func TestUserDelete(t *testing.T) {
	db, mock := getDB(t)
	defer db.Close()

	ap := &mockAuthProvider{}
	mazc := mockAuthzClient{}
	us := NewUserService(ap, db, &mazc, nil, common.CliConfigDownloadData{}, getLogger(), true)

	uuuid := uuid.New().String()
	puuid := uuid.New().String()
	ouuid := uuid.New().String()

	mock.ExpectQuery(`SELECT "identities"."id" FROM "identities" WHERE .*traits ->> 'email' = 'user-` + uuuid + `'`).
		WithArgs().WillReturnRows(sqlmock.NewRows([]string{"id", "traits"}).AddRow(uuuid, []byte(`{"email":"johndoe@provider.com"}`)))
	mock.ExpectBegin()
	_ = addUserRoleMappingsUpdateExpectation(mock, uuuid)
	// User delete is via kratos
	addUserGroupMappingsUpdateExpectation(mock, uuuid)
	mock.ExpectCommit()

	user := &userv3.User{
		Metadata: &v3.Metadata{Partner: "partner-" + puuid, Organization: "org-" + ouuid, Name: "user-" + uuuid},
	}
	_, err := us.Delete(context.Background(), user)
	if err != nil {
		t.Fatal("could not delete user:", err)
	}

	performBasicAuthProviderChecks(t, *ap, 0, 0, 0, 1)
}