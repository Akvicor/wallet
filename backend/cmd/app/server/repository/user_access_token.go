package repository

import (
	"context"
	"wallet/cmd/app/server/common/resp"
	"wallet/cmd/app/server/model"
)

var UserAccessToken = new(userAccessTokenRepository)

type userAccessTokenRepository struct {
	base[*model.UserAccessToken]
}

/**
查找
*/

// FindAll 获取所有的记录, page为nil时不分页
func (u *userAccessTokenRepository) FindAll(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderAccessToken) (tokens []*model.UserAccessToken, err error) {
	return u.paging(page, u.preload(c, alive, preloader).Order("disabled ASC"))
}

// FindAllByUID 获取所有的记录, page为nil时不分页
func (u *userAccessTokenRepository) FindAllByUID(c context.Context, page *resp.PageModel, alive bool, preloader *model.PreloaderAccessToken, uid int64) (tokens []*model.UserAccessToken, err error) {
	return u.paging(page, u.preload(c, alive, preloader).Where("uid = ?", uid).Order("disabled ASC"))
}

// FindByToken 通过Token获取记录
func (u *userAccessTokenRepository) FindByToken(c context.Context, alive bool, preloader *model.PreloaderAccessToken, tokenStr string) (token *model.UserAccessToken, err error) {
	token = new(model.UserAccessToken)
	err = u.preload(c, alive, preloader).Where("token = ?", tokenStr).First(token).Error
	return token, err
}

/**
创建
*/

// Create 创建记录
func (u *userAccessTokenRepository) Create(c context.Context, token *model.UserAccessToken) error {
	return u.WrapResultErr(u.db(c).Create(token))
}

/**
更新
*/

// UpdateLastUsedByToken 更新上次使用时间
func (u *userAccessTokenRepository) UpdateLastUsedByToken(c context.Context, alive bool, tokenStr string, timestamp int64) error {
	return u.WrapResultErr(u.alive(c, alive).Where("token = ?", tokenStr).UpdateColumn("last_used", timestamp))
}

// UpdateDisabledByUID 更新停用时间
func (u *userAccessTokenRepository) UpdateDisabledByUID(c context.Context, alive bool, uid, timestamp int64) error {
	return u.WrapResultErr(u.alive(c, alive).Where("uid = ?", uid).UpdateColumn("disabled", timestamp))
}
