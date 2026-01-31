package service

import (
	"context"
	"strconv"

	"course_select/internal/domain/model"
	"course_select/internal/domain/repository"
	"course_select/internal/infrastructure/encrypt"
	"course_select/internal/pkg/errcode"
)

// MemberService 成员服务
// 提供成员的增删改查等业务逻辑
type MemberService struct {
	memberRepo repository.IMemberRepo
}

// NewMemberService 创建成员服务
func NewMemberService(memberRepo repository.IMemberRepo) *MemberService {
	return &MemberService{
		memberRepo: memberRepo,
	}
}

// Create 创建成员
// 参数:
//   - ctx: 上下文
//   - req: 创建成员请求
//
// 返回:
//   - *model.Member: 创建的成员
//   - error: 错误信息
func (s *MemberService) Create(ctx context.Context, req *model.CreateMemberRequest) (*model.Member, error) {
	// Step 1: 检查用户名是否已存在
	existing, err := s.memberRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errcode.UserHasExisted
	}

	// Step 2: 加密密码
	hashedPassword, err := encrypt.HashPassword(req.Password)
	if err != nil {
		return nil, errcode.UnknownError.WithMsg("密码加密失败")
	}

	// Step 3: 创建成员
	member := &model.Member{
		Username:  req.Username,
		Password:  hashedPassword,
		Nickname:  req.Nickname,
		UserType:  req.UserType,
		IsDeleted: false,
	}

	if err := s.memberRepo.Create(ctx, member); err != nil {
		return nil, err
	}

	return member, nil
}

// Get 获取成员
// 参数:
//   - ctx: 上下文
//   - userID: 用户ID
//
// 返回:
//   - *model.Member: 成员信息
//   - error: 错误信息
func (s *MemberService) Get(ctx context.Context, userID string) (*model.Member, error) {
	// Step 1: 解析用户ID
	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil, errcode.ParamInvalid
	}

	// Step 2: 获取成员信息
	member, err := s.memberRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errcode.UserNotExisted
	}
	if member.IsDeleted {
		return nil, errcode.UserHasDeleted
	}
	return member, nil
}

// List 获取成员列表
// 参数:
//   - ctx: 上下文
//   - offset: 偏移量
//   - limit: 限制数量
//
// 返回:
//   - []*model.Member: 成员列表
//   - int64: 总数
//   - error: 错误信息
func (s *MemberService) List(ctx context.Context, offset, limit int) ([]*model.Member, int64, error) {
	members, err := s.memberRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.memberRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return members, count, nil
}

// Update 更新成员
// 参数:
//   - ctx: 上下文
//   - userID: 用户ID
//   - nickname: 新昵称
//
// 返回:
//   - error: 错误信息
func (s *MemberService) Update(ctx context.Context, userID, nickname string) error {
	// Step 1: 解析用户ID
	id, err := strconv.Atoi(userID)
	if err != nil {
		return errcode.ParamInvalid
	}

	// Step 2: 验证成员是否存在
	member, err := s.memberRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if member == nil {
		return errcode.UserNotExisted
	}
	if member.IsDeleted {
		return errcode.UserHasDeleted
	}

	// Step 3: 更新成员信息
	return s.memberRepo.Update(ctx, id, map[string]interface{}{
		"nickname": nickname,
	})
}

// Delete 删除成员 (软删除)
// 参数:
//   - ctx: 上下文
//   - userID: 用户ID
//
// 返回:
//   - error: 错误信息
func (s *MemberService) Delete(ctx context.Context, userID string) error {
	// Step 1: 解析用户ID
	id, err := strconv.Atoi(userID)
	if err != nil {
		return errcode.ParamInvalid
	}

	// Step 2: 验证成员是否存在
	member, err := s.memberRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if member == nil {
		return errcode.UserNotExisted
	}

	// Step 3: 软删除成员
	return s.memberRepo.Delete(ctx, id)
}

// IsUserExist 检查用户是否存在
// 参数:
//   - ctx: 上下文
//   - userID: 用户ID
//
// 返回:
//   - bool: 是否存在且未删除
//   - error: 错误信息
func (s *MemberService) IsUserExist(ctx context.Context, userID string) (bool, error) {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return false, nil
	}
	member, err := s.memberRepo.GetByID(ctx, id)
	if err != nil {
		return false, err
	}
	return member != nil && !member.IsDeleted, nil
}
