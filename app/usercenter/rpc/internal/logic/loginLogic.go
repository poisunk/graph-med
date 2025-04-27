package logic

import (
	"context"
	"github.com/pkg/errors"
	"graph-med/app/usercenter/model"
	"graph-med/app/usercenter/rpc/usercenter"
	"graph-med/pkg/tool"
	"graph-med/pkg/xerr"

	"graph-med/app/usercenter/rpc/internal/svc"
	"graph-med/app/usercenter/rpc/pd"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrGenerateTokenError = xerr.NewErrMsg("生成token失败")
var ErrUsernamePwdError = xerr.NewErrMsg("账号或密码不正确")

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pd.LoginReq) (*pd.LoginResp, error) {
	var user *model.User
	var err error
	switch in.AuthType {
	case model.UserAuthTypeSystem:
		user, err = l.loginByEmail(in.AuthKey, in.Password)
	default:
		return nil, xerr.NewErrCode(xerr.SERVER_COMMON_ERROR)
	}
	if err != nil {
		return nil, err
	}

	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&usercenter.GenerateTokenReq{
		UserId: user.UserId,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "GenerateToken fail")
	}

	return &usercenter.LoginResp{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}

func (l *LoginLogic) loginByEmail(email, password string) (*model.User, error) {
	user, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, email)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "根据邮箱查询用户信息失败，email:%s,err:%v", email, err)
	}

	if user == nil {
		return nil, errors.Wrapf(ErrUserNoExistsError, "email:%s", email)
	}

	if !(tool.Md5ByString(password) == user.Password) {
		return nil, errors.Wrap(ErrUsernamePwdError, "密码匹配出错")
	}

	return user, nil
}
