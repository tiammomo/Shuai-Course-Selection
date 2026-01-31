package service

import (
	"context"
	"strconv"

	"course_select/internal/domain/model"
	"course_select/internal/domain/repository"
	"course_select/internal/pkg/errcode"
)

// CourseService 课程服务
// 提供课程的增删改查以及教师绑定等功能
type CourseService struct {
	courseRepo repository.ICourseRepo
	bindRepo   repository.IBindRepo
	choiceRepo repository.IChoiceRepo
}

// NewCourseService 创建课程服务
func NewCourseService(courseRepo repository.ICourseRepo, bindRepo repository.IBindRepo, choiceRepo repository.IChoiceRepo) *CourseService {
	return &CourseService{
		courseRepo: courseRepo,
		bindRepo:   bindRepo,
		choiceRepo: choiceRepo,
	}
}

// Create 创建课程
// 参数:
//   - ctx: 上下文
//   - req: 创建课程请求
//
// 返回:
//   - *model.Course: 创建的课程
//   - error: 错误信息
func (s *CourseService) Create(ctx context.Context, req *model.CreateCourseRequest) (*model.Course, error) {
	course := &model.Course{
		Name:        req.Name,
		Capacity:    req.Cap,
		CapSelected: 0,
		TeacherID:   nil,
	}

	if err := s.courseRepo.Create(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

// Get 获取课程
// 参数:
//   - ctx: 上下文
//   - courseID: 课程ID
//
// 返回:
//   - *model.Course: 课程信息
//   - error: 错误信息
func (s *CourseService) Get(ctx context.Context, courseID string) (*model.Course, error) {
	// Step 1: 解析课程ID
	id, err := strconv.Atoi(courseID)
	if err != nil {
		return nil, errcode.ParamInvalid
	}

	// Step 2: 获取课程信息
	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errcode.CourseNotExisted
	}

	return course, nil
}

// List 获取课程列表
// 参数:
//   - ctx: 上下文
//   - offset: 偏移量
//   - limit: 限制数量
//
// 返回:
//   - []*model.Course: 课程列表
//   - error: 错误信息
func (s *CourseService) List(ctx context.Context, offset, limit int) ([]*model.Course, error) {
	return s.courseRepo.List(ctx, offset, limit)
}

// BindCourse 绑定课程到教师
// 参数:
//   - ctx: 上下文
//   - courseID: 课程ID
//   - teacherID: 教师ID
//
// 返回:
//   - error: 错误信息
func (s *CourseService) BindCourse(ctx context.Context, courseID, teacherID string) error {
	// Step 1: 解析参数
	cID, err := strconv.Atoi(courseID)
	if err != nil {
		return errcode.ParamInvalid
	}
	tID, err := strconv.Atoi(teacherID)
	if err != nil {
		return errcode.ParamInvalid
	}

	// Step 2: 检查课程是否存在
	course, err := s.courseRepo.GetByID(ctx, cID)
	if err != nil {
		return err
	}
	if course == nil {
		return errcode.CourseNotExisted
	}
	if course.TeacherID != nil {
		return errcode.CourseHasBound
	}

	// Step 3: 检查课程是否已被其他教师绑定
	existingTeacher, err := s.bindRepo.GetByCourseID(ctx, cID)
	if err != nil {
		return err
	}
	if existingTeacher != nil {
		return errcode.CourseHasBound
	}

	// Step 4: 创建绑定关系
	bind := &model.Bind{
		TeacherID: tID,
		CourseID:  cID,
	}
	if err := s.bindRepo.Create(ctx, bind); err != nil {
		return err
	}

	// Step 5: 更新课程的教师ID
	return s.courseRepo.Update(ctx, cID, map[string]interface{}{
		"teacher_id": tID,
	})
}

// UnbindCourse 解绑课程
// 参数:
//   - ctx: 上下文
//   - courseID: 课程ID
//   - teacherID: 教师ID
//
// 返回:
//   - error: 错误信息
func (s *CourseService) UnbindCourse(ctx context.Context, courseID, teacherID string) error {
	// Step 1: 解析参数
	cID, err := strconv.Atoi(courseID)
	if err != nil {
		return errcode.ParamInvalid
	}
	tID, err := strconv.Atoi(teacherID)
	if err != nil {
		return errcode.ParamInvalid
	}

	// Step 2: 检查绑定是否存在
	teacher, err := s.bindRepo.GetByCourseID(ctx, cID)
	if err != nil {
		return err
	}
	if teacher == nil {
		return errcode.CourseNotBind
	}
	if *teacher != tID {
		return errcode.PermDenied
	}

	// Step 3: 删除绑定关系
	if err := s.bindRepo.DeleteByCourseID(ctx, cID); err != nil {
		return err
	}

	// Step 4: 更新课程的教师ID为nil
	return s.courseRepo.Update(ctx, cID, map[string]interface{}{
		"teacher_id": nil,
	})
}

// GetTeacherCourses 获取教师的课程列表
// 参数:
//   - ctx: 上下文
//   - teacherID: 教师ID
//
// 返回:
//   - []*model.Course: 课程列表
//   - error: 错误信息
func (s *CourseService) GetTeacherCourses(ctx context.Context, teacherID string) ([]*model.Course, error) {
	tID, err := strconv.Atoi(teacherID)
	if err != nil {
		return nil, errcode.ParamInvalid
	}
	return s.bindRepo.GetByTeacherID(ctx, tID)
}

// IsCourseExist 检查课程是否存在
// 参数:
//   - ctx: 上下文
//   - courseID: 课程ID
//
// 返回:
//   - bool: 是否存在
//   - error: 错误信息
func (s *CourseService) IsCourseExist(ctx context.Context, courseID string) (bool, error) {
	id, err := strconv.Atoi(courseID)
	if err != nil {
		return false, nil
	}
	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return false, err
	}
	return course != nil, nil
}
