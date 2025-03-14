package user_models

import (
	"github.com/google/uuid"
	"go-zero2/common/models"
)

type SysUser struct {
	models.Model
	UUID         uuid.UUID `json:"uuid" gorm:"default:null;comment:用户UUID"`                                              // 用户UUID
	Username     string    `json:"userName" gorm:"index;comment:用户登录名"`                                                  // 用户登录名
	Password     string    `json:"-"  gorm:"comment:用户登录密码"`                                                             // 用户登录密码
	NickName     string    `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`                                            // 用户昵称
	SideMode     string    `json:"sideMode" gorm:"default:dark;comment:用户侧边主题"`                                          // 用户侧边主题
	HeaderImg    string    `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	BaseColor    string    `json:"baseColor" gorm:"default:#fff;comment:基础颜色"`                                           // 基础颜色
	AuthorityId  uint      `json:"authorityId" gorm:"default:888;comment:用户角色ID"`                                        // 用户角色ID
	Phone        string    `json:"phone"  gorm:"comment:用户手机号"`                                                          // 用户手机号
	Email        string    `json:"email"  gorm:"comment:用户邮箱"`                                                           // 用户邮箱
	Enable       int       `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"`                                      //用户是否被冻结 1正常 2冻结
	ParentID     uint      `json:"parentId" gorm:"default:0;not null;comment:上级代理用户ID"`                                  // 上级代理用户ID
	Balance      int       `json:"balance" gorm:"default:0;not null;comment:余额"`
	Commission   int       `json:"commission" gorm:"default:0;not null;comment:佣金"`
	Score        int       `json:"score" gorm:"default:0;not null;comment:积分"`
	Openid       string    `json:"openid" gorm:"default:'';not null;comment:openid:用户唯一标识"`
	ReferralCode string    `json:"referral_code" gorm:"default:'';not null;comment:用户的推荐码"`
	Ratio        int       `json:"ratio" gorm:"default:0;not null;comment:分润比例"` // 代理人分润比例字段
}
