package handler

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
	"github.com/MamangRust/monolith-point-of-sale-user/repository"
)

func mapPaginationMeta(meta *pbcommon.PaginationMeta) *pbcommon.PaginationMeta {
	if meta == nil {
		return nil
	}
	return &pbcommon.PaginationMeta{
		CurrentPage:  meta.CurrentPage,
		PageSize:     meta.PageSize,
		TotalPages:   meta.TotalPages,
		TotalRecords: meta.TotalRecords,
	}
}

func formatOptionalTime(t *time.Time) string {
	return convert.FormatTimePtr(t)
}

func mapUserToProto(user *salemongo.User) *pbusers.UserResponse {
	if user == nil {
		return nil
	}
	createdAtStr := convert.FormatTimePtr(user.CreatedAt)
	updatedAtStr := convert.FormatTimePtr(user.UpdatedAt)
	return &pbusers.UserResponse{
		Id:        user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
	}
}

// mapUserWithPasswordToProto is the auth-only projection: it carries the
// password hash and must never be served through the public gateway.
func mapUserWithPasswordToProto(user *salemongo.User) *pbusers.UserResponseWithPassword {
	if user == nil {
		return nil
	}
	return &pbusers.UserResponseWithPassword{
		Id:        user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: convert.FormatTimePtr(user.CreatedAt),
		UpdatedAt: convert.FormatTimePtr(user.UpdatedAt),
	}
}

func mapGetUsersRowToProto(user *repository.UserRow) *pbusers.UserResponse {
	if user == nil {
		return nil
	}
	createdAtStr := convert.FormatTimePtr(user.CreatedAt)
	updatedAtStr := convert.FormatTimePtr(user.UpdatedAt)
	return &pbusers.UserResponse{
		Id:        user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
	}
}

func mapGetUsersRowsToProto(users []*repository.UserRow) []*pbusers.UserResponse {
	var responseUsers []*pbusers.UserResponse
	for _, u := range users {
		responseUsers = append(responseUsers, mapGetUsersRowToProto(u))
	}
	return responseUsers
}

func mapUserDeleteAtToProto(user *salemongo.User) *pbusers.UserResponseDeleteAt {
	if user == nil {
		return nil
	}
	createdAtStr := convert.FormatTimePtr(user.CreatedAt)
	updatedAtStr := convert.FormatTimePtr(user.UpdatedAt)
	var deletedAt *wrapperspb.StringValue
	if user.DeletedAt != nil {
		deletedAt = convert.TimeToWrappers(user.DeletedAt)
	}
	return &pbusers.UserResponseDeleteAt{
		Id:        user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
		DeletedAt: deletedAt,
	}
}

func mapActiveUserToProto(user *repository.UserActiveRow) *pbusers.UserResponseDeleteAt {
	if user == nil {
		return nil
	}
	createdAtStr := convert.FormatTimePtr(user.CreatedAt)
	updatedAtStr := convert.FormatTimePtr(user.UpdatedAt)
	deletedAtVal := convert.TimeToWrappers(user.DeletedAt)
	return &pbusers.UserResponseDeleteAt{
		Id:        user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
		DeletedAt: deletedAtVal,
	}
}

func mapActiveUsersToProto(users []*repository.UserActiveRow) []*pbusers.UserResponseDeleteAt {
	var responseUsers []*pbusers.UserResponseDeleteAt
	for _, u := range users {
		responseUsers = append(responseUsers, mapActiveUserToProto(u))
	}
	return responseUsers
}

func mapTrashedUserToProto(user *repository.UserTrashedRow) *pbusers.UserResponseDeleteAt {
	if user == nil {
		return nil
	}
	createdAtStr := convert.FormatTimePtr(user.CreatedAt)
	updatedAtStr := convert.FormatTimePtr(user.UpdatedAt)
	deletedAtVal := convert.TimeToWrappers(user.DeletedAt)
	return &pbusers.UserResponseDeleteAt{
		Id:        user.UserID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
		DeletedAt: deletedAtVal,
	}
}

func mapTrashedUsersToProto(users []*repository.UserTrashedRow) []*pbusers.UserResponseDeleteAt {
	var responseUsers []*pbusers.UserResponseDeleteAt
	for _, u := range users {
		responseUsers = append(responseUsers, mapTrashedUserToProto(u))
	}
	return responseUsers
}
