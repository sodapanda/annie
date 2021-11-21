package extractors

import (
	"net/url"
	"strings"

	"github.com/sodapanda/annie/extractors/acfun"
	"github.com/sodapanda/annie/extractors/bcy"
	"github.com/sodapanda/annie/extractors/bilibili"
	"github.com/sodapanda/annie/extractors/douyin"
	"github.com/sodapanda/annie/extractors/douyu"
	"github.com/sodapanda/annie/extractors/eporner"
	"github.com/sodapanda/annie/extractors/facebook"
	"github.com/sodapanda/annie/extractors/geekbang"
	"github.com/sodapanda/annie/extractors/haokan"
	"github.com/sodapanda/annie/extractors/instagram"
	"github.com/sodapanda/annie/extractors/iqiyi"
	"github.com/sodapanda/annie/extractors/mgtv"
	"github.com/sodapanda/annie/extractors/miaopai"
	"github.com/sodapanda/annie/extractors/netease"
	"github.com/sodapanda/annie/extractors/pixivision"
	"github.com/sodapanda/annie/extractors/pornhub"
	"github.com/sodapanda/annie/extractors/qq"
	"github.com/sodapanda/annie/extractors/streamtape"
	"github.com/sodapanda/annie/extractors/tangdou"
	"github.com/sodapanda/annie/extractors/tiktok"
	"github.com/sodapanda/annie/extractors/tumblr"
	"github.com/sodapanda/annie/extractors/twitter"
	"github.com/sodapanda/annie/extractors/types"
	"github.com/sodapanda/annie/extractors/udn"
	"github.com/sodapanda/annie/extractors/universal"
	"github.com/sodapanda/annie/extractors/vimeo"
	"github.com/sodapanda/annie/extractors/weibo"
	"github.com/sodapanda/annie/extractors/xvideos"
	"github.com/sodapanda/annie/extractors/yinyuetai"
	"github.com/sodapanda/annie/extractors/youku"
	"github.com/sodapanda/annie/extractors/youtube"
	"github.com/sodapanda/annie/utils"
)

var extractorMap map[string]types.Extractor

func init() {
	douyinExtractor := douyin.New()
	youtubeExtractor := youtube.New()
	stExtractor := streamtape.New()

	extractorMap = map[string]types.Extractor{
		"": universal.New(), // universal extractor

		"douyin":     douyinExtractor,
		"iesdouyin":  douyinExtractor,
		"bilibili":   bilibili.New(),
		"bcy":        bcy.New(),
		"pixivision": pixivision.New(),
		"youku":      youku.New(),
		"youtube":    youtubeExtractor,
		"youtu":      youtubeExtractor, // youtu.be
		"iqiyi":      iqiyi.New(iqiyi.SiteTypeIqiyi),
		"iq":         iqiyi.New(iqiyi.SiteTypeIQ),
		"mgtv":       mgtv.New(),
		"tangdou":    tangdou.New(),
		"tumblr":     tumblr.New(),
		"vimeo":      vimeo.New(),
		"facebook":   facebook.New(),
		"douyu":      douyu.New(),
		"miaopai":    miaopai.New(),
		"163":        netease.New(),
		"weibo":      weibo.New(),
		"instagram":  instagram.New(),
		"twitter":    twitter.New(),
		"qq":         qq.New(),
		"yinyuetai":  yinyuetai.New(),
		"geekbang":   geekbang.New(),
		"pornhub":    pornhub.New(),
		"xvideos":    xvideos.New(),
		"udn":        udn.New(),
		"tiktok":     tiktok.New(),
		"haokan":     haokan.New(),
		"acfun":      acfun.New(),
		"eporner":    eporner.New(),
		"streamtape": stExtractor,
		"streamta":   stExtractor, // streamta.pe
	}
}

// Extract is the main function to extract the data.
func Extract(u string, option types.Options) ([]*types.Data, error) {
	u = strings.TrimSpace(u)
	var domain string

	bilibiliShortLink := utils.MatchOneOf(u, `^(av|BV|ep)\w+`)
	if len(bilibiliShortLink) > 1 {
		bilibiliURL := map[string]string{
			"av": "https://www.bilibili.com/video/",
			"BV": "https://www.bilibili.com/video/",
			"ep": "https://www.bilibili.com/bangumi/play/",
		}
		domain = "bilibili"
		u = bilibiliURL[bilibiliShortLink[1]] + u
	} else {
		u, err := url.ParseRequestURI(u)
		if err != nil {
			return nil, err
		}
		if u.Host == "haokan.baidu.com" {
			domain = "haokan"
		} else {
			domain = utils.Domain(u.Host)
		}
	}
	extractor := extractorMap[domain]
	if extractor == nil {
		extractor = extractorMap[""]
	}
	videos, err := extractor.Extract(u, option)
	if err != nil {
		return nil, err
	}
	for _, v := range videos {
		v.FillUpStreamsData()
	}
	return videos, nil
}
