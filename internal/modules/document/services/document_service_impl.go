package services

import (
	"context"
	"fmt"
	"time"

	"testcase/internal/modules/document/dto"
	"testcase/internal/modules/document/entities"
	"testcase/internal/modules/document/repositories"
	userEntity "testcase/internal/modules/user/entities"
	"testcase/internal/utils"
)

type documentServiceImpl struct {
	repo repositories.DocumentRepo
}

func NewDocumentService(repo repositories.DocumentRepo) DocumentService {
	return &documentServiceImpl{
		repo: repo,
	}
}

func (d *documentServiceImpl) CreateDocument(ctx context.Context, input *dto.CreateDocumentDTO) (*entities.Document, error) {
	document := &entities.Document{
		Title:           input.Title,
		Status:          entities.StatusPending,
		CurrentApprover: 1,
		CreatedAt:       time.Now(),
	}

	if err := d.repo.CreateDocument(ctx, document); err != nil {
		return nil, utils.NewAppError(utils.ErrInternalServer, fmt.Errorf("failed to create document: %w", err))
	}

	return document, nil
}

func (d *documentServiceImpl) FindById(ctx context.Context, id string) (*entities.Document, error) {
	if id == "" {
		return nil, utils.NewAppError(utils.ErrInvalidRequest, fmt.Errorf("document ID is required"))
	}

	document, err := d.repo.FindById(id)
	if err != nil {
		return nil, utils.NewAppError(utils.ErrNotFound, fmt.Errorf("document not found: %w", err))
	}

	return document, nil
}

func (d *documentServiceImpl) SubmitAction(ctx context.Context, id string, input *dto.UpdateDocumentDTO) (*entities.Document, error) {
	if err := d.validateSubmitActionInput(id, input); err != nil {
		return nil, err
	}
	role := ctx.Value(utils.RoleContextKey).(string)
	fmt.Println("User role from context:", role)
	document, err := d.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if !d.validateRoleApproval(document.CurrentApprover, role) {
		return nil, utils.NewAppError(utils.ErrForbiddenAccess, fmt.Errorf("user role %s not authorized for approver level %d", role, document.CurrentApprover))
	}
	if err := d.validateDocumentState(document); err != nil {
		return nil, err
	}
	if err := d.processApprovalAction(document, input); err != nil {
		return nil, err
	}
	if err := d.repo.UpdateDocument(ctx, document); err != nil {
		return nil, utils.NewAppError(utils.ErrInternalServer, fmt.Errorf("failed to update document: %w", err))
	}
	return document, nil
}

func (d *documentServiceImpl) ResubmitAction(ctx context.Context, id string) (*entities.Document, error) {
	document, err := d.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if document.Status != entities.StatusRejected {
		return nil, utils.NewAppError(utils.ErrInvalidRequest, fmt.Errorf("only rejected documents can be resubmitted"))
	}

	now := time.Now()

	document.Status = entities.StatusNeedRevision
	document.CurrentApprover = 1
	document.UpdatedAt = now
	document.Approver1Action = nil
	document.Approver1Comment = nil
	document.Approver1Date = nil
	document.Approver2Action = nil
	document.Approver2Comment = nil
	document.Approver2Date = nil
	document.Approver3Action = nil
	document.Approver3Comment = nil
	document.Approver3Date = nil

	if err := d.repo.UpdateDocument(ctx, document); err != nil {
		return nil, utils.NewAppError(utils.ErrInternalServer, fmt.Errorf("failed to resubmit document: %w", err))
	}

	return document, nil
}

func (d *documentServiceImpl) validateSubmitActionInput(id string, input *dto.UpdateDocumentDTO) error {
	if id == "" {
		return utils.NewAppError(utils.ErrInvalidRequest, fmt.Errorf("document ID is required"))
	}

	if input == nil {
		return utils.NewAppError(utils.ErrInvalidRequest, fmt.Errorf("input data is required"))
	}

	if input.Action != entities.ActionApprove && input.Action != entities.ActionReject {
		return utils.NewAppError(utils.ErrInvalidRequest, fmt.Errorf("invalid action: must be approve or reject"))
	}

	return nil
}

func (d *documentServiceImpl) validateDocumentState(document *entities.Document) error {
	if document.Status == entities.StatusApproved {
		return utils.NewAppError(utils.ErrInvalidRequest, fmt.Errorf("document is already approved"))
	}

	if document.Status == entities.StatusRejected {
		return utils.NewAppError(utils.ErrInvalidRequest, fmt.Errorf("document is rejected, use resubmit instead"))
	}

	if document.CurrentApprover < 1 || document.CurrentApprover > 3 {
		return utils.NewAppError(utils.ErrInvalidRequest, fmt.Errorf("invalid approver level"))
	}

	return nil
}

func (d *documentServiceImpl) processApprovalAction(document *entities.Document, input *dto.UpdateDocumentDTO) error {
	now := time.Now()

	switch input.Action {
	case entities.ActionReject:
		return d.processRejection(document, input.Comment, &now)
	case entities.ActionApprove:
		return d.processApproval(document, input.Comment, &now)
	default:
		return utils.NewAppError(utils.ErrInvalidRequest, fmt.Errorf("invalid action"))
	}
}

func (d *documentServiceImpl) validateRoleApproval(currentApprover int, userRole string) bool {
	approverRoleMap := map[int]userEntity.RoleEnum{
		1: userEntity.RoleAdmin1,
		2: userEntity.RoleAdmin2,
		3: userEntity.RoleAdmin3,
	}
	requiredRole, exists := approverRoleMap[currentApprover]
	if !exists {
		return false
	}
	return string(requiredRole) == userRole
}

func (d *documentServiceImpl) processRejection(document *entities.Document, comment *string, timestamp *time.Time) error {
	document.Status = entities.StatusRejected

	rejectAction := entities.ActionReject
	switch document.CurrentApprover {
	case 1:
		document.Approver1Action = &rejectAction
		document.Approver1Comment = comment
		document.Approver1Date = timestamp
	case 2:
		document.Approver2Action = &rejectAction
		document.Approver2Comment = comment
		document.Approver2Date = timestamp
		document.Approver1Action = nil
		document.Approver1Comment = nil
		document.Approver1Date = nil
	case 3:
		document.Approver3Action = &rejectAction
		document.Approver3Comment = comment
		document.Approver3Date = timestamp
		document.Approver1Action = nil
		document.Approver1Comment = nil
		document.Approver1Date = nil
		document.Approver2Action = nil
		document.Approver2Comment = nil
		document.Approver2Date = nil

	default:
		return utils.NewAppError(utils.ErrInvalidRequest, fmt.Errorf("invalid approver level: %d", document.CurrentApprover))
	}

	document.CurrentApprover = 1

	return nil
}
func (d *documentServiceImpl) processApproval(document *entities.Document, comment *string, timestamp *time.Time) error {
	approveAction := entities.ActionApprove
	switch document.CurrentApprover {
	case 1:
		document.Approver1Action = &approveAction
		document.Approver1Comment = comment
		document.Approver1Date = timestamp
		document.CurrentApprover = 2
		document.Status = entities.StatusPending

	case 2:
		document.Approver2Action = &approveAction
		document.Approver2Comment = comment
		document.Approver2Date = timestamp
		document.CurrentApprover = 3
		document.Status = entities.StatusPending

	case 3:
		document.Approver3Action = &approveAction
		document.Approver3Comment = comment
		document.Approver3Date = timestamp
		document.Status = entities.StatusApproved

	default:
		return utils.NewAppError(utils.ErrInvalidRequest, fmt.Errorf("invalid approver level: %d", document.CurrentApprover))
	}
	return nil
}
