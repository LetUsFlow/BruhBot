package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tcolgate/mp3"
)

var joinedServers []string
var db *sql.DB
var sounds []Sound

func main() {

	// Load configuration file
	content, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
		return
	}
	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
		return
	}

	// Connect to statistics database
	db, err = sql.Open("sqlite3", "./statistics.sqlite")
	if err != nil {
		log.Fatal("Error opening database", err)
	}
	defer func() { _ = db.Close() }()

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection,", err)
		return
	}

	// creating commands
	for i := 21; i >= 1; i-- {
		sounds = append(sounds, Sound{fmt.Sprintf("%s%d", "moan", i), fmt.Sprintf("%s%d%s", "sounds/moans/moan", i, ".mp3"), false, time.Minute})
	}

	sounds = append(sounds,
		Sound{"mo", "sounds/mo.mp3", true, time.Minute},
		Sound{"ma", "sounds/ma.mp3", true, time.Minute},
		Sound{"ginf", "sounds/ginf.mp3", false, time.Minute},
		Sound{"depression", "sounds/warum_bin_ich_so_fröhlich.mp3", false, time.Minute},
		Sound{"teams", "sounds/teams.mp3", false, time.Minute},
		Sound{"okay", "sounds/okay.mp3", false, time.Minute},
		Sound{"yeet", "sounds/yeet.mp3", false, time.Minute},
		Sound{"marcos", "sounds/marcos.mp3", false, time.Minute},
		Sound{"outlook", "sounds/outlook.mp3", false, time.Minute},
		Sound{"bonk", "sounds/bonk.mp3", false, time.Minute},
		Sound{"moan", "sounds/bonk.mp3", false, time.Minute},
		Sound{"bruh", "sounds/bruh.mp3", false, time.Minute},
		Sound{"bann", "sounds/ban_den_weg.mp3", false, time.Minute},
		Sound{"jamoin", "sounds/ja_moin.mp3", false, time.Minute},
		Sound{"megalovania", "sounds/megalovania.mp3", false, time.Minute},
		Sound{"ough", "sounds/ough.mp3", false, time.Minute},
		Sound{"yooo", "sounds/yooooooooooo.mp3", false, time.Minute},
		Sound{"haha", "sounds/haha.mp3", false, time.Minute},
		Sound{"letsgo", "sounds/letsgo.mp3", false, time.Minute},
		Sound{"lugner", "sounds/minze.mp3", false, time.Minute},
		Sound{"minze", "sounds/minze.mp3", false, time.Minute},
		Sound{"electro", "sounds/electroboom.mp3", false, time.Minute},
		Sound{"boom", "sounds/full_bridge_rectifier_song.mp3", false, time.Minute},
		Sound{"rectify", "sounds/full_bridge_rectifier_song.mp3", false, time.Minute},
		Sound{"fichtl", "sounds/fichtl.mp3", false, time.Minute},
		Sound{"ara", "sounds/ara_ara.mp3", false, time.Minute},
	)

	for i, sound := range sounds {
		t := 0.0

		r, err := os.Open(sound.filename)
		if err != nil {
			fmt.Println(err)
			return
		}

		d := mp3.NewDecoder(r)
		var f mp3.Frame
		skipped := 0

		for {
			if err := d.Decode(&f, &skipped); err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println(err)
				return
			}

			t = t + f.Duration().Seconds()
		}

		err = r.Close()
		if err != nil {
			return
		}

		sounds[i].duration = time.Duration(t)*time.Second + time.Second
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

	if m.Content == "bruhelp" {
		var helpstring = sounds[0].message

		for i := 1; i < len(sounds); i++ {
			helpstring += fmt.Sprintf(", %s", sounds[i].message)
		}

		_, _ = s.ChannelMessageSend(m.ChannelID, helpstring)
		return
	}

	// check for commands
	for _, sound := range sounds {
		if voiceMessageHandler(s, m, sound) {
			return
		}
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

func voiceMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate, sound Sound) bool {
	message := sound.message
	filename := sound.filename

	m.Content = strings.ToLower(m.Content)
	if m.Content == message {
		playSound(s, m, filename, message, true, sound)
		return true
	}
	if strings.Contains(m.Content, message) && !sound.handleOnlyFullMessage {
		playSound(s, m, filename, message, false, sound)
		return true
	}
	return false
}

func playSound(s *discordgo.Session, m *discordgo.MessageCreate, filename string, commandString string, sendErrMsg bool, sound Sound) {

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

	go removeGuildAfterTimeout(m.GuildID, sound.duration)

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

func removeGuildAfterTimeout(gid string, duration time.Duration) {
	time.Sleep(duration)
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

type Sound struct {
	message, filename     string
	handleOnlyFullMessage bool
	duration              time.Duration
}

type Config struct {
	Token string
}
