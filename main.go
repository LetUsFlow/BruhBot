package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

var joinedServers []string
var db *sql.DB

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + "ODI2MTgxMTc0NjczNjcwMTQ0.YGIvLQ.SKLXUWuiZ9ZqvddDZ5iTkbIssrg")
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Connect to statistics database
	db, err = sql.Open("sqlite3", "./statistics.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	_ = dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	for i := 21; i >= 1; i-- {
		if voiceMessageHandler(s, m, fmt.Sprintf("%s%d", "moan", i), fmt.Sprintf("%s%d%s", "sounds/moans/moan", i, ".mp3"), false) {
			return
		}
	}
	if voiceMessageHandler(s, m, "mo", "sounds/mo.mp3", true) {
		return
	}
	if voiceMessageHandler(s, m, "ma", "sounds/ma.mp3", true) {
		return
	}
	if voiceMessageHandler(s, m, "ginf", "sounds/ginf.mp3", false) {
		return
	}
	if voiceMessageHandler(s, m, "teams", "sounds/teams.mp3", false) {
		return
	}
	if voiceMessageHandler(s, m, "okay", "sounds/okay.mp3", false) {
		return
	}
	if voiceMessageHandler(s, m, "yeet", "sounds/yeet.mp3", false) {
		return
	}
	if voiceMessageHandler(s, m, "marcos", "sounds/marcos.mp3", false) {
		return
	}
	if voiceMessageHandler(s, m, "outlook", "sounds/outlook.mp3", false) {
		return
	}
	if voiceMessageHandler(s, m, "bonk", "sounds/bonk.mp3", false) {
		return
	}
	if voiceMessageHandler(s, m, "moan", "sounds/bonk.mp3", false) {
		return
	}
	if voiceMessageHandler(s, m, "bruh", "sounds/bruh.mp3", false) {
		return
	}
	if voiceMessageHandler(s, m, "bann", "sounds/ban_den_weg.mp3", false) {
		return
	}
	if voiceMessageHandler(s, m, "jamoin", "sounds/ja_moin.mp3", false) {
		return
	}
	if voiceMessageHandler(s, m, "megalovania", "sounds/megalovania.mp3", false) {
		return
	}
	if voiceMessageHandler(s, m, "ough", "sounds/ough.mp3", false) {
		return
	}
	if voiceMessageHandler(s, m, "yooo", "sounds/yooooooooooo.mp3", false) {
		return
	}
	if voiceMessageHandler(s, m, "haha", "sounds/haha.mp3", false) {
		return
	}
	if voiceMessageHandler(s, m, "letsgo", "sounds/letsgo.mp3", false) {
		return
	}

	// If the message is "bing" reply with "Bong!"
	if strings.ToLower(m.Content) == "bing" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Bong!")

		return
	}

	// If the message is "bong" reply with "Bing!"
	if strings.ToLower(m.Content) == "bong" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Bing!")

		return
	}

}

func voiceMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate, message string, filename string, handleOnlyFullMessage bool) bool {
	m.Content = strings.ToLower(m.Content)
	if m.Content == message {
		playSound(s, m, filename, message, true)
		return true
	}
	if strings.Contains(m.Content, message) && !handleOnlyFullMessage {
		playSound(s, m, filename, message, false)
		return true
	}
	return false
}

func playSound(s *discordgo.Session, m *discordgo.MessageCreate, filename string, commandString string, sendErrMsg bool) {

	// s.Guild() funktioniert hier nicht, weil die VoiceStates nur in "state-cached guilds" verfügbar sind,
	// deshalb s.State.Guild()
	st, _ := s.State.Guild(m.GuildID)

	vc := func() (vc *discordgo.Channel) {
		for _, state := range st.VoiceStates {
			if state.UserID == m.Author.ID {
				channel, _ := s.State.Channel(state.ChannelID)
				return channel
			}
		}
		return nil
	}()
	if vc == nil {
		if sendErrMsg {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Bischte dumm oder was? Du muss schon in nem Channel sein kek alda")
		}
		return
	}

	if contains(joinedServers, m.GuildID) {
		return
	}
	joinedServers = append(joinedServers, m.GuildID)

	go removeGuildAfterTimeout(m.GuildID)

	dvc, err := s.ChannelVoiceJoin(vc.GuildID, vc.ID, false, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	updateUserCommand(db, m.Author.ID, commandString)
	dgvoice.PlayAudioFile(dvc, filename, make(chan bool))
	_ = dvc.Disconnect()

	joinedServers = remove(joinedServers, m.GuildID)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func remove(s []string, e string) []string {
	for i, v := range s {
		if v == e {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func removeGuildAfterTimeout(gid string) {
	time.Sleep(time.Minute)
	if contains(joinedServers, gid) {
		joinedServers = remove(joinedServers, gid)
	}
}

func updateUserCommand(db *sql.DB, userid string, command string) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	var stmt *sql.Stmt
	if checkUserEntryExists(db, userid, command) { // UPDATE
		stmt, err = tx.Prepare("UPDATE `statistics` SET `count` = `count` + 1 WHERE `userid` = ? AND `command` = ?")
		if err != nil {
			log.Fatal(err)
		}
	} else { // INSERT
		stmt, err = tx.Prepare("INSERT INTO `statistics` (`userid`, `command`, `count`) VALUES (?, ?, 1)")
		if err != nil {
			log.Fatal(err)
		}
	}
	defer func() { _ = stmt.Close() }()
	_, err = stmt.Exec(userid, command)
	if err != nil {
		log.Fatal(err)
	}
	_ = tx.Commit()
}

func checkUserEntryExists(db *sql.DB, userid string, command string) bool {
	stmt, err := db.Prepare("SELECT COUNT(*) AS `count` FROM `statistics` WHERE `userid` = ? AND `command` = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = stmt.Close() }()
	rows, err := stmt.Query(&userid, &command)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = rows.Close() }()
	rows.Next()
	var count int
	err = rows.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count == 1
}
