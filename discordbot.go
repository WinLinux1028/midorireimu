package main

import (
	"bufio"
	srand "crypto/rand"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	mrand "math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/v3/cpu"

	"github.com/WinLinux1028/dgconv"
	"github.com/WinLinux1028/typeconv"
	"github.com/bwmarrin/discordgo"
)

//グローバル変数定義
var (
	prefix     string = "*;"
	sc         chan os.Signal
	adminid    []string = []string{"704702259665043476"}
	bugrep     string   = "777459366555156501"
	globalname string   = "青霊夢_test"
)

func init() {
	token, _ := os.Executable()
	f, err := os.Open(filepath.Dir(token) + "/discordtoken.txt")
	if err != nil {
		fmt.Println("Input your bot token to (executable file directory)/../discordtoken.txt")
		return
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		token = scanner.Text()
		break
	}
	f.Close()

	var dg, err2 = discordgo.New("Bot " + token)
	if err2 != nil {
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
	go func(s *discordgo.Session) {
		for {
			now := time.Now()
			var future time.Time
			if now.Hour() > 8 {
				future = time.Date(now.Year(), now.Month(), now.Day()+1, 8, 10, 0, 0, time.Local)
			} else if now.Hour() < 8 {
				future = time.Date(now.Year(), now.Month(), now.Day(), 8, 10, 0, 0, time.Local)
			} else if now.Hour() == 8 {
				if now.Minute() >= 10 {
					future = time.Date(now.Year(), now.Month(), now.Day()+1, 8, 10, 0, 0, time.Local)
				} else if now.Minute() < 10 {
					future = time.Date(now.Year(), now.Month(), now.Day(), 8, 10, 0, 0, time.Local)
				}
			}
			time.Sleep(future.Sub(now))
			for _, send := range s.State.Guilds {
				if send.SystemChannelID != "" {
					s.ChannelMessageSend(send.SystemChannelID, "野獣の時間だよぉ!")
				}
			}
		}
	}(dg)

	sc = make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
}

func main() {
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Type != discordgo.MessageTypeDefault {
		return
	}
	channel, _ := s.State.Channel(m.ChannelID)
	if channel.Name == globalname {
		globalchat(s, m)
		return
	}
	if len(m.Content) < len(prefix) {
		return
	}
	if strings.Contains(strings.ToLower(m.Content), "ypa") == true {
		s.ChannelMessageSend(m.ChannelID, "お前スパイだろ､粛清(正しくはУраね)")
	}
	if m.Content[0:len(prefix)] != prefix {
		return
	}
	var command []string = readcmd(m)
	switch command[0] {
	case "ping":
		ping(s, m)
	case "ping2":
		ping2(s, m)
	case "野生":
		yasei(s, m)
	case "DM":
		anonmsg(s, m, command)
	case "チャンネルトピック":
		chtopic(s, m, command)
	case "チャンネル":
		chsend(s, m, command)
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
	case "ピン":
		pin(s, m, command)
	case "時間確認":
		timecheck(s, m)
	case "脱出":
		bye(s, m)
	case "リンク":
		link(s, m, command)
	case "uuser":
		uuser(s, m, command)
	case "使用率":
		memorycheck(s, m)
	case "全体人数":
		amountofmember(s, m, command)
	case "ユーザー人数":
		amountofhuman(s, m, command)
	case "bot人数":
		amountofbot(s, m, command)
	case "チャンネル確認":
		chcheck(s, m, command)
	case "役職持ち確認":
		rolecheck(s, m, command)
	case "vcから切断":
		vcremove(s, m, command)
	case "end":
		shutdown(s, m)
	case "ランダムユーザー":
		randuser(s, m)
	case "導入数確認":
		botcheck(s, m)
	case "ui":
		ui(s, m, command)
	case "鯖知りたい":
		guildstate(s, m, command)
	case "sh":
		shell(s, m, command)
	case "バグ報告":
		bugreport(s, m, command)
	case "help":
		help(s, m)
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
	for c[0:1] == " " {
		c = c[1:]
	}
	return
}

func username(user *discordgo.User) (a string) {
	a = "<@" + user.ID + ">(" + user.Username + "#" + user.Discriminator + ")"
	return
}

func formattime(time time.Time) (s string) {
	s = time.Format("2006年1月2日 15時4分5秒.999999999")
	i := strings.Index(s, "秒")
	s = s[0:i] + "秒" + s[i+len("秒")+1:]
	return
}

//グローバルチャット
func globalchat(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer cmderror(s, m)
	nowguild, _ := s.State.Guild(m.GuildID)
	embed := make([]*discordgo.MessageEmbed, 0, 10)
	embed = append(embed, &discordgo.MessageEmbed{
		Description: m.Content,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    m.Author.Username + "#" + m.Author.Discriminator + "@" + nowguild.Name,
			IconURL: m.Author.AvatarURL("4096"),
		},
	})
	if len(m.Attachments) == 0 {
	} else {
		for i, att := range m.Attachments {
			if i == 0 {
				embed[0].Image = &discordgo.MessageEmbedImage{
					URL: att.ProxyURL,
				}
			} else {
				embed = append(embed, &discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{
						Name:    m.Author.Username + "#" + m.Author.Discriminator + "@" + nowguild.Name,
						IconURL: m.Author.AvatarURL("4096"),
					},
					Image: &discordgo.MessageEmbedImage{
						URL: att.ProxyURL,
					},
				})
			}
		}
	}
	for _, guilds := range s.State.Guilds {
		for _, channels := range guilds.Channels {
			if channels.Type == discordgo.ChannelTypeGuildText {
				if channels.Name == globalname {
					for _, embeds := range embed {
						s.ChannelMessageSendEmbed(channels.ID, embeds)
					}
					break
				}
			}
		}
	}
	s.ChannelMessageDelete(m.ChannelID, m.ID)
}

//ここからコマンド
func ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer cmderror(s, m)
	s.ChannelMessageSend(m.ChannelID, "pong!")
}

func ping2(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer cmderror(s, m)
	var b *discordgo.Message
	var a = time.Now()
	b, _ = s.ChannelMessageSend(m.ChannelID, "計測中……!")
	var c = time.Since(a)
	s.ChannelMessageEdit(m.ChannelID, b.ID, "pong！\n結果:**"+typeconv.Stringc(math.Round(float64(c)/1000000)/1000)+"**秒ですฅ✧！")
}

func yasei(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer cmderror(s, m)
	mrand.Seed(time.Now().UnixNano())
	var embed = &discordgo.MessageEmbed{
		Title: "あ！",
		Color: mrand.Intn(0xffffff),
	}
	var a *discordgo.Channel
	a, _ = s.State.Channel(m.ChannelID)
	if a.Type != 0 {
		embed.Description = ("野生の" + m.Author.Username + "が飛び出してきた！")
	} else {
		if m.Member.Nick != "" {
			embed.Description = ("野生の" + m.Author.Username + "(" + m.Member.Nick + ")が飛び出してきた！")
		} else {
			embed.Description = ("野生の" + m.Author.Username + "が飛び出してきた！")
		}
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func anonmsg(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	b, _ := s.UserChannelCreate(dgconv.Getuser(s, command[1]))
	f := strfukugen(command, 2)
	_, err := s.ChannelMessageSend(b.ID, f)
	if err != nil {
		panic(err)
	}
	s.ChannelMessageSend(m.ChannelID, "あなたのメッセージ､届けましたよ")
}

func chtopic(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	var ch string
	var mode bool
	_, err := s.State.Channel(dgconv.Getchannel(s, command[1]))
	if err != nil {
		ch = m.ChannelID
		mode = false
	} else {
		ch = dgconv.Getchannel(s, command[1])
		mode = true
	}
	a, _ := s.State.UserChannelPermissions(m.Author.ID, ch)
	if a&discordgo.PermissionManageChannels == discordgo.PermissionManageChannels {
		c := &discordgo.ChannelEdit{}
		if mode {
			c.Topic = strfukugen(command, 2)
		} else {
			c.Topic = strfukugen(command, 1)
		}
		_, err = s.ChannelEditComplex(ch, c)
		if err != nil {
			panic(err)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "何様のつもりですか...?")
	}
}

func chsend(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	var mode int
	channel := dgconv.Getchannel(s, command[1])
	if channel == "" {
		channel = m.ChannelID
		mode = 1
	} else {
		mode = 2
	}
	a, _ := s.State.UserChannelPermissions(m.Author.ID, channel)
	check, _ := s.State.Channel(channel)
	if a&discordgo.PermissionSendMessages == discordgo.PermissionSendMessages || check.Type != discordgo.ChannelTypeGuildText {
		_, err := s.ChannelMessageSend(channel, strfukugen(command, mode))
		if err != nil {
			panic(err)
		}
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	} else {
		s.ChannelMessageSend(m.ChannelID, "このチャンネルにメッセージを送信する権限がありません")
	}
}

func follow(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	a, _ := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if a&discordgo.PermissionManageWebhooks == discordgo.PermissionManageWebhooks {
		b, _ := s.State.Channel(command[1])
		if b.Type == discordgo.ChannelTypeGuildNews {
			_, err := s.ChannelNewsFollow(command[1], m.ChannelID)
			if err != nil {
				panic(err)
			}
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
		b, _ := s.User(dgconv.Getuser(s, command[1]))
		if reason == "" {
			reason = "未指定"
		}
		err := s.GuildMemberDeleteWithReason(m.GuildID, dgconv.Getuser(s, command[1]), reason)
		if err != nil {
			panic(err)
		}
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
		b, _ := s.User(dgconv.Getuser(s, command[1]))
		if reason == "" {
			reason = "未指定"
		}
		err := s.GuildBanCreateWithReason(m.GuildID, dgconv.Getuser(s, command[1]), reason, 0)
		if err != nil {
			panic(err)
		}
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
	check, _ := s.State.Channel(m.ChannelID)
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
	list := make([]int, 2, 2)
	for a, i := range liststr {
		list[a] = typeconv.Intc(i)
	}
	number := typeconv.Stringc(mrand.Intn(list[1]-list[0]+1) + list[0])
	s.ChannelMessageSend(m.ChannelID, number)
}

func giverole(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	a, _ := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if a&discordgo.PermissionManageRoles == discordgo.PermissionManageRoles {
		err := s.GuildMemberRoleAdd(m.GuildID, dgconv.Getuser(s, command[1]), dgconv.Getrole(s, m, strfukugen(command, 2)))
		if err != nil {
			panic(err)
		}
		user, _ := s.User(dgconv.Getuser(s, command[1]))
		role, _ := s.State.Role(m.GuildID, dgconv.Getrole(s, m, strfukugen(command, 2)))
		s.ChannelMessageSend(m.ChannelID, username(user)+"さんに"+"<@&"+role.ID+">("+role.Name+")を付与しました")
	} else {
		s.ChannelMessageSend(m.ChannelID, "ロール管理権限がありません")
	}
}

func pin(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	a, _ := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	check, _ := s.State.Channel(m.ChannelID)
	if a&discordgo.PermissionManageChannels == discordgo.PermissionManageChannels || check.Type != discordgo.ChannelTypeGuildText {
		msgs, _ := s.ChannelMessagesPinned(m.ChannelID)
		var mode int
		for _, msg := range msgs {
			if msg.ID == command[1] {
				mode = 1
				break
			}
		}
		if mode == 0 {
			err := s.ChannelMessagePin(m.ChannelID, command[1])
			if err != nil {
				panic(err)
			}
			s.ChannelMessageSend(m.ChannelID, "ピンをしました。："+username(m.Author))
		} else {
			err := s.ChannelMessageUnpin(m.ChannelID, command[1])
			if err != nil {
				panic(err)
			}
			s.ChannelMessageSend(m.ChannelID, "ピンを外しました。："+username(m.Author))
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "メッセージ管理権限がありません")
	}
}

func timecheck(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer cmderror(s, m)
	mrand.Seed(time.Now().UnixNano())
	embed := &discordgo.MessageEmbed{
		Title:       "時間です。よく見ておいてくださいね。",
		Color:       mrand.Intn(0xffffff),
		Description: formattime(time.Now()),
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func bye(s *discordgo.Session, m *discordgo.MessageCreate) {
	a, _ := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if a&discordgo.PermissionManageChannels == discordgo.PermissionManageChannels {
		guild, _ := s.State.Guild(m.GuildID)
		s.ChannelMessageSend(m.ChannelID, guild.Name+"("+guild.ID+")から退室しました。")
		s.GuildLeave(m.GuildID)
	} else {
		s.ChannelMessageSend(m.ChannelID, "何様のつもりですか...?")
	}
}

func link(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	resp, err := http.Get("https://is.gd/create.php?format=simple&url=" + strfukugen(command, 1))
	if err != nil {
		resp.Body.Close()
		panic(err)
	}
	defer resp.Body.Close()
	respbyte, _ := ioutil.ReadAll(resp.Body)
	mrand.Seed(time.Now().UnixNano())
	embed := &discordgo.MessageEmbed{
		Title:       "短縮リンク",
		Color:       mrand.Intn(0xffffff),
		Description: string(respbyte),
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	s.ChannelMessageDelete(m.ChannelID, m.ID)
}

func uuser(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	mrand.Seed(time.Now().UnixNano())
	var guildnames string
	thisuser := dgconv.Getuser(s, command[1])
	for _, guild := range s.State.Guilds {
		for _, mem := range guild.Members {
			if mem.User.ID == thisuser {
				guildnames = guildnames + guild.Name + "\n"
				break
			}
		}
	}
	embed := &discordgo.MessageEmbed{
		Title:       "該当ユーザーが居る場所",
		Description: guildnames,
		Color:       mrand.Intn(0xffffff),
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func memorycheck(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer cmderror(s, m)
	if searchslice(adminid, m.Author.ID) {
		mem, _ := mem.VirtualMemory()
		cpuper, _ := cpu.Percent(1000000000, false)
		cpupers := typeconv.Stringc(math.Round(cpuper[0]*100) / 100)
		allmem := typeconv.Stringc(math.Round(float64(mem.Total/10000000)) / 100)
		used := typeconv.Stringc(math.Round(float64(mem.Used/10000000)) / 100)
		usedp := typeconv.Stringc(math.Round((float64(mem.Used)/float64(mem.Total))*10000) / 100)
		free := typeconv.Stringc(math.Round(float64(mem.Available/10000000)) / 100)
		freep := typeconv.Stringc(100 - typeconv.Float64c(usedp))
		s.ChannelMessageSend(m.ChannelID, "CPU使用率:"+cpupers+"%\n全てのメモリ容量:"+allmem+"GB\n使用量:"+used+"GB("+usedp+"%)\n空き容量:"+free+"GB("+freep+"%)")
	} else {
		s.ChannelMessageSend(m.ChannelID, "だが断る()")
	}
}

func amountofmember(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	var a *discordgo.Guild
	if len(command) == 1 {
		a, _ = s.State.Guild(m.GuildID)
	} else {
		a, _ = s.State.Guild(command[1])
	}
	s.ChannelMessageSend(m.ChannelID, "メンバー数:"+typeconv.Stringc(len(a.Members)))
}

func amountofhuman(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	var a *discordgo.Guild
	if len(command) == 1 {
		a, _ = s.State.Guild(m.GuildID)
	} else {
		a, _ = s.State.Guild(command[1])
	}
	var i int
	for _, member := range a.Members {
		if member.User.Bot != true {
			i = i + 1
		}
	}
	s.ChannelMessageSend(m.ChannelID, "ユーザー数:"+typeconv.Stringc(i))
}

func amountofbot(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	var a *discordgo.Guild
	if len(command) == 1 {
		a, _ = s.State.Guild(m.GuildID)
	} else {
		a, _ = s.State.Guild(command[1])
	}
	var i int
	for _, member := range a.Members {
		if member.User.Bot == true {
			i = i + 1
		}
	}
	s.ChannelMessageSend(m.ChannelID, "BOT数:"+typeconv.Stringc(i))
}

func chcheck(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	mrand.Seed(time.Now().UnixNano())
	var i string
	if len(command) == 1 {
		i = m.GuildID
	} else {
		i = command[1]
	}
	g, _ := s.State.Guild(i)
	embed := &discordgo.MessageEmbed{
		Title:       "チャンネル数:" + typeconv.Stringc(len(g.Channels)),
		Description: "500になったら作れません",
		Color:       mrand.Intn(0xffffff),
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func rolecheck(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	mrand.Seed(time.Now().UnixNano())
	var member string
	var kazu int
	role := dgconv.Getrole(s, m, strfukugen(command, 1))
	guild, _ := s.State.Guild(m.GuildID)
	for _, m := range guild.Members {
		if searchslice(m.Roles, role) {
			member = member + m.User.Username + "\n"
			kazu = kazu + 1
		}
	}
	role2, _ := s.State.Role(m.GuildID, role)
	embed := &discordgo.MessageEmbed{
		Title:       role2.Name + "を持つメンバー一覧\n人数:" + typeconv.Stringc(kazu),
		Description: member,
		Color:       mrand.Intn(0xffffff),
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func vcremove(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	a, _ := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if a&discordgo.PermissionVoiceMoveMembers == discordgo.PermissionVoiceMoveMembers {
		err := s.GuildMemberMove(m.GuildID, dgconv.Getuser(s, strfukugen(command, 1)), nil)
		if err != nil {
			panic(err)
		}
		user, _ := s.User(dgconv.Getuser(s, strfukugen(command, 1)))
		s.ChannelMessageSend(m.ChannelID, username(user)+"をボイスチャンネルから切断しました")
	} else {
		s.ChannelMessageSend(m.ChannelID, "メンバー移動権限がありません")
	}
}

func shutdown(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer cmderror(s, m)
	if searchslice(adminid, m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "Shutting down...")
		sc <- syscall.SIGINT
	} else {
		s.ChannelMessageSend(m.ChannelID, "何様のつもりですか...?")
	}
}

func randuser(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer cmderror(s, m)
	mrand.Seed(time.Now().UnixNano())
	var users []*discordgo.Member
	for _, guild := range s.State.Guilds {
		users = append(users, guild.Members...)
	}
	user := users[mrand.Intn(len(users))].User
	embed := &discordgo.MessageEmbed{
		Title:       "誰が出るかな?",
		Description: username(user),
		Color:       mrand.Intn(0xffffff),
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: user.AvatarURL("4096")},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func botcheck(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer cmderror(s, m)
	mrand.Seed(time.Now().UnixNano())
	var members int
	var guilds int
	for _, guild := range s.State.Guilds {
		guilds = guilds + 1
		members = members + guild.MemberCount
	}
	embed := &discordgo.MessageEmbed{
		Title: "サーバー数:" + typeconv.Stringc(guilds) + "\nメンバー数:" + typeconv.Stringc(members),
		Color: mrand.Intn(0xffffff),
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func ui(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	var id string
	if len(command) != 2 {
		id = m.Author.ID
	} else {
		id = dgconv.Getuser(s, command[1])
	}
	member := dgconv.Getmember(s, id)
	user, _ := s.User(id)
	ch, _ := s.State.Channel(m.ChannelID)
	createdtime, _ := discordgo.SnowflakeTimestamp(id)
	status, _ := s.State.Presence(member.GuildID, id)
	fields := make([]*discordgo.MessageEmbedField, 0, 10)
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "名前",
		Value: user.Username,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "ID",
		Value: user.ID,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "タグ",
		Value: user.Discriminator,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "BOT",
		Value: typeconv.Stringc(user.Bot),
	})
	if ch.Type == discordgo.ChannelTypeGuildText {
		if member.Nick != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  "ニックネーム",
				Value: member.Nick,
			})
		}
	}
	if status.Game != nil {
		if status.Game.State != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  "アクティビティ",
				Value: status.Game.State,
			})
		} else if status.Game.Name != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  "アクティビティ",
				Value: status.Game.Name,
			})
		}
	}
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "ステータス",
		Value: string(status.Status),
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "アカウント作成日",
		Value: formattime(createdtime),
	})
	if ch.Type == discordgo.ChannelTypeGuildText {
		var role string
		var role2 *discordgo.Role
		for _, roles := range member.Roles {
			role2, _ = s.State.Role(m.GuildID, roles)
			role = role + role2.Name + ","
		}
		role = role[:len(role)-1]
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "ロール",
			Value: role,
		})
	}

	mrand.Seed(time.Now().UnixNano())
	embed := &discordgo.MessageEmbed{
		Title: user.Username + "の情報",
		Color: mrand.Intn(0xffffff),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: user.AvatarURL("4096"),
		},
		Fields: fields,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func guildstate(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	var err error
	var guild string
	if len(command) != 2 {
		guild = m.GuildID
	} else {
		guild = command[1]
	}
	guild2, _ := s.State.Guild(guild)
	mrand.Seed(time.Now().UnixNano())
	fields := make([]*discordgo.MessageEmbedField, 0, 13)
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "名前",
		Value: guild2.Name,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "ID",
		Value: guild,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "地域",
		Value: guild2.Region,
	})
	time, _ := discordgo.SnowflakeTimestamp(guild)
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "作成日",
		Value: formattime(time),
	})
	owner, _ := s.User(guild2.OwnerID)
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "オーナー",
		Value: username(owner),
	})
	var textch int
	var voicech int
	var catch int
	for _, ch := range guild2.Channels {
		if ch.Type == discordgo.ChannelTypeGuildText || ch.Type == discordgo.ChannelTypeGuildNews {
			textch++
		} else if ch.Type == discordgo.ChannelTypeGuildVoice {
			voicech++
		} else if ch.Type == discordgo.ChannelTypeGuildCategory {
			catch++
		}
	}
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "テキストチャンネル数",
		Value:  typeconv.Stringc(textch),
		Inline: true,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "ボイスチャンネル数",
		Value:  typeconv.Stringc(voicech),
		Inline: true,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "カテゴリ数",
		Value:  typeconv.Stringc(catch),
		Inline: true,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "合計チャンネル数",
		Value:  typeconv.Stringc(textch + voicech),
		Inline: true,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "MFA必須",
		Value:  typeconv.Stringc(int(guild2.MfaLevel)),
		Inline: true,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "認証レベル",
		Value:  typeconv.Stringc(int(guild2.VerificationLevel)),
		Inline: true,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "ブーストレベル",
		Value:  typeconv.Stringc(int(guild2.PremiumTier)),
		Inline: true,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "ブーストした人数",
		Value:  typeconv.Stringc(guild2.PremiumSubscriptionCount),
		Inline: true,
	})
	embed := &discordgo.MessageEmbed{
		Title: guild2.Name + "の情報",
		Color: mrand.Intn(0xffffff),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: guild2.IconURL(),
		},
		Fields: fields,
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		panic(err)
	}
}

func shell(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	if searchslice(adminid, m.Author.ID) {
		cmd := exec.Command(os.Getenv("SHELL"), "-c", strfukugen(command, 1))
		a, _ := cmd.Output()
		if string(a) == "" {
			s.ChannelMessageSend(m.ChannelID, "出力がありませんでした")
			return
		}
		var sw bool
		var str string
		for _, strs := range strings.Split(string(a), "\n") {
			if len(strs) <= 2000 {
				if len(str+strs) <= 2000 {
					str = str + strs + "\n"
					sw = true
				} else {
					s.ChannelMessageSend(m.ChannelID, str)
					str = strs + "\n"
					sw = true
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, str)
				s.ChannelMessageSend(m.ChannelID, strs[0:1997]+"...")
				str = ""
				sw = false
			}
		}
		if sw == true {
			s.ChannelMessageSend(m.ChannelID, str)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "何様のつもりですか...?")
	}
}

func bugreport(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	mrand.Seed(time.Now().UnixNano())
	check, _ := s.State.Channel(m.ChannelID)
	desc := "報告内容:\n" + strfukugen(command, 1) + "\n報告者: " + username(m.Author) + "\nサーバー: "
	if check.Type == discordgo.ChannelTypeGuildText {
		guild, _ := s.State.Guild(m.GuildID)
		desc = desc + guild.Name + "(" + m.GuildID + ")"
	} else {
		desc = desc + "DM(" + m.ChannelID + ")"
	}
	embed := &discordgo.MessageEmbed{
		Title:       "意見ありがとうございます",
		Description: desc,
		Color:       mrand.Intn(0xffffff),
	}
	s.ChannelMessageSendEmbed(bugrep, embed)
	s.ChannelMessageSend(m.ChannelID, "参考にします")
}

func help(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title: "ヘルプコマンドです",
		Description: `___**パスワード**___:指定した桁数のパスワードを生成します､詳しい使い方は'パスワード --help'を実行してください
		___**DM**___:指定した人に匿名でDMを送ります(使い方:DM @ユーザー メッセージ)
		___**ピン**___：ピン留めします(使い方: ピン メッセージID)
		___**ping**___：動作確認します
		___**ping2**___：応答速度を確認します
		___**ui**___：ユーザーを調べます(使い方：ui @ユーザー)
		___**uuser**___：ユーザー調査です！(使い方：uuser @ユーザー)
		___**サイコロをふる**___：サイコロをふります！(使い方：サイコロをふる 1d6)
		___**チャンネル**___：指定したチャンネルに書き込みます！(使い方：チャンネル #(チャンネル名) てすと)(チャンネル名の指定がなければ現在のチャンネルに書き込まれます)
		___**チャンネル確認**___：チャンネル数を確認します！
		___**全体人数**___：サーバー全体の人数を調べます！(使い方: 全体人数 サーバーID)(サーバーIDがない場合は現在のサーバーになります)
		___**ユーザー人数**___：ユーザーの人数を調べます！(使い方: ユーザー人数 サーバーID)(サーバーIDがない場合は現在のサーバーになります)
		___**bot人数**___：BOTの数を調べます (使い方: bot人数 サーバーID)(サーバーIDがない場合は現在のサーバーになります)
		___**リンク**___：短縮リンクを作ります！(使い方：リンク URL)
		___**フォロー**___：アナウンスチャンネルをフォローします
		___**役職持ち確認**___：ロール所持者を確認します！(使い方：役職持ち確認 @ロール)
		___**時間確認**___：時間を確認します！
		___**野生**___：ネタコマンドです！
		___**鯖知りたい**___：サーバーの情報を知ることができます！(使い方：鯖知りたい サーバーID)
		___**ランダムユーザー**___: ネタコマンドその2
		___**バグ報告**___: バグ報告します
		`,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func help2(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title: "運営用ヘルプコマンドです",
		Description: `___**チャンネルトピック**___：チャンネルトピックをいじります！(使い方：チャンネルトピック #チャンネル名 てすと)(チャンネル名がない場合は現在のチャンネル)
		___**kick**___：対象者を蹴ります！(使い方：kick ユーザーID 理由)(理由がない場合は未指定となります)
		___**ban**___：対象者をBANします！(使い方：ban ユーザーID 理由)(理由がない場合は未指定となります)
		___**役職付与**___：ロールを付与します！(使い方：役職付与 @メンバー @ロール)
		___**脱出**___：サーバーから退室します！
		___**使用率**___：使用率を調べます！
		___**vcから切断**___: VCから強制切断します！(使い方：vcから切断 @メンバー)
		___**end**___：BOTを終了させます！
		___**導入数確認**___: BOTがいるサーバーの数を調べます
		`,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
