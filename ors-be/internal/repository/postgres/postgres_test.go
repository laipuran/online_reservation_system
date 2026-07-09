package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

const repositoryTestDatabaseURLEnv = "ORS_REPOSITORY_TEST_DATABASE_URL"

type repoFixture struct {
	ctx context.Context

	pool *pgxpool.Pool

	users         repository.UserRepository
	providers     repository.ServiceProviderRepository
	categories    repository.CategoryRepository
	tags          repository.TagRepository
	serviceTags   repository.ServiceTagRepository
	interests     repository.UserInterestRepository
	services      repository.ServiceRepository
	reservations  repository.ReservationRepository
	reviews       repository.ReviewRepository
	notifications repository.NotificationRepository
}

func TestConnect_InvalidDSN(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pool, err := Connect(ctx, "://bad")
	if err == nil {
		if pool != nil {
			pool.Close()
		}
		t.Fatal("Connect() error = nil, want error")
	}
}

func TestConnect_Success(t *testing.T) {
	url := testDatabaseURL(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := Connect(ctx, url)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	t.Cleanup(pool.Close)
}

func TestUserProviderCategoryTagRepositories(t *testing.T) {
	f := newRepoFixture(t)
	ctx := f.ctx

	customer := f.createUser(t, "customer@example.com", "customer")
	providerUser := f.createUser(t, "provider@example.com", "provider")

	if customer.ID == 0 || customer.CreatedAt.IsZero() || customer.UpdatedAt.IsZero() {
		t.Fatalf("created user missing generated fields: %+v", customer)
	}

	byEmail, err := f.users.GetByEmail(ctx, customer.Email)
	requireNoError(t, err)
	if byEmail == nil || byEmail.ID != customer.ID || byEmail.Phone != "" || byEmail.AvatarURL != "" {
		t.Fatalf("GetByEmail() = %+v, want user with empty optional fields", byEmail)
	}

	byID, err := f.users.GetByID(ctx, customer.ID)
	requireNoError(t, err)
	if byID == nil || byID.Email != customer.Email {
		t.Fatalf("GetByID() = %+v, want email %s", byID, customer.Email)
	}

	missingUser, err := f.users.GetByEmail(ctx, "missing@example.com")
	requireNoError(t, err)
	if missingUser != nil {
		t.Fatalf("GetByEmail(missing) = %+v, want nil", missingUser)
	}
	missingUser, err = f.users.GetByID(ctx, 999999)
	requireNoError(t, err)
	if missingUser != nil {
		t.Fatalf("GetByID(missing) = %+v, want nil", missingUser)
	}

	err = f.users.Create(ctx, &model.User{
		Name:         "Duplicate",
		Email:        customer.Email,
		PasswordHash: "hash",
		Role:         "customer",
	})
	if err == nil {
		t.Fatal("Create(duplicate user email) error = nil, want error")
	}

	provider := &model.ServiceProvider{
		UserID:       providerUser.ID,
		BusinessName: "Provider One",
	}
	requireNoError(t, f.providers.Create(ctx, provider))
	if provider.ID == 0 || provider.CreatedAt.IsZero() || provider.UpdatedAt.IsZero() {
		t.Fatalf("created provider missing generated fields: %+v", provider)
	}

	providerByID, err := f.providers.GetByID(ctx, provider.ID)
	requireNoError(t, err)
	if providerByID == nil || providerByID.Description != "" || providerByID.Address != "" || providerByID.Phone != "" || providerByID.Email != "" || providerByID.LogoURL != "" {
		t.Fatalf("GetByID(provider empty optionals) = %+v", providerByID)
	}

	providerByUser, err := f.providers.GetByUserID(ctx, providerUser.ID)
	requireNoError(t, err)
	if providerByUser == nil || providerByUser.ID != provider.ID {
		t.Fatalf("GetByUserID() = %+v, want provider %d", providerByUser, provider.ID)
	}

	provider.BusinessName = "Provider Updated"
	provider.Description = "Description"
	provider.Address = "Address"
	provider.Phone = "1234567890"
	provider.Email = "provider@example.com"
	provider.LogoURL = "https://example.com/logo.png"
	requireNoError(t, f.providers.Update(ctx, provider))

	updatedProvider, err := f.providers.GetByID(ctx, provider.ID)
	requireNoError(t, err)
	if updatedProvider == nil || updatedProvider.BusinessName != "Provider Updated" || updatedProvider.Description != "Description" || updatedProvider.LogoURL == "" {
		t.Fatalf("updated provider = %+v", updatedProvider)
	}

	missingProvider, err := f.providers.GetByID(ctx, 999999)
	requireNoError(t, err)
	if missingProvider != nil {
		t.Fatalf("GetByID(missing provider) = %+v, want nil", missingProvider)
	}
	missingProvider, err = f.providers.GetByUserID(ctx, 999999)
	requireNoError(t, err)
	if missingProvider != nil {
		t.Fatalf("GetByUserID(missing provider) = %+v, want nil", missingProvider)
	}
	err = f.providers.Create(ctx, &model.ServiceProvider{UserID: providerUser.ID, BusinessName: "Duplicate Provider"})
	if err == nil {
		t.Fatal("Create(duplicate provider user_id) error = nil, want error")
	}

	parent := f.createCategory(t, "Wellness", "", nil)
	parentID := parent.ID
	child := f.createCategory(t, "Massage", "Deep tissue", &parentID)

	category, err := f.categories.GetByID(ctx, child.ID)
	requireNoError(t, err)
	if category == nil || category.ParentID == nil || *category.ParentID != parent.ID || category.Description != "Deep tissue" {
		t.Fatalf("GetByID(category child) = %+v", category)
	}

	topCategory, err := f.categories.GetByID(ctx, parent.ID)
	requireNoError(t, err)
	if topCategory == nil || topCategory.ParentID != nil || topCategory.Description != "" {
		t.Fatalf("GetByID(category parent) = %+v", topCategory)
	}

	missingCategory, err := f.categories.GetByID(ctx, 999999)
	requireNoError(t, err)
	if missingCategory != nil {
		t.Fatalf("GetByID(missing category) = %+v, want nil", missingCategory)
	}

	categories, err := f.categories.List(ctx)
	requireNoError(t, err)
	if len(categories) != 2 || categories[0].ID != parent.ID || categories[1].ID != child.ID {
		t.Fatalf("List(categories) IDs = %v, want [%d %d]", categoryIDs(categories), parent.ID, child.ID)
	}

	tagA := f.createTag(t, "relax")
	tagB := f.createTag(t, "strength")

	tagByID, err := f.tags.GetByID(ctx, tagA.ID)
	requireNoError(t, err)
	if tagByID == nil || tagByID.Name != "relax" {
		t.Fatalf("GetByID(tag) = %+v", tagByID)
	}

	tagByName, err := f.tags.GetByName(ctx, "strength")
	requireNoError(t, err)
	if tagByName == nil || tagByName.ID != tagB.ID {
		t.Fatalf("GetByName(tag) = %+v, want ID %d", tagByName, tagB.ID)
	}

	missingTag, err := f.tags.GetByID(ctx, 999999)
	requireNoError(t, err)
	if missingTag != nil {
		t.Fatalf("GetByID(missing tag) = %+v, want nil", missingTag)
	}
	missingTag, err = f.tags.GetByName(ctx, "missing")
	requireNoError(t, err)
	if missingTag != nil {
		t.Fatalf("GetByName(missing tag) = %+v, want nil", missingTag)
	}

	tags, err := f.tags.List(ctx)
	requireNoError(t, err)
	if len(tags) != 2 || tags[0].ID != tagA.ID || tags[1].ID != tagB.ID {
		t.Fatalf("List(tags) IDs = %v, want [%d %d]", tagIDs(tags), tagA.ID, tagB.ID)
	}

	err = f.tags.Create(ctx, &model.Tag{Name: "relax"})
	if err == nil {
		t.Fatal("Create(duplicate tag name) error = nil, want error")
	}
}

func TestAssociationRepositories_ReplaceListAndRollback(t *testing.T) {
	f := newRepoFixture(t)
	ctx := f.ctx

	customer := f.createUser(t, "assoc-customer@example.com", "customer")
	providerUser := f.createUser(t, "assoc-provider@example.com", "provider")
	provider := f.createProvider(t, providerUser.ID, "Association Provider")
	category := f.createCategory(t, "Association Category", "", nil)
	service := f.createService(t, provider.ID, category.ID, "Association Service", "tag test", 100, 60, "active")
	tagA := f.createTag(t, "assoc-a")
	tagB := f.createTag(t, "assoc-b")

	requireNoError(t, f.serviceTags.ReplaceByServiceID(ctx, service.ID, []int64{tagB.ID, tagA.ID}))
	serviceTags, err := f.serviceTags.ListByServiceID(ctx, service.ID)
	requireNoError(t, err)
	if got, want := tagIDs(serviceTags), []int64{tagA.ID, tagB.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByServiceID() IDs = %v, want %v", got, want)
	}

	requireNoError(t, f.serviceTags.ReplaceByServiceID(ctx, service.ID, nil))
	serviceTags, err = f.serviceTags.ListByServiceID(ctx, service.ID)
	requireNoError(t, err)
	if len(serviceTags) != 0 {
		t.Fatalf("ListByServiceID(after clear) len = %d, want 0", len(serviceTags))
	}

	requireNoError(t, f.serviceTags.ReplaceByServiceID(ctx, service.ID, []int64{tagA.ID}))
	err = f.serviceTags.ReplaceByServiceID(ctx, service.ID, []int64{tagB.ID, 999999})
	if err == nil {
		t.Fatal("ReplaceByServiceID(invalid tag) error = nil, want error")
	}
	serviceTags, err = f.serviceTags.ListByServiceID(ctx, service.ID)
	requireNoError(t, err)
	if got, want := tagIDs(serviceTags), []int64{tagA.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByServiceID(after rollback) IDs = %v, want %v", got, want)
	}

	requireNoError(t, f.interests.ReplaceByUserID(ctx, customer.ID, []int64{tagB.ID, tagA.ID}))
	interests, err := f.interests.ListByUserID(ctx, customer.ID)
	requireNoError(t, err)
	if got, want := tagIDs(interests), []int64{tagA.ID, tagB.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByUserID() IDs = %v, want %v", got, want)
	}

	requireNoError(t, f.interests.ReplaceByUserID(ctx, customer.ID, nil))
	interests, err = f.interests.ListByUserID(ctx, customer.ID)
	requireNoError(t, err)
	if len(interests) != 0 {
		t.Fatalf("ListByUserID(after clear) len = %d, want 0", len(interests))
	}

	requireNoError(t, f.interests.ReplaceByUserID(ctx, customer.ID, []int64{tagA.ID}))
	err = f.interests.ReplaceByUserID(ctx, customer.ID, []int64{tagB.ID, 999999})
	if err == nil {
		t.Fatal("ReplaceByUserID(invalid tag) error = nil, want error")
	}
	interests, err = f.interests.ListByUserID(ctx, customer.ID)
	requireNoError(t, err)
	if got, want := tagIDs(interests), []int64{tagA.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByUserID(after rollback) IDs = %v, want %v", got, want)
	}
}

func TestServiceRepository_CreateGetListUpdateAndStatus(t *testing.T) {
	f := newRepoFixture(t)
	ctx := f.ctx

	providerUserA := f.createUser(t, "service-provider-a@example.com", "provider")
	providerUserB := f.createUser(t, "service-provider-b@example.com", "provider")
	providerA := f.createProvider(t, providerUserA.ID, "Provider A")
	providerB := f.createProvider(t, providerUserB.ID, "Provider B")
	categoryA := f.createCategory(t, "Beauty", "", nil)
	categoryB := f.createCategory(t, "Health", "", nil)

	serviceA := f.createService(t, providerA.ID, categoryA.ID, "Massage Therapy", "Deep relaxation", 99, 60, "active")
	serviceB := f.createService(t, providerA.ID, categoryB.ID, "Dental Cleaning", "", 180, 45, "inactive")
	serviceC := f.createService(t, providerB.ID, categoryA.ID, "Personal Training", "Strength plan", 150, 30, "active")

	if _, err := f.pool.Exec(ctx, `UPDATE services SET avg_rating = $1 WHERE id = $2`, 4.8, serviceA.ID); err != nil {
		t.Fatalf("seed serviceA rating: %v", err)
	}
	if _, err := f.pool.Exec(ctx, `UPDATE services SET avg_rating = $1 WHERE id = $2`, 3.2, serviceB.ID); err != nil {
		t.Fatalf("seed serviceB rating: %v", err)
	}
	if _, err := f.pool.Exec(ctx, `UPDATE services SET avg_rating = $1 WHERE id = $2`, 4.1, serviceC.ID); err != nil {
		t.Fatalf("seed serviceC rating: %v", err)
	}

	got, err := f.services.GetByID(ctx, serviceB.ID)
	requireNoError(t, err)
	if got == nil || got.Description != "" || got.ImageURL != "" || got.Status != "inactive" {
		t.Fatalf("GetByID(serviceB) = %+v", got)
	}

	missingService, err := f.services.GetByID(ctx, 999999)
	requireNoError(t, err)
	if missingService != nil {
		t.Fatalf("GetByID(missing service) = %+v, want nil", missingService)
	}

	view, err := f.services.GetViewByID(ctx, serviceA.ID)
	requireNoError(t, err)
	if view == nil || view.Provider.ID != providerA.ID || view.Category.ID != categoryA.ID || view.Provider.BusinessName != "Provider A" {
		t.Fatalf("GetViewByID() = %+v", view)
	}

	missingView, err := f.services.GetViewByID(ctx, 999999)
	requireNoError(t, err)
	if missingView != nil {
		t.Fatalf("GetViewByID(missing) = %+v, want nil", missingView)
	}

	items, total, err := f.services.List(ctx, model.ServiceFilter{Page: 1, PageSize: 10})
	requireNoError(t, err)
	if total != 3 || len(items) != 3 {
		t.Fatalf("List(all) total/len = %d/%d, want 3/3", total, len(items))
	}

	items, total, err = f.services.List(ctx, model.ServiceFilter{Keyword: "massage", Page: 1, PageSize: 10})
	requireNoError(t, err)
	if total != 1 || len(items) != 1 || items[0].ID != serviceA.ID {
		t.Fatalf("List(keyword) total/items = %d/%v, want service %d", total, serviceViewIDs(items), serviceA.ID)
	}

	items, total, err = f.services.List(ctx, model.ServiceFilter{CategoryID: &categoryA.ID, Page: 1, PageSize: 10})
	requireNoError(t, err)
	if total != 2 || len(items) != 2 {
		t.Fatalf("List(category) total/len = %d/%d, want 2/2", total, len(items))
	}

	items, total, err = f.services.List(ctx, model.ServiceFilter{ProviderID: &providerA.ID, Page: 1, PageSize: 10})
	requireNoError(t, err)
	if total != 2 || len(items) != 2 {
		t.Fatalf("List(provider) total/len = %d/%d, want 2/2", total, len(items))
	}

	minPrice := 120.0
	maxPrice := 160.0
	items, total, err = f.services.List(ctx, model.ServiceFilter{MinPrice: &minPrice, MaxPrice: &maxPrice, Page: 1, PageSize: 10})
	requireNoError(t, err)
	if total != 1 || len(items) != 1 || items[0].ID != serviceC.ID {
		t.Fatalf("List(price range) total/items = %d/%v, want service %d", total, serviceViewIDs(items), serviceC.ID)
	}

	items, total, err = f.services.List(ctx, model.ServiceFilter{Status: "active", Page: 1, PageSize: 10})
	requireNoError(t, err)
	if total != 2 || len(items) != 2 {
		t.Fatalf("List(status active) total/len = %d/%d, want 2/2", total, len(items))
	}

	items, total, err = f.services.List(ctx, model.ServiceFilter{SortBy: "price", SortOrder: "asc", Page: 1, PageSize: 2})
	requireNoError(t, err)
	if total != 3 || len(items) != 2 || items[0].ID != serviceA.ID || items[1].ID != serviceC.ID {
		t.Fatalf("List(price asc page 1) total/items = %d/%v", total, serviceViewIDs(items))
	}

	items, total, err = f.services.List(ctx, model.ServiceFilter{SortBy: "price", SortOrder: "desc", Page: 2, PageSize: 1})
	requireNoError(t, err)
	if total != 3 || len(items) != 1 || items[0].ID != serviceC.ID {
		t.Fatalf("List(price desc page 2) total/items = %d/%v, want service %d", total, serviceViewIDs(items), serviceC.ID)
	}

	items, total, err = f.services.List(ctx, model.ServiceFilter{SortBy: "rating", SortOrder: "desc", Page: 1, PageSize: 1})
	requireNoError(t, err)
	if total != 3 || len(items) != 1 || items[0].ID != serviceA.ID {
		t.Fatalf("List(rating desc) total/items = %d/%v, want service %d", total, serviceViewIDs(items), serviceA.ID)
	}

	serviceB.CategoryID = categoryA.ID
	serviceB.Title = "Dental Cleaning Updated"
	serviceB.Description = ""
	serviceB.Price = 175
	serviceB.DurationMinutes = 50
	serviceB.ImageURL = ""
	requireNoError(t, f.services.Update(ctx, serviceB))

	updated, err := f.services.GetByID(ctx, serviceB.ID)
	requireNoError(t, err)
	if updated == nil || updated.Title != "Dental Cleaning Updated" || updated.Description != "" || updated.ImageURL != "" || updated.CategoryID != categoryA.ID {
		t.Fatalf("updated service = %+v", updated)
	}

	err = f.services.Update(ctx, &model.Service{
		ID:              999999,
		CategoryID:      categoryA.ID,
		Title:           "Missing",
		Price:           1,
		DurationMinutes: 1,
	})
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("Update(missing) error = %v, want pgx.ErrNoRows", err)
	}

	requireNoError(t, f.services.UpdateStatus(ctx, serviceA.ID, "inactive"))
	statusUpdated, err := f.services.GetByID(ctx, serviceA.ID)
	requireNoError(t, err)
	if statusUpdated == nil || statusUpdated.Status != "inactive" {
		t.Fatalf("UpdateStatus(serviceA) = %+v, want inactive", statusUpdated)
	}

	requireNoError(t, f.services.UpdateStatus(ctx, 999999, "active"))
}

func TestReservationRepository_CreateQueryStatusAndConflict(t *testing.T) {
	f := newRepoFixture(t)
	ctx := f.ctx

	customerA := f.createUser(t, "reservation-customer-a@example.com", "customer")
	customerB := f.createUser(t, "reservation-customer-b@example.com", "customer")
	providerUserA := f.createUser(t, "reservation-provider-a@example.com", "provider")
	providerUserB := f.createUser(t, "reservation-provider-b@example.com", "provider")
	providerA := f.createProvider(t, providerUserA.ID, "Reservation Provider A")
	providerB := f.createProvider(t, providerUserB.ID, "Reservation Provider B")
	category := f.createCategory(t, "Reservation Category", "", nil)
	serviceA := f.createService(t, providerA.ID, category.ID, "Reservation Service A", "", 100, 60, "active")
	serviceB := f.createService(t, providerB.ID, category.ID, "Reservation Service B", "", 120, 60, "active")

	base := time.Date(2026, 7, 10, 10, 0, 0, 0, time.UTC)
	reservationA := f.createReservation(t, customerA.ID, serviceA.ID, base, "", "")
	reservationB := f.createReservation(t, customerA.ID, serviceA.ID, base.Add(2*time.Hour), "confirmed", "bring water")
	reservationC := f.createReservation(t, customerB.ID, serviceA.ID, base.Add(4*time.Hour), "completed", "done")
	reservationD := f.createReservation(t, customerA.ID, serviceB.ID, base.Add(6*time.Hour), "pending", "other provider")

	if reservationA.Status != "pending" || reservationA.Note != "" || reservationA.CreatedAt.IsZero() || reservationA.UpdatedAt.IsZero() {
		t.Fatalf("default reservation = %+v", reservationA)
	}

	err := f.reservations.Create(ctx, &model.Reservation{
		UserID:    customerA.ID,
		ServiceID: serviceA.ID,
		StartTime: base,
		EndTime:   base.Add(time.Hour),
	})
	if !errors.Is(err, repository.ErrReservationTimeConflict) {
		t.Fatalf("Create(duplicate service/start) error = %v, want ErrReservationTimeConflict", err)
	}

	got, err := f.reservations.GetByID(ctx, reservationB.ID)
	requireNoError(t, err)
	if got == nil || got.Note != "bring water" || got.Status != "confirmed" {
		t.Fatalf("GetByID(reservationB) = %+v", got)
	}

	missing, err := f.reservations.GetByID(ctx, 999999)
	requireNoError(t, err)
	if missing != nil {
		t.Fatalf("GetByID(missing reservation) = %+v, want nil", missing)
	}

	got, err = f.reservations.GetByIDForUser(ctx, reservationB.ID, customerA.ID)
	requireNoError(t, err)
	if got == nil || got.ID != reservationB.ID {
		t.Fatalf("GetByIDForUser() = %+v, want %d", got, reservationB.ID)
	}
	got, err = f.reservations.GetByIDForUser(ctx, reservationB.ID, customerB.ID)
	requireNoError(t, err)
	if got != nil {
		t.Fatalf("GetByIDForUser(wrong user) = %+v, want nil", got)
	}

	got, err = f.reservations.GetByIDForProvider(ctx, reservationB.ID, providerA.ID)
	requireNoError(t, err)
	if got == nil || got.ID != reservationB.ID {
		t.Fatalf("GetByIDForProvider() = %+v, want %d", got, reservationB.ID)
	}
	got, err = f.reservations.GetByIDForProvider(ctx, reservationB.ID, providerB.ID)
	requireNoError(t, err)
	if got != nil {
		t.Fatalf("GetByIDForProvider(wrong provider) = %+v, want nil", got)
	}

	userReservations, err := f.reservations.ListByUserID(ctx, customerA.ID, "", 10, 0)
	requireNoError(t, err)
	if got, want := reservationIDs(userReservations), []int64{reservationD.ID, reservationB.ID, reservationA.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByUserID(all) IDs = %v, want %v", got, want)
	}

	pendingReservations, err := f.reservations.ListByUserID(ctx, customerA.ID, "pending", 10, 0)
	requireNoError(t, err)
	if got, want := reservationIDs(pendingReservations), []int64{reservationD.ID, reservationA.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByUserID(pending) IDs = %v, want %v", got, want)
	}

	pagedReservations, err := f.reservations.ListByUserID(ctx, customerA.ID, "", 1, 1)
	requireNoError(t, err)
	if got, want := reservationIDs(pagedReservations), []int64{reservationB.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByUserID(page) IDs = %v, want %v", got, want)
	}

	providerReservations, err := f.reservations.ListByProviderID(ctx, providerA.ID, "", 10, 0)
	requireNoError(t, err)
	if got, want := reservationIDs(providerReservations), []int64{reservationC.ID, reservationB.ID, reservationA.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByProviderID(all) IDs = %v, want %v", got, want)
	}

	completedReservations, err := f.reservations.ListByProviderID(ctx, providerA.ID, "completed", 10, 0)
	requireNoError(t, err)
	if got, want := reservationIDs(completedReservations), []int64{reservationC.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByProviderID(completed) IDs = %v, want %v", got, want)
	}

	defaultLimitReservations, err := f.reservations.ListByProviderID(ctx, providerA.ID, "", 0, -1)
	requireNoError(t, err)
	if len(defaultLimitReservations) != 3 {
		t.Fatalf("ListByProviderID(default limit) len = %d, want 3", len(defaultLimitReservations))
	}

	cancelled, err := f.reservations.UpdateStatus(ctx, reservationA.ID, "cancelled")
	requireNoError(t, err)
	if cancelled == nil || cancelled.Status != "cancelled" {
		t.Fatalf("UpdateStatus() = %+v, want cancelled", cancelled)
	}

	missing, err = f.reservations.UpdateStatus(ctx, 999999, "cancelled")
	requireNoError(t, err)
	if missing != nil {
		t.Fatalf("UpdateStatus(missing) = %+v, want nil", missing)
	}
}

func TestReviewRepository_CreateAndLists(t *testing.T) {
	f := newRepoFixture(t)
	ctx := f.ctx

	customerA := f.createUser(t, "review-customer-a@example.com", "customer")
	customerB := f.createUser(t, "review-customer-b@example.com", "customer")
	providerUserA := f.createUser(t, "review-provider-a@example.com", "provider")
	providerUserB := f.createUser(t, "review-provider-b@example.com", "provider")
	providerA := f.createProvider(t, providerUserA.ID, "Review Provider A")
	providerB := f.createProvider(t, providerUserB.ID, "Review Provider B")
	category := f.createCategory(t, "Review Category", "", nil)
	serviceA := f.createService(t, providerA.ID, category.ID, "Review Service A", "", 100, 60, "active")
	serviceB := f.createService(t, providerB.ID, category.ID, "Review Service B", "", 120, 60, "active")

	base := time.Date(2026, 7, 11, 9, 0, 0, 0, time.UTC)
	reservationA := f.createReservation(t, customerA.ID, serviceA.ID, base, "completed", "")
	reservationB := f.createReservation(t, customerB.ID, serviceA.ID, base.Add(2*time.Hour), "completed", "")
	reservationC := f.createReservation(t, customerA.ID, serviceB.ID, base.Add(4*time.Hour), "completed", "")

	reviewA := f.createReview(t, reservationA.ID, customerA.ID, serviceA.ID, 5, "")
	reviewB := f.createReview(t, reservationB.ID, customerB.ID, serviceA.ID, 4, "solid")
	reviewC := f.createReview(t, reservationC.ID, customerA.ID, serviceB.ID, 3, "other")

	got, err := f.reviews.GetByID(ctx, reviewA.ID)
	requireNoError(t, err)
	if got == nil || got.Comment != "" || got.Rating != 5 {
		t.Fatalf("GetByID(reviewA) = %+v", got)
	}

	got, err = f.reviews.GetByReservationID(ctx, reservationB.ID)
	requireNoError(t, err)
	if got == nil || got.ID != reviewB.ID {
		t.Fatalf("GetByReservationID() = %+v, want %d", got, reviewB.ID)
	}

	missing, err := f.reviews.GetByID(ctx, 999999)
	requireNoError(t, err)
	if missing != nil {
		t.Fatalf("GetByID(missing review) = %+v, want nil", missing)
	}
	missing, err = f.reviews.GetByReservationID(ctx, 999999)
	requireNoError(t, err)
	if missing != nil {
		t.Fatalf("GetByReservationID(missing review) = %+v, want nil", missing)
	}

	serviceReviews, err := f.reviews.ListByServiceID(ctx, serviceA.ID, 10, 0)
	requireNoError(t, err)
	if got, want := reviewIDs(serviceReviews), []int64{reviewB.ID, reviewA.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByServiceID() IDs = %v, want %v", got, want)
	}

	serviceReviews, err = f.reviews.ListByServiceID(ctx, serviceA.ID, 1, 1)
	requireNoError(t, err)
	if got, want := reviewIDs(serviceReviews), []int64{reviewA.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByServiceID(page) IDs = %v, want %v", got, want)
	}

	userReviews, err := f.reviews.ListByUserID(ctx, customerA.ID, 10, 0)
	requireNoError(t, err)
	if got, want := reviewIDs(userReviews), []int64{reviewC.ID, reviewA.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByUserID() IDs = %v, want %v", got, want)
	}

	providerReviews, err := f.reviews.ListByProviderID(ctx, providerA.ID, 10, 0)
	requireNoError(t, err)
	if got, want := reviewIDs(providerReviews), []int64{reviewB.ID, reviewA.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByProviderID(providerA) IDs = %v, want %v", got, want)
	}

	providerReviews, err = f.reviews.ListByProviderID(ctx, providerB.ID, 10, 0)
	requireNoError(t, err)
	if got, want := reviewIDs(providerReviews), []int64{reviewC.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByProviderID(providerB) IDs = %v, want %v", got, want)
	}
}

func TestNotificationRepository_CreateListCountAndMarkRead(t *testing.T) {
	f := newRepoFixture(t)
	ctx := f.ctx

	userA := f.createUser(t, "notification-user-a@example.com", "customer")
	userB := f.createUser(t, "notification-user-b@example.com", "customer")

	noticeA := f.createNotification(t, userA.ID, "First", "first content", "system", false)
	noticeB := f.createNotification(t, userA.ID, "Second", "second content", "reservation_confirmed", true)
	noticeC := f.createNotification(t, userA.ID, "Third", "third content", "system", false)
	noticeD := f.createNotification(t, userB.ID, "Other", "other content", "system", false)

	got, err := f.notifications.GetByID(ctx, noticeA.ID)
	requireNoError(t, err)
	if got == nil || got.Title != "First" || got.IsRead {
		t.Fatalf("GetByID(noticeA) = %+v", got)
	}

	missing, err := f.notifications.GetByID(ctx, 999999)
	requireNoError(t, err)
	if missing != nil {
		t.Fatalf("GetByID(missing notification) = %+v, want nil", missing)
	}

	allForUser, err := f.notifications.ListByUserID(ctx, userA.ID, nil, 10, 0)
	requireNoError(t, err)
	if got, want := notificationIDs(allForUser), []int64{noticeC.ID, noticeB.ID, noticeA.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByUserID(nil) IDs = %v, want %v", got, want)
	}

	allForUser, err = f.notifications.ListByUserID(ctx, userA.ID, nil, 1, 1)
	requireNoError(t, err)
	if got, want := notificationIDs(allForUser), []int64{noticeB.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByUserID(page) IDs = %v, want %v", got, want)
	}

	isRead := true
	readForUser, err := f.notifications.ListByUserID(ctx, userA.ID, &isRead, 10, 0)
	requireNoError(t, err)
	if got, want := notificationIDs(readForUser), []int64{noticeB.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByUserID(read) IDs = %v, want %v", got, want)
	}

	isRead = false
	unreadForUser, err := f.notifications.ListByUserID(ctx, userA.ID, &isRead, 10, 0)
	requireNoError(t, err)
	if got, want := notificationIDs(unreadForUser), []int64{noticeC.ID, noticeA.ID}; !sameInt64s(got, want) {
		t.Fatalf("ListByUserID(unread) IDs = %v, want %v", got, want)
	}

	unreadCount, err := f.notifications.CountUnread(ctx, userA.ID)
	requireNoError(t, err)
	if unreadCount != 2 {
		t.Fatalf("CountUnread(userA) = %d, want 2", unreadCount)
	}

	updated, err := f.notifications.MarkRead(ctx, noticeA.ID, userA.ID)
	requireNoError(t, err)
	if updated == nil || !updated.IsRead {
		t.Fatalf("MarkRead() = %+v, want read notification", updated)
	}

	crossUser, err := f.notifications.MarkRead(ctx, noticeD.ID, userA.ID)
	requireNoError(t, err)
	if crossUser != nil {
		t.Fatalf("MarkRead(cross user) = %+v, want nil", crossUser)
	}

	missing, err = f.notifications.MarkRead(ctx, 999999, userA.ID)
	requireNoError(t, err)
	if missing != nil {
		t.Fatalf("MarkRead(missing) = %+v, want nil", missing)
	}

	unreadCount, err = f.notifications.CountUnread(ctx, userA.ID)
	requireNoError(t, err)
	if unreadCount != 1 {
		t.Fatalf("CountUnread(after MarkRead) = %d, want 1", unreadCount)
	}

	updatedCount, err := f.notifications.MarkAllRead(ctx, userA.ID)
	requireNoError(t, err)
	if updatedCount != 1 {
		t.Fatalf("MarkAllRead(userA) = %d, want 1", updatedCount)
	}

	updatedCount, err = f.notifications.MarkAllRead(ctx, userA.ID)
	requireNoError(t, err)
	if updatedCount != 0 {
		t.Fatalf("MarkAllRead(userA again) = %d, want 0", updatedCount)
	}

	unreadCount, err = f.notifications.CountUnread(ctx, userB.ID)
	requireNoError(t, err)
	if unreadCount != 1 {
		t.Fatalf("CountUnread(userB) = %d, want 1", unreadCount)
	}
}

func newRepoFixture(t *testing.T) *repoFixture {
	t.Helper()

	url := testDatabaseURL(t)
	ctx := context.Background()
	schema := fmt.Sprintf("ors_repo_test_%d_%d", os.Getpid(), time.Now().UnixNano())

	controlPool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatalf("create control pool: %v", err)
	}
	if _, err := controlPool.Exec(ctx, "CREATE SCHEMA "+quoteIdentifier(schema)); err != nil {
		controlPool.Close()
		t.Fatalf("create schema %s: %v", schema, err)
	}
	t.Cleanup(func() {
		if _, err := controlPool.Exec(context.Background(), "DROP SCHEMA IF EXISTS "+quoteIdentifier(schema)+" CASCADE"); err != nil {
			t.Errorf("drop schema %s: %v", schema, err)
		}
		controlPool.Close()
	})

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatalf("parse test database url: %v", err)
	}
	if cfg.ConnConfig.RuntimeParams == nil {
		cfg.ConnConfig.RuntimeParams = make(map[string]string)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = schema

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("create schema pool: %v", err)
	}
	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping schema pool: %v", err)
	}
	runUpMigrations(t, ctx, pool)

	return &repoFixture{
		ctx:           ctx,
		pool:          pool,
		users:         NewUserRepo(pool),
		providers:     NewServiceProviderRepo(pool),
		categories:    NewCategoryRepo(pool),
		tags:          NewTagRepo(pool),
		serviceTags:   NewServiceTagRepo(pool),
		interests:     NewUserInterestRepo(pool),
		services:      NewServiceRepo(pool),
		reservations:  NewReservationRepo(pool),
		reviews:       NewReviewRepo(pool),
		notifications: NewNotificationRepo(pool),
	}
}

func testDatabaseURL(t *testing.T) string {
	t.Helper()

	url := os.Getenv(repositoryTestDatabaseURLEnv)
	if url == "" {
		t.Skipf("%s is not set; skipping PostgreSQL repository integration test", repositoryTestDatabaseURLEnv)
	}
	return url
}

func runUpMigrations(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve test file path")
	}

	migrationsDir := filepath.Join(filepath.Dir(file), "..", "..", "..", "migrations")
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
	if err != nil {
		t.Fatalf("glob migrations: %v", err)
	}

	filtered := make([]string, 0, 10)
	for _, path := range files {
		base := filepath.Base(path)
		if base >= "001_" && base <= "010_zzzzzzzzzzzzzzzz.up.sql" {
			filtered = append(filtered, path)
		}
	}
	sort.Strings(filtered)
	if len(filtered) != 10 {
		t.Fatalf("found %d migrations in 001..010, want 10: %v", len(filtered), filtered)
	}

	for _, path := range filtered {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read migration %s: %v", path, err)
		}
		for _, stmt := range strings.Split(string(content), ";") {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := pool.Exec(ctx, stmt); err != nil {
				t.Fatalf("execute migration %s statement %q: %v", filepath.Base(path), stmt, err)
			}
		}
	}
}

func quoteIdentifier(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}

func (f *repoFixture) createUser(t *testing.T, email, role string) *model.User {
	t.Helper()
	if role == "" {
		role = "customer"
	}
	user := &model.User{
		Name:         "User " + email,
		Email:        email,
		PasswordHash: "hash",
		Role:         role,
	}
	requireNoError(t, f.users.Create(f.ctx, user))
	return user
}

func (f *repoFixture) createProvider(t *testing.T, userID int64, businessName string) *model.ServiceProvider {
	t.Helper()
	provider := &model.ServiceProvider{
		UserID:       userID,
		BusinessName: businessName,
	}
	requireNoError(t, f.providers.Create(f.ctx, provider))
	return provider
}

func (f *repoFixture) createCategory(t *testing.T, name, description string, parentID *int64) *model.Category {
	t.Helper()
	category := &model.Category{
		Name:        name,
		Description: description,
		ParentID:    parentID,
	}
	requireNoError(t, f.categories.Create(f.ctx, category))
	return category
}

func (f *repoFixture) createTag(t *testing.T, name string) *model.Tag {
	t.Helper()
	tag := &model.Tag{Name: name}
	requireNoError(t, f.tags.Create(f.ctx, tag))
	return tag
}

func (f *repoFixture) createService(t *testing.T, providerID, categoryID int64, title, description string, price float64, duration int, status string) *model.Service {
	t.Helper()
	if status == "" {
		status = "active"
	}
	service := &model.Service{
		ProviderID:      providerID,
		CategoryID:      categoryID,
		Title:           title,
		Description:     description,
		Price:           price,
		DurationMinutes: duration,
		Status:          status,
	}
	requireNoError(t, f.services.Create(f.ctx, service))
	return service
}

func (f *repoFixture) createReservation(t *testing.T, userID, serviceID int64, start time.Time, status, note string) *model.Reservation {
	t.Helper()
	reservation := &model.Reservation{
		UserID:    userID,
		ServiceID: serviceID,
		StartTime: start,
		EndTime:   start.Add(time.Hour),
		Status:    status,
		Note:      note,
	}
	requireNoError(t, f.reservations.Create(f.ctx, reservation))
	return reservation
}

func (f *repoFixture) createReview(t *testing.T, reservationID, userID, serviceID int64, rating int16, comment string) *model.Review {
	t.Helper()
	review := &model.Review{
		ReservationID: reservationID,
		UserID:        userID,
		ServiceID:     serviceID,
		Rating:        rating,
		Comment:       comment,
	}
	requireNoError(t, f.reviews.Create(f.ctx, review))
	return review
}

func (f *repoFixture) createNotification(t *testing.T, userID int64, title, content, notificationType string, isRead bool) *model.Notification {
	t.Helper()
	notification := &model.Notification{
		UserID:  userID,
		Title:   title,
		Content: content,
		Type:    notificationType,
		IsRead:  isRead,
	}
	requireNoError(t, f.notifications.Create(f.ctx, notification))
	return notification
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func categoryIDs(categories []*model.Category) []int64 {
	ids := make([]int64, 0, len(categories))
	for _, category := range categories {
		ids = append(ids, category.ID)
	}
	return ids
}

func tagIDs(tags []*model.Tag) []int64 {
	ids := make([]int64, 0, len(tags))
	for _, tag := range tags {
		ids = append(ids, tag.ID)
	}
	return ids
}

func serviceViewIDs(services []*model.ServiceView) []int64 {
	ids := make([]int64, 0, len(services))
	for _, service := range services {
		ids = append(ids, service.ID)
	}
	return ids
}

func reservationIDs(reservations []*model.Reservation) []int64 {
	ids := make([]int64, 0, len(reservations))
	for _, reservation := range reservations {
		ids = append(ids, reservation.ID)
	}
	return ids
}

func reviewIDs(reviews []*model.Review) []int64 {
	ids := make([]int64, 0, len(reviews))
	for _, review := range reviews {
		ids = append(ids, review.ID)
	}
	return ids
}

func notificationIDs(notifications []*model.Notification) []int64 {
	ids := make([]int64, 0, len(notifications))
	for _, notification := range notifications {
		ids = append(ids, notification.ID)
	}
	return ids
}

func sameInt64s(got, want []int64) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
