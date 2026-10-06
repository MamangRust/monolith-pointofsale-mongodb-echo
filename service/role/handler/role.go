package handler

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-role/repository"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
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

func mapResponseRole(role *salemongo.Role) *pbroles.RoleResponse {
	if role == nil {
		return nil
	}
	createdAtStr := convert.FormatTimePtr(role.CreatedAt)
	updatedAtStr := convert.FormatTimePtr(role.UpdatedAt)
	return &pbroles.RoleResponse{
		Id:        role.RoleID,
		Name:      role.RoleName,
		CreatedAt: createdAtStr,
		UpdatedAt: updatedAtStr,
	}
}

func mapResponsesRole(roles []*repository.RoleRow) []*pbroles.RoleResponse {
	var responseRoles []*pbroles.RoleResponse
	for _, role := range roles {
		if role == nil {
			continue
		}
		responseRoles = append(responseRoles, &pbroles.RoleResponse{
			Id:        role.RoleID,
			Name:      role.RoleName,
			CreatedAt: convert.FormatTimePtr(role.CreatedAt),
			UpdatedAt: convert.FormatTimePtr(role.UpdatedAt),
		})
	}
	return responseRoles
}

func mapResponsesRoleFromDB(roles []*salemongo.Role) []*pbroles.RoleResponse {
	var responseRoles []*pbroles.RoleResponse
	for _, role := range roles {
		responseRoles = append(responseRoles, mapResponseRole(role))
	}
	return responseRoles
}

func formatOptionalTime(t *time.Time) string {
	return convert.FormatTimePtr(t)
}

func mapResponsesRoleFromActive(roles []*repository.RoleActiveRow) []*pbroles.RoleResponseDeleteAt {
	var responseRoles []*pbroles.RoleResponseDeleteAt
	for _, role := range roles {
		if role == nil {
			continue
		}
		responseRoles = append(responseRoles, &pbroles.RoleResponseDeleteAt{
			Id:        role.RoleID,
			Name:      role.RoleName,
			CreatedAt: convert.FormatTimePtr(role.CreatedAt),
			UpdatedAt: convert.FormatTimePtr(role.UpdatedAt),
			DeletedAt: formatOptionalTime(role.DeletedAt),
		})
	}
	return responseRoles
}

func mapResponsesRoleFromTrashed(roles []*repository.RoleTrashedRow) []*pbroles.RoleResponseDeleteAt {
	var responseRoles []*pbroles.RoleResponseDeleteAt
	for _, role := range roles {
		if role == nil {
			continue
		}
		responseRoles = append(responseRoles, &pbroles.RoleResponseDeleteAt{
			Id:        role.RoleID,
			Name:      role.RoleName,
			CreatedAt: convert.FormatTimePtr(role.CreatedAt),
			UpdatedAt: convert.FormatTimePtr(role.UpdatedAt),
			DeletedAt: formatOptionalTime(role.DeletedAt),
		})
	}
	return responseRoles
}
