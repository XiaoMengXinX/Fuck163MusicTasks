package utils

// APIConfig API配置 , 用于设置自定义API（参见 https://github.com/Binaryify/NeteaseCloudMusicApi ）
type APIConfig struct {
	NeteaseAPI string
}

// Cookies 用户 cookies 数据类型
type Cookies []struct {
	Key   string
	Value string
}

// Headers 自定义 Headers 数据类型 (仅对于非 eapi 有效)
type Headers []struct {
	Key   string
	Value string
}

// RequestData 传入请求数据类型
type RequestData struct {
	Cookies Cookies
	Headers Headers
	Body    string
}

// eapi 请求所需要的参数
type eapiOption struct {
	Json string
	Path string
	Url  string
}

// Config 配置文件结构
type Config struct {
	NeteaseAPI string `json:"NeteaseAPI"`
	DEBUG      bool   `json:"DEBUG"`
	Users      []struct {
		Cookies Cookies `json:"Cookies"`
	} `json:"Users"`
	EventSendConfig struct {
		LagConfig LagConfig `json:"LagConfig"`
	} `json:"EventSendConfig"`
	CommentReplyConfig struct {
		RepliedComment []struct {
			ID        int `json:"ID"`
			CommentId int `json:"CommentId"`
		} `json:"RepliedComment"`
		LagConfig LagConfig `json:"LagConfig"`
	} `json:"CommentReplyConfig"`
	SendMsgConfig struct {
		UserID    [][]int   `json:"UserID"`
		LagConfig LagConfig `json:"LagConfig"`
	}
	Content []string `json:"Content"`
	Cron    struct {
		Enabled    bool      `json:"Enabled"`
		Expression string    `json:"Expression"`
		EnableLag  bool      `json:"EnableLag"`
		LagConfig  LagConfig `json:"LagConfig"`
	} `json:"Cron"`
}

// LagConfig 延迟设置
type LagConfig struct {
	LagBetweenSendAndDelete bool `json:"LagBetweenSendAndDelete"`
	RandomLag               bool `json:"RandomLag"`
	DefaultLag              int  `json:"DefaultLag"`
	LagMin                  int  `json:"LagMin"`
	LagMax                  int  `json:"LagMax"`
}

// UserDetailData 用户详细数据
type UserDetailData struct {
	Identify struct {
		ImageUrl  string      `json:"imageUrl"`
		ImageDesc string      `json:"imageDesc"`
		ActionUrl interface{} `json:"actionUrl"`
	} `json:"identify"`
	CurrentExpert struct {
		RoleId   int         `json:"roleId"`
		RoleName string      `json:"roleName"`
		Level    interface{} `json:"level"`
	} `json:"currentExpert"`
	ExpertArray []struct {
		RoleId   int         `json:"roleId"`
		RoleName string      `json:"roleName"`
		Level    interface{} `json:"level"`
	} `json:"expertArray"`
	CurrentProduct interface{} `json:"currentProduct"`
	Products       []struct {
		ProductionTypeName string `json:"productionTypeName"`
		ProductionTypeId   int    `json:"productionTypeId"`
	} `json:"products"`
	Level       int `json:"level"`
	ListenSongs int `json:"listenSongs"`
	UserPoint   struct {
		UserId       int `json:"userId"`
		Balance      int `json:"balance"`
		UpdateTime   int `json:"updateTime"`
		Version      int `json:"version"`
		Status       int `json:"status"`
		BlockBalance int `json:"blockBalance"`
	} `json:"userPoint"`
	MobileSign bool `json:"mobileSign"`
	PcSign     bool `json:"pcSign"`
	Profile    struct {
		AvatarDetail struct {
			UserType        interface{} `json:"userType"`
			IdentityLevel   int         `json:"identityLevel"`
			IdentityIconUrl string      `json:"identityIconUrl"`
		} `json:"avatarDetail"`
		BackgroundImgIdStr string      `json:"backgroundImgIdStr"`
		UserId             int         `json:"userId"`
		DjStatus           int         `json:"djStatus"`
		AccountStatus      int         `json:"accountStatus"`
		Province           int         `json:"province"`
		VipType            int         `json:"vipType"`
		Followed           bool        `json:"followed"`
		CreateTime         int         `json:"createTime"`
		AvatarImgId        int         `json:"avatarImgId"`
		Birthday           int         `json:"birthday"`
		Gender             int         `json:"gender"`
		Nickname           string      `json:"nickname"`
		AvatarImgIdStr     string      `json:"avatarImgIdStr"`
		Description        string      `json:"description"`
		Mutual             bool        `json:"mutual"`
		RemarkName         interface{} `json:"remarkName"`
		UserType           int         `json:"userType"`
		AuthStatus         int         `json:"authStatus"`
		DetailDescription  string      `json:"detailDescription"`
		Experts            struct {
		} `json:"experts"`
		ExpertTags      interface{} `json:"expertTags"`
		City            int         `json:"city"`
		DefaultAvatar   bool        `json:"defaultAvatar"`
		BackgroundImgId int         `json:"backgroundImgId"`
		BackgroundUrl   string      `json:"backgroundUrl"`
		AvatarUrl       string      `json:"avatarUrl"`
		Signature       string      `json:"signature"`
		Authority       int         `json:"authority"`
		AllAuthTypes    []struct {
			Type int      `json:"type"`
			Desc string   `json:"desc"`
			Tags []string `json:"tags"`
		} `json:"allAuthTypes"`
		Followeds                 int  `json:"followeds"`
		Follows                   int  `json:"follows"`
		Blacklist                 bool `json:"blacklist"`
		ArtistId                  int  `json:"artistId"`
		EventCount                int  `json:"eventCount"`
		AllSubscribedCount        int  `json:"allSubscribedCount"`
		PlaylistBeSubscribedCount int  `json:"playlistBeSubscribedCount"`
		MainAuthType              struct {
			Type int      `json:"type"`
			Desc string   `json:"desc"`
			Tags []string `json:"tags"`
		} `json:"mainAuthType"`
		AvatarImgIdStr1 string      `json:"avatarImgId_str"`
		FollowTime      interface{} `json:"followTime"`
		FollowMe        bool        `json:"followMe"`
		ArtistIdentity  []int       `json:"artistIdentity"`
		CCount          int         `json:"cCount"`
		SDJPCount       int         `json:"sDJPCount"`
		ArtistName      string      `json:"artistName"`
		PlaylistCount   int         `json:"playlistCount"`
		SCount          int         `json:"sCount"`
		NewFollows      int         `json:"newFollows"`
	} `json:"profile"`
	PeopleCanSeeMyPlayRecord bool `json:"peopleCanSeeMyPlayRecord"`
	Bindings                 []struct {
		UserId       int         `json:"userId"`
		Url          string      `json:"url"`
		ExpiresIn    int         `json:"expiresIn"`
		RefreshTime  int         `json:"refreshTime"`
		BindingTime  int         `json:"bindingTime"`
		TokenJsonStr interface{} `json:"tokenJsonStr"`
		Expired      bool        `json:"expired"`
		Id           int         `json:"id"`
		Type         int         `json:"type"`
	} `json:"bindings"`
	AdValid    bool `json:"adValid"`
	Code       int  `json:"code"`
	CreateTime int  `json:"createTime"`
	CreateDays int  `json:"createDays"`
}

// LoginStatData 登录状态数据
type LoginStatData struct {
	Data struct {
		Code    int `json:"code"`
		Account struct {
			Id                 int    `json:"id"`
			UserName           string `json:"userName"`
			Type               int    `json:"type"`
			Status             int    `json:"status"`
			WhitelistAuthority int    `json:"whitelistAuthority"`
			CreateTime         int    `json:"createTime"`
			TokenVersion       int    `json:"tokenVersion"`
			Ban                int    `json:"ban"`
			BaoyueVersion      int    `json:"baoyueVersion"`
			DonateVersion      int    `json:"donateVersion"`
			VipType            int    `json:"vipType"`
			AnonimousUser      bool   `json:"anonimousUser"`
			PaidFee            bool   `json:"paidFee"`
		} `json:"account"`
		Profile struct {
			UserId              int         `json:"userId"`
			UserType            int         `json:"userType"`
			Nickname            string      `json:"nickname"`
			AvatarImgId         int         `json:"avatarImgId"`
			AvatarUrl           string      `json:"avatarUrl"`
			BackgroundImgId     int         `json:"backgroundImgId"`
			BackgroundUrl       string      `json:"backgroundUrl"`
			Signature           string      `json:"signature"`
			CreateTime          int         `json:"createTime"`
			UserName            string      `json:"userName"`
			AccountType         int         `json:"accountType"`
			ShortUserName       string      `json:"shortUserName"`
			Birthday            int         `json:"birthday"`
			Authority           int         `json:"authority"`
			Gender              int         `json:"gender"`
			AccountStatus       int         `json:"accountStatus"`
			Province            int         `json:"province"`
			City                int         `json:"city"`
			AuthStatus          int         `json:"authStatus"`
			Description         string      `json:"description"`
			DetailDescription   string      `json:"detailDescription"`
			DefaultAvatar       bool        `json:"defaultAvatar"`
			ExpertTags          interface{} `json:"expertTags"`
			Experts             interface{} `json:"experts"`
			DjStatus            int         `json:"djStatus"`
			LocationStatus      int         `json:"locationStatus"`
			VipType             int         `json:"vipType"`
			Followed            bool        `json:"followed"`
			Mutual              bool        `json:"mutual"`
			Authenticated       bool        `json:"authenticated"`
			LastLoginTime       int         `json:"lastLoginTime"`
			LastLoginIP         string      `json:"lastLoginIP"`
			RemarkName          interface{} `json:"remarkName"`
			ViptypeVersion      int         `json:"viptypeVersion"`
			AuthenticationTypes int         `json:"authenticationTypes"`
			AvatarDetail        interface{} `json:"avatarDetail"`
			Anchor              bool        `json:"anchor"`
		} `json:"profile"`
	} `json:"data"`
}

// TasksData 用于解析音乐人任务 api 的返回数据
type TasksData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List []struct {
			Business        string `json:"business"`
			MissionId       int    `json:"missionId"`
			UserId          int    `json:"userId"`
			MissionEntityId int    `json:"missionEntityId"`
			RewardId        int    `json:"rewardId"`
			ProgressRate    int    `json:"progressRate"`
			Description     string `json:"description"`
			Type            int    `json:"type"`
			Tag             int    `json:"tag"`
			ActionType      int    `json:"actionType"`
			Platform        int    `json:"platform"`
			Status          int    `json:"status"`
			Button          string `json:"button"`
			SortValue       int    `json:"sortValue"`
			StartTime       int    `json:"startTime"`
			EndTime         int    `json:"endTime"`
			ExtendInfo      string `json:"extendInfo"`
			Period          int    `json:"period"`
			NeedToReceive   int    `json:"needToReceive,omitempty"`
			TargetCount     int    `json:"targetCount"`
			RewardWorth     string `json:"rewardWorth"`
			RewardType      int    `json:"rewardType"`
			UserMissionId   int    `json:"userMissionId,omitempty"`
			UpdateTime      int    `json:"updateTime,omitempty"`
		} `json:"list"`
	} `json:"data"`
}

// CloudbeanData 用于解析云豆数量 api 的返回数据
type CloudbeanData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ArtistId      int         `json:"artistId"`
		CloudBean     int         `json:"cloudBean"`
		MusicianLevel int         `json:"musicianLevel"`
		MaxCloudBean  int         `json:"maxCloudBean"`
		NeedPop       bool        `json:"needPop"`
		FlowStage     interface{} `json:"flowStage"`
		DynamicNews   interface{} `json:"dynamicNews"`
		Signed        bool        `json:"signed"`
	} `json:"data"`
}

// MusicianSignResult 用于解析音乐人签到 api 的返回数据
type MusicianSignResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    bool   `json:"data"`
}

// UserSignResult 用于解析用户签到 api 的返回数据
type UserSignResult struct {
	Point int    `json:"point"`
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
}

// SendMsgResult 用于解析发送私信 api 的返回数据
type SendMsgResult struct {
	Code    int `json:"code"`
	NewMsgs []struct {
		FromUser struct {
			Description        string      `json:"description"`
			BackgroundUrl      string      `json:"backgroundUrl"`
			BackgroundImgId    int         `json:"backgroundImgId"`
			Birthday           int         `json:"birthday"`
			AccountStatus      int         `json:"accountStatus"`
			City               int         `json:"city"`
			DetailDescription  string      `json:"detailDescription"`
			DefaultAvatar      bool        `json:"defaultAvatar"`
			DjStatus           int         `json:"djStatus"`
			Followed           bool        `json:"followed"`
			UserId             int         `json:"userId"`
			AvatarDetail       interface{} `json:"avatarDetail"`
			UserType           int         `json:"userType"`
			AvatarImgId        int         `json:"avatarImgId"`
			Gender             int         `json:"gender"`
			BackgroundImgIdStr string      `json:"backgroundImgIdStr"`
			AvatarImgIdStr     string      `json:"avatarImgIdStr"`
			AvatarUrl          string      `json:"avatarUrl"`
			AuthStatus         int         `json:"authStatus"`
			ExpertTags         interface{} `json:"expertTags"`
			VipType            int         `json:"vipType"`
			Experts            interface{} `json:"experts"`
			Province           int         `json:"province"`
			RemarkName         interface{} `json:"remarkName"`
			Nickname           string      `json:"nickname"`
			Mutual             bool        `json:"mutual"`
			Signature          string      `json:"signature"`
			Authority          int         `json:"authority"`
		} `json:"fromUser"`
		ToUser struct {
			Description        string      `json:"description"`
			BackgroundUrl      string      `json:"backgroundUrl"`
			BackgroundImgId    int         `json:"backgroundImgId"`
			Birthday           int         `json:"birthday"`
			AccountStatus      int         `json:"accountStatus"`
			City               int         `json:"city"`
			DetailDescription  string      `json:"detailDescription"`
			DefaultAvatar      bool        `json:"defaultAvatar"`
			DjStatus           int         `json:"djStatus"`
			Followed           bool        `json:"followed"`
			UserId             int         `json:"userId"`
			AvatarDetail       interface{} `json:"avatarDetail"`
			UserType           int         `json:"userType"`
			AvatarImgId        int         `json:"avatarImgId"`
			Gender             int         `json:"gender"`
			BackgroundImgIdStr string      `json:"backgroundImgIdStr"`
			AvatarImgIdStr     string      `json:"avatarImgIdStr"`
			AvatarUrl          string      `json:"avatarUrl"`
			AuthStatus         int         `json:"authStatus"`
			ExpertTags         interface{} `json:"expertTags"`
			VipType            int         `json:"vipType"`
			Experts            interface{} `json:"experts"`
			Province           int         `json:"province"`
			RemarkName         interface{} `json:"remarkName"`
			Nickname           string      `json:"nickname"`
			Mutual             bool        `json:"mutual"`
			Signature          string      `json:"signature"`
			Authority          int         `json:"authority"`
		} `json:"toUser"`
		RealFromUser interface{} `json:"realFromUser"`
		Msg          string      `json:"msg"`
		Time         int         `json:"time"`
		BatchId      int         `json:"batchId"`
		Id           int         `json:"id"`
	} `json:"newMsgs"`
	Id            int           `json:"id"`
	Sendblacklist []interface{} `json:"sendblacklist"`
	Blacklist     []interface{} `json:"blacklist"`
}

// CommentConfig 发送评论参数
type CommentConfig struct {
	CommentType int    // 评论类型, 1 发送, 2 回复
	ResType     int    // 资源类型, 0: 歌曲, 1: mv, 2: 歌单, 3: 专辑, 4: 电台, 5: 视频, 6: 动态
	ID          int    // 对应资源 id
	ThreadId    int    // 给动态发送评论，则不需要传 id，需要传动态的 threadId
	Content     string // 要发送的内容
	CommentId   int    // 回复的评论id (回复评论时必填)
}

// CommentResult 用于解析发送评论 api 的返回数据
type CommentResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Comment struct {
		User struct {
			LocationInfo interface{} `json:"locationInfo"`
			LiveInfo     interface{} `json:"liveInfo"`
			Anonym       int         `json:"anonym"`
			AvatarDetail interface{} `json:"avatarDetail"`
			UserType     int         `json:"userType"`
			AuthStatus   int         `json:"authStatus"`
			AvatarUrl    string      `json:"avatarUrl"`
			Nickname     string      `json:"nickname"`
			ExpertTags   interface{} `json:"expertTags"`
			VipType      int         `json:"vipType"`
			RemarkName   interface{} `json:"remarkName"`
			Experts      interface{} `json:"experts"`
			VipRights    struct {
				Associator   interface{} `json:"associator"`
				MusicPackage struct {
					VipCode int  `json:"vipCode"`
					Rights  bool `json:"rights"`
				} `json:"musicPackage"`
				RedVipAnnualCount int `json:"redVipAnnualCount"`
				RedVipLevel       int `json:"redVipLevel"`
			} `json:"vipRights"`
			UserId int `json:"userId"`
		} `json:"user"`
		BeRepliedUser struct {
			LocationInfo interface{} `json:"locationInfo"`
			LiveInfo     interface{} `json:"liveInfo"`
			Anonym       int         `json:"anonym"`
			AvatarDetail interface{} `json:"avatarDetail"`
			UserType     int         `json:"userType"`
			AuthStatus   int         `json:"authStatus"`
			AvatarUrl    string      `json:"avatarUrl"`
			Nickname     string      `json:"nickname"`
			ExpertTags   interface{} `json:"expertTags"`
			VipType      int         `json:"vipType"`
			RemarkName   interface{} `json:"remarkName"`
			Experts      interface{} `json:"experts"`
			VipRights    interface{} `json:"vipRights"`
			UserId       int         `json:"userId"`
		} `json:"beRepliedUser"`
		CommentLocationType int         `json:"commentLocationType"`
		CommentId           int         `json:"commentId"`
		ExpressionUrl       interface{} `json:"expressionUrl"`
		Time                int         `json:"time"`
		Content             string      `json:"content"`
	} `json:"comment"`
}

// EventConfig 发送动态的参数
type EventConfig struct {
	Msg        string `json:"msg"`
	Type       string `json:"type"`
	UUID       string `json:"uuid"`
	Pics       string `json:"pics"`
	AddComment string `json:"addComment"`
	Header     string `json:"header"`
	ER         string `json:"e_r"`
}

// SendEventResult 用于解析 eapi 发送动态的返回数据
type SendEventResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	UserId  int    `json:"userId"`
	Id      int    `json:"id"`
	Event   struct {
		ActName          interface{} `json:"actName"`
		PendantData      interface{} `json:"pendantData"`
		ForwardCount     int         `json:"forwardCount"`
		LotteryEventData interface{} `json:"lotteryEventData"`
		Json             string      `json:"json"`
		User             struct {
			DefaultAvatar      bool        `json:"defaultAvatar"`
			Province           int         `json:"province"`
			AuthStatus         int         `json:"authStatus"`
			Followed           bool        `json:"followed"`
			AvatarUrl          string      `json:"avatarUrl"`
			AccountStatus      int         `json:"accountStatus"`
			Gender             int         `json:"gender"`
			City               int         `json:"city"`
			Birthday           int         `json:"birthday"`
			UserId             int         `json:"userId"`
			UserType           int         `json:"userType"`
			Nickname           string      `json:"nickname"`
			Signature          string      `json:"signature"`
			Description        string      `json:"description"`
			DetailDescription  string      `json:"detailDescription"`
			AvatarImgId        int         `json:"avatarImgId"`
			BackgroundImgId    int         `json:"backgroundImgId"`
			BackgroundUrl      string      `json:"backgroundUrl"`
			Authority          int         `json:"authority"`
			Mutual             bool        `json:"mutual"`
			ExpertTags         interface{} `json:"expertTags"`
			Experts            interface{} `json:"experts"`
			DjStatus           int         `json:"djStatus"`
			VipType            int         `json:"vipType"`
			RemarkName         interface{} `json:"remarkName"`
			UrlAnalyze         bool        `json:"urlAnalyze"`
			Followeds          int         `json:"followeds"`
			AvatarImgIdStr     string      `json:"avatarImgId_str"`
			AvatarImgIdStr1    string      `json:"avatarImgIdStr"`
			BackgroundImgIdStr string      `json:"backgroundImgIdStr"`
			VipRights          interface{} `json:"vipRights"`
			AvatarDetail       struct {
				UserType        int    `json:"userType"`
				IdentityLevel   int    `json:"identityLevel"`
				IdentityIconUrl string `json:"identityIconUrl"`
			} `json:"avatarDetail"`
			CommonIdentity interface{} `json:"commonIdentity"`
		} `json:"user"`
		Uuid               string        `json:"uuid"`
		ExpireTime         int           `json:"expireTime"`
		RcmdInfo           interface{}   `json:"rcmdInfo"`
		EventTime          int           `json:"eventTime"`
		ActId              int           `json:"actId"`
		Pics               []interface{} `json:"pics"`
		TmplId             int           `json:"tmplId"`
		ShowTime           int           `json:"showTime"`
		InsertTime         int           `json:"insertTime"`
		Id                 int           `json:"id"`
		Type               int           `json:"type"`
		TopEvent           bool          `json:"topEvent"`
		InsiteForwardCount int           `json:"insiteForwardCount"`
		Info               struct {
			CommentThread struct {
				Id               string        `json:"id"`
				ResourceInfo     interface{}   `json:"resourceInfo"`
				ResourceType     int           `json:"resourceType"`
				CommentCount     int           `json:"commentCount"`
				LikedCount       int           `json:"likedCount"`
				ShareCount       int           `json:"shareCount"`
				HotCount         int           `json:"hotCount"`
				LatestLikedUsers []interface{} `json:"latestLikedUsers"`
				ResourceOwnerId  int           `json:"resourceOwnerId"`
				ResourceId       int           `json:"resourceId"`
			} `json:"commentThread"`
			LatestLikedUsers []interface{} `json:"latestLikedUsers"`
			Liked            bool          `json:"liked"`
			Comments         []interface{} `json:"comments"`
			ResourceType     int           `json:"resourceType"`
			ResourceId       int           `json:"resourceId"`
			ThreadId         string        `json:"threadId"`
			ShareCount       int           `json:"shareCount"`
			CommentCount     int           `json:"commentCount"`
			LikedCount       int           `json:"likedCount"`
		} `json:"info"`
		TailMark    interface{} `json:"tailMark"`
		ExtJsonInfo struct {
			ActId          int           `json:"actId"`
			ActIds         []interface{} `json:"actIds"`
			Uuid           string        `json:"uuid"`
			ExtType        string        `json:"extType"`
			ExtId          string        `json:"extId"`
			CircleId       interface{}   `json:"circleId"`
			CirclePubType  interface{}   `json:"circlePubType"`
			TailMark       interface{}   `json:"tailMark"`
			TypeDesc       interface{}   `json:"typeDesc"`
			PrivacySetting int           `json:"privacySetting"`
			ExtParams      struct {
			} `json:"extParams"`
		} `json:"extJsonInfo"`
		PrivacySetting int         `json:"privacySetting"`
		ExtPageParam   interface{} `json:"extPageParam"`
		LogInfo        interface{} `json:"logInfo"`
	} `json:"event"`
	Sns struct {
	} `json:"sns"`
	ResUrl string `json:"resUrl"`
}

// DelEventConfig 用于解析 eapi 删除评论的返回数据
type DelEventConfig struct {
	ID     string `json:"id"`
	Header string `json:"header"`
	ER     string `json:"e_r"`
}

// PlainResult 用于解析一些操作类 api 的返回值
type PlainResult struct {
	Message string `json:"message"`
	Msg     string `json:"msg"`
	Code    int    `json:"code"`
}

// RandomNum 随机数设置
type RandomNum struct {
	IsRandom   bool
	DefaultNum int
	MinNum     int
	MaxNum     int
}
