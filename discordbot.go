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
	prefix  string = "*;"
	adminid        = []string{"704702259665043476"}
)

func main() {
	token, _ := os.Executable()
	f, err := os.Open(filepath.Dir(token) + "/../discordtoken.txt")
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
		amountofmember(s, m)
	case "ユーザー人数":
		amountofhuman(s, m)
	case "bot人数":
		amountofbot(s, m)
	case "test":
		test(s, m, command)
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
	var embed = discordgo.MessageEmbed{
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
	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
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
	a, _ := s.State.UserChannelPermissions(m.Author.ID, command[1])
	if a&discordgo.PermissionManageChannels == discordgo.PermissionManageChannels {
		c := discordgo.ChannelEdit{}
		c.Topic = strfukugen(command, 2)
		_, err := s.ChannelEditComplex(command[1], &c)
		if err != nil {
			panic(err)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "何様のつもりですか...?")
	}
}

func chsend(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	defer cmderror(s, m)
	a, _ := s.State.UserChannelPermissions(m.Author.ID, command[1])
	check, _ := s.State.Channel(m.ChannelID)
	if a&discordgo.PermissionSendMessages == discordgo.PermissionSendMessages || check.Type == discordgo.ChannelTypeDM {
		_, err := s.ChannelMessageSend(command[1], strfukugen(command, 2))
		if err != nil {
			panic(err)
		}
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
	var list []int
	for _, i := range liststr {
		list = append(list, typeconv.Intc(i))
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
	if a&discordgo.PermissionManageChannels == discordgo.PermissionManageChannels || check.Type == discordgo.ChannelTypeDM {
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
	embed := discordgo.MessageEmbed{
		Title:       "時間です。よく見ておいてくださいね。",
		Color:       mrand.Intn(0xffffff),
		Description: formattime(time.Now()),
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
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
	embed := discordgo.MessageEmbed{
		Title:       "短縮リンク",
		Color:       mrand.Intn(0xffffff),
		Description: string(respbyte),
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
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
	embed := discordgo.MessageEmbed{
		Title:       "該当ユーザーが居る場所",
		Description: guildnames,
		Color:       mrand.Intn(0xffffff),
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
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

func amountofmember(s *discordgo.Session, m *discordgo.MessageCreate) {
	a, _ := s.State.Guild(m.GuildID)
	s.ChannelMessageSend(m.ChannelID, "メンバー数:"+typeconv.Stringc(len(a.Members)))
}

func amountofhuman(s *discordgo.Session, m *discordgo.MessageCreate) {
	a, _ := s.State.Guild(m.GuildID)
	var i int
	for _, member := range a.Members {
		if member.User.Bot != true {
			i = i + 1
		}
	}
	s.ChannelMessageSend(m.ChannelID, "ユーザー数:"+typeconv.Stringc(i))
}

func amountofbot(s *discordgo.Session, m *discordgo.MessageCreate) {
	a, _ := s.State.Guild(m.GuildID)
	var i int
	for _, member := range a.Members {
		if member.User.Bot == true {
			i = i + 1
		}
	}
	s.ChannelMessageSend(m.ChannelID, "BOT数:"+typeconv.Stringc(i))
}

func test(s *discordgo.Session, m *discordgo.MessageCreate, command []string) {
	fmt.Println(strfukugen(command, 0))
}
