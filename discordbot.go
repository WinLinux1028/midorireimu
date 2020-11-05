package main

import (
	srand "crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/WinLinux1028/typeconv"
	"github.com/bwmarrin/discordgo"
)

//グローバル変数定義
var (
	prefix  string = "*;"
	adminid        = []string{"704702259665043476"}
)

func main() {
	var Token string = "YOUR TOKEN"

	var dg, err = discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	//ここから起動後に行いたい処理
	fmt.Println("logged in as " + dg.State.User.Username)
	for _, a := range adminid {
		b, _ := dg.UserChannelCreate(a)
		dg.ChannelMessageSend(b.ID, "起動完了")
	}
	dg.UpdateStreamingStatus(1, "Go版テスト", "https://www.youtube.com/watch?v=KcDED7_f258")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Type != discordgo.MessageTypeDefault {
		return
	}
	if len(m.Content) < len(prefix) {
		return
	}
	if m.Content[0:len(prefix)] != prefix {
		return
	}
	var command []string = readcmd(m)
	switch command[0] {
	case "ping":
		ping(s, m, command)
	case "ping2":
		ping2(s, m, command)
	case "野生":
		yasei(s, m, command)
	case "DM":
		anonmsg(s, m, command)
	case "チャンネルトピック":
		chtopic(s, m, command)
	case "チャンネル":
		chsend(s, m, command)
	case "チャンネル2":
		chsend2(s, m, command)
	case "フォロー":
		follow(s, m, command)
	case "kick":
		kick(s, m, command)
	case "ban":
		ban(s, m, command)
	case "パスワード":
		passwd(s, m, command)
	case "サイコロをふる":
		dice(s, m, command)
	case "役職付与":
		giverole(s, m, command)
	}
}

func readcmd(m *discordgo.MessageCreate) (a []string) {
	a = strings.Split(m.Content[len(prefix):], " ")
	return
}

func cmderror(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := recover()
	if err != nil {
		a := searchslice(adminid, m.Author.ID)
		if a == true {
			errstr := fmt.Sprintln(err)
			s.ChannelMessageSend(m.ChannelID, errstr)
		} else {
			s.ChannelMessageSend(m.ChannelID, "何かしらのエラーが起きたようです､構文ミスでないことを確認してから開発者にご相談ください\nまた､DMでは使えないコマンドも存在します")
		}
	}
}

func searchslice(a interface{}, b interface{}) (d bool) {
	switch a.(type) {
	case []string:
		fukugen := a.([]string)
		search := b.(string)
		for _, c := range fukugen {
			if c == search {
				d = true
				break
			}
		}
		return
	case []int:
		fukugen := a.([]int)
		search := b.(int)
		for _, c := range fukugen {
			if c == search {
				d = true
				break
			}
		}
		return
	default:
		return
	}
}

func strfukugen(command []string, b int) (c string) {
	for d, e := range command {
		if d >= b {
			c = c + " " + e
		}
	}
	return
}

func username(user *discordgo.User) (a string) {
	a = "<@" + user.ID + ">(" + user.Username + "#" + user.Discriminator + ")"
	return
}

func ping(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	s.ChannelMessageSend(m.ChannelID, "pong!")
}

func ping2(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	var b *discordgo.Message
	var a = time.Now()
	b, _ = s.ChannelMessageSend(m.ChannelID, "計測中……!")
	var c = time.Since(a)
	s.ChannelMessageEdit(m.ChannelID, b.ID, "pong！\n結果:**"+typeconv.Stringc(float64(c)/1000000000)+"**秒ですฅ✧！")
}

func yasei(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	mrand.Seed(time.Now().UnixNano())
	var embed = discordgo.MessageEmbed{
		Title: "あ！",
		Color: mrand.Intn(0xffffff),
	}
	var a *discordgo.Channel
	a, _ = s.Channel(m.ChannelID)
	if a.Type != 0 {
		embed.Description = ("野生の" + m.Author.Username + "が飛び出してきた！")
	} else {
		if m.Member.Nick != "" {
			embed.Description = ("野生の" + m.Author.Username + "(" + m.Member.Nick + ")が飛び出してきた！")
		} else {
			embed.Description = ("野生の" + m.Author.Username + "が飛び出してきた！")
		}
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}

func anonmsg(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	b, _ := s.UserChannelCreate(command[1])
	f := strfukugen(command, 2)
	s.ChannelMessageSend(b.ID, f)
	s.ChannelMessageSend(m.ChannelID, "あなたのメッセージ､届けましたよ")
}

func chtopic(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	a, _ := s.State.UserChannelPermissions(m.Author.ID, command[1])
	if a&discordgo.PermissionManageChannels == discordgo.PermissionManageChannels {
		c := discordgo.ChannelEdit{}
		c.Topic = strfukugen(command, 2)
		s.ChannelEditComplex(command[1], &c)
	} else {
		s.ChannelMessageSend(m.ChannelID, "何様のつもりですか...?")
	}
}

func chsend(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	a, _ := s.State.UserChannelPermissions(m.Author.ID, command[1])
	if a&discordgo.PermissionSendMessages == discordgo.PermissionSendMessages {
		s.ChannelMessageSend(command[1], strfukugen(command, 2))
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	} else {
		s.ChannelMessageSend(m.ChannelID, "このチャンネルにメッセージを送信する権限がありません")
	}
}

func chsend2(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	s.ChannelMessageSend(m.ChannelID, strfukugen(command, 1))
	s.ChannelMessageDelete(m.ChannelID, m.ID)
}

func follow(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	a, _ := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if a&discordgo.PermissionManageWebhooks == discordgo.PermissionManageWebhooks {
		b, _ := s.Channel(command[1])
		if b.Type == discordgo.ChannelTypeGuildNews {
			s.ChannelNewsFollow(command[1], m.ChannelID)
			s.ChannelMessageSend(m.ChannelID, "アナウンスチャンネルをフォローした。\nいらなくなったら運営に頼んで消して貰ってね。")
		} else {
			s.ChannelMessageSend(m.ChannelID, "アナウンスチャンネルじゃないよー")
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Webhookの操作権限がありません")
	}
}

func kick(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	a, _ := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if a&discordgo.PermissionKickMembers == discordgo.PermissionKickMembers {
		reason := strfukugen(command, 2)
		b, _ := s.User(command[1])
		if reason == "" {
			reason = "未指定"
		}
		s.GuildMemberDeleteWithReason(m.GuildID, command[1], reason)
		s.ChannelMessageSend(m.ChannelID, "実行者："+username(m.Author)+"\n"+username(b)+"をキックした。\n理由："+reason)
	} else {
		s.ChannelMessageSend(m.ChannelID, "キック権限がありません")
	}
}

func ban(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	a, _ := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if a&discordgo.PermissionBanMembers == discordgo.PermissionBanMembers {
		reason := strfukugen(command, 2)
		b, _ := s.User(command[1])
		if reason == "" {
			reason = "未指定"
		}
		s.GuildBanCreateWithReason(m.GuildID, command[1], reason, 0)
		s.ChannelMessageSend(m.ChannelID, "実行者："+username(m.Author)+"\n"+username(b)+"をBANした。\n理由："+reason)
	} else {
		s.ChannelMessageSend(m.ChannelID, "BAN権限がありません")
	}
}

func passwd(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	if command[1] == "--help" {
		s.ChannelMessageSend(m.ChannelID, "パスワード生成コマンドです\n基本的な使い方: パスワード 桁数\n桁数の前に入れるオプション:\n--no-spchar 記号を含まない\n--only-number 数字のみ")
		return
	}
	check, _ := s.Channel(m.ChannelID)
	if check.Type != discordgo.ChannelTypeDM {
		s.ChannelMessageSend(m.ChannelID, "このコマンドはDM以外で実行するべきではありません､このBOTとのDM上で実行してください")
		return
	}
	var letters string
	var a int
	var max int
	var passwd string
	if command[1] == "--only-number" {
		letters = "0123456789"
		a = typeconv.Intc(command[2])
		max = 2000
	} else if command[1] == "--no-spchar" {
		letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		a = typeconv.Intc(command[2])
		max = 2000
	} else {
		passwd = "`"
		letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\"#$%&'\\()*+,-./:;<=>?@[]^_{|}~"
		a = typeconv.Intc(command[1])
		max = 1998
	}
	if a > 0 && a <= max {
		i := 0
		for i < a {
			letterslen := big.NewInt(int64(len(letters) - 1))
			b, _ := srand.Int(srand.Reader, letterslen)
			c := b.Int64()
			passwd = passwd + letters[c:c+1]
			i = i + 1
		}
		if passwd[0:1] == "`" {
			passwd = passwd + "`"
		}
		s.ChannelMessageSend(m.ChannelID, "パスワードの生成が完了しました､生成されたパスワードは")
		s.ChannelMessageSend(m.ChannelID, passwd)
		s.ChannelMessageSend(m.ChannelID, "です")
	} else {
		maxstr := typeconv.Stringc(max)
		s.ChannelMessageSend(m.ChannelID, "1文字以上､"+maxstr+"文字以下にしてください")
	}
}

func dice(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	mrand.Seed(time.Now().UnixNano())
	liststr := strings.Split(command[1], "d")
	var list []int
	for _, i := range liststr {
		list = append(list, typeconv.Intc(i))
	}
	number := typeconv.Stringc(mrand.Intn(list[1]-list[0]+1) + list[0])
	s.ChannelMessageSend(m.ChannelID, number)
}

func giverole(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	a, _ := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if a&discordgo.PermissionManageRoles == discordgo.PermissionManageRoles {
		s.GuildMemberRoleAdd(m.GuildID, command[1], command[2])
		s.ChannelMessageSend(m.ChannelID, "TEST")
	} else {
		s.ChannelMessageSend(m.ChannelID, "ロール管理権限がありません")
	}
}
