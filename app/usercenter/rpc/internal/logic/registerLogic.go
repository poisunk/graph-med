package logic

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"graph-med/app/usercenter/model"
	"graph-med/app/usercenter/rpc/usercenter"
	"graph-med/pkg/tool"
	"graph-med/pkg/xerr"

	"graph-med/app/usercenter/rpc/internal/svc"
	"graph-med/app/usercenter/rpc/pd"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUserAlreadyRegisterError = xerr.NewErrMsg("user has been registered")

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pd.RegisterReq) (*pd.RegisterResp, error) {

	// 1. 判断用户是否存在
	user, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, in.Email)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "email:%s,err:%v", in.Email, err)
	}
	if user != nil {
		return nil, errors.Wrapf(ErrUserAlreadyRegisterError, "Register user exists email:%s,err:%v", in.Email, err)
	}

	// 2. 插入用户
	var id int64
	if err := l.svcCtx.UserModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		user := new(model.User)
		// email
		err = user.Email.Scan(in.Email)
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "email:%s,err:%v", in.Email, err)
		}
		// userId
		user.UserId = uuid.New().String()
		// nickname
		if len(in.Nickname) == 0 {
			user.Nickname = tool.RandomString(8, tool.RAND_STRING_KIND_ALL)
		}
		// password
		if len(in.Password) > 0 {
			user.Password = tool.Md5ByString(in.Password)
		}
		// insert
		insertResult, err := l.svcCtx.UserModel.Insert(ctx, user)
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Register db user Insert err:%v,user:%+v", err, user)
		}
		lastId, err := insertResult.LastInsertId()
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Register db user insertResult.LastInsertId err:%v,user:%+v", err, user)
		}
		id = lastId
		return nil
	}); err != nil {
		return nil, err
	}

	// 3. 生成token
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&usercenter.GenerateTokenReq{
		Id: id,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "GenerateToken fail")
	}

	return &usercenter.RegisterResp{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}
